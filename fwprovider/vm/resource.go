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

	"github.com/bpg/terraform-provider-proxmox/fwprovider/vm/cpu"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/vm/vga"
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
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithConfigure   = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

// Resource implements the resource.Resource interface for managing VMs.
type Resource struct {
	client proxmox.Client
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
		id, err := r.client.Cluster().GetVMID(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Failed to get VM ID", err.Error())
			return
		}

		plan.ID = types.Int64Value(int64(*id))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Clone != nil {
		r.clone(ctx, plan, &resp.Diagnostics)
	} else {
		r.create(ctx, plan, &resp.Diagnostics)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// read back the VM from the PVE API to populate computed fields
	exists := read(ctx, r.client, &plan, &resp.Diagnostics)
	if !exists {
		resp.Diagnostics.AddError("VM does not exist after creation", "")
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
		Template:    proxmoxtypes.CustomBoolPtr(plan.Template.ValueBoolPointer()),
		VMID:        int(plan.ID.ValueInt64()),
	}

	// fill out create body fields with values from other resource blocks
	cpu.FillCreateBody(ctx, plan.CPU, createBody, diags)
	vga.FillCreateBody(ctx, plan.VGA, createBody, diags)

	if diags.HasError() {
		return
	}

	// .VM(0) is used to create a new VM, the VM ID is not used in the API URL
	vmAPI := r.client.Node(plan.NodeName.ValueString()).VM(0)

	err := vmAPI.CreateVM(ctx, createBody)
	if err != nil {
		diags.AddError("Failed to create VM", err.Error())
	}
}

func (r *Resource) clone(ctx context.Context, plan Model, diags *diag.Diagnostics) {
	if plan.Clone == nil {
		diags.AddError("Clone configuration is missing", "")
		return
	}

	sourceID := int(plan.Clone.ID.ValueInt64())
	vmAPI := r.client.Node(plan.NodeName.ValueString()).VM(sourceID)

	// name and description for the clone are optional, but they are not copied from the source VM.
	cloneBody := &vms.CloneRequestBody{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueStringPointer(),
		VMIDNew:     int(plan.ID.ValueInt64()),
	}

	err := vmAPI.CloneVM(ctx, int(plan.Clone.Retries.ValueInt64()), cloneBody)
	if err != nil {
		diags.AddError("Failed to clone VM", err.Error())
	}

	if diags.HasError() {
		return
	}

	// now load the clone's configuration into a temporary model and update what is needed comparing to the plan
	clone := Model{
		ID:          plan.ID,
		CPU:         plan.CPU,
		Name:        plan.Name,
		Description: plan.Description,
		NodeName:    plan.NodeName,
	}

	read(ctx, r.client, &clone, diags)

	if diags.HasError() {
		return
	}

	r.update(ctx, plan, clone, true, diags)
}

//nolint:dupl
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
		tflog.Info(ctx, "VM does not exist, removing from the state", map[string]interface{}{
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

	r.update(ctx, plan, state, false, &resp.Diagnostics)

	// read back the VM from the PVE API to populate computed fields
	exists := read(ctx, r.client, &plan, &resp.Diagnostics)
	if !exists {
		resp.Diagnostics.AddError("VM does not exist after update", "")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// set state to the updated plan data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// update updates the VM with the new configuration.
//
// The isClone parameter is used to determine if the VM is being updated as part of a clone operation.
// During a clone operation, the attributes are copied from the source VM to the clone, so for computed attributes
// that are optional, we need to handle them differently. If they are not set in the clone configuration, we keep the
// source VM's values.
// During the normal update operation, if a computed attribute is not set in the plan, we remove it from the VM, so it
// can assume its default PVE-provided value.
func (r *Resource) update(ctx context.Context, plan, state Model, isClone bool, diags *diag.Diagnostics) {
	vmAPI := r.client.Node(plan.NodeName.ValueString()).VM(int(plan.ID.ValueInt64()))

	updateBody := &vms.UpdateRequestBody{}

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

	// For optional computed fields only:
	// The first condition is for the clone case, where the tags (captured in `state.Tags`)
	// have already been copied from the source VM to the clone during the cloning process.
	// Then, if the clone config does not have tags, we keep the cloned ones in the VM.
	// Otherwise, if the clone config has empty tags we remove them from the VM.
	// And finally, if the clone config has tags we update them in th VM
	if !plan.Tags.Equal(state.Tags) && !plan.Tags.IsUnknown() {
		if plan.Tags.IsNull() || len(plan.Tags.Elements()) == 0 {
			del("Tags")
		} else {
			updateBody.Tags = plan.Tags.ValueStringPointer(ctx, diags)
		}
	}

	cpu.FillUpdateBody(ctx, plan.CPU, state.CPU, updateBody, isClone, diags)
	vga.FillUpdateBody(ctx, plan.VGA, state.VGA, updateBody, isClone, diags)

	if !updateBody.IsEmpty() {
		updateBody.VMID = int(plan.ID.ValueInt64())

		err := vmAPI.UpdateVM(ctx, updateBody)
		if err != nil {
			diags.AddError("Failed to update VM", err.Error())
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

	exists := read(ctx, r.client, &state, &resp.Diagnostics)
	if !exists {
		resp.Diagnostics.AddError(fmt.Sprintf("VM %d does not exist on node %s", id, nodeName), "")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
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
