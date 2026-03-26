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

type standardRepositoryResourceShort struct{ standardRepositoryResource }

var (
	_ resource.Resource                = &standardRepositoryResourceShort{}
	_ resource.ResourceWithConfigure   = &standardRepositoryResourceShort{}
	_ resource.ResourceWithImportState = &standardRepositoryResourceShort{}
	_ resource.ResourceWithMoveState   = &standardRepositoryResourceShort{}
)

// NewShortStandardRepositoryResource creates the short-name alias proxmox_apt_standard_repository.
func NewShortStandardRepositoryResource() resource.Resource {
	return &standardRepositoryResourceShort{}
}

func (r *standardRepositoryResourceShort) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_apt_standard_repository"
}

func (r *standardRepositoryResourceShort) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	r.standardRepositoryResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *standardRepositoryResourceShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.standardRepositoryResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_apt_standard_repository", &schemaResp.Schema),
	}
}
