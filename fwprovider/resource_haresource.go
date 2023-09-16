/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fwprovider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/structure"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	haresources "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/resources"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// haResourceResource contains the resource's internal data.
// NOTE: the naming is horrible, but this is the convention used by the framework.
// and the entity name in the API is "ha resource", so...
type haResourceResource struct {
	// The HA resources API client
	client haresources.Client
}

// Ensure the resource implements the expected interfaces.
var (
	_ resource.Resource                = &haResourceResource{}
	_ resource.ResourceWithConfigure   = &haResourceResource{}
	_ resource.ResourceWithImportState = &haResourceResource{}
)

// NewHAResourceResource returns a new resource for managing High Availability resources.
func NewHAResourceResource() resource.Resource {
	return &haResourceResource{}
}

// Metadata defines the name of the resource.
func (r *haResourceResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_haresource"
}

// Schema defines the schema for the resource.
func (r *haResourceResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages Proxmox HA resources.",
		Attributes: map[string]schema.Attribute{
			"id": structure.IDAttribute(),
			"resource_id": schema.StringAttribute{
				Description: "The Proxmox HA resource identifier",
				Required:    true,
				Validators: []validator.String{
					validators.HAResourceIDValidator(),
				},
			},
			"state": schema.StringAttribute{
				Description: "The desired state of the resource.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("started"),
				Validators: []validator.String{
					validators.HAResourceStateValidator(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of HA resources to create. If unset, it will be deduced from the `resource_id`.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					validators.HAResourceTypeValidator(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "The comment associated with this resource.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[^\s]|^$`), "must not start with whitespace"),
					stringvalidator.RegexMatches(regexp.MustCompile(`[^\s]$|^$`), "must not end with whitespace"),
				},
			},
			"group": schema.StringAttribute{
				Description: "The identifier of the High Availability group this resource is a member of.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-_\.]*[a-zA-Z0-9]$`),
						"must start with a letter, end with a letter or number, be composed of "+
							"letters, numbers, '-', '_' and '.', and must be at least 2 characters long",
					),
				},
			},
			"max_relocate": schema.Int64Attribute{
				Description: "The maximal number of relocation attempts.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 10),
				},
			},
			"max_restart": schema.Int64Attribute{
				Description: "The maximal number of restart attempts.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 10),
				},
			},
		},
	}
}

// Configure adds the provider-configured client to the resource.
func (r *haResourceResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(proxmox.Client)
	if ok {
		r.client = *client.Cluster().HA().Resources()
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData),
		)
	}
}

// Create creates a new HA resource.
func (r *haResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data haResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resID, err := proxmoxtypes.ParseHAResourceID(data.ResourceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error parsing Proxmox HA resource identifier",
			fmt.Sprintf("Couldn't parse the Terraform resource ID into a valid HA resource identifier: %s. "+
				"Please report this issue to the provider developers.", err),
		)

		return
	}

	createRequest := data.toCreateRequest(resID)

	err = r.client.Create(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not create HA resource '%v'.", resID),
			err.Error(),
		)

		return
	}

	data.ID = types.StringValue(resID.String())

	r.readBack(ctx, &data, &resp.Diagnostics, &resp.State)
}

// Update updates an existing HA resource.
func (r *haResourceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data, state haResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resID, err := proxmoxtypes.ParseHAResourceID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error parsing Proxmox HA resource identifier",
			fmt.Sprintf("Couldn't parse the Terraform resource ID into a valid HA resource identifier: %s. "+
				"Please report this issue to the provider developers.", err),
		)

		return
	}

	updateRequest := data.toUpdateRequest(&state)

	err = r.client.Update(ctx, resID, updateRequest)
	if err == nil {
		r.readBack(ctx, &data, &resp.Diagnostics, &resp.State)
	} else {
		resp.Diagnostics.AddError(
			"Error updating HA resource",
			fmt.Sprintf("Could not update HA resource '%s', unexpected error: %s",
				state.Group.ValueString(), err.Error()),
		)
	}
}

// Delete deletes an existing HA resource.
func (r *haResourceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data haResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resID, err := proxmoxtypes.ParseHAResourceID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error parsing Proxmox HA resource identifier",
			fmt.Sprintf("Couldn't parse the Terraform resource ID into a valid HA resource identifier: %s. "+
				"Please report this issue to the provider developers.", err),
		)

		return
	}

	err = r.client.Delete(ctx, resID)
	if err != nil {
		if strings.Contains(err.Error(), "no such resource") {
			resp.Diagnostics.AddWarning(
				"HA resource does not exist",
				fmt.Sprintf(
					"Could not delete HA resource '%v', it does not exist or has been deleted outside of Terraform.",
					resID,
				),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error deleting HA resource",
				fmt.Sprintf("Could not delete HA resource '%v', unexpected error: %s",
					resID, err.Error()),
			)
		}
	}
}

// Read reads the HA resource.
func (r *haResourceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data haResourceModel

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

// ImportState imports a HA resource from the Proxmox cluster.
func (r *haResourceResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	reqID := req.ID
	data := haResourceModel{
		ID:         types.StringValue(reqID),
		ResourceID: types.StringValue(reqID),
	}
	r.readBack(ctx, &data, &resp.Diagnostics, &resp.State)
}

// read reads information about a HA resource from the cluster. The Terraform resource identifier must have been set
// in the model before this function is called.
func (r *haResourceResource) read(ctx context.Context, data *haResourceModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	resID, err := proxmoxtypes.ParseHAResourceID(data.ID.ValueString())
	if err != nil {
		diags.AddError(
			"Unexpected error parsing Proxmox HA resource identifier",
			fmt.Sprintf("Couldn't parse the Terraform resource ID into a valid HA resource identifier: %s. "+
				"Please report this issue to the provider developers.", err),
		)

		return false, diags
	}

	res, err := r.client.Get(ctx, resID)
	if err != nil {
		if !strings.Contains(err.Error(), "no such resource") {
			diags.AddError("Could not read HA resource", err.Error())
		}

		return false, diags
	}

	data.importFromAPI(res)

	return true, nil
}

// readBack reads information about a created or modified HA resource from the cluster then updates the response
// state accordingly. It is assumed that the `state`'s identifier is set.
func (r *haResourceResource) readBack(
	ctx context.Context,
	data *haResourceModel,
	respDiags *diag.Diagnostics,
	respState *tfsdk.State,
) {
	found, diags := r.read(ctx, data)

	respDiags.Append(diags...)

	if !found {
		respDiags.AddError(
			"HA resource not found after update",
			"Failed to find the resource when trying to read back the updated HA resource's data.",
		)
	}

	if !respDiags.HasError() {
		respDiags.Append(respState.Set(ctx, *data)...)
	}
}
