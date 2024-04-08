/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

//nolint:nolintlint,gofumpt,wsl // wsl linter is random-linter-loop-broken and only causes code to become unreadable.
package fwprovider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/structure"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	mappings "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/mapping"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// Ensure the resource implements the required interfaces.
var (
	_ resource.Resource                = &hardwareMappingPCIResource{}
	_ resource.ResourceWithConfigure   = &hardwareMappingPCIResource{}
	_ resource.ResourceWithImportState = &hardwareMappingPCIResource{}
)

// HardwareMappingResourceErrMessageInvalidPath is the error message for an invalid Linux device path for a hardware
// mapping of the specified type.
// Extracting the message helps to reduce duplicated code and allows to use it in automated unit and acceptance tests.
//
//nolint:gochecknoglobals
var HardwareMappingResourceErrMessageInvalidPath = func(hmType proxmoxtypes.HardwareMappingType) string {
	return fmt.Sprintf("not a valid Linux device path for hardware mapping of type %q", hmType)
}

// hardwareMappingPCIResource contains the PCI hardware mapping resource's internal data.
type hardwareMappingPCIResource struct {
	// client is the hardware mapping API client.
	client mappings.Client
}

// read reads information about a PCI hardware mapping from the Proxmox VE API.
func (r *hardwareMappingPCIResource) read(ctx context.Context, hm *hardwareMappingPCIModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	hmName := hm.Name.ValueString()
	data, err := r.client.Get(ctx, proxmoxtypes.HardwareMappingTypePCI, hmName)
	if err != nil {
		if strings.Contains(err.Error(), "no such resource") {
			diags.AddError("Could not read PCI hardware mapping", err.Error())
		}

		return false, diags
	}

	hm.importFromAPI(ctx, data)

	return true, nil
}

// readBack reads information about a created or modified PCI hardware mapping from the Proxmox API then updates the
// response state accordingly.
// The Terraform resource identifier must have been set in the state before this method is called!
func (r *hardwareMappingPCIResource) readBack(
	ctx context.Context,
	hm *hardwareMappingPCIModel,
	respDiags *diag.Diagnostics,
	respState *tfsdk.State,
) {
	found, diags := r.read(ctx, hm)

	respDiags.Append(diags...)

	if !found {
		respDiags.AddError(
			"PCI hardware mapping resource not found after update",
			"Failed to find the resource when trying to read back the updated PCI hardware mapping's data.",
		)
	}

	if !respDiags.HasError() {
		respDiags.Append(respState.Set(ctx, *hm)...)
	}
}

// Configure adds the provider-configured client to the resource.
func (r *hardwareMappingPCIResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
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
	}

	r.client = *client.Cluster().HardwareMapping()
}

// Create creates a new PCI hardware mapping.
func (r *hardwareMappingPCIResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var hm hardwareMappingPCIModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmName := hm.Name.ValueString()
	// Ensure to keep both in sync since the name represents the ID.
	hm.ID = hm.Name

	if err := r.client.Create(ctx, proxmoxtypes.HardwareMappingTypePCI, hm.toCreateRequest()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not create PCI hardware mapping %q.", hmName),
			err.Error(),
		)

		return
	}

	r.readBack(ctx, &hm, &resp.Diagnostics, &resp.State)
}

