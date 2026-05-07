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
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cdrom"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cpu"
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

	// set state to the updated plan data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *Resource) create(ctx context.Context, plan Model, diags *diag.Diagnostics) {
	createBody := &vms.CreateRequestBody{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueStringPointer(),
		Tags:        plan.Tags.ValueStringPointer(ctx, diags),
		VMID:        int(plan.ID.ValueInt64()),
	}

	// fill out create body fields with values from other resource blocks
	cdrom.FillCreateBody(ctx, plan.CDROM, createBody, diags)
	cpu.FillCreateBody(ctx, plan.CPU, createBody, diags)
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

	// Convert to template if requested
	if !plan.Template.IsNull() && plan.Template.ValueBool() {
		tflog.Info(ctx, fmt.Sprintf("Converting VM %d to template", plan.ID.ValueInt64()))

		vmAPI = r.client.Node(plan.NodeName.ValueString()).VM(int(plan.ID.ValueInt64()))
		vmAPI.ConvertToTemplate(ctx).AddDiags(diags, fmt.Sprintf("Unable to Convert VM %d to Template", plan.ID.ValueInt64()))
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
	cdrom.FillUpdateBody(ctx, plan.CDROM, state.CDROM, updateBody, diags)
	cpu.FillUpdateBody(ctx, plan.CPU, state.CPU, updateBody, diags)
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

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
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
