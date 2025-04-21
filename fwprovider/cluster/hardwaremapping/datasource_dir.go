/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package hardwaremapping

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types/hardwaremapping"
	mappings "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/mapping"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/hardwaremapping"
)

// Ensure the implementation satisfies the required interfaces.
var (
	_ datasource.DataSource              = &dirDataSource{}
	_ datasource.DataSourceWithConfigure = &dirDataSource{}
)

// dirDataSource is the data source implementation for a directory mapping.
type dirDataSource struct {
	client *mappings.Client
}

// Configure adds the provider-configured client to the data source.
func (d *dirDataSource) Configure(
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

	d.client = cfg.Client.Cluster().HardwareMapping()
}

// Metadata returns the data source type name.
func (d *dirDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hardware_mapping_dir"
}

// Read fetches the specified directory mapping from the Proxmox VE API.
func (d *dirDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var hm modelDir

	resp.Diagnostics.Append(req.Config.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmID := hm.Name.ValueString()
	// Ensure to keep both in sync since the name represents the ID.
	hm.ID = hm.Name

	data, err := d.client.Get(ctx, proxmoxtypes.TypeDir, hmID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read directory mapping %q", hmID),
			err.Error(),
		)

		return
	}

	hm.importFromAPI(ctx, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &hm)...)
}

// Schema defines the schema for the directory mapping.
func (d *dirDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	comment := dataSourceSchemaBaseAttrComment
	comment.Optional = false
	comment.Computed = true
	comment.Description = "The comment of this directory mapping."

	resp.Schema = schema.Schema{
		Description: "Retrieves a directory mapping from a Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			schemaAttrNameComment: comment,
			schemaAttrNameMap: schema.SetNestedAttribute{
				Computed:    true,
				Description: "The actual map of devices for the directory mapping.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaAttrNameMapNode: schema.StringAttribute{
							Computed:    true,
							Description: "The node name attribute of the map.",
						},
						schemaAttrNameMapPath: schema.StringAttribute{
							// For directory mappings the path is required and refers
							// to the POSIX path of the directory as visible from the node.
							Computed:    true,
							CustomType:  customtypes.PathType{},
							Description: "The path attribute of the map.",
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			schemaAttrNameName: schema.StringAttribute{
				Description: "The name of this directory mapping.",
				Required:    true,
			},
			schemaAttrNameTerraformID: attribute.ResourceID(
				"The unique identifier of this directory mapping data source.",
			),
		},
	}
}

// NewDirDataSource returns a new data source for a directory mapping.
// This is a helper function to simplify the provider implementation.
func NewDirDataSource() datasource.DataSource {
	return &dirDataSource{}
}
