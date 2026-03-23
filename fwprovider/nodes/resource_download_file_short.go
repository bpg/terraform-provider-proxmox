/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

//nolint:dupl // short-name alias wrappers share the same boilerplate structure by design
package nodes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

type downloadFileResourceShort struct{ downloadFileResource }

var (
	_ resource.Resource              = &downloadFileResourceShort{}
	_ resource.ResourceWithConfigure = &downloadFileResourceShort{}
	_ resource.ResourceWithMoveState = &downloadFileResourceShort{}
)

// NewShortDownloadFileResource creates the short-name alias proxmox_download_file.
func NewShortDownloadFileResource() resource.Resource {
	return &downloadFileResourceShort{}
}

func (r *downloadFileResourceShort) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_download_file"
}

func (r *downloadFileResourceShort) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	r.downloadFileResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *downloadFileResourceShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.downloadFileResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_download_file", &schemaResp.Schema),
	}
}
