/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

//nolint:dupl // short-name alias wrappers share the same boilerplate structure by design
package apt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

type repositoryResourceShort struct{ repositoryResource }

var (
	_ resource.Resource                = &repositoryResourceShort{}
	_ resource.ResourceWithConfigure   = &repositoryResourceShort{}
	_ resource.ResourceWithImportState = &repositoryResourceShort{}
	_ resource.ResourceWithMoveState   = &repositoryResourceShort{}
)

// NewShortRepositoryResource creates the short-name alias proxmox_apt_repository.
func NewShortRepositoryResource() resource.Resource {
	return &repositoryResourceShort{}
}

func (r *repositoryResourceShort) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_apt_repository"
}

func (r *repositoryResourceShort) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	r.repositoryResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *repositoryResourceShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.repositoryResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_apt_repository", &schemaResp.Schema),
	}
}
