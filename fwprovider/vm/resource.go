package vm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

const (
	defaultCreateTimeout = 30 * time.Minute
	defaultReadTimeout   = 5 * time.Minute
	defaultUpdateTimeout = 30 * time.Minute
	defaultDeleteTimeout = 10 * time.Minute

	// these timeouts are for individual PVE operations.
	defaultShutdownTimeout = 5 * time.Minute
)

var (
	_ resource.Resource                = &vmResource{}
	_ resource.ResourceWithConfigure   = &vmResource{}
	_ resource.ResourceWithImportState = &vmResource{}
)

type vmResource struct {
	client proxmox.Client
}

// NewVMResource creates a new resource for managing VMs.
func NewVMResource() resource.Resource {
	return &vmResource{}
}

// Metadata defines the name of the resource.
func (r *vmResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_vm2"
}

func (r *vmResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(proxmox.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got: %T", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *vmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vmModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, d := plan.Timeouts.Create(ctx, defaultCreateTimeout)
	resp.Diagnostics.Append(d...)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// vmID := int(plan.VMID.ValueInt64())
	vmID := int(plan.ID.ValueInt64())
	if vmID == 0 {
		id, err := r.client.Cluster().GetVMID(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Failed to get VM ID", err.Error())
			return
		}

		vmID = *id
	}

	vmAPI := r.client.Node(plan.NodeName.ValueString()).VM(vmID)

	createBody := &vms.CreateRequestBody{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueStringPointer(),
		VMID:        &vmID,
	}

	err := vmAPI.CreateVM(ctx, createBody)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create VM", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// read back the VM from the PVE API to populate computed fields
	exists, err := r.read(ctx, vmAPI, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read VM", err.Error())
	}

	if !exists {
		resp.Diagnostics.AddError("VM does not exist after creation", "")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// set state to the updated plan data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *vmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vmModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, d := state.Timeouts.Read(ctx, defaultReadTimeout)
	resp.Diagnostics.Append(d...)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	vmAPI := r.client.Node(state.NodeName.ValueString()).VM(int(state.ID.ValueInt64()))

	exists, err := r.read(ctx, vmAPI, &state)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read VM", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if !exists {
		tflog.Info(ctx, "VM does not exist, removing from the state", map[string]interface{}{
			"id": state.ID.ValueInt64(),
		})
		resp.State.RemoveResource(ctx)

		return
	}

	// store updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *vmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state vmModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, d := plan.Timeouts.Update(ctx, defaultUpdateTimeout)
	resp.Diagnostics.Append(d...)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	vmAPI := r.client.Node(plan.NodeName.ValueString()).VM(int(plan.ID.ValueInt64()))

	updateBody := &vms.UpdateRequestBody{
		VMID: proxmoxtypes.Int64PtrToIntPtr(plan.ID.ValueInt64Pointer()),
	}

	var errs []error

	del := func(field string) {
		errs = append(errs, updateBody.ToDelete(field))
	}

	if !plan.Description.Equal(state.Description) {
		if plan.Description.IsNull() {
			del("Description")
		} else {
			updateBody.Description = plan.Description.ValueStringPointer()
		}
	}

	if !plan.Name.Equal(state.Name) {
		if plan.Name.IsNull() {
			del("Name")
		} else {
			updateBody.Name = plan.Name.ValueStringPointer()
		}
	}

	err := vmAPI.UpdateVM(ctx, updateBody)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update VM", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// read back the VM from the PVE API to populate computed fields
	exists, err := r.read(ctx, vmAPI, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read VM", err.Error())
	}

	if !exists {
		resp.Diagnostics.AddError("VM does not exist after update", "")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// set state to the updated plan data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *vmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vmModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, d := state.Timeouts.Delete(ctx, defaultDeleteTimeout)
	resp.Diagnostics.Append(d...)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	vmAPI := r.client.Node(state.NodeName.ValueString()).VM(int(state.ID.ValueInt64()))

	// Stop or shut down the virtual machine before deleting it.
	status, err := vmAPI.GetVMStatus(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get VM status", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// stop := d.Get(mkStopOnDestroy).(bool)
	stop := false

	if status.Status != "stopped" {
		if stop {
			if e := vmStop(ctx, vmAPI); e != nil {
				resp.Diagnostics.AddWarning("Failed to stop VM", e.Error())
			}
		} else {
			if e := vmShutdown(ctx, vmAPI); e != nil {
				resp.Diagnostics.AddWarning("Failed to shut down VM", e.Error())
			}
		}
	}

	err = vmAPI.DeleteVM(ctx)
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError("Failed to delete VM", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *vmResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	ctx, cancel := context.WithTimeout(ctx, defaultReadTimeout)
	defer cancel()

	nodeName, vmid, found := strings.Cut(req.ID, "/")

	id, err := strconv.Atoi(vmid)
	if !found || err != nil || id == 0 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: `node_name/id`. Got: %q", req.ID),
		)

		return
	}

	vmAPI := r.client.Node(nodeName).VM(id)

	var ts timeouts.Value

	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("timeouts"), &ts)...)

	if resp.Diagnostics.HasError() {
		return
	}

	state := vmModel{
		ID:       types.Int64Value(int64(id)),
		NodeName: types.StringValue(nodeName),
		Timeouts: ts,
	}

	exists, err := r.read(ctx, vmAPI, &state)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read VM", err.Error())
	}

	if !exists {
		resp.Diagnostics.AddError(fmt.Sprintf("VM %d does not exist on node %s", id, nodeName), "")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// read retrieves the current state of the resource from the API and updates the state.
// Returns false if the resource does not exist, so the caller can remove it from the state if necessary.
func (r *vmResource) read(ctx context.Context, vmAPI *vms.Client, model *vmModel) (bool, error) {
	// Retrieve the entire configuration in order to compare it to the state.
	vmConfig, err := vmAPI.GetVM(ctx)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			tflog.Info(ctx, "VM does not exist, removing from the state", map[string]interface{}{
				"vm_id": vmAPI.VMID,
			})

			return false, nil
		}

		return false, fmt.Errorf("failed to get VM: %w", err)
	}

	vmStatus, err := vmAPI.GetVMStatus(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get VM status: %w", err)
	}

	err = model.updateFromAPI(*vmConfig, *vmStatus)
	if err != nil {
		return false, fmt.Errorf("failed to update VM from API: %w", err)
	}

	return true, nil
}

// Shutdown the VM, then wait for it to actually shut down (it may not be shut down immediately if
// running in HA mode).
func vmShutdown(ctx context.Context, vmAPI *vms.Client) error {
	tflog.Debug(ctx, "Shutting down VM")

	shutdownTimeoutSec := int(defaultShutdownTimeout.Seconds())

	if dl, ok := ctx.Deadline(); ok {
		time.Until(dl)
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
