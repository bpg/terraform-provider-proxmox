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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	_ resource.Resource                = &pciResource{}
	_ resource.ResourceWithConfigure   = &pciResource{}
	_ resource.ResourceWithImportState = &pciResource{}
)

// pciResource contains the PCI hardware mapping resource's internal data.
type pciResource struct {
	// client is the hardware mapping API client.
	client *mappings.Client
}

// read reads information about a PCI hardware mapping from the Proxmox VE API.
func (r *pciResource) read(ctx context.Context, hm *modelPCI) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	hmName := hm.Name.ValueString()

	data, err := r.client.Get(ctx, proxmoxtypes.TypePCI, hmName)
	if err != nil {
		if strings.Contains(err.Error(), "no such resource") {
			diags.AddError("Could not read PCI hardware mapping", err.Error())
		}

		return false, diags
	}

	hm.importFromAPI(ctx, data)

	return true, nil
}

// readBack reads information about a created or modified PCI hardware mapping from the Proxmox VE API then updates the
// response state accordingly.
// The Terraform resource identifier must have been set in the state before this method is called!
func (r *pciResource) readBack(ctx context.Context, hm *modelPCI, respDiags *diag.Diagnostics, respState *tfsdk.State) {
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
func (r *pciResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new PCI hardware mapping.
func (r *pciResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var hm modelPCI

	resp.Diagnostics.Append(req.Plan.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmName := hm.Name.ValueString()
	// Ensure to keep both in sync since the name represents the ID.
	hm.ID = hm.Name

	if err := r.client.Create(ctx, proxmoxtypes.TypePCI, hm.toCreateRequest()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not create PCI hardware mapping %q.", hmName),
			err.Error(),
		)

		return
	}

	r.readBack(ctx, &hm, &resp.Diagnostics, &resp.State)
}

// Delete deletes an existing PCI hardware mapping.
func (r *pciResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var hm modelPCI

	resp.Diagnostics.Append(req.State.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmID := hm.Name.ValueString()

	if err := r.client.Delete(ctx, proxmoxtypes.TypePCI, hmID); err != nil {
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
func (r *pciResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	data := modelPCI{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}

	resource.ImportStatePassthroughID(ctx, path.Root(schemaAttrNameTerraformID), req, resp)
	r.readBack(ctx, &data, &resp.Diagnostics, &resp.State)
}

// Metadata defines the name of the PCI hardware mapping.
func (r *pciResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hardware_mapping_pci"
}

// Read reads the PCI hardware mapping.
func (r *pciResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data modelPCI

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
func (r *pciResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	comment := resourceSchemaBaseAttrComment
	comment.Description = "The comment of this PCI hardware mapping."
	commentMap := comment
	commentMap.Description = "The comment of the mapped PCI device."

	resp.Schema = schema.Schema{
		Description: "Manages a PCI hardware mapping in a Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			schemaAttrNameComment: comment,
			schemaAttrNameMap: schema.SetNestedAttribute{
				Description: "The actual map of devices for the PCI hardware mapping.",
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
						schemaAttrNameMapIOMMUGroup: schema.Int64Attribute{
							Description: "The IOMMU group of the map. Not mandatory for the Proxmox VE API call, " +
								"but causes a PCI hardware mapping to be incomplete when not set",
							Optional: true,
						},
						schemaAttrNameMapNode: schema.StringAttribute{
							Description: "The node name of the map.",
							Required:    true,
						},
						schemaAttrNameMapPath: schema.StringAttribute{
							CustomType:  customtypes.PathType{},
							Description: "The path of the map.",
							// For hardware mappings of type PCI, the path is required while it is optional for USB.
							Required: true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									customtypes.PathPCIValueRegEx,
									ErrResourceMessageInvalidPath(proxmoxtypes.TypePCI),
								),
							},
						},
						schemaAttrNameMapSubsystemID: schema.StringAttribute{
							Description: "The subsystem ID group of the map. Not mandatory for the Proxmox VE API call, " +
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
			schemaAttrNameMediatedDevices: schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Indicates whether to enable mediated devices.",
			},
			schemaAttrNameName: schema.StringAttribute{
				Description: "The name of this PCI hardware mapping.",
				Required:    true,
			},
			schemaAttrNameTerraformID: attribute.ID(
				"The unique identifier of this PCI hardware mapping resource.",
			),
		},
	}
}

// Update updates an existing PCI hardware mapping.
func (r *pciResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var hmCurrent, hmPlan modelPCI

	resp.Diagnostics.Append(req.Plan.Get(ctx, &hmPlan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &hmCurrent)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmName := hmPlan.Name.ValueString()

	if err := r.client.Update(
		ctx,
		proxmoxtypes.TypePCI,
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

// NewPCIResource returns a new resource for managing a PCI hardware mapping.
// This is a helper function to simplify the provider implementation.
func NewPCIResource() resource.Resource {
	return &pciResource{}
}
