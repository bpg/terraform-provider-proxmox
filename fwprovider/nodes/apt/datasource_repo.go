/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package apt

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
)

// Ensure the implementation satisfies the required interfaces.
var (
	_ datasource.DataSource              = &repositoryDataSource{}
	_ datasource.DataSourceWithConfigure = &repositoryDataSource{}
)

// repositoryDataSource is the data source implementation for an APT repository.
type repositoryDataSource struct {
	// client is the Proxmox VE API client.
	client proxmox.Client
}

// Configure adds the provider-configured client to the data source.
func (d *repositoryDataSource) Configure(
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
func (d *repositoryDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_apt_repository"
}

// Read fetches the specified APT repository from the Proxmox VE API.
func (d *repositoryDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var rp modelRepo

	resp.Diagnostics.Append(req.Config.Get(ctx, &rp)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.Node(rp.Node.ValueString()).APT().Repositories().Get(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Could not read APT repository", err.Error())

		return
	}

	resp.Diagnostics.Append(rp.importFromAPI(ctx, data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &rp)...)
}

// Schema defines the schema for the APT repository.
func (d *repositoryDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves an APT repository from a Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			SchemaAttrNameComment: schema.StringAttribute{
				Computed:    true,
				Description: "The associated comment.",
			},
			SchemaAttrNameComponents: schema.ListAttribute{
				Computed:    true,
				Description: "The list of components.",
				ElementType: types.StringType,
			},
			SchemaAttrNameEnabled: schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates the activation status.",
			},
			SchemaAttrNameFilePath: schema.StringAttribute{
				Description: "The absolute path of the source list file that contains this repository.",
				Required:    true,
				Validators: []validator.String{
					validators.NonEmptyString(),
				},
			},
			SchemaAttrNameFileType: schema.StringAttribute{
				Computed:    true,
				Description: "The format of the defining source list file.",
			},
			SchemaAttrNameIndex: schema.Int64Attribute{
				Description: "The index within the defining source list file.",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			SchemaAttrNameNode: schema.StringAttribute{
				Description: "The name of the target Proxmox VE node.",
				Required:    true,
				Validators: []validator.String{
					validators.NonEmptyString(),
				},
			},
			SchemaAttrNamePackageTypes: schema.ListAttribute{
				Computed:    true,
				Description: "The list of package types.",
				ElementType: types.StringType,
			},
			SchemaAttrNameSuites: schema.ListAttribute{
				Computed:    true,
				Description: "The list of package distributions.",
				ElementType: types.StringType,
			},
			SchemaAttrNameTerraformID: attribute.ID("The unique identifier of this APT repository data source."),
			SchemaAttrNameURIs: schema.ListAttribute{
				Computed:    true,
				Description: "The list of repository URIs.",
				ElementType: types.StringType,
			},
		},
	}
}

// NewRepositoryDataSource returns a new data source for an APT repository.
// This is a helper function to simplify the provider implementation.
func NewRepositoryDataSource() datasource.DataSource {
	return &repositoryDataSource{}
}
