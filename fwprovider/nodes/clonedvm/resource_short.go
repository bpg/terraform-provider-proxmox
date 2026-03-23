/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package clonedvm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

type resourceShort struct{ Resource }

var (
	_ resource.Resource                = &resourceShort{}
	_ resource.ResourceWithConfigure   = &resourceShort{}
	_ resource.ResourceWithImportState = &resourceShort{}
	_ resource.ResourceWithMoveState   = &resourceShort{}
)

// NewShortResource creates the short-name alias proxmox_cloned_vm.
func NewShortResource() resource.Resource {
	return &resourceShort{}
}

func (r *resourceShort) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_cloned_vm"
}

func (r *resourceShort) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	r.Resource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *resourceShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.Resource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_cloned_vm", &schemaResp.Schema),
	}
}
