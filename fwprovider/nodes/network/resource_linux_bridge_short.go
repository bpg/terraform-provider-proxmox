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

type linuxBridgeResourceShort struct{ linuxBridgeResource }

var (
	_ resource.Resource                = &linuxBridgeResourceShort{}
	_ resource.ResourceWithConfigure   = &linuxBridgeResourceShort{}
	_ resource.ResourceWithImportState = &linuxBridgeResourceShort{}
	_ resource.ResourceWithMoveState   = &linuxBridgeResourceShort{}
)

// NewShortLinuxBridgeResource creates the short-name alias proxmox_network_linux_bridge.
func NewShortLinuxBridgeResource() resource.Resource {
	return &linuxBridgeResourceShort{}
}

func (r *linuxBridgeResourceShort) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_network_linux_bridge"
}

func (r *linuxBridgeResourceShort) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	r.linuxBridgeResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *linuxBridgeResourceShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.linuxBridgeResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_network_linux_bridge", &schemaResp.Schema),
	}
}
