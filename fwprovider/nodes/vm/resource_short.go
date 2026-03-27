/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

const shortResourceTypeName = "proxmox_vm"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &resourceShort{}
	_ resource.ResourceWithConfigure   = &resourceShort{}
	_ resource.ResourceWithImportState = &resourceShort{}
	_ resource.ResourceWithMoveState   = &resourceShort{}
)

// resourceShort is the short-name alias for the VM2 resource (ADR-007).
type resourceShort struct {
	Resource
}

// NewShortResource creates a new short-named VM2 resource.
func NewShortResource() resource.Resource {
	return &resourceShort{}
}

// Metadata defines the short resource type name.
func (r *resourceShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = shortResourceTypeName
}

// Schema returns the schema with no deprecation message (this is the canonical name).
func (r *resourceShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.Resource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

// MoveState supports migrating from the old long-name resource.
func (r *resourceShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse

	r.Resource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_vm2", &schemaResp.Schema),
	}
}
