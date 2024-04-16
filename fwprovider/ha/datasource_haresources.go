/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ha

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/structure"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	haresources "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/resources"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &haResourcesDatasource{}
	_ datasource.DataSourceWithConfigure = &haResourcesDatasource{}
)

// NewHAResourcesDataSource is a helper function to simplify the provider implementation.
func NewHAResourcesDataSource() datasource.DataSource {
	return &haResourcesDatasource{}
}

// haResourcesDatasource is the data source implementation for High Availability resources.
type haResourcesDatasource struct {
	client *haresources.Client
}

// haResourcesModel maps the schema data for the High Availability resources data source.
type haResourcesModel struct {
	// The Terraform resource identifier
	ID types.String `tfsdk:"id"`
	// The type of HA resources to fetch. If unset, all resources will be fetched.
	Type types.String `tfsdk:"type"`
	// The set of HA resource identifiers
	Resources types.Set `tfsdk:"resource_ids"`
}

// Metadata returns the data source type name.
func (d *haResourcesDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_haresources"
}

// Schema returns the schema for the data source.
func (d *haResourcesDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of High Availability resources.",
		Attributes: map[string]schema.Attribute{
			"id": structure.IDAttribute(),
			"type": schema.StringAttribute{
				Description: "The type of High Availability resources to fetch (`vm` or `ct`). All resources " +
					"will be fetched if this option is unset.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ct", "vm"),
				},
			},
			"resource_ids": schema.SetAttribute{
				Description: "The identifiers of the High Availability resources.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider-configured client to the data source.
func (d *haResourcesDatasource) Configure(
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
			fmt.Sprintf("Expected *proxmox.Client, got: %T",
				req.ProviderData),
		)
	}
}

// Read fetches the list of HA resources from the Proxmox cluster then converts it to a list of strings.
func (d *haResourcesDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		data      haResourcesModel
		fetchType *proxmoxtypes.HAResourceType
	)

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Type.IsNull() {
		data.ID = types.StringValue("haresources")
	} else {
		confType, err := proxmoxtypes.ParseHAResourceType(data.Type.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected HA resource type",
				fmt.Sprintf(
					"Couldn't parse configuration into a valid HA resource type: %s. Please report this issue to the "+
						"provider developers.", err.Error(),
				),
			)

			return
		}

		fetchType = &confType
		data.ID = types.StringValue(fmt.Sprintf("haresources:%v", confType))
	}

	list, err := d.client.List(ctx, fetchType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read High Availability resources",
			err.Error(),
		)

		return
	}

	resources := make([]attr.Value, len(list))
	for i, v := range list {
		resources[i] = types.StringValue(v.ID.String())
	}

	resourcesValue, diags := types.SetValue(types.StringType, resources)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Resources = resourcesValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
