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
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types/hardwaremapping"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	mappings "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/mapping"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/hardwaremapping"
)

// Ensure the implementation satisfies the required interfaces.
var (
	_ datasource.DataSource              = &datasourceUSB{}
	_ datasource.DataSourceWithConfigure = &datasourceUSB{}
)

// datasourceUSB is the data source implementation for a USB hardware mapping.
type datasourceUSB struct {
	client *mappings.Client
}

// Configure adds the provider-configured client to the data source.
func (d *datasourceUSB) Configure(
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
			fmt.Sprintf("Expected *proxmox.Client, got: %T", req.ProviderData),
		)

		return
	}

	d.client = client.Cluster().HardwareMapping()
}

// Metadata returns the data source type name.
func (d *datasourceUSB) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hardware_mapping_usb"
}

// Read fetches the specified USB hardware mapping from the Proxmox VE API.
func (d *datasourceUSB) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var hm modelUSB

	resp.Diagnostics.Append(req.Config.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmID := hm.Name.ValueString()
	// Ensure to keep both in sync since the name represents the ID.
	hm.ID = hm.Name

	data, err := d.client.Get(ctx, proxmoxtypes.TypeUSB, hmID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read USB hardware mapping %q", hmID),
			err.Error(),
		)

		return
	}

	hm.importFromAPI(ctx, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &hm)...)
}

// Schema defines the schema for the USB hardware mapping.
func (d *datasourceUSB) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	comment := dataSourceSchemaBaseAttrComment
	comment.Optional = false
	comment.Computed = true
	comment.Description = "The comment of this USB hardware mapping."
	commentMap := comment
	commentMap.Description = "The comment of the mapped USB device."

	resp.Schema = schema.Schema{
		Description: "Retrieves a USB hardware mapping from a Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			schemaAttrNameComment: comment,
			schemaAttrNameMap: schema.SetNestedAttribute{
				Computed:    true,
				Description: "The actual map of devices for the hardware mapping.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaAttrNameComment: commentMap,
						schemaAttrNameMapDeviceID: schema.StringAttribute{
							Computed:    true,
							Description: "The ID attribute of the map.",
							Validators: []validator.String{
								validators.HardwareMappingDeviceIDValidator(),
							},
						},
						schemaAttrNameMapNode: schema.StringAttribute{
							Computed:    true,
							Description: "The node name attribute of the map.",
						},
						schemaAttrNameMapPath: schema.StringAttribute{
							// For hardware mappings of type USB the path is optional and indicates that the device is mapped through
							// the device ID instead of ports.
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
				Description: "The name of this USB hardware mapping.",
				Required:    true,
			},
			schemaAttrNameTerraformID: attribute.ID(
				"The unique identifier of this USB hardware mapping data source.",
			),
		},
	}
}

// NewDataSourceUSB returns a new data source for a USB hardware mapping.
// This is a helper function to simplify the provider implementation.
func NewDataSourceUSB() datasource.DataSource {
	return &datasourceUSB{}
}
