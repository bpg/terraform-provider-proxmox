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

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/agent"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cdrom"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/clone"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cpu"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/initialization"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/memory"
	network_device "github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/network_device"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/rng"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/vga"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Model represents the VM model.
//
// Note: for computed fields / blocks we have to use an Object type (or an alias),
// or a custom type in order to hold an unknown value.
type Model struct {
	Agent                            agent.Value          `tfsdk:"agent"`
	CDROM                            cdrom.Value          `tfsdk:"cdrom"`
	Clone                            clone.Value          `tfsdk:"clone"`
	CPU                              cpu.Value            `tfsdk:"cpu"`
	Description                      types.String         `tfsdk:"description"`
	ID                               types.Int64          `tfsdk:"id"`
	Initialization                   initialization.Value `tfsdk:"initialization"`
	Memory                           memory.Value         `tfsdk:"memory"`
	Name                             types.String         `tfsdk:"name"`
	NetworkDevice                    network_device.Value `tfsdk:"network_device"`
	NodeName                         types.String         `tfsdk:"node_name"`
	RNG                              rng.Value            `tfsdk:"rng"`
	Started                          types.Bool           `tfsdk:"started"`
	StopOnDestroy                    types.Bool           `tfsdk:"stop_on_destroy"`
	PurgeOnDestroy                   types.Bool           `tfsdk:"purge_on_destroy"`
	DeleteUnreferencedDisksOnDestroy types.Bool           `tfsdk:"delete_unreferenced_disks_on_destroy"`
	Tags                             stringset.Value      `tfsdk:"tags"`
	Template                         types.Bool           `tfsdk:"template"`
	Timeouts                         timeouts.Value       `tfsdk:"timeouts"`
	VGA                              vga.Value            `tfsdk:"vga"`
}

// DatasourceModel represents the VM datasource model.
// It excludes resource-only lifecycle fields (stop_on_destroy, purge_on_destroy,
// delete_unreferenced_disks_on_destroy) that have no API representation.
type DatasourceModel struct {
	Agent          agent.Value                    `tfsdk:"agent"`
	CDROM          cdrom.Value                    `tfsdk:"cdrom"`
	CPU            cpu.Value                      `tfsdk:"cpu"`
	Description    types.String                   `tfsdk:"description"`
	ID             types.Int64                    `tfsdk:"id"`
	Initialization initialization.DataSourceValue `tfsdk:"initialization"`
	Memory         memory.Value                   `tfsdk:"memory"`
	Name           types.String                   `tfsdk:"name"`
	NetworkDevice  network_device.Value           `tfsdk:"network_device"`
	NodeName       types.String                   `tfsdk:"node_name"`
	RNG            rng.Value                      `tfsdk:"rng"`
	Started        types.Bool                     `tfsdk:"started"`
	Status         types.String                   `tfsdk:"status"`
	Tags           stringset.Value                `tfsdk:"tags"`
	Template       types.Bool                     `tfsdk:"template"`
	Timeouts       timeouts.Value                 `tfsdk:"timeouts"`
	VGA            vga.Value                      `tfsdk:"vga"`
}

// readForDatasource retrieves the VM from the API and populates the datasource model.
// Returns false if the resource does not exist.
func readForDatasource(ctx context.Context, client proxmox.Client, model *DatasourceModel, diags *diag.Diagnostics) bool {
	vmAPI := client.Node(model.NodeName.ValueString()).VM(int(model.ID.ValueInt64()))

	config, err := vmAPI.GetVM(ctx)
	if err != nil {
		if !errors.Is(err, api.ErrResourceDoesNotExist) {
			diags.AddError(fmt.Sprintf("Unable to Read VM %d", model.ID.ValueInt64()), err.Error())
		}

		return false
	}

	status, err := vmAPI.GetVMStatus(ctx)
	if err != nil {
		diags.AddError(fmt.Sprintf("Unable to Read VM %d Status", model.ID.ValueInt64()), err.Error())
		return false
	}

	if status.VMID == nil {
		diags.AddError(
			fmt.Sprintf("Unable to Read VM %d Status", model.ID.ValueInt64()),
			"VM ID is missing in the status API response",
		)

		return false
	}

	model.ID = types.Int64Value(int64(*status.VMID))
	model.Status = types.StringValue(status.Status)
	model.Started = types.BoolValue(status.Status == "running")
	model.Tags = stringset.NewValueString(config.Tags, diags)

	model.Description = attribute.StringValueFromPtr(config.Description)

	model.Name = attribute.StringValueFromPtr(config.Name)

	model.Template = attribute.BoolValueFromCustomBoolPtr(config.Template)

	model.Agent = agent.NewValue(ctx, config, diags)
	model.CPU = cpu.NewValue(ctx, config, diags)
	model.Memory = memory.NewValue(ctx, config, diags)
	model.NetworkDevice = network_device.NewValue(ctx, config, diags)
	model.RNG = rng.NewValue(ctx, config, diags)
	model.VGA = vga.NewValue(ctx, config, diags)
	model.CDROM = cdrom.NewValue(ctx, config, diags)
	model.Initialization = initialization.NewDataSourceValue(ctx, config, diags)

	return true
}

// read retrieves the current state of the resource from the API and updates the state.
// Returns false if the resource does not exist, so the caller can remove it from the state if necessary.
func read(ctx context.Context, client proxmox.Client, model *Model, diags *diag.Diagnostics) bool {
	vmAPI := client.Node(model.NodeName.ValueString()).VM(int(model.ID.ValueInt64()))

	// Retrieve the entire configuration in order to compare it to the state.
	config, err := vmAPI.GetVM(ctx)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			tflog.Info(ctx, "VM does not exist, removing from the state", map[string]any{
				"vm_id": vmAPI.VMID,
			})
		} else {
			diags.AddError(fmt.Sprintf("Unable to Read VM %d", model.ID.ValueInt64()), err.Error())
		}

		return false
	}

	status, err := vmAPI.GetVMStatus(ctx)
	if err != nil {
		diags.AddError(fmt.Sprintf("Unable to Read VM %d Status", model.ID.ValueInt64()), err.Error())
		return false
	}

	if status.VMID == nil {
		diags.AddError(
			fmt.Sprintf("Unable to Read VM %d Status", model.ID.ValueInt64()),
			"VM ID is missing in the status API response",
		)

		return false
	}

	model.ID = types.Int64Value(int64(*status.VMID))
	model.Started = types.BoolValue(status.Status == "running")

	// Optional fields can be removed from the model, use StringPointerValue to handle removal on nil
	model.Description = types.StringPointerValue(config.Description)
	model.Name = types.StringPointerValue(config.Name)
	model.Tags = stringset.NewValueString(config.Tags, diags)
	model.Template = types.BoolPointerValue(config.Template.PointerBool())

	// Blocks
	model.Agent = agent.NewValue(ctx, config, diags)
	model.CPU = cpu.NewValue(ctx, config, diags)
	model.Memory = memory.NewValue(ctx, config, diags)
	model.NetworkDevice = network_device.NewValue(ctx, config, diags)
	model.RNG = rng.NewValue(ctx, config, diags)
	model.VGA = vga.NewValue(ctx, config, diags)
	model.CDROM = cdrom.NewValue(ctx, config, diags)
	model.Initialization = initialization.NewValue(ctx, config, diags)

	// clone is write-only at create time; preserve what was in state on subsequent reads.
	// model.Clone is left as-is (UseStateForUnknown handles it via the plan modifier).

	return true
}
