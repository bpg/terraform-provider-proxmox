/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package fwprovider

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/structure"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/mapping"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// Ensure the implementation satisfies the required interfaces.
var (
	_ datasource.DataSource              = &hardwareMappingsDatasource{}
	_ datasource.DataSourceWithConfigure = &hardwareMappingsDatasource{}
)

// hardwareMappingsDatasource is the data source implementation for a hardware mapping.
type hardwareMappingsDatasource struct {
	client *mapping.Client
}

// Configure adds the provider-configured client to the data source.
func (d *hardwareMappingsDatasource) Configure(
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
func (d *hardwareMappingsDatasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_hardware_mappings"
}

// Read fetches the list of hardware mappings from the Proxmox VE API then converts it to a list of strings.
func (d *hardwareMappingsDatasource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var hm hardwareMappingsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmType, err := proxmoxtypes.ParseHardwareMappingType(hm.Type.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Could not parse hardware mapping type", err.Error())
		return
	}

	list, err := d.client.List(ctx, hmType, hm.CheckNode.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read hardware mappings",
			err.Error(),
		)

		return
	}

	createCheckDiagnostics := func(
		hmID string,
		input []mapping.HardwareMappingNodeCheckDiagnostic,
	) []hardwareMappingsNodeCheckDiagnosticModel {
		checks := make([]hardwareMappingsNodeCheckDiagnosticModel, len(input))

		for idx, check := range input {
			model := hardwareMappingsNodeCheckDiagnosticModel{
				MappingID: types.StringValue(hmID),
				Severity:  types.StringPointerValue(check.Severity),
			}
			// Strip the unnecessary new line control character (\n) from the end of the message that is, for whatever reason,
			// returned this way by the Proxmox VE API.
			msg := strings.TrimSuffix(types.StringPointerValue(check.Message).ValueString(), "\n")
			model.Message = types.StringPointerValue(&msg)
			checks[idx] = model
		}

		return checks
	}

	mappings := make([]attr.Value, len(list))
	for idx, data := range list {
		mappings[idx] = types.StringValue(data.ID)
		// One of the fields only exists when the "check-node" option was passed to the Proxmox VE API with a valid node
		// name.
		// Note that the Proxmox VE API, for whatever reason, only returns one error at a time, even though the field is an
		// array.
		if (data.ChecksPCI != nil && len(data.ChecksPCI) > 0) || (data.ChecksUSB != nil && len(data.ChecksUSB) > 0) {
			switch data.Type {
			case proxmoxtypes.HardwareMappingTypePCI:
				hm.Checks = append(hm.Checks, createCheckDiagnostics(data.ID, data.ChecksPCI)...)
			case proxmoxtypes.HardwareMappingTypeUSB:
				hm.Checks = append(hm.Checks, createCheckDiagnostics(data.ID, data.ChecksUSB)...)
			}
		}
		// Ensure to keep the order of the diagnostic entries to prevent random plan changes.
		slices.SortStableFunc(
			hm.Checks, func(a, b hardwareMappingsNodeCheckDiagnosticModel) int {
				return strings.Compare(a.MappingID.ValueString(), b.MappingID.ValueString())
			},
		)
	}

	values, diags := types.SetValue(types.StringType, mappings)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	hm.MappingIDs = values
	hm.ID = types.StringValue("hardware_mappings")

	resp.Diagnostics.Append(resp.State.Set(ctx, &hm)...)
}

// Schema returns the schema for the data source.
func (d *hardwareMappingsDatasource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of hardware mapping resources.",
		Attributes: map[string]schema.Attribute{
			hardwareMappingsSchemaAttrNameChecks: schema.ListNestedAttribute{
				Computed:    true,
				Description: `Might contain relevant diagnostics about incorrect configurations.`,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						hardwareMappingsSchemaAttrNameChecksDiagnosticsMappingID: schema.StringAttribute{
							Computed:    true,
							Description: "The corresponding hardware mapping ID of the node check diagnostic entry.",
						},
						hardwareMappingsSchemaAttrNameChecksDiagnosticsMessage: schema.StringAttribute{
							Computed:    true,
							Description: "The message of the node check diagnostic entry.",
						},
						hardwareMappingsSchemaAttrNameChecksDiagnosticsSeverity: schema.StringAttribute{
							Computed:    true,
							Description: "The severity of the node check diagnostic entry.",
						},
					},
				},
			},
			hardwareMappingsSchemaAttrNameCheckNode: schema.StringAttribute{
				Description: "The name of the node whose configurations should be checked for correctness.",
				Optional:    true,
			},
			hardwareMappingsSchemaAttrNameHardwareMappingIDs: schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The identifiers of the hardware mappings.",
			},
			hardwareMappingSchemaAttrNameTerraformID: structure.IDAttribute(
				"The unique identifier of this hardware mappings data source.",
			),
			hardwareMappingSchemaAttrNameType: schema.StringAttribute{
				Description: "The type of the hardware mappings.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						[]string{
							proxmoxtypes.HardwareMappingTypePCI.String(),
							proxmoxtypes.HardwareMappingTypeUSB.String(),
						}...,
					),
				},
			},
		},
	}
}

// NewHardwareMappingsDatasource returns a new data source for hardware mappings.
// This is a helper function to simplify the provider implementation.
func NewHardwareMappingsDatasource() datasource.DataSource {
	return &hardwareMappingsDatasource{}
}
