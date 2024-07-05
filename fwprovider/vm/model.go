/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vm

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/vm/cdrom"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/vm/cpu"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/vm/vga"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Model represents the VM model.
//
// Note: for computed fields / blocks we have to use an Object type (or an alias),
// or a custom type in order to hold an unknown value.
type Model struct {
	Description types.String `tfsdk:"description"`
	CDROM       cdrom.Value  `tfsdk:"cdrom"`
	CPU         cpu.Value    `tfsdk:"cpu"`
	Clone       *struct {
		ID      types.Int64 `tfsdk:"id"`
		Retries types.Int64 `tfsdk:"retries"`
	} `tfsdk:"clone"`
	ID            types.Int64     `tfsdk:"id"`
	Name          types.String    `tfsdk:"name"`
	NodeName      types.String    `tfsdk:"node_name"`
	StopOnDestroy types.Bool      `tfsdk:"stop_on_destroy"`
	Tags          stringset.Value `tfsdk:"tags"`
	Template      types.Bool      `tfsdk:"template"`
	Timeouts      timeouts.Value  `tfsdk:"timeouts"`
	VGA           vga.Value       `tfsdk:"vga"`
}

// read retrieves the current state of the resource from the API and updates the state.
// Returns false if the resource does not exist, so the caller can remove it from the state if necessary.
func read(ctx context.Context, client proxmox.Client, model *Model, diags *diag.Diagnostics) bool {
	vmAPI := client.Node(model.NodeName.ValueString()).VM(int(model.ID.ValueInt64()))

	// Retrieve the entire configuration in order to compare it to the state.
	config, err := vmAPI.GetVM(ctx)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			tflog.Info(ctx, "VM does not exist, removing from the state", map[string]interface{}{
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

	// Optional fields can be removed from the model, use StringPointerValue to handle removal on nil
	model.Description = types.StringPointerValue(config.Description)
	model.Name = types.StringPointerValue(config.Name)
	model.Tags = stringset.NewValue(config.Tags, diags)
	model.Template = types.BoolPointerValue(config.Template.PointerBool())

	// Blocks
	model.CPU = cpu.NewValue(ctx, config, diags)
	model.VGA = vga.NewValue(ctx, config, diags)

	model.CDROM = cdrom.NewValue(ctx, config, diags)

	return true
}
