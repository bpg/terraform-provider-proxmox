/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package applier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

var (
	_ resource.Resource              = &sdnApplierShort{}
	_ resource.ResourceWithConfigure = &sdnApplierShort{}
	_ resource.ResourceWithMoveState = &sdnApplierShort{}
)

type sdnApplierShort struct{ Resource }

// NewShortResource creates the short-name alias proxmox_sdn_applier.
func NewShortResource() resource.Resource {
	return &sdnApplierShort{}
}

func (r *sdnApplierShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_applier"
}

func (r *sdnApplierShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.Resource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *sdnApplierShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.Resource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_sdn_applier", &schemaResp.Schema),
	}
}
