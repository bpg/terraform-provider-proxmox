/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

type nodeConfigResourceShort struct{ nodeConfigResource }

var (
	_ resource.Resource                = &nodeConfigResourceShort{}
	_ resource.ResourceWithConfigure   = &nodeConfigResourceShort{}
	_ resource.ResourceWithImportState = &nodeConfigResourceShort{}
	_ resource.ResourceWithMoveState   = &nodeConfigResourceShort{}
)

func NewShortNodeConfigResource() resource.Resource {
	return &nodeConfigResourceShort{}
}

func (r *nodeConfigResourceShort) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_node_config"
}

func (r *nodeConfigResourceShort) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	r.nodeConfigResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *nodeConfigResourceShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.nodeConfigResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_node_config", &schemaResp.Schema),
	}
}
