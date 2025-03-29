/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package hardwaremapping

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

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types/hardwaremapping"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	mappings "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/mapping"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/hardwaremapping"
)

// Ensure the resource implements the required interfaces.
var (
	_ resource.Resource                = &usbResource{}
	_ resource.ResourceWithConfigure   = &usbResource{}
	_ resource.ResourceWithImportState = &usbResource{}
)

// usbResource contains the USB hardware mapping resource's internal data.
type usbResource struct {
	// client is the hardware mapping API client.
	client *mappings.Client
}

// read reads information about a USB hardware mapping from the Proxmox VE API.
func (r *usbResource) read(ctx context.Context, hm *modelUSB) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	hmName := hm.Name.ValueString()

	data, err := r.client.Get(ctx, proxmoxtypes.TypeUSB, hmName)
	if err != nil {
		if strings.Contains(err.Error(), "no such resource") {
			diags.AddError("Could not read USB hardware mapping", err.Error())
		}

		return false, diags
	}

	hm.importFromAPI(ctx, data)

	return true, nil
}

// readBack reads information about a created or modified USB hardware mapping from the Proxmox VE API then updates the
// response state accordingly.
// The Terraform resource identifier must have been set in the state before this method is called!
func (r *usbResource) readBack(ctx context.Context, hm *modelUSB, respDiags *diag.Diagnostics, respState *tfsdk.State) {
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
func (r *usbResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.Resource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected config.Resource, got: %T", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client.Cluster().HardwareMapping()
}

// Create creates a new USB hardware mapping.
func (r *usbResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var hm modelUSB

	resp.Diagnostics.Append(req.Plan.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmName := hm.Name.ValueString()
	// Ensure to keep both in sync since the name represents the ID.
	hm.ID = hm.Name

	apiReq := hm.toCreateRequest()

	if err := r.client.Create(ctx, proxmoxtypes.TypeUSB, apiReq); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not create USB hardware mapping %q.", hmName),
			err.Error(),
		)

		return
	}

	r.readBack(ctx, &hm, &resp.Diagnostics, &resp.State)
}

// Delete deletes an existing USB hardware mapping.
func (r *usbResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var hm modelUSB

	resp.Diagnostics.Append(req.State.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmID := hm.Name.ValueString()

	if err := r.client.Delete(ctx, proxmoxtypes.TypeUSB, hmID); err != nil {
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
func (r *usbResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	data := modelUSB{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}

	resource.ImportStatePassthroughID(ctx, path.Root(schemaAttrNameTerraformID), req, resp)
	r.readBack(ctx, &data, &resp.Diagnostics, &resp.State)
}

// Metadata defines the name of the USB hardware mapping.
func (r *usbResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hardware_mapping_usb"
}

// Read reads the USB hardware mapping.
//

func (r *usbResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data modelUSB

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
func (r *usbResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	comment := resourceSchemaBaseAttrComment
	comment.Description = "The comment of this USB hardware mapping."
	commentMap := comment
	commentMap.Description = "The comment of the mapped USB device."

	resp.Schema = schema.Schema{
		Description: "Manages a USB hardware mapping in a Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			schemaAttrNameComment: comment,
			schemaAttrNameMap: schema.SetNestedAttribute{
				Description: "The actual map of devices for the hardware mapping.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaAttrNameComment: commentMap,
						schemaAttrNameMapDeviceID: schema.StringAttribute{
							Description: "The ID of the map.",
							Required:    true,
							Validators: []validator.String{
								validators.HardwareMappingDeviceIDValidator(),
							},
						},
						schemaAttrNameMapNode: schema.StringAttribute{
							Description: "The node name of the map.",
							Required:    true,
						},
						schemaAttrNameMapPath: schema.StringAttribute{
							CustomType: customtypes.PathType{},
							Description: "The path of the map. For hardware mappings of type USB the path is optional and indicates" +
								" that the device is mapped through the device ID instead of ports.",
							Optional: true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									customtypes.PathUSBValueRegEx,
									ErrResourceMessageInvalidPath(proxmoxtypes.TypeUSB),
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
			schemaAttrNameName: schema.StringAttribute{
				Description: "The name of this hardware mapping.",
				Required:    true,
			},
			schemaAttrNameTerraformID: attribute.ID(
				"The unique identifier of this USB hardware mapping resource.",
			),
		},
	}
}

// Update updates an existing USB hardware mapping.
func (r *usbResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var hmCurrent, hmPlan modelUSB

	resp.Diagnostics.Append(req.Plan.Get(ctx, &hmPlan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &hmCurrent)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmName := hmPlan.Name.ValueString()

	apiReq := hmPlan.toUpdateRequest(&hmCurrent)

	if err := r.client.Update(ctx, proxmoxtypes.TypeUSB, hmName, apiReq); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not update USB hardware mapping %q.", hmName),
			err.Error(),
		)

		return
	}

	r.readBack(ctx, &hmPlan, &resp.Diagnostics, &resp.State)
}

// NewUSBResource returns a new resource for managing a USB hardware mapping.
// This is a helper function to simplify the provider implementation.
func NewUSBResource() resource.Resource {
	return &usbResource{}
}
