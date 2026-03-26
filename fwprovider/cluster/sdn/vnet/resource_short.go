/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vnet

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

var (
	_ resource.Resource                = &sdnVnetShort{}
	_ resource.ResourceWithConfigure   = &sdnVnetShort{}
	_ resource.ResourceWithImportState = &sdnVnetShort{}
	_ resource.ResourceWithMoveState   = &sdnVnetShort{}
)

type sdnVnetShort struct{ Resource }

// NewShortResource creates the short-name alias proxmox_sdn_vnet.
func NewShortResource() resource.Resource {
	return &sdnVnetShort{}
}

func (r *sdnVnetShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_vnet"
}

func (r *sdnVnetShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.Resource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *sdnVnetShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.Resource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_sdn_vnet", &schemaResp.Schema),
	}
}
