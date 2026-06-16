/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zfs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	zfsapi "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/disks/zfs"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// draidConfigModel maps the draid_config nested attribute.
type draidConfigModel struct {
	Data   types.Int64 `tfsdk:"data"`
	Spares types.Int64 `tfsdk:"spares"`
}

// zfsPoolModel maps the proxmox_disks_zfs schema.
type zfsPoolModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`
	Name     types.String `tfsdk:"name"`

	// Write-only creation parameters (RequiresReplace).
	Devices     types.List        `tfsdk:"devices"`
	RaidLevel   types.String      `tfsdk:"raidlevel"`
	AShift      types.Int64       `tfsdk:"ashift"`
	Compression types.String      `tfsdk:"compression"`
	DraidConfig *draidConfigModel `tfsdk:"draid_config"`

	// Local-only create-time flag.
	AddStorage types.Bool `tfsdk:"add_storage"`

	// Local-only delete-time flags.
	CleanupConfig types.Bool `tfsdk:"cleanup_config"`
	CleanupDisks  types.Bool `tfsdk:"cleanup_disks"`

	// Computed from the detail API.
	State  types.String `tfsdk:"state"`
	Errors types.String `tfsdk:"errors"`
}

// toCreateBody builds the POST request body for creating a ZFS pool.
func (m *zfsPoolModel) toCreateBody(ctx context.Context, diags *diag.Diagnostics) *zfsapi.CreateRequestBody {
	var deviceList []string

	d := m.Devices.ElementsAs(ctx, &deviceList, false)
	diags.Append(d...)

	if d.HasError() {
		return nil
	}

	body := &zfsapi.CreateRequestBody{
		Name:      m.Name.ValueString(),
		Devices:   strings.Join(deviceList, ","),
		RaidLevel: m.RaidLevel.ValueString(),
	}

	body.AddStorage = attribute.CustomBoolPtrFromValue(m.AddStorage)
	body.AShift = attribute.Int64PtrFromValue(m.AShift)
	body.Compression = attribute.StringPtrFromValue(m.Compression)

	if m.DraidConfig != nil {
		cfg := fmt.Sprintf("data=%d,spares=%d",
			m.DraidConfig.Data.ValueInt64(),
			m.DraidConfig.Spares.ValueInt64(),
		)
		body.DraidConfig = &cfg
	}

	return body
}

// toDeleteParams builds the DELETE query parameters from local-only state.
func (m *zfsPoolModel) toDeleteParams() *zfsapi.DeleteRequestParams {
	return &zfsapi.DeleteRequestParams{
		CleanupConfig: proxmoxtypes.CustomBool(m.CleanupConfig.ValueBool()).Pointer(),
		CleanupDisks:  proxmoxtypes.CustomBool(m.CleanupDisks.ValueBool()).Pointer(),
	}
}

// fromAPI populates computed fields from the detail API response.
// Write-only fields (devices, raidlevel, ashift, compression, draid_config,
// add_storage, cleanup_config, cleanup_disks) are left untouched so plan
// values survive Read without causing "inconsistent result" errors.
func (m *zfsPoolModel) fromAPI(data *zfsapi.GetResponseData) {
	m.ID = types.StringValue(m.NodeName.ValueString() + "/" + data.Name)
	m.Name = types.StringValue(data.Name)
	m.State = types.StringValue(data.State)
	m.Errors = types.StringValue(data.Errors)
}
