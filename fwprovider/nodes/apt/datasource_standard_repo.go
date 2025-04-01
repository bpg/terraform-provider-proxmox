/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package apt

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types/nodes/apt"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
)

// Ensure the implementation satisfies the required interfaces.
var (
	_ datasource.DataSource              = &standardRepositoryDataSource{}
	_ datasource.DataSourceWithConfigure = &standardRepositoryDataSource{}
)

// standardRepositoryDataSource is the data source implementation for an APT standard repository.
type standardRepositoryDataSource struct {
	// client is the Proxmox VE API client.
	client proxmox.Client
}

// Configure adds the provider-configured client to the data source.
func (d *standardRepositoryDataSource) Configure(
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

	d.client = cfg.Client
}

// Metadata returns the data source type name.
func (d *standardRepositoryDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_apt_standard_repository"
}

// Read fetches the specified APT standard repository from the Proxmox VE API.
func (d *standardRepositoryDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var srp modelStandardRepo

	resp.Diagnostics.Append(req.Config.Get(ctx, &srp)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.Node(srp.Node.ValueString()).APT().Repositories().Get(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Could not read APT standard repository", err.Error())

		return
	}

	srp.importFromAPI(ctx, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &srp)...)
}

// Schema defines the schema for the APT standard repository.
func (d *standardRepositoryDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves an APT standard repository from a Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			SchemaAttrNameStandardDescription: schema.StringAttribute{
				Computed:    true,
				Description: "The description of the APT standard repository.",
			},
			SchemaAttrNameFilePath: schema.StringAttribute{
				Computed:    true,
				Description: "The absolute path of the source list file that contains this standard repository.",
			},
			SchemaAttrNameStandardHandle: schema.StringAttribute{
				CustomType:  customtypes.StandardRepoHandleType{},
				Description: "The handle of the APT standard repository.",
				Required:    true,
			},
			SchemaAttrNameIndex: schema.Int64Attribute{
				Computed:    true,
				Description: "The index within the defining source list file.",
			},
			SchemaAttrNameStandardName: schema.StringAttribute{
				Computed:    true,
				Description: "The name of the APT standard repository.",
			},
			SchemaAttrNameNode: schema.StringAttribute{
				Description: "The name of the target Proxmox VE node.",
				Required:    true,
				Validators: []validator.String{
					validators.NonEmptyString(),
				},
			},
			SchemaAttrNameStandardStatus: schema.Int64Attribute{
				Computed:    true,
				Description: "Indicates the activation status.",
			},
			SchemaAttrNameTerraformID: attribute.ResourceID(
				"The unique identifier of this APT standard repository data source.",
			),
		},
	}
}

// NewStandardRepositoryDataSource returns a new data source for an APT standard repository.
// This is a helper function to simplify the provider implementation.
func NewStandardRepositoryDataSource() datasource.DataSource {
	return &standardRepositoryDataSource{}
}
