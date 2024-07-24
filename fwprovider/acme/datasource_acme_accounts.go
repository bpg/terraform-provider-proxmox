/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package acme

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/acme/account"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &acmeAccountsDatasource{}
	_ datasource.DataSourceWithConfigure = &acmeAccountsDatasource{}
)

// NewACMEAccountsDataSource is a helper function to simplify the provider implementation.
func NewACMEAccountsDataSource() datasource.DataSource {
	return &acmeAccountsDatasource{}
}

// acmeAccountsDatasource is the data source implementation for ACME accounts.
type acmeAccountsDatasource struct {
	client *account.Client
}

// acmeAccountsModel maps the schema data for the ACME accounts data source.
type acmeAccountsModel struct {
	Accounts types.Set `tfsdk:"accounts"`
}

// Metadata returns the data source type name.
func (d *acmeAccountsDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_acme_accounts"
}

// Schema returns the schema for the data source.
func (d *acmeAccountsDatasource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of ACME accounts.",
		Attributes: map[string]schema.Attribute{
			"accounts": schema.SetAttribute{
				Description: "The identifiers of the ACME accounts.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider-configured client to the data source.
func (d *acmeAccountsDatasource) Configure(
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

// Read fetches the list of ACME Accounts from the Proxmox cluster then converts it to a list of strings.
func (d *acmeAccountsDatasource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state acmeAccountsModel

	list, err := d.client.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read ACME accounts",
			err.Error(),
		)

		return
	}

	accounts := make([]attr.Value, len(list))
	for i, v := range list {
		accounts[i] = types.StringValue(v.Name)
	}

	accountsValue, diags := types.SetValue(types.StringType, accounts)
	resp.Diagnostics.Append(diags...)

	state.Accounts = accountsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
