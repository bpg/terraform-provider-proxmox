/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

type acmeCertificateResourceShort struct{ acmeCertificateResource }

var (
	_ resource.Resource                = &acmeCertificateResourceShort{}
	_ resource.ResourceWithConfigure   = &acmeCertificateResourceShort{}
	_ resource.ResourceWithImportState = &acmeCertificateResourceShort{}
	_ resource.ResourceWithMoveState   = &acmeCertificateResourceShort{}
)

// NewShortACMECertificateResource creates the short-name alias proxmox_acme_certificate.
func NewShortACMECertificateResource() resource.Resource {
	return &acmeCertificateResourceShort{}
}

func (r *acmeCertificateResourceShort) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_acme_certificate"
}

func (r *acmeCertificateResourceShort) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	r.acmeCertificateResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *acmeCertificateResourceShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.acmeCertificateResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_acme_certificate", &schemaResp.Schema),
	}
}
