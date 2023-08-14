/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"context"
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	hagroups "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/groups"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &hagroupResource{}
	_ resource.ResourceWithConfigure   = &hagroupResource{}
	_ resource.ResourceWithImportState = &hagroupResource{}
)

// hagroupResourceModel is the model used to represent a High Availability group.
type hagroupResourceModel struct {
	// Identifier used by Terrraform
	ID types.String `tfsdk:"id"`
}

// NewHAGroupResource creates a new resource for managing Linux Bridge network interfaces.
func NewHAGroupResource() resource.Resource {
	return &hagroupResource{}
}

// hagroupResource contains the resource's internal data.
type hagroupResource struct {
	// The HA groups API client
	client hagroups.Client
}

// Metadata defines the name of the resource.
func (r *hagroupResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_hagroup"
}

// Schema defines the schema for the resource.
func (r *hagroupResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{}
}

// Configure accesses the provider-configured Proxmox API client on behalf of the resource.
func (r *hagroupResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(proxmox.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData),
		)

		return
	}

	r.client = *client.Cluster().HA().Groups()
}

// Create creates a new HA group on the Proxmox cluster.
func (r *hagroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

// Read reads a HA group definition from the Proxmox cluster.
func (r *hagroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates a HA group definition on the Proxmox cluster.
func (r *hagroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes a HA group definition.
func (r *hagroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// ImportState imports a HA group from the Proxmox cluster.
func (r *hagroupResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
}
