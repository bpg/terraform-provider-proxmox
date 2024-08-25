/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package acme

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/acme/plugins"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &acmePluginsDatasource{}
	_ datasource.DataSourceWithConfigure = &acmePluginsDatasource{}
)

// NewACMEPluginsDataSource is a helper function to simplify the provider implementation.
func NewACMEPluginsDataSource() datasource.DataSource {
	return &acmePluginsDatasource{}
}

// acmePluginsDatasource is the data source implementation for ACME plugins.
type acmePluginsDatasource struct {
	client *plugins.Client
}

// Metadata returns the data source type name.
func (d *acmePluginsDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_acme_plugins"
}

// Schema returns the schema for the data source.
func (d *acmePluginsDatasource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of ACME plugins.",
		Attributes: map[string]schema.Attribute{
			"plugins": schema.ListNestedAttribute{
				Description: "List of ACME plugins",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"api": schema.StringAttribute{
							Description: "API plugin name.",
							Computed:    true,
						},
						"data": schema.MapAttribute{
							Description: "DNS plugin data.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"digest": schema.StringAttribute{
							Description: "Prevent changes if current configuration file has a different digest. " +
								"This can be used to prevent concurrent modifications.",
							Computed: true,
						},
						"plugin": schema.StringAttribute{
							Description: "ACME Plugin ID name.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "ACME challenge type (dns, standalone).",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("dns", "standalone"),
							},
						},
						"validation_delay": schema.Int64Attribute{
							Description: "Extra delay in seconds to wait before requesting validation. " +
								"Allows to cope with a long TTL of DNS records (0 - 172800).",
							Computed: true,
							Validators: []validator.Int64{
								int64validator.Between(0, 172800),
							},
						},
					},
				},
				Computed: true,
			},
		},
	}
}

// Configure adds the provider-configured client to the data source.
func (d *acmePluginsDatasource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(proxmox.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got: %T",
				req.ProviderData),
		)

		return
	}

	d.client = client.Cluster().ACME().Plugins()
}

// Read fetches the list of ACME plugins from the Proxmox cluster.
func (d *acmePluginsDatasource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state acmePluginsModel

	list, err := d.client.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read ACME plugins",
			err.Error(),
		)

		return
	}

	for _, plugin := range list {
		mapValue, diags := types.MapValueFrom(ctx, types.StringType, plugin.Data)
		resp.Diagnostics.Append(diags...)

		state.Plugins = append(state.Plugins, acmePluginModel{
			baseACMEPluginModel: baseACMEPluginModel{
				API:             types.StringValue(plugin.API),
				Data:            mapValue,
				Digest:          types.StringValue(plugin.Digest),
				Plugin:          types.StringValue(plugin.Plugin),
				ValidationDelay: types.Int64Value(plugin.ValidationDelay),
			},
			Type: types.StringValue(plugin.Type),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
