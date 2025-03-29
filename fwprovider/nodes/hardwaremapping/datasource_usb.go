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
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	mappings "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/mapping"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/hardwaremapping"
)

// Ensure the implementation satisfies the required interfaces.
var (
	_ datasource.DataSource              = &usbDataSource{}
	_ datasource.DataSourceWithConfigure = &usbDataSource{}
)

// usbDataSource is the data source implementation for a USB hardware mapping.
type usbDataSource struct {
	client *mappings.Client
}

// Configure adds the provider-configured client to the data source.
func (d *usbDataSource) Configure(
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
func (d *usbDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hardware_mapping_usb"
}

// Read fetches the specified USB hardware mapping from the Proxmox VE API.
func (d *usbDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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
func (d *usbDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

// NewUSBDataSource returns a new data source for a USB hardware mapping.
// This is a helper function to simplify the provider implementation.
func NewUSBDataSource() datasource.DataSource {
	return &usbDataSource{}
}
