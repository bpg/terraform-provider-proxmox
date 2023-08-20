/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/internal/structure"
	customtypes "github.com/bpg/terraform-provider-proxmox/internal/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	haresources "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/resources"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &haresourceDatasource{}
	_ datasource.DataSourceWithConfigure = &haresourceDatasource{}
)

// NewHAResourceDataSource is a helper function to simplify the provider implementation.
func NewHAResourceDataSource() datasource.DataSource {
	return &haresourceDatasource{}
}

// haresourceDatasource is the data source implementation for High Availability resources.
type haresourceDatasource struct {
	client *haresources.Client
}

// Metadata returns the data source type name.
func (d *haresourceDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_haresource"
}

// Schema returns the schema for the data source.
func (d *haresourceDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of High Availability resources.",
		Attributes: map[string]schema.Attribute{
			"id": structure.IDAttribute(),
			"resource_id": schema.StringAttribute{
				Description: "The identifier of the Proxmox HA resource to read.",
				Required:    true,
				Validators: []validator.String{
					customtypes.HAResourceIDValidator(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of High Availability resource (`vm` or `ct`).",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "The comment associated with this resource.",
				Computed:    true,
			},
			"group": schema.StringAttribute{
				Description: "The identifier of the High Availability group this resource is a member of.",
				Computed:    true,
			},
			"max_relocate": schema.Int64Attribute{
				Description: "The maximal number of relocation attempts.",
				Computed:    true,
			},
			"max_restart": schema.Int64Attribute{
				Description: "The maximal number of restart attempts.",
				Computed:    true,
			},
			"state": schema.StringAttribute{
				Description: "The desired state of the resource.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider-configured client to the data source.
func (d *haresourceDatasource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(proxmox.Client)
	if ok {
		d.client = client.Cluster().HA().Resources()
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData),
		)
	}
}

// Read fetches the specified HA resource.
func (d *haresourceDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data haresourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resID, err := customtypes.ParseHAResourceID(data.ResourceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error parsing Proxmox HA resource identifier",
			fmt.Sprintf("Couldn't parse configuration into a valid HA resource identifier: %s. "+
				"Please report this issue to the provider developers.", err.Error()),
		)

		return
	}

	resource, err := d.client.Get(ctx, resID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read High Availability resource %v", resID),
			err.Error(),
		)

		return
	}

	data.importFromAPI(resource)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
