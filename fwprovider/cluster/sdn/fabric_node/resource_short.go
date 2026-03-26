/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric_node

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

// --- OpenFabric Node ---

var (
	_ resource.ResourceWithConfigure        = &openFabricNodeShort{}
	_ resource.ResourceWithImportState      = &openFabricNodeShort{}
	_ resource.ResourceWithConfigValidators = &openFabricNodeShort{}
	_ resource.ResourceWithMoveState        = &openFabricNodeShort{}
)

type openFabricNodeShort struct{ *OpenFabricResource }

// NewOpenFabricShortResource creates the short-name alias proxmox_sdn_fabric_node_openfabric.
func NewOpenFabricShortResource() resource.Resource {
	inner := NewOpenFabricResource().(*OpenFabricResource)
	return &openFabricNodeShort{OpenFabricResource: inner}
}

func (r *openFabricNodeShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_fabric_node_openfabric"
}

func (r *openFabricNodeShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.OpenFabricResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *openFabricNodeShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.OpenFabricResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_sdn_fabric_node_openfabric", &schemaResp.Schema),
	}
}

// --- OSPF Node ---

var (
	_ resource.ResourceWithConfigure   = &ospfFabricNodeShort{}
	_ resource.ResourceWithImportState = &ospfFabricNodeShort{}
	_ resource.ResourceWithMoveState   = &ospfFabricNodeShort{}
)

type ospfFabricNodeShort struct{ *OSPFResource }

// NewOSPFShortResource creates the short-name alias proxmox_sdn_fabric_node_ospf.
func NewOSPFShortResource() resource.Resource {
	inner := NewOSPFResource().(*OSPFResource)
	return &ospfFabricNodeShort{OSPFResource: inner}
}

func (r *ospfFabricNodeShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_fabric_node_ospf"
}

func (r *ospfFabricNodeShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.OSPFResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *ospfFabricNodeShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.OSPFResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_sdn_fabric_node_ospf", &schemaResp.Schema),
	}
}
