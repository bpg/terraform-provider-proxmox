/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

// --- OpenFabric ---

var (
	_ resource.Resource                     = &openFabricShort{}
	_ resource.ResourceWithConfigure        = &openFabricShort{}
	_ resource.ResourceWithImportState      = &openFabricShort{}
	_ resource.ResourceWithConfigValidators = &openFabricShort{}
	_ resource.ResourceWithMoveState        = &openFabricShort{}
)

type openFabricShort struct{ *OpenFabricResource }

// NewOpenFabricShortResource creates the short-name alias proxmox_sdn_fabric_openfabric.
func NewOpenFabricShortResource() resource.Resource {
	inner := NewOpenFabricResource().(*OpenFabricResource)
	return &openFabricShort{OpenFabricResource: inner}
}

func (r *openFabricShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_fabric_openfabric"
}

func (r *openFabricShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.OpenFabricResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *openFabricShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.OpenFabricResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_sdn_fabric_openfabric", &schemaResp.Schema),
	}
}

// --- OSPF ---

var (
	_ resource.Resource                = &ospfFabricShort{}
	_ resource.ResourceWithConfigure   = &ospfFabricShort{}
	_ resource.ResourceWithImportState = &ospfFabricShort{}
	_ resource.ResourceWithMoveState   = &ospfFabricShort{}
)

type ospfFabricShort struct{ *OSPFResource }

// NewOSPFShortResource creates the short-name alias proxmox_sdn_fabric_ospf.
func NewOSPFShortResource() resource.Resource {
	inner := NewOSPFResource().(*OSPFResource)
	return &ospfFabricShort{OSPFResource: inner}
}

func (r *ospfFabricShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_fabric_ospf"
}

func (r *ospfFabricShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.OSPFResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *ospfFabricShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.OSPFResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_sdn_fabric_ospf", &schemaResp.Schema),
	}
}
