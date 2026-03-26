/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

type nodeFirewallOptionsResourceShort struct{ nodeFirewallOptionsResource }

var (
	_ resource.Resource                = &nodeFirewallOptionsResourceShort{}
	_ resource.ResourceWithConfigure   = &nodeFirewallOptionsResourceShort{}
	_ resource.ResourceWithImportState = &nodeFirewallOptionsResourceShort{}
	_ resource.ResourceWithMoveState   = &nodeFirewallOptionsResourceShort{}
)

// NewShortNodeFirewallOptionsResource creates the short-name alias proxmox_node_firewall.
func NewShortNodeFirewallOptionsResource() resource.Resource {
	return &nodeFirewallOptionsResourceShort{}
}

func (r *nodeFirewallOptionsResourceShort) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_node_firewall"
}

func (r *nodeFirewallOptionsResourceShort) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	r.nodeFirewallOptionsResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *nodeFirewallOptionsResourceShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.nodeFirewallOptionsResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_node_firewall", &schemaResp.Schema),
	}
}
