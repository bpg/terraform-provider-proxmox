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
	_ datasource.DataSource              = &hardwareMappingPCIDatasource{}
	_ datasource.DataSourceWithConfigure = &hardwareMappingPCIDatasource{}
)

// hardwareMappingPCIDatasource is the data source implementation for a PCI hardware mapping.
type hardwareMappingPCIDatasource struct {
	// client is the hardware mapping API client.
	client *mappings.Client
}

// Configure adds the provider-configured client to the data source.
func (d *hardwareMappingPCIDatasource) Configure(
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
func (d *hardwareMappingPCIDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_hardware_mapping_pci"
}

// Read fetches the specified PCI hardware mapping from the Proxmox VE API.
func (d *hardwareMappingPCIDatasource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var hm hardwareMappingPCIModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmID := hm.Name.ValueString()
	// Ensure to keep both in sync since the name represents the ID.
	hm.ID = hm.Name

	data, err := d.client.Get(ctx, proxmoxtypes.HardwareMappingTypePCI, hmID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read PCI hardware mapping %q", hmID),
			err.Error(),
		)

		return
	}

	hm.importFromAPI(ctx, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &hm)...)
}

// Schema defines the schema for the PCI hardware mapping.
func (d *hardwareMappingPCIDatasource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	comment := hardwareMappingDataSourceSchemaWithBaseAttrComment
	comment.Optional = false
	comment.Computed = true
	comment.Description = "The comment of this PCI hardware mapping."
	commentMap := comment
	commentMap.Description = "The comment of the mapped PCI device."

	resp.Schema = schema.Schema{
		Description: "Retrieves a PCI hardware mapping from a Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			hardwareMappingSchemaAttrNameComment: comment,
			hardwareMappingSchemaAttrNameMap: schema.SetNestedAttribute{
				Computed:    true,
				Description: "The actual map of devices for the hardware mapping.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						hardwareMappingSchemaAttrNameComment: commentMap,
						hardwareMappingSchemaAttrNameMapIOMMUGroup: schema.Int64Attribute{
							Computed:    true,
							Description: "The IOMMU group attribute of the map.",
						},
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
							// For hardware mappings of type PCI, the path is required while it is optional for USB.
							Computed:    true,
							CustomType:  customtypes.HardwareMappingPathType{},
							Description: "The path attribute of the map.",
						},
						hardwareMappingSchemaAttrNameMapSubsystemID: schema.StringAttribute{
							Computed: true,
							Description: "The subsystem ID attribute of the map." +
								"Not mandatory for the Proxmox API call, but causes a PCI hardware mapping to be incomplete when not " +
								"set.",
							Validators: []validator.String{
								validators.HardwareMappingDeviceIDValidator(),
							},
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			hardwareMappingSchemaAttrNameMediatedDevices: schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether to use with mediated devices.",
			},
			hardwareMappingSchemaAttrNameName: schema.StringAttribute{
				Description: "The name of this PCI hardware mapping.",
				Required:    true,
			},
			hardwareMappingSchemaAttrNameTerraformID: structure.IDAttribute(
				"The unique identifier of this PCI hardware mapping data source.",
			),
		},
	}
}

// NewHardwareMappingPCIDatasource returns a new data source for a PCI hardware mapping.
// This is a helper function to simplify the provider implementation.
func NewHardwareMappingPCIDatasource() datasource.DataSource {
	return &hardwareMappingPCIDatasource{}
}
