/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pools

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

const shortPoolMembershipTypeName = "proxmox_pool_membership"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                     = (*poolMembershipShortResource)(nil)
	_ resource.ResourceWithConfigure        = (*poolMembershipShortResource)(nil)
	_ resource.ResourceWithImportState      = (*poolMembershipShortResource)(nil)
	_ resource.ResourceWithConfigValidators = (*poolMembershipShortResource)(nil)
	_ resource.ResourceWithMoveState        = (*poolMembershipShortResource)(nil)
)

// poolMembershipShortResource is the short-name alias for the pool membership resource (ADR-007).
type poolMembershipShortResource struct {
	poolMembershipResource
}

// NewShortPoolMembershipResource creates a new short-named pool membership resource.
func NewShortPoolMembershipResource() resource.Resource {
	return &poolMembershipShortResource{}
}

// Metadata defines the short resource type name.
func (r *poolMembershipShortResource) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = shortPoolMembershipTypeName
}

// Schema returns the schema with no deprecation message (this is the canonical name).
func (r *poolMembershipShortResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	r.poolMembershipResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

// MoveState supports migrating from the old long-name resource.
func (r *poolMembershipShortResource) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse

	r.poolMembershipResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_pool_membership", &schemaResp.Schema),
	}
}
