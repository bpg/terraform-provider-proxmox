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

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/acme/plugins"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &acmePluginDatasource{}
	_ datasource.DataSourceWithConfigure = &acmePluginDatasource{}
)

// NewACMEPluginDataSource is a helper function to simplify the provider implementation.
func NewACMEPluginDataSource() datasource.DataSource {
	return &acmePluginDatasource{}
}

// acmePluginDatasource is the data source implementation for ACME plugin.
type acmePluginDatasource struct {
	client *plugins.Client
}

// Metadata returns the data source type name.
func (d *acmePluginDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_acme_plugin"
}

// Schema returns the schema for the data source.
func (d *acmePluginDatasource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a single ACME plugin by plugin ID name.",
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
				Required:    true,
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
	}
}

// Configure adds the provider-configured client to the data source.
func (d *acmePluginDatasource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.DataSource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected config.DataSource, got: %T", req.ProviderData),
		)

		return
	}

	d.client = cfg.Client.Cluster().ACME().Plugins()
}

// Read fetches the ACME plugin from the Proxmox cluster.
func (d *acmePluginDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state acmePluginModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Plugin.ValueString()

	plugin, err := d.client.Get(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read ACME plugin",
			err.Error(),
		)

		return
	}

	state.API = types.StringValue(plugin.API)

	mapValue, diags := types.MapValueFrom(ctx, types.StringType, plugin.Data)
	resp.Diagnostics.Append(diags...)

	state.Data = mapValue
	state.Digest = types.StringValue(plugin.Digest)
	state.Plugin = types.StringValue(plugin.Plugin)
	state.Type = types.StringValue(plugin.Type)
	state.ValidationDelay = types.Int64Value(plugin.ValidationDelay)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
