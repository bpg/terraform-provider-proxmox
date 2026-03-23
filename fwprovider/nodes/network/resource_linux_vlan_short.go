/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

type linuxVLANResourceShort struct{ linuxVLANResource }

var (
	_ resource.Resource                = &linuxVLANResourceShort{}
	_ resource.ResourceWithConfigure   = &linuxVLANResourceShort{}
	_ resource.ResourceWithImportState = &linuxVLANResourceShort{}
	_ resource.ResourceWithMoveState   = &linuxVLANResourceShort{}
)

// NewShortLinuxVLANResource creates the short-name alias proxmox_network_linux_vlan.
func NewShortLinuxVLANResource() resource.Resource {
	return &linuxVLANResourceShort{}
}

func (r *linuxVLANResourceShort) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_network_linux_vlan"
}

func (r *linuxVLANResourceShort) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	r.linuxVLANResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *linuxVLANResourceShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.linuxVLANResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_network_linux_vlan", &schemaResp.Schema),
	}
}
