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

type ociImageResourceShort struct{ ociImageResource }

var (
	_ resource.Resource              = &ociImageResourceShort{}
	_ resource.ResourceWithConfigure = &ociImageResourceShort{}
	_ resource.ResourceWithMoveState = &ociImageResourceShort{}
)

// NewShortOCIImageResource creates the short-name alias proxmox_oci_image.
func NewShortOCIImageResource() resource.Resource {
	return &ociImageResourceShort{}
}

func (r *ociImageResourceShort) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_oci_image"
}

func (r *ociImageResourceShort) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	r.ociImageResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *ociImageResourceShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.ociImageResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_oci_image", &schemaResp.Schema),
	}
}
