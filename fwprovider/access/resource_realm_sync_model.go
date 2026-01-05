/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type realmSyncModel struct {
	ID             types.String `tfsdk:"id"`
	Realm          types.String `tfsdk:"realm"`
	Scope          types.String `tfsdk:"scope"`
	RemoveVanished types.String `tfsdk:"remove_vanished"`
	EnableNew      types.Bool   `tfsdk:"enable_new"`
	Full           types.Bool   `tfsdk:"full"`
	Purge          types.Bool   `tfsdk:"purge"`
	DryRun         types.Bool   `tfsdk:"dry_run"`
}

func (m *realmSyncModel) toSyncRequest() *access.RealmSyncRequestBody {
	body := &access.RealmSyncRequestBody{}

	if !m.Scope.IsNull() {
		body.Scope = m.Scope.ValueStringPointer()
	}

	if !m.RemoveVanished.IsNull() && m.RemoveVanished.ValueString() != "" {
		body.RemoveVanished = m.RemoveVanished.ValueStringPointer()
	}

	if !m.EnableNew.IsNull() {
		body.EnableNew = proxmoxtypes.CustomBoolPtr(m.EnableNew.ValueBoolPointer())
	}

	// Full and Purge are deprecated by Proxmox in favor of RemoveVanished.
	// They are still sent to the API for backward compatibility, but may be
	// removed in future Proxmox versions.
	if !m.Full.IsNull() {
		body.Full = proxmoxtypes.CustomBoolPtr(m.Full.ValueBoolPointer())
	}

	if !m.Purge.IsNull() {
		body.Purge = proxmoxtypes.CustomBoolPtr(m.Purge.ValueBoolPointer())
	}

	if !m.DryRun.IsNull() {
		body.DryRun = proxmoxtypes.CustomBoolPtr(m.DryRun.ValueBoolPointer())
	}

	return body
}
