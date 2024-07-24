/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package acme

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/acme/account"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &acmeAccountDatasource{}
	_ datasource.DataSourceWithConfigure = &acmeAccountDatasource{}
)

// NewACMEAccountDataSource is a helper function to simplify the provider implementation.
func NewACMEAccountDataSource() datasource.DataSource {
	return &acmeAccountDatasource{}
}

// acmeAccountDatasource is the data source implementation for ACME accounts.
type acmeAccountDatasource struct {
	client *account.Client
}

// accountModel is the model used to represent an ACME account.
type accountModel struct {
	// Name is the ACME account config file name.
	Name types.String `tfsdk:"name"`
	// Account is the ACME account information.
	// Account types.Map `tfsdk:"account"` // XXX
	// Directory is the URL of the ACME CA directory endpoint.
	Directory types.String `tfsdk:"directory"`
	// Location is the location of the ACME account.
	Location types.String `tfsdk:"location"`
	// URL of CA TermsOfService - setting this indicates agreement.
	TOS types.String `tfsdk:"tos"`
}

// Metadata returns the data source type name.
func (d *acmeAccountDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_acme_account"
}

// Schema returns the schema for the data source.
func (d *acmeAccountDatasource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific ACME account.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The identifier of the ACME account to read.",
				Optional:    true,
			},
			// "account": schema.MapAttribute{
			// 	Description: "The ACME account information.",
			// 	ElementType: types.StringType,
			// 	Computed:    true,
			// }, // XXX
			"directory": schema.StringAttribute{
				Description: "The directory URL of the ACME account.",
				Computed:    true,
			},
			"location": schema.StringAttribute{
				Description: "The location URL of the ACME account.",
				Computed:    true,
			},
			"tos": schema.StringAttribute{
				Description: "The URL of the terms of service of the ACME account.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider-configured client to the data source.
func (d *acmeAccountDatasource) Configure(
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

	d.client = client.Cluster().ACME().Account()
}

// Read retrieves the ACME account information.
func (d *acmeAccountDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state accountModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name.ValueString()

	account, err := d.client.Get(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read ACME account '%s'", name),
			err.Error(),
		)

		return
	}

	// mapValue, diags := types.MapValueFrom(ctx, types.StringType, account.Account)
	// state.Account = mapValue
	// resp.Diagnostics.Append(diags...)
	// XXX

	state.Directory = types.StringValue(account.Directory)
	state.Location = types.StringValue(account.Location)
	state.TOS = types.StringValue(account.TOS)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
