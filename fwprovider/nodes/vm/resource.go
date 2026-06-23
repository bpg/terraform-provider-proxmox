/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/agent"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cdrom"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/clone"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cpu"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/initialization"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/memory"
	network_device "github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/network_device"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/rng"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/vga"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
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
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithConfigure   = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

// Resource implements the resource.Resource interface for managing VMs.
type Resource struct {
	client      proxmox.Client
	idGenerator cluster.IDGenerator
}

// NewResource creates a new resource for managing VMs.
func NewResource() resource.Resource {
	return &Resource{}
}

// Metadata defines the name of the resource.
func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_vm2"
}

// Configure sets the client for the resource.
func (r *Resource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.Resource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got: %T", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client
	r.idGenerator = cfg.IDGenerator
}

// Create creates a new VM.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, d := plan.Timeouts.Create(ctx, defaultCreateTimeout)
	resp.Diagnostics.Append(d...)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if plan.ID.ValueInt64() == 0 {
		id, err := r.idGenerator.NextID(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Unable to Generate VM ID", err.Error())
			return
		}

		plan.ID = types.Int64Value(int64(id))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Snapshot plan values for Optional-only sub-blocks before create() or read()
	// can overwrite them. Used below to restore them after read() in the clone case.
	cloneWasUsed := !plan.Clone.IsNull() && !plan.Clone.IsUnknown()
	planCPU := plan.CPU
	planCDROM := plan.CDROM
	planInit := plan.Initialization

	r.create(ctx, plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// read back the VM from the PVE API to populate computed fields
	exists := read(ctx, r.client, &plan, &resp.Diagnostics)
	if !exists {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Read VM %d After Creation", plan.ID.ValueInt64()),
			"VM does not exist after creation",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// After clone, read() picks up template-inherited values for Optional-only attributes
	// (e.g. cpu.sockets=1, cpu.numa=false, initialization keys/vendor). These contradict
	// the plan (which had them null) and cause "inconsistent result after apply" errors.
	// Restore plan snapshot so state matches the plan. The inherited values will be
	// cleaned up on the next apply via FillUpdateBody's delete-on-null logic.
	if cloneWasUsed {
		plan.CPU = planCPU
		plan.CDROM = planCDROM
		plan.Initialization = planInit
	}

	// set state to the updated plan data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *Resource) create(ctx context.Context, plan Model, diags *diag.Diagnostics) {
	if !plan.Clone.IsNull() && !plan.Clone.IsUnknown() {
		r.createFromClone(ctx, plan, diags)
		return
	}

	createBody := &vms.CreateRequestBody{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueStringPointer(),
		Tags:        plan.Tags.ValueStringPointer(ctx, diags),
		VMID:        int(plan.ID.ValueInt64()),
	}

	// fill out create body fields with values from other resource blocks
	agent.FillCreateBody(ctx, plan.Agent, createBody, diags)
	cdrom.FillCreateBody(ctx, plan.CDROM, createBody, diags)
	cpu.FillCreateBody(ctx, plan.CPU, createBody, diags)
	initialization.FillCreateBody(ctx, plan.Initialization, createBody, diags)
	memory.FillCreateBody(ctx, plan.Memory, createBody, diags)
	network_device.FillCreateBody(ctx, plan.NetworkDevice, createBody, diags)
	rng.FillCreateBody(ctx, plan.RNG, createBody, diags)
	vga.FillCreateBody(ctx, plan.VGA, createBody, diags)

	if diags.HasError() {
		return
	}

	// .VM(0) is used to create a new VM, the VM ID is not used in the API URL
	vmAPI := r.client.Node(plan.NodeName.ValueString()).VM(0)

	if vmAPI.CreateVM(ctx, createBody).AddDiags(diags, fmt.Sprintf("Unable to Create VM %d", plan.ID.ValueInt64())) {
		return
	}

	r.postCreateVM(ctx, plan, diags)
}

// createFromClone clones an existing VM and applies the desired config via UpdateVM.
func (r *Resource) createFromClone(ctx context.Context, plan Model, diags *diag.Diagnostics) {
	targetNodeName := plan.NodeName.ValueString()
	newVMID := int(plan.ID.ValueInt64())

	sourceNodeName := clone.SourceNodeName(ctx, plan.Clone, targetNodeName, diags)
	sourceVMID := clone.SourceVMID(ctx, plan.Clone, diags)
	retries := clone.Retries(ctx, plan.Clone, diags)
	cloneBody := clone.BuildCloneBody(ctx, plan.Clone, newVMID, diags)

	if diags.HasError() {
		return
	}

	// Set target node if cloning cross-node (PVE handles migration when target != source)
	if sourceNodeName != targetNodeName {
		cloneBody.TargetNodeName = &targetNodeName
	}

	tflog.Info(ctx, fmt.Sprintf("Cloning VM %d from %s/%d", newVMID, sourceNodeName, sourceVMID))

	sourceVMAPI := r.client.Node(sourceNodeName).VM(sourceVMID)

	if sourceVMAPI.CloneVM(ctx, retries, cloneBody).AddDiags(diags, fmt.Sprintf("Unable to Clone VM %d", newVMID)) {
		return
	}

	// Wait for the cloned VM to become ready
	targetVMAPI := r.client.Node(targetNodeName).VM(newVMID)

	if err := targetVMAPI.WaitForVMConfigUnlock(ctx, true); err != nil {
		diags.AddError(fmt.Sprintf("Unable to Clone VM %d", newVMID), err.Error())
		return
	}

	// Apply the desired config on top of the cloned VM via UpdateVM
	updateBody := &vms.CreateRequestBody{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueStringPointer(),
		Tags:        plan.Tags.ValueStringPointer(ctx, diags),
	}

	agent.FillCreateBody(ctx, plan.Agent, updateBody, diags)
	cdrom.FillCreateBody(ctx, plan.CDROM, updateBody, diags)
	cpu.FillCreateBody(ctx, plan.CPU, updateBody, diags)
	initialization.FillCreateBody(ctx, plan.Initialization, updateBody, diags)
	memory.FillCreateBody(ctx, plan.Memory, updateBody, diags)
	network_device.FillCreateBody(ctx, plan.NetworkDevice, updateBody, diags)
	rng.FillCreateBody(ctx, plan.RNG, updateBody, diags)
	vga.FillCreateBody(ctx, plan.VGA, updateBody, diags)

	// For any Optional-only CPU attribute that is null in the plan, explicitly delete it
	// so the template's inherited value does not persist and cause drift on subsequent plans.
	// Only runs when the user set at least one CPU attribute (plan.CPU non-null).
	if !plan.CPU.IsNull() && !plan.CPU.IsUnknown() {
		cpu.AddCloneCleanupDeletes(ctx, plan.CPU, updateBody, diags)
	}

	if diags.HasError() {
		return
	}

	if err := targetVMAPI.UpdateVM(ctx, updateBody); err != nil {
		diags.AddError(fmt.Sprintf("Unable to Configure Cloned VM %d", newVMID), err.Error())
		return
	}

	r.postCreateVM(ctx, plan, diags)
}

// postCreateVM handles common post-create actions: starting the VM and converting to template.
func (r *Resource) postCreateVM(ctx context.Context, plan Model, diags *diag.Diagnostics) {
	vmID := int(plan.ID.ValueInt64())
	vmAPI := r.client.Node(plan.NodeName.ValueString()).VM(vmID)

	// Start failure is a warning, not an error: the VM was created successfully.
	// State is saved with started=false; the update path retries on the next apply.
	if !plan.Started.IsNull() && !plan.Started.IsUnknown() && plan.Started.ValueBool() {
		vmStart(ctx, vmAPI).AddDiagsAsWarnings(diags, fmt.Sprintf("Unable to Start VM %d", vmID))
	}

	// Convert to template if requested
	if !plan.Template.IsNull() && plan.Template.ValueBool() {
		tflog.Info(ctx, fmt.Sprintf("Converting VM %d to template", vmID))
		vmAPI.ConvertToTemplate(ctx).AddDiags(diags, fmt.Sprintf("Unable to Convert VM %d to Template", vmID))
	}
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, d := state.Timeouts.Read(ctx, defaultReadTimeout)
	resp.Diagnostics.Append(d...)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	exists := read(ctx, r.client, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	if !exists {
		tflog.Info(ctx, "VM does not exist, removing from the state", map[string]any{
			"id": state.ID.ValueInt64(),
		})
		resp.State.RemoveResource(ctx)

		return
	}

	// store updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update updates the VM with the new configuration.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, d := plan.Timeouts.Update(ctx, defaultUpdateTimeout)
	resp.Diagnostics.Append(d...)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	r.update(ctx, plan, state, &resp.Diagnostics)

	// read back the VM from the PVE API to populate computed fields
	exists := read(ctx, r.client, &plan, &resp.Diagnostics)
	if !exists {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Read VM %d After Update", plan.ID.ValueInt64()),
			"VM does not exist after update",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// set state to the updated plan data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// update updates the VM with the new configuration.
func (r *Resource) update(ctx context.Context, plan, state Model, diags *diag.Diagnostics) {
	vmAPI := r.client.Node(plan.NodeName.ValueString()).VM(int(plan.ID.ValueInt64()))

	updateBody := &vms.UpdateRequestBody{}

	attribute.CheckDeleteBody(plan.Description, state.Description, updateBody, "description")
	attribute.CheckDeleteBody(plan.Name, state.Name, updateBody, "name")
	attribute.CheckDeleteBody(plan.Tags, state.Tags, updateBody, "tags")

	if attribute.IsDefined(plan.Description) && !plan.Description.Equal(state.Description) {
		updateBody.Description = plan.Description.ValueStringPointer()
	}

	if attribute.IsDefined(plan.Name) && !plan.Name.Equal(state.Name) {
		updateBody.Name = plan.Name.ValueStringPointer()
	}

	if attribute.IsDefined(plan.Tags) && !plan.Tags.Equal(state.Tags) && len(plan.Tags.Elements()) > 0 {
		updateBody.Tags = plan.Tags.ValueStringPointer(ctx, diags)
	}

	// fill out update body fields with values from other resource blocks
	agent.FillUpdateBody(ctx, plan.Agent, state.Agent, updateBody, diags)
	cdrom.FillUpdateBody(ctx, plan.CDROM, state.CDROM, updateBody, diags)
	cpu.FillUpdateBody(ctx, plan.CPU, state.CPU, updateBody, diags)
	initialization.FillUpdateBody(ctx, plan.Initialization, state.Initialization, updateBody, diags)
	memory.FillUpdateBody(ctx, plan.Memory, state.Memory, updateBody, diags)
	network_device.FillUpdateBody(ctx, plan.NetworkDevice, state.NetworkDevice, updateBody, diags)
	rng.FillUpdateBody(ctx, plan.RNG, state.RNG, updateBody, diags)
	vga.FillUpdateBody(ctx, plan.VGA, state.VGA, updateBody, diags)

	if !updateBody.IsEmpty() {
		updateBody.VMID = int(plan.ID.ValueInt64())

		err := vmAPI.UpdateVM(ctx, updateBody)
		if err != nil {
			diags.AddError(fmt.Sprintf("Unable to Update VM %d", plan.ID.ValueInt64()), err.Error())
			return
		}
	}

	// Handle started transitions
	planStarted := !plan.Started.IsNull() && !plan.Started.IsUnknown() && plan.Started.ValueBool()
	stateStarted := !state.Started.IsNull() && !state.Started.IsUnknown() && state.Started.ValueBool()

	if planStarted && !stateStarted {
		vmStart(ctx, vmAPI).AddDiags(diags, fmt.Sprintf("Unable to Start VM %d", plan.ID.ValueInt64()))
	} else if !planStarted && !plan.Started.IsNull() && stateStarted {
		vmShutdown(ctx, vmAPI).AddDiags(diags, fmt.Sprintf("Unable to Shutdown VM %d", plan.ID.ValueInt64()))
	}

	// Handle template conversion if the template flag changed to true
	if !plan.Template.IsNull() && !state.Template.IsNull() {
		oldTemplate := state.Template.ValueBool()
		newTemplate := plan.Template.ValueBool()

		if !oldTemplate && newTemplate {
			tflog.Info(ctx, fmt.Sprintf("Converting VM %d to template", plan.ID.ValueInt64()))

			status, err := vmAPI.GetVMStatus(ctx)
			if err != nil {
				diags.AddError(fmt.Sprintf("Unable to Read VM %d Status", plan.ID.ValueInt64()), err.Error())
				return
			}

			if status != nil && status.Status != "stopped" {
				tflog.Info(ctx, fmt.Sprintf("Stopping VM %d before converting to template", plan.ID.ValueInt64()))

				if vmStop(ctx, vmAPI).AddDiags(diags, fmt.Sprintf("Unable to Stop VM %d Before Template Conversion", plan.ID.ValueInt64())) {
					return
				}
			}

			if vmAPI.ConvertToTemplate(ctx).AddDiags(diags, fmt.Sprintf("Unable to Convert VM %d to Template", plan.ID.ValueInt64())) {
				return
			}
		} else if oldTemplate && !newTemplate {
			diags.AddError(fmt.Sprintf("Unable to Convert Template Back to VM %d", plan.ID.ValueInt64()), "Templates cannot be converted back to regular VMs")
			return
		}
	}
}

// Delete deletes the VM.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Model

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
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Read VM %d Status", state.ID.ValueInt64()), err.Error())
	}

	if resp.Diagnostics.HasError() || status == nil {
		return
	}

	if status.Status != "stopped" {
		// Stop/shutdown failures during delete are non-fatal — reported as warnings.
		if state.StopOnDestroy.ValueBool() {
			vmStop(ctx, vmAPI).AddDiagsAsWarnings(&resp.Diagnostics, fmt.Sprintf("Unable to Stop VM %d", state.ID.ValueInt64()))
		} else {
			vmShutdown(ctx, vmAPI).AddDiagsAsWarnings(&resp.Diagnostics, fmt.Sprintf("Unable to Shutdown VM %d", state.ID.ValueInt64()))
		}
	}

	purge := state.PurgeOnDestroy.ValueBool()
	deleteUnreferencedDisks := state.DeleteUnreferencedDisksOnDestroy.ValueBool()

	result := vmAPI.DeleteVM(ctx, purge, deleteUnreferencedDisks)
	if result.Err() != nil && !errors.Is(result.Err(), api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Delete VM %d", state.ID.ValueInt64()), result.Err().Error())
	}

	for _, w := range result.Warnings() {
		resp.Diagnostics.AddWarning(fmt.Sprintf("Unable to Delete VM %d", state.ID.ValueInt64()), w)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports the state of the VM from the API.
func (r *Resource) ImportState(
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
			"Unable to Import VM",
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

	exists := read(ctx, r.client, &state, &resp.Diagnostics)
	if !exists {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Import VM %d", id),
			fmt.Sprintf("VM does not exist on node %q", nodeName),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// not clear why this is needed, but ImportStateVerify fails without it
	state.StopOnDestroy = types.BoolValue(false)
	state.PurgeOnDestroy = types.BoolValue(true)
	state.DeleteUnreferencedDisksOnDestroy = types.BoolValue(true)
	// clone is write-only; imported VMs have no tracked clone origin.
	state.Clone = clone.NullValue()

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Start the VM and wait for it to reach running status.
func vmStart(ctx context.Context, vmAPI *vms.Client) tasks.TaskResult {
	tflog.Debug(ctx, "Starting VM")

	startTimeoutSec := 300
	if dl, ok := ctx.Deadline(); ok {
		startTimeoutSec = int(time.Until(dl).Seconds())
	}

	result := vmAPI.StartVM(ctx, startTimeoutSec)
	if result.Err() != nil {
		return result
	}

	if err := vmAPI.WaitForVMStatus(ctx, "running"); err != nil {
		return tasks.TaskFailedWithWarnings(
			fmt.Errorf("failed to wait for VM to start: %w", err),
			result.Warnings(),
		)
	}

	return result
}

// Shutdown the VM, then wait for it to actually shut down (it may not be shut down immediately if
// running in HA mode).
func vmShutdown(ctx context.Context, vmAPI *vms.Client) tasks.TaskResult {
	tflog.Debug(ctx, "Shutting down VM")

	shutdownTimeoutSec := int(defaultShutdownTimeout.Seconds())

	if dl, ok := ctx.Deadline(); ok {
		shutdownTimeoutSec = int(time.Until(dl).Seconds())
	}

	result := vmAPI.ShutdownVM(ctx, &vms.ShutdownRequestBody{
		ForceStop: proxmoxtypes.CustomBool(true).Pointer(),
		Timeout:   &shutdownTimeoutSec,
	})
	if result.Err() != nil {
		return result
	}

	if err := vmAPI.WaitForVMStatus(ctx, "stopped"); err != nil {
		return tasks.TaskFailedWithWarnings(
			fmt.Errorf("failed to wait for VM to shut down: %w", err),
			result.Warnings(),
		)
	}

	return result
}

// Forcefully stop the VM, then wait for it to actually stop.
func vmStop(ctx context.Context, vmAPI *vms.Client) tasks.TaskResult {
	tflog.Debug(ctx, "Stopping VM")

	result := vmAPI.StopVM(ctx)
	if result.Err() != nil {
		return result
	}

	if err := vmAPI.WaitForVMStatus(ctx, "stopped"); err != nil {
		return tasks.TaskFailedWithWarnings(
			fmt.Errorf("failed to wait for VM to stop: %w", err),
			result.Warnings(),
		)
	}

	return result
}
