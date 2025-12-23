/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package clonedvm

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cdrom"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cpu"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/memory"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/rng"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/vga"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

const (
	defaultCreateTimeout = 30 * time.Minute
	defaultReadTimeout   = 5 * time.Minute
	defaultUpdateTimeout = 30 * time.Minute
	defaultDeleteTimeout = 10 * time.Minute

	defaultShutdownTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithConfigure   = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

// Resource implements a cloned VM managed resource.
type Resource struct {
	client      proxmox.Client
	idGenerator cluster.IDGenerator
}

// NewResource creates the cloned VM resource.
func NewResource() resource.Resource {
	return &Resource{}
}

// Metadata sets the resource name.
func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_cloned_vm"
}

// Configure wires provider data.
func (r *Resource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	// leave Delete as-is (may be null) to avoid forcing empty object into state

	cfg, ok := req.ProviderData.(config.Resource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected config.Resource, got: %T", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client
	r.idGenerator = cfg.IDGenerator
}

// Create clones and configures the VM.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, d := plan.Timeouts.Create(ctx, defaultCreateTimeout)
	resp.Diagnostics.Append(d...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if plan.ID.ValueInt64() == 0 {
		id, err := r.idGenerator.NextID(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Failed to generate VM ID", err.Error())
			return
		}

		plan.ID = types.Int64Value(int64(id))
	}

	sourceNode := plan.Clone.SourceNodeName.ValueString()
	if sourceNode == "" {
		sourceNode = plan.NodeName.ValueString()
	}

	if sourceNode == "" {
		resp.Diagnostics.AddError("Missing source node", "either clone.source_node_name or node_name must be set")
		return
	}

	if plan.Clone.SourceVMID.IsUnknown() || plan.Clone.SourceVMID.IsNull() {
		resp.Diagnostics.AddError("Missing source_vm_id", "clone.source_vm_id must be set")
		return
	}

	targetNode := plan.NodeName.ValueString()

	cloneBody := buildCloneBody(plan)

	if resp.Diagnostics.HasError() {
		return
	}

	sourceVM := r.client.Node(sourceNode).VM(int(plan.Clone.SourceVMID.ValueInt64()))

	retries := int(plan.Clone.Retries.ValueInt64())
	if retries == 0 {
		retries = 1
	}

	err := sourceVM.CloneVM(ctx, retries, cloneBody)
	if err != nil {
		resp.Diagnostics.AddError("Failed to clone VM", err.Error())
		return
	}

	vmAPI := r.client.Node(targetNode).VM(int(plan.ID.ValueInt64()))

	// Read current VM config to get existing disk file paths before updating
	// This is needed because updating existing disks requires the file path
	currentConfig, err := vmAPI.GetVM(ctx)
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError("Failed to get VM config", err.Error())
		return
	}

	applyManaged(ctx, vmAPI, plan, currentConfig, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	exists := read(ctx, vmAPI, &plan, &resp.Diagnostics)
	if !exists {
		resp.Diagnostics.AddError("VM does not exist after creation", "")
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes state.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, d := state.Timeouts.Read(ctx, defaultReadTimeout)
	resp.Diagnostics.Append(d...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	vmAPI := r.client.Node(state.NodeName.ValueString()).VM(int(state.ID.ValueInt64()))

	exists := read(ctx, vmAPI, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if !exists {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update applies managed config changes.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan Model
	var state Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, d := plan.Timeouts.Update(ctx, defaultUpdateTimeout)
	resp.Diagnostics.Append(d...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if plan.ID.IsUnknown() || plan.ID.IsNull() {
		plan.ID = state.ID
	}

	if plan.NodeName.IsUnknown() || plan.NodeName.IsNull() {
		plan.NodeName = state.NodeName
	}

	vmAPI := r.client.Node(plan.NodeName.ValueString()).VM(int(plan.ID.ValueInt64()))

	// Read current VM config to get existing disk file paths before updating
	currentConfig, err := vmAPI.GetVM(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get VM config", err.Error())
		return
	}

	applyManaged(ctx, vmAPI, plan, currentConfig, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	exists := read(ctx, vmAPI, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if !exists {
		resp.Diagnostics.AddError("VM no longer exists", "")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete removes the VM.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, d := state.Timeouts.Delete(ctx, defaultDeleteTimeout)
	resp.Diagnostics.Append(d...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	vmAPI := r.client.Node(state.NodeName.ValueString()).VM(int(state.ID.ValueInt64()))

	status, err := vmAPI.GetVMStatus(ctx)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Failed to get VM status", err.Error())

		return
	}

	if status != nil && status.Status != "stopped" {
		if state.StopOnDestroy.ValueBool() {
			if e := vmStop(ctx, vmAPI); e != nil {
				resp.Diagnostics.AddWarning("Failed to stop VM", e.Error())
			}
		} else {
			if e := vmShutdown(ctx, vmAPI); e != nil {
				resp.Diagnostics.AddWarning("Failed to shut down VM", e.Error())
			}
		}
	}

	err = vmAPI.DeleteVM(ctx, state.PurgeOnDestroy.ValueBool(), state.DeleteUnreferencedDisksOnDestroy.ValueBool())
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError("Failed to delete VM", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState supports import using node/vmid.
func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	nodeName, vmid, found := strings.Cut(req.ID, "/")

	id, err := strconv.Atoi(vmid)
	if !found || err != nil || id == 0 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: `node_name/id`. Got: %q", req.ID),
		)

		return
	}

	var ts timeouts.Value
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &ts)...)

	if resp.Diagnostics.HasError() {
		return
	}

	state := Model{
		ID:       types.Int64Value(int64(id)),
		NodeName: types.StringValue(nodeName),
		Timeouts: ts,
	}

	vmAPI := r.client.Node(nodeName).VM(id)

	exists := read(ctx, vmAPI, &state, &resp.Diagnostics)
	if !exists {
		resp.Diagnostics.AddError(fmt.Sprintf("VM %d does not exist on node %s", id, nodeName), "")
		return
	}

	state.StopOnDestroy = types.BoolValue(false)
	state.PurgeOnDestroy = types.BoolValue(true)
	state.DeleteUnreferencedDisksOnDestroy = types.BoolValue(true)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func buildCloneBody(plan Model) *vms.CloneRequestBody {
	body := &vms.CloneRequestBody{
		VMIDNew: int(plan.ID.ValueInt64()),
	}

	if !plan.Clone.Full.IsUnknown() && !plan.Clone.Full.IsNull() {
		body.FullCopy = proxmoxtypes.CustomBoolPtr(plan.Clone.Full.ValueBoolPointer())
	}

	if !plan.Clone.TargetDatastore.IsUnknown() && !plan.Clone.TargetDatastore.IsNull() {
		body.TargetStorage = plan.Clone.TargetDatastore.ValueStringPointer()
	}

	if !plan.Clone.TargetFormat.IsUnknown() && !plan.Clone.TargetFormat.IsNull() {
		body.TargetStorageFormat = plan.Clone.TargetFormat.ValueStringPointer()
	}

	if !plan.Clone.SnapshotName.IsUnknown() && !plan.Clone.SnapshotName.IsNull() {
		body.SnapshotName = plan.Clone.SnapshotName.ValueStringPointer()
	}

	if !plan.Clone.PoolID.IsUnknown() && !plan.Clone.PoolID.IsNull() {
		body.PoolID = plan.Clone.PoolID.ValueStringPointer()
	}

	if !plan.Clone.BandwidthLimit.IsUnknown() && !plan.Clone.BandwidthLimit.IsNull() {
		v := int(plan.Clone.BandwidthLimit.ValueInt64())
		body.BandwidthLimit = &v
	}

	if !plan.Description.IsUnknown() && !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}

	if !plan.Name.IsUnknown() && !plan.Name.IsNull() {
		body.Name = plan.Name.ValueStringPointer()
	}

	if !plan.NodeName.IsUnknown() && !plan.NodeName.IsNull() {
		body.TargetNodeName = plan.NodeName.ValueStringPointer()
	}

	return body
}

func applyManaged(ctx context.Context, vmAPI *vms.Client, plan Model, currentConfig *vms.GetResponseData, diags *diag.Diagnostics) {
	updateBody := &vms.UpdateRequestBody{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueStringPointer(),
	}
	if !plan.Tags.IsUnknown() && !plan.Tags.IsNull() {
		updateBody.Tags = plan.Tags.ValueStringPointer(ctx, diags)
	}

	cdrom.FillCreateBody(ctx, plan.CDROM, updateBody, diags)
	cpu.FillCreateBody(ctx, plan.CPU, updateBody, diags)
	memory.FillUpdateBody(ctx, plan.Memory, updateBody, diags)
	rng.FillCreateBody(ctx, plan.RNG, updateBody, diags)
	vga.FillCreateBody(ctx, plan.VGA, updateBody, diags)

	applyNetwork(ctx, plan.Network, updateBody, diags)

	if diags.HasError() {
		return
	}

	diskResizes := applyDisks(plan.Disk, currentConfig, updateBody, diags)

	if diags.HasError() {
		return
	}

	if plan.Delete != nil {
		for _, slot := range plan.Delete.Network {
			if slot.IsUnknown() || slot.IsNull() {
				continue
			}

			updateBody.Delete = append(updateBody.Delete, slot.ValueString())
		}

		for _, slot := range plan.Delete.Disk {
			if slot.IsUnknown() || slot.IsNull() {
				continue
			}

			updateBody.Delete = append(updateBody.Delete, slot.ValueString())
		}
	}

	if !updateBody.IsEmpty() {
		if err := vmAPI.UpdateVM(ctx, updateBody); err != nil {
			diags.AddError("Failed to update VM", err.Error())
			return
		}
	}

	for _, resize := range diskResizes {
		if err := vmAPI.ResizeVMDisk(ctx, resize); err != nil {
			diags.AddError("Failed to resize VM disk", err.Error())
			return
		}
	}
}

func applyNetwork(ctx context.Context, nets map[string]NetworkModel, body *vms.UpdateRequestBody, diags *diag.Diagnostics) {
	if len(nets) == 0 {
		return
	}

	maxIdx := -1

	for slot := range nets {
		idx, ok := slotIndex(slot, "net")
		if !ok {
			diags.AddError("Invalid network slot", fmt.Sprintf("Unsupported network slot key %q", slot))
			return
		}

		if idx > maxIdx {
			maxIdx = idx
		}
	}

	devices := make(vms.CustomNetworkDevices, maxIdx+1)

	for slot, cfg := range nets {
		idx, _ := slotIndex(slot, "net")

		dev := vms.CustomNetworkDevice{
			Enabled: true,
			Model:   "virtio",
		}

		if !cfg.Model.IsUnknown() && !cfg.Model.IsNull() {
			dev.Model = cfg.Model.ValueString()
		}

		if !cfg.Bridge.IsUnknown() && !cfg.Bridge.IsNull() {
			dev.Bridge = cfg.Bridge.ValueStringPointer()
		}

		if !cfg.Firewall.IsUnknown() && !cfg.Firewall.IsNull() {
			dev.Firewall = proxmoxtypes.CustomBoolPtr(cfg.Firewall.ValueBoolPointer())
		}

		if !cfg.LinkDown.IsUnknown() && !cfg.LinkDown.IsNull() {
			dev.LinkDown = proxmoxtypes.CustomBoolPtr(cfg.LinkDown.ValueBoolPointer())
		}

		if !cfg.MACAddress.IsUnknown() && !cfg.MACAddress.IsNull() {
			dev.MACAddress = cfg.MACAddress.ValueStringPointer()
		}

		if !cfg.MTU.IsUnknown() && !cfg.MTU.IsNull() {
			v := int(cfg.MTU.ValueInt64())
			dev.MTU = &v
		}

		if !cfg.Queues.IsUnknown() && !cfg.Queues.IsNull() {
			v := int(cfg.Queues.ValueInt64())
			dev.Queues = &v
		}

		if !cfg.RateLimit.IsUnknown() && !cfg.RateLimit.IsNull() {
			v := cfg.RateLimit.ValueFloat64()
			dev.RateLimit = &v
		}

		if !cfg.Tag.IsUnknown() && !cfg.Tag.IsNull() {
			v := int(cfg.Tag.ValueInt64())
			dev.Tag = &v
		}

		if !cfg.Trunks.IsUnknown() && !cfg.Trunks.IsNull() {
			var trunks []int64
			d := cfg.Trunks.ElementsAs(ctx, &trunks, false)
			diags.Append(d...)

			if !d.HasError() {
				dev.Trunks = make([]int, len(trunks))
				for i, v := range trunks {
					dev.Trunks[i] = int(v)
				}
			}
		}

		devices[idx] = dev
	}

	body.NetworkDevices = devices
}

func applyDisks(
	disks map[string]DiskModel,
	currentConfig *vms.GetResponseData,
	body *vms.UpdateRequestBody,
	diags *diag.Diagnostics,
) []*vms.ResizeDiskRequestBody {
	if len(disks) == 0 {
		return nil
	}

	var resizes []*vms.ResizeDiskRequestBody

	for slot, cfg := range disks {
		currentDevice := (*vms.CustomStorageDevice)(nil)
		if currentConfig != nil {
			currentDevice = currentConfig.StorageDevices[slot]
		}

		isNewDisk := currentDevice == nil

		device := vms.CustomStorageDevice{}

		if !cfg.File.IsUnknown() && !cfg.File.IsNull() {
			device.FileVolume = cfg.File.ValueString()
		} else if !isNewDisk {
			device.FileVolume = currentDevice.FileVolume
		}

		if device.FileVolume == "" {
			if cfg.DatastoreID.IsUnknown() || cfg.DatastoreID.IsNull() {
				diags.AddError("Missing datastore_id", fmt.Sprintf("Disk %q requires either file or datastore_id+size_gb", slot))
				return nil
			}

			if cfg.SizeGB.IsUnknown() || cfg.SizeGB.IsNull() || cfg.SizeGB.ValueInt64() == 0 {
				diags.AddError("Missing size_gb", fmt.Sprintf("Disk %q requires size_gb when file is not provided", slot))
				return nil
			}

			device.FileVolume = fmt.Sprintf("%s:%d", cfg.DatastoreID.ValueString(), cfg.SizeGB.ValueInt64())
		}

		if !cfg.AIO.IsUnknown() && !cfg.AIO.IsNull() {
			device.AIO = cfg.AIO.ValueStringPointer()
		}

		if !cfg.SizeGB.IsUnknown() && !cfg.SizeGB.IsNull() && cfg.SizeGB.ValueInt64() > 0 {
			desiredGB := cfg.SizeGB.ValueInt64()
			if isNewDisk {
				device.Size = proxmoxtypes.DiskSizeFromGigabytes(desiredGB)
			} else if currentDevice.Size != nil {
				currentGB := currentDevice.Size.InGigabytes()
				if desiredGB < currentGB {
					diags.AddError(
						"Disk resize failure",
						fmt.Sprintf("Disk %q: requested size_gb (%d) is lower than current size_gb (%d)", slot, desiredGB, currentGB),
					)

					return nil
				}

				if desiredGB > currentGB {
					resizes = append(resizes, &vms.ResizeDiskRequestBody{
						Disk: slot,
						Size: *proxmoxtypes.DiskSizeFromGigabytes(desiredGB),
					})
				}
			}
		}

		if !cfg.Backup.IsUnknown() && !cfg.Backup.IsNull() {
			device.Backup = proxmoxtypes.CustomBoolPtr(cfg.Backup.ValueBoolPointer())
		}

		if !cfg.Discard.IsUnknown() && !cfg.Discard.IsNull() {
			device.Discard = cfg.Discard.ValueStringPointer()
		}

		if !cfg.Cache.IsUnknown() && !cfg.Cache.IsNull() {
			device.Cache = cfg.Cache.ValueStringPointer()
		}

		if !cfg.IOThread.IsUnknown() && !cfg.IOThread.IsNull() {
			device.IOThread = proxmoxtypes.CustomBoolPtr(cfg.IOThread.ValueBoolPointer())
		}

		if !cfg.Replicate.IsUnknown() && !cfg.Replicate.IsNull() {
			device.Replicate = proxmoxtypes.CustomBoolPtr(cfg.Replicate.ValueBoolPointer())
		}

		if !cfg.Serial.IsUnknown() && !cfg.Serial.IsNull() {
			device.Serial = cfg.Serial.ValueStringPointer()
		}

		if isNewDisk && !cfg.SSD.IsUnknown() && !cfg.SSD.IsNull() {
			device.SSD = proxmoxtypes.CustomBoolPtr(cfg.SSD.ValueBoolPointer())
		}

		if !cfg.ImportFrom.IsUnknown() && !cfg.ImportFrom.IsNull() {
			device.ImportFrom = cfg.ImportFrom.ValueStringPointer()
		}

		if !cfg.Format.IsUnknown() && !cfg.Format.IsNull() {
			device.Format = cfg.Format.ValueStringPointer()
		}

		if !cfg.Media.IsUnknown() && !cfg.Media.IsNull() {
			device.Media = cfg.Media.ValueStringPointer()
		}

		body.AddCustomStorageDevice(slot, device)
	}

	return resizes
}

func read(ctx context.Context, vmAPI *vms.Client, model *Model, diags *diag.Diagnostics) bool {
	config, err := vmAPI.GetVM(ctx)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			tflog.Info(ctx, "VM does not exist, removing from state", map[string]any{
				"vm_id": vmAPI.VMID,
			})
		} else {
			diags.AddError("Failed to get VM", err.Error())
		}

		return false
	}

	status, err := vmAPI.GetVMStatus(ctx)
	if err != nil {
		diags.AddError("Failed to get VM status", err.Error())
		return false
	}

	if status.VMID == nil {
		diags.AddError("VM ID is missing in status API response", "")
		return false
	}

	model.ID = types.Int64Value(int64(*status.VMID))
	if !model.Description.IsUnknown() && !model.Description.IsNull() {
		model.Description = types.StringPointerValue(config.Description)
	}

	if !model.Name.IsUnknown() && !model.Name.IsNull() {
		model.Name = types.StringPointerValue(config.Name)
	}

	if !model.Tags.IsUnknown() && !model.Tags.IsNull() {
		model.Tags = stringset.NewValueString(config.Tags, diags)
	}

	if model.Network != nil {
		for slot := range model.Network {
			model.Network[slot] = readNetworkSlot(config, slot, model.Network[slot])
		}
	}

	if model.Disk != nil {
		for slot := range model.Disk {
			model.Disk[slot] = readDiskSlot(config, slot, model.Disk[slot])
		}
	}

	return true
}

func readNetworkSlot(config *vms.GetResponseData, slot string, current NetworkModel) NetworkModel {
	nm := current

	dev := networkDeviceBySlot(config, slot)
	if dev == nil {
		return current
	}

	if !nm.Bridge.IsUnknown() && !nm.Bridge.IsNull() && dev.Bridge != nil {
		nm.Bridge = types.StringValue(*dev.Bridge)
	}

	if !nm.Firewall.IsUnknown() && !nm.Firewall.IsNull() && dev.Firewall != nil {
		nm.Firewall = types.BoolPointerValue(dev.Firewall.PointerBool())
	}

	if !nm.LinkDown.IsUnknown() && !nm.LinkDown.IsNull() && dev.LinkDown != nil {
		nm.LinkDown = types.BoolPointerValue(dev.LinkDown.PointerBool())
	}

	if !nm.MACAddress.IsUnknown() && !nm.MACAddress.IsNull() && dev.MACAddress != nil {
		nm.MACAddress = types.StringValue(*dev.MACAddress)
	}

	if !nm.Model.IsUnknown() && !nm.Model.IsNull() && dev.Model != "" {
		nm.Model = types.StringValue(dev.Model)
	}

	if !nm.MTU.IsUnknown() && !nm.MTU.IsNull() && dev.MTU != nil {
		nm.MTU = types.Int64Value(int64(*dev.MTU))
	}

	if !nm.Queues.IsUnknown() && !nm.Queues.IsNull() && dev.Queues != nil {
		nm.Queues = types.Int64Value(int64(*dev.Queues))
	}

	if !nm.RateLimit.IsUnknown() && !nm.RateLimit.IsNull() && dev.RateLimit != nil {
		nm.RateLimit = types.Float64Value(*dev.RateLimit)
	}

	if !nm.Tag.IsUnknown() && !nm.Tag.IsNull() && dev.Tag != nil {
		nm.Tag = types.Int64Value(int64(*dev.Tag))
	}

	if !nm.Trunks.IsUnknown() && !nm.Trunks.IsNull() && len(dev.Trunks) > 0 {
		vals := make([]attr.Value, len(dev.Trunks))
		for i, v := range dev.Trunks {
			vals[i] = types.Int64Value(int64(v))
		}

		nm.Trunks = types.SetValueMust(types.Int64Type, vals)
	}

	return nm
}

func readDiskSlot(config *vms.GetResponseData, slot string, current DiskModel) DiskModel {
	dm := current

	device := config.StorageDevices[slot]
	if device == nil {
		return current
	}

	if !dm.File.IsUnknown() && !dm.File.IsNull() {
		dm.File = types.StringValue(device.FileVolume)
	}

	if !dm.DatastoreID.IsUnknown() && !dm.DatastoreID.IsNull() {
		if parts := strings.SplitN(device.FileVolume, ":", 2); len(parts) == 2 && parts[0] != "" {
			dm.DatastoreID = types.StringValue(parts[0])
		}
	}

	if !dm.SizeGB.IsUnknown() && !dm.SizeGB.IsNull() && device.Size != nil {
		dm.SizeGB = types.Int64Value(device.Size.InGigabytes())
	}

	if !dm.Format.IsUnknown() && !dm.Format.IsNull() && device.Format != nil {
		dm.Format = types.StringValue(*device.Format)
	}

	if !dm.AIO.IsUnknown() && !dm.AIO.IsNull() && device.AIO != nil {
		dm.AIO = types.StringValue(*device.AIO)
	}

	if !dm.Backup.IsUnknown() && !dm.Backup.IsNull() && device.Backup != nil {
		dm.Backup = types.BoolPointerValue(device.Backup.PointerBool())
	}

	if !dm.Discard.IsUnknown() && !dm.Discard.IsNull() && device.Discard != nil {
		dm.Discard = types.StringValue(*device.Discard)
	}

	if !dm.Cache.IsUnknown() && !dm.Cache.IsNull() && device.Cache != nil {
		dm.Cache = types.StringValue(*device.Cache)
	}

	if !dm.IOThread.IsUnknown() && !dm.IOThread.IsNull() && device.IOThread != nil {
		dm.IOThread = types.BoolPointerValue(device.IOThread.PointerBool())
	}

	if !dm.Replicate.IsUnknown() && !dm.Replicate.IsNull() && device.Replicate != nil {
		dm.Replicate = types.BoolPointerValue(device.Replicate.PointerBool())
	}

	if !dm.Serial.IsUnknown() && !dm.Serial.IsNull() && device.Serial != nil {
		dm.Serial = types.StringValue(*device.Serial)
	}

	if !dm.SSD.IsUnknown() && !dm.SSD.IsNull() && device.SSD != nil {
		dm.SSD = types.BoolPointerValue(device.SSD.PointerBool())
	}

	if !dm.ImportFrom.IsUnknown() && !dm.ImportFrom.IsNull() && device.ImportFrom != nil {
		dm.ImportFrom = types.StringValue(*device.ImportFrom)
	}

	if !dm.Media.IsUnknown() && !dm.Media.IsNull() && device.Media != nil {
		dm.Media = types.StringValue(*device.Media)
	}

	return dm
}

func networkDeviceBySlot(config *vms.GetResponseData, slot string) *vms.CustomNetworkDevice {
	idx, ok := slotIndex(slot, "net")
	if !ok {
		return nil
	}

	val := reflect.ValueOf(config)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	field := val.FieldByName(fmt.Sprintf("NetworkDevice%d", idx))
	if !field.IsValid() || field.IsNil() {
		return nil
	}

	device, ok := field.Interface().(*vms.CustomNetworkDevice)
	if !ok {
		return nil
	}

	return device
}

func slotIndex(slot string, prefix string) (int, bool) {
	if !strings.HasPrefix(slot, prefix) {
		return 0, false
	}

	idx, err := strconv.Atoi(strings.TrimPrefix(slot, prefix))
	if err != nil || idx < 0 {
		return 0, false
	}

	return idx, true
}

// Shutdown the VM, then wait for it to actually shut down.
func vmShutdown(ctx context.Context, vmAPI *vms.Client) error {
	tflog.Debug(ctx, "Shutting down VM")

	shutdownTimeoutSec := int(defaultShutdownTimeout.Seconds())

	if dl, ok := ctx.Deadline(); ok {
		shutdownTimeoutSec = int(time.Until(dl).Seconds())
	}

	err := vmAPI.ShutdownVM(ctx, &vms.ShutdownRequestBody{
		ForceStop: proxmoxtypes.CustomBool(true).Pointer(),
		Timeout:   &shutdownTimeoutSec,
	})
	if err != nil {
		return fmt.Errorf("failed to initiate shut down of VM: %w", err)
	}

	err = vmAPI.WaitForVMStatus(ctx, "stopped")
	if err != nil {
		return fmt.Errorf("failed to wait for VM to shut down: %w", err)
	}

	return nil
}

// Forcefully stop the VM, then wait for it to actually stop.
func vmStop(ctx context.Context, vmAPI *vms.Client) error {
	tflog.Debug(ctx, "Stopping VM")

	err := vmAPI.StopVM(ctx)
	if err != nil {
		return fmt.Errorf("failed to initiate stop of VM: %w", err)
	}

	err = vmAPI.WaitForVMStatus(ctx, "stopped")
	if err != nil {
		return fmt.Errorf("failed to wait for VM to stop: %w", err)
	}

	return nil
}
