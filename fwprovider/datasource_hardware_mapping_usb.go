/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package fwprovider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/structure"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	mappings "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/mapping"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// Ensure the implementation satisfies the required interfaces.
var (
	_ datasource.DataSource              = &hardwareMappingUSBDatasource{}
	_ datasource.DataSourceWithConfigure = &hardwareMappingUSBDatasource{}
)

// hardwareMappingUSBDatasource is the data source implementation for a USB hardware mapping.
type hardwareMappingUSBDatasource struct {
	client *mappings.Client
}

// Configure adds the provider-configured client to the data source.
func (d *hardwareMappingUSBDatasource) Configure(
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
func (d *hardwareMappingUSBDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_hardware_mapping_usb"
}

// Read fetches the specified USB hardware mapping from the Proxmox VE API.
func (d *hardwareMappingUSBDatasource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var hm hardwareMappingUSBModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmID := hm.Name.ValueString()
	// Ensure to keep both in sync since the name represents the ID.
	hm.ID = hm.Name

	data, err := d.client.Get(ctx, proxmoxtypes.HardwareMappingTypeUSB, hmID)
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
func (d *hardwareMappingUSBDatasource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	comment := hardwareMappingDataSourceSchemaWithBaseAttrComment
	comment.Optional = false
	comment.Computed = true
	comment.Description = "The comment of this USB hardware mapping."
	commentMap := comment
	commentMap.Description = "The comment of the mapped USB device."

	resp.Schema = schema.Schema{
		Description: "Retrieves a USB hardware mapping from a Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			hardwareMappingSchemaAttrNameComment: comment,
			hardwareMappingSchemaAttrNameMap: schema.SetNestedAttribute{
				Computed:    true,
				Description: "The actual map of devices for the hardware mapping.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						hardwareMappingSchemaAttrNameComment: commentMap,
						hardwareMappingSchemaAttrNameMapDeviceID: schema.StringAttribute{
							Computed:    true,
							Description: "The ID attribute of the map.",
							Validators: []validator.String{
								validators.HardwareMappingDeviceIDValidator(),
							},
						},
						hardwareMappingSchemaAttrNameMapNode: schema.StringAttribute{
							Computed:    true,
							Description: "The node name attribute of the map.",
						},
						hardwareMappingSchemaAttrNameMapPath: schema.StringAttribute{
							// For hardware mappings of type USB the path is optional and indicates that the device is mapped through
							// the device ID instead of ports.
							Computed:    true,
							CustomType:  customtypes.HardwareMappingPathType{},
							Description: "The path attribute of the map.",
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			hardwareMappingSchemaAttrNameName: schema.StringAttribute{
				Description: "The name of this USB hardware mapping.",
				Required:    true,
			},
			hardwareMappingSchemaAttrNameTerraformID: structure.IDAttribute(
				"The unique identifier of this USB hardware mapping data source.",
			),
		},
	}
}

// NewHardwareMappingUSBDatasource returns a new data source for a USB hardware mapping.
// This is a helper function to simplify the provider implementation.
func NewHardwareMappingUSBDatasource() datasource.DataSource {
	return &hardwareMappingUSBDatasource{}
}