// Delete deletes an existing PCI hardware mapping.
func (r *hardwareMappingPCIResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var hm hardwareMappingPCIModel

	resp.Diagnostics.Append(req.State.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmID := hm.Name.ValueString()

	if err := r.client.Delete(ctx, proxmoxtypes.HardwareMappingTypePCI, hmID); err != nil {
		if strings.Contains(err.Error(), "no such resource") {
			resp.Diagnostics.AddWarning(
				"PCI hardware mapping does not exist",
				fmt.Sprintf(
					"Could not delete PCI hardware mapping %q, it does not exist or has been deleted outside of Terraform.",
					hmID,
				),
			)
		} else {
			resp.Diagnostics.AddError(fmt.Sprintf("Could not delete PCI hardware mapping %q.", hmID), err.Error())
		}
	}
}

// ImportState imports a PCI hardware mapping from the Proxmox VE API.
func (r *hardwareMappingPCIResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	data := hardwareMappingPCIModel{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}

	resource.ImportStatePassthroughID(ctx, path.Root(hardwareMappingSchemaAttrNameTerraformID), req, resp)
	r.readBack(ctx, &data, &resp.Diagnostics, &resp.State)
}

// Metadata defines the name of the PCI hardware mapping.
func (r *hardwareMappingPCIResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_hardware_mapping_pci"
}

// Read reads the PCI hardware mapping.
//

func (r *hardwareMappingPCIResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data hardwareMappingPCIModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	found, diags := r.read(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if !resp.Diagnostics.HasError() {
		if found {
			resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
		} else {
			resp.State.RemoveResource(ctx)
		}
	}
}

// Schema defines the schema for the PCI hardware mapping.
func (r *hardwareMappingPCIResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	comment := hardwareMappingResourceSchemaWithBaseAttrComment
	comment.Description = "The comment of this PCI hardware mapping."
	commentMap := comment
	commentMap.Description = "The comment of the mapped PCI device."

	resp.Schema = schema.Schema{
		Description: "Manages a PCI hardware mapping in a Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			hardwareMappingSchemaAttrNameComment: comment,
			hardwareMappingSchemaAttrNameMap: schema.SetNestedAttribute{
				Description: "The actual map of devices for the PCI hardware mapping.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						hardwareMappingSchemaAttrNameComment: commentMap,
						hardwareMappingSchemaAttrNameMapDeviceID: schema.StringAttribute{
							Description: "The ID of the map.",
							Required:    true,
							Validators: []validator.String{
								validators.HardwareMappingDeviceIDValidator(),
							},
						},
						hardwareMappingSchemaAttrNameMapIOMMUGroup: schema.Int64Attribute{
							Description: "The IOMMU group of the map. Not mandatory for the Proxmox API call, " +
								"but causes a PCI hardware mapping to be incomplete when not set",
							Optional: true,
						},
						hardwareMappingSchemaAttrNameMapNode: schema.StringAttribute{
							Description: "The node name of the map.",
							Required:    true,
						},
						hardwareMappingSchemaAttrNameMapPath: schema.StringAttribute{
							CustomType:  customtypes.HardwareMappingPathType{},
							Description: "The path of the map.",
							// For hardware mappings of type PCI, the path is required while it is optional for USB.
							Required: true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									customtypes.HardwareMappingPathPCIValueRegEx,
									HardwareMappingResourceErrMessageInvalidPath(proxmoxtypes.HardwareMappingTypePCI),
								),
							},
						},
						hardwareMappingSchemaAttrNameMapSubsystemID: schema.StringAttribute{
							Description: "The subsystem ID group of the map. Not mandatory for the Proxmox API call, " +
								"but causes a PCI hardware mapping to be incomplete when not set",
							Optional: true,
							Validators: []validator.String{
								validators.HardwareMappingDeviceIDValidator(),
							},
						},
					},
				},
				Required: true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			hardwareMappingSchemaAttrNameMediatedDevices: schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Indicates whether to enable mediated devices.",
			},
			hardwareMappingSchemaAttrNameName: schema.StringAttribute{
				Description: "The name of this PCI hardware mapping.",
				Required:    true,
			},
			hardwareMappingSchemaAttrNameTerraformID: structure.IDAttribute(
				"The unique identifier of this PCI hardware mapping resource.",
			),
		},
	}
}

// Update updates an existing PCI hardware mapping.
func (r *hardwareMappingPCIResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var hmCurrent, hmPlan hardwareMappingPCIModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &hmPlan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &hmCurrent)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmName := hmPlan.Name.ValueString()

	if err := r.client.Update(
		ctx,
		proxmoxtypes.HardwareMappingTypePCI,
		hmName,
		hmPlan.toUpdateRequest(&hmCurrent),
	); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not update PCI hardware mapping %q.", hmName),
			err.Error(),
		)

		return
	}

	r.readBack(ctx, &hmPlan, &resp.Diagnostics, &resp.State)
}

// NewHardwareMappingPCIResource returns a new resource for managing a PCI hardware mapping.
// This is a helper function to simplify the provider implementation.
func NewHardwareMappingPCIResource() resource.Resource {
	return &hardwareMappingPCIResource{}
}
