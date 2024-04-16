/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

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
	_ resource.Resource                = &hardwareMappingUSBResource{}
	_ resource.ResourceWithConfigure   = &hardwareMappingUSBResource{}
	_ resource.ResourceWithImportState = &hardwareMappingUSBResource{}
)

// hardwareMappingUSBResource contains the USB hardware mapping resource's internal data.
type hardwareMappingUSBResource struct {
	// client is the hardware mapping API client.
	client mappings.Client
}

// read reads information about a USB hardware mapping from the Proxmox API.
func (r *hardwareMappingUSBResource) read(ctx context.Context, hm *hardwareMappingUSBModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	hmName := hm.Name.ValueString()
	data, err := r.client.Get(ctx, proxmoxtypes.HardwareMappingTypeUSB, hmName)
	if err != nil {
		if strings.Contains(err.Error(), "no such resource") {
			diags.AddError("Could not read USB hardware mapping", err.Error())
		}

		return false, diags
	}

	hm.importFromAPI(ctx, data)

	return true, nil
}

// readBack reads information about a created or modified USB hardware mapping from the Proxmox API then updates the
// response state accordingly.
// The Terraform resource identifier must have been set in the state before this method is called!
func (r *hardwareMappingUSBResource) readBack(
	ctx context.Context,
	hm *hardwareMappingUSBModel,
	respDiags *diag.Diagnostics,
	respState *tfsdk.State,
) {
	found, diags := r.read(ctx, hm)

	respDiags.Append(diags...)

	if !found {
		respDiags.AddError(
			"USB hardware mapping resource not found after update",
			"Failed to find the resource when trying to read back the updated USB hardware mapping's data.",
		)
	}

	if !respDiags.HasError() {
		respDiags.Append(respState.Set(ctx, *hm)...)
	}
}

// Configure adds the provider-configured client to the resource.
func (r *hardwareMappingUSBResource) Configure(
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

// Create creates a new USB hardware mapping.
func (r *hardwareMappingUSBResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var hm hardwareMappingUSBModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmName := hm.Name.ValueString()
	// Ensure to keep both in sync since the name represents the ID.
	hm.ID = hm.Name

	apiReq := hm.toCreateRequest()

	if err := r.client.Create(ctx, proxmoxtypes.HardwareMappingTypeUSB, apiReq); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not create USB hardware mapping %q.", hmName),
			err.Error(),
		)

		return
	}

	r.readBack(ctx, &hm, &resp.Diagnostics, &resp.State)
}

// Delete deletes an existing USB hardware mapping.
func (r *hardwareMappingUSBResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var hm hardwareMappingUSBModel

	resp.Diagnostics.Append(req.State.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmID := hm.Name.ValueString()

	if err := r.client.Delete(ctx, proxmoxtypes.HardwareMappingTypeUSB, hmID); err != nil {
		if strings.Contains(err.Error(), "no such resource") {
			resp.Diagnostics.AddWarning(
				"USB hardware mapping does not exist",
				fmt.Sprintf(
					"Could not delete USB hardware mapping %q, it does not exist or has been deleted outside of Terraform.",
					hmID,
				),
			)
		} else {
			resp.Diagnostics.AddError(fmt.Sprintf("Could not delete USB hardware mapping %q.", hmID), err.Error())
		}
	}
}

// ImportState imports a USB hardware mapping from the Proxmox VE API.
func (r *hardwareMappingUSBResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	data := hardwareMappingUSBModel{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}

	resource.ImportStatePassthroughID(ctx, path.Root(hardwareMappingSchemaAttrNameTerraformID), req, resp)
	r.readBack(ctx, &data, &resp.Diagnostics, &resp.State)
}

// Metadata defines the name of the USB hardware mapping.
func (r *hardwareMappingUSBResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_hardware_mapping_usb"
}

// Read reads the USB hardware mapping.
//

func (r *hardwareMappingUSBResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data hardwareMappingUSBModel

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

// Schema defines the schema for the USB hardware mapping.
func (r *hardwareMappingUSBResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	comment := hardwareMappingResourceSchemaWithBaseAttrComment
	comment.Description = "The comment of this USB hardware mapping."
	commentMap := comment
	commentMap.Description = "The comment of the mapped USB device."

	resp.Schema = schema.Schema{
		Description: "Manages a USB hardware mapping in a Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			hardwareMappingSchemaAttrNameComment: comment,
			hardwareMappingSchemaAttrNameMap: schema.SetNestedAttribute{
				Description: "The actual map of devices for the hardware mapping.",
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
						hardwareMappingSchemaAttrNameMapNode: schema.StringAttribute{
							Description: "The node name of the map.",
							Required:    true,
						},
						hardwareMappingSchemaAttrNameMapPath: schema.StringAttribute{
							CustomType: customtypes.HardwareMappingPathType{},
							Description: "The path of the map. For hardware mappings of type USB the path is optional and indicates" +
								" that the device is mapped through the device ID instead of ports.",
							Optional: true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									customtypes.HardwareMappingPathUSBValueRegEx,
									HardwareMappingResourceErrMessageInvalidPath(proxmoxtypes.HardwareMappingTypeUSB),
								),
							},
						},
					},
				},
				Required: true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			hardwareMappingSchemaAttrNameName: schema.StringAttribute{
				Description: "The name of this hardware mapping.",
				Required:    true,
			},
			hardwareMappingSchemaAttrNameTerraformID: structure.IDAttribute(
				"The unique identifier of this USB hardware mapping resource.",
			),
		},
	}
}

// Update updates an existing USB hardware mapping.
func (r *hardwareMappingUSBResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var hmCurrent, hmPlan hardwareMappingUSBModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &hmPlan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &hmCurrent)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmName := hmPlan.Name.ValueString()

	apiReq := hmPlan.toUpdateRequest(&hmCurrent)

	if err := r.client.Update(ctx, proxmoxtypes.HardwareMappingTypeUSB, hmName, apiReq); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not update USB hardware mapping %q.", hmName),
			err.Error(),
		)

		return
	}

	r.readBack(ctx, &hmPlan, &resp.Diagnostics, &resp.State)
}

// NewHardwareMappingUSBResource returns a new resource for managing a USB hardware mapping.
// This is a helper function to simplify the provider implementation.
func NewHardwareMappingUSBResource() resource.Resource {
	return &hardwareMappingUSBResource{}
}
