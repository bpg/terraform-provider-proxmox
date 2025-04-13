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
	mappings "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/mapping"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/hardwaremapping"
)

// Ensure the resource implements the required interfaces.
var (
	_ resource.Resource                = &dirResource{}
	_ resource.ResourceWithConfigure   = &dirResource{}
	_ resource.ResourceWithImportState = &dirResource{}
)

// dirResource contains the directory mapping resource's internal data.
type dirResource struct {
	// client is the hardware mapping API client.
	client *mappings.Client
}

// read reads information about a directory mapping from the Proxmox VE API.
func (r *dirResource) read(ctx context.Context, hm *modelDir) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	hmName := hm.Name.ValueString()

	data, err := r.client.Get(ctx, proxmoxtypes.TypeDir, hmName)
	if err != nil {
		if strings.Contains(err.Error(), "no such resource") {
			diags.AddError("Could not read directory mapping", err.Error())
		}

		return false, diags
	}

	hm.importFromAPI(ctx, data)

	return true, nil
}

// readBack reads information about a created or modified directory mapping from the Proxmox VE API then updates the
// response state accordingly.
// The Terraform resource identifier must have been set in the state before this method is called!
func (r *dirResource) readBack(ctx context.Context, hm *modelDir, respDiags *diag.Diagnostics, respState *tfsdk.State) {
	found, diags := r.read(ctx, hm)

	respDiags.Append(diags...)

	if !found {
		respDiags.AddError(
			"directory mapping resource not found after update",
			"Failed to find the resource when trying to read back the updated directory mapping's data.",
		)
	}

	if !respDiags.HasError() {
		respDiags.Append(respState.Set(ctx, *hm)...)
	}
}

// Configure adds the provider-configured client to the resource.
func (r *dirResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new directory mapping.
func (r *dirResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var hm modelDir

	resp.Diagnostics.Append(req.Plan.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmName := hm.Name.ValueString()
	// Ensure to keep both in sync since the name represents the ID.
	hm.ID = hm.Name

	apiReq := hm.toCreateRequest()

	if err := r.client.Create(ctx, proxmoxtypes.TypeDir, apiReq); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not create directory mapping %q.", hmName),
			err.Error(),
		)

		return
	}

	r.readBack(ctx, &hm, &resp.Diagnostics, &resp.State)
}

// Delete deletes an existing directory mapping.
func (r *dirResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var hm modelDir

	resp.Diagnostics.Append(req.State.Get(ctx, &hm)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmID := hm.Name.ValueString()

	if err := r.client.Delete(ctx, proxmoxtypes.TypeDir, hmID); err != nil {
		if strings.Contains(err.Error(), "no such resource") {
			resp.Diagnostics.AddWarning(
				"directory mapping does not exist",
				fmt.Sprintf(
					"Could not delete directory mapping %q, it does not exist or has been deleted outside of Terraform.",
					hmID,
				),
			)
		} else {
			resp.Diagnostics.AddError(fmt.Sprintf("Could not delete directory mapping %q.", hmID), err.Error())
		}
	}
}

// ImportState imports a directory mapping from the Proxmox VE API.
func (r *dirResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	data := modelDir{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}

	resource.ImportStatePassthroughID(ctx, path.Root(schemaAttrNameTerraformID), req, resp)
	r.readBack(ctx, &data, &resp.Diagnostics, &resp.State)
}

// Metadata defines the name of the directory mapping.
func (r *dirResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hardware_mapping_dir"
}

// Read reads the directory mapping.
//

func (r *dirResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data modelDir

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

// Schema defines the schema for the directory mapping.
func (r *dirResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	comment := resourceSchemaBaseAttrComment
	comment.Description = "The comment of this directory mapping."

	resp.Schema = schema.Schema{
		Description: "Manages a directory mapping in a Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			schemaAttrNameComment: comment,
			schemaAttrNameMap: schema.SetNestedAttribute{
				Description: "The actual map of devices for the hardware mapping.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaAttrNameMapNode: schema.StringAttribute{
							Description: "The node this mapping applies to.",
							Required:    true,
						},
						schemaAttrNameMapPath: schema.StringAttribute{
							CustomType: customtypes.PathType{},
							Description: "The path of the map. For directory mappings the path is required and refers" +
								" to the POSIX path of the directory as visible from the node.",
							Required: true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									customtypes.PathDirValueRegEx,
									ErrResourceMessageInvalidPath(proxmoxtypes.TypeDir),
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
				Description: "The name of this directory mapping.",
				Required:    true,
			},
			schemaAttrNameTerraformID: attribute.ResourceID(
				"The unique identifier of this directory mapping resource.",
			),
		},
	}
}

// Update updates an existing directory mapping.
func (r *dirResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var hmCurrent, hmPlan modelDir

	resp.Diagnostics.Append(req.Plan.Get(ctx, &hmPlan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &hmCurrent)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hmName := hmPlan.Name.ValueString()

	apiReq := hmPlan.toUpdateRequest(&hmCurrent)

	if err := r.client.Update(ctx, proxmoxtypes.TypeDir, hmName, apiReq); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not update directory mapping %q.", hmName),
			err.Error(),
		)

		return
	}

	r.readBack(ctx, &hmPlan, &resp.Diagnostics, &resp.State)
}

// NewDirResource returns a new resource for managing a directory mapping.
// This is a helper function to simplify the provider implementation.
func NewDirResource() resource.Resource {
	return &dirResource{}
}
