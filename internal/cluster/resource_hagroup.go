/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/internal/tffwk"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	hagroups "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/groups"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &hagroupResource{}
	_ resource.ResourceWithConfigure   = &hagroupResource{}
	_ resource.ResourceWithImportState = &hagroupResource{}
)

// NewHAGroupResource creates a new resource for managing Linux Bridge network interfaces.
func NewHAGroupResource() resource.Resource {
	return &hagroupResource{}
}

// hagroupResource contains the resource's internal data.
type hagroupResource struct {
	// The HA groups API client
	client hagroups.Client
}

// Metadata defines the name of the resource.
func (r *hagroupResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_hagroup"
}

// Schema defines the schema for the resource.
func (r *hagroupResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a High Availability group in a Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			"id": tffwk.IDAttribute(),
			"group": schema.StringAttribute{
				Description: "The identifier of the High Availability group to manage.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-_\.]*[a-zA-Z0-9]$`),
						"must start with a letter, end with a letter or number, be composed of "+
							"letters, numbers, '-', '_' and '.', and must be at least 2 characters long.",
					),
				},
			},
			"digest": schema.StringAttribute{
				Description: "The SHA-1 digest of the group's configuration",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "The comment associated with this group",
				Optional:    true,
			},
			"members": schema.MapAttribute{
				Description: "The member nodes for this group, associated with their priority or to null if no priority is set.",
				Required:    true,
				ElementType: types.Int64Type,
				Validators: []validator.Map{
					mapvalidator.SizeAtLeast(1),
					mapvalidator.KeysAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]*[a-zA-Z0-9])?$`),
							"must be a valid Proxmox node name",
						),
					),
					mapvalidator.ValueInt64sAre(int64validator.Between(0, 1000)),
				},
			},
			"no_failback": schema.BoolAttribute{
				Description: "A flag that indicates that failing back to a higher priority node is disabled for this HA " +
					"group. Defaults to `false`.",
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
			},
			"restricted": schema.BoolAttribute{
				Description: "A flag that indicates that other nodes may not be used to run resources associated to this HA " +
					"group. Defaults to `false`.",
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
			},
		},
	}
}

// Configure accesses the provider-configured Proxmox API client on behalf of the resource.
func (r *hagroupResource) Configure(
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
			fmt.Sprintf("Expected *proxmox.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData),
		)

		return
	}

	r.client = *client.Cluster().HA().Groups()
}

// Create creates a new HA group on the Proxmox cluster.
func (r *hagroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data hagroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupID := data.Group.ValueString()
	createRequest := &hagroups.HAGroupCreateRequestBody{}
	createRequest.ID = groupID
	createRequest.Comment = tffwk.OptStringFromModel(data.Comment)
	createRequest.Nodes = r.groupMembersToString(data.Members)
	createRequest.NoFailback = tffwk.BoolintFromModel(data.NoFailback)
	createRequest.Restricted = tffwk.BoolintFromModel(data.Restricted)
	createRequest.Type = "group"

	err := r.client.Create(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not create HA group '%s'.", groupID),
			err.Error(),
		)

		return
	}

	data.ID = types.StringValue(groupID)

	r.readBack(ctx, &data, &resp.Diagnostics, &resp.State)
}

// Read reads a HA group definition from the Proxmox cluster.
func (r *hagroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data hagroupModel

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

// Update updates a HA group definition on the Proxmox cluster.
func (r *hagroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state hagroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &hagroups.HAGroupUpdateRequestBody{}
	updateRequest.Comment = tffwk.OptStringFromModel(data.Comment)
	updateRequest.Digest = tffwk.OptStringFromModel(state.Digest)
	updateRequest.Nodes = r.groupMembersToString(data.Members)
	updateRequest.NoFailback = tffwk.BoolintFromModel(data.NoFailback)
	updateRequest.Restricted = tffwk.BoolintFromModel(data.Restricted)

	if updateRequest.Comment == nil && !state.Comment.IsNull() {
		updateRequest.Delete = "comment"
	}

	err := r.client.Update(ctx, state.Group.ValueString(), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating HA group",
			fmt.Sprintf("Could not update HA group '%s', unexpected error: %s",
				state.Group.ValueString(), err.Error()),
		)
		return
	}

	r.readBack(ctx, &data, &resp.Diagnostics, &resp.State)
}

// Delete deletes a HA group definition.
func (r *hagroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data hagroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupID := data.Group.ValueString()

	err := r.client.Delete(ctx, groupID)
	if err != nil {
		if strings.Contains(err.Error(), "no such ha group") {
			resp.Diagnostics.AddWarning(
				"HA group does not exist",
				fmt.Sprintf(
					"Could not delete HA group '%s', it does not exist or has been deleted outside of Terraform.",
					groupID,
				),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error deleting HA group",
				fmt.Sprintf("Could not delete HA group '%s', unexpected error: %s",
					groupID, err.Error()),
			)
		}
	}
}

// ImportState imports a HA group from the Proxmox cluster.
func (r *hagroupResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
}

// readBack reads information about a created or modified HA group from the cluster then updates the response
// state accordingly. It is assumed that the `state`'s identifier is set.
func (r *hagroupResource) readBack(ctx context.Context, data *hagroupModel, respDiags *diag.Diagnostics, respState *tfsdk.State) {
	found, diags := r.read(ctx, data)

	respDiags.Append(diags...)

	if !found {
		respDiags.AddError(
			"HA group not found after update",
			"Failed to find the group when trying to read back the updated HA group's data.",
		)
	}

	if !respDiags.HasError() {
		respDiags.Append(respState.Set(ctx, *data)...)
	}
}

// read reads information about a HA group from the cluster. The group identifier must have been set in the
// `data`.
func (r *hagroupResource) read(ctx context.Context, data *hagroupModel) (bool, diag.Diagnostics) {
	name := data.Group.ValueString()

	group, err := r.client.Get(ctx, name)
	if err != nil {
		var diags diag.Diagnostics

		if !strings.Contains(err.Error(), "no such ha group") {
			diags.AddError("Could not read HA group", err.Error())
		}

		return false, diags
	}

	return true, data.importFromAPI(*group)
}

// groupMembersToString converts the map of group member nodes into a string.
func (r *hagroupResource) groupMembersToString(members types.Map) string {
	mbElements := members.Elements()
	mbNodes := make([]string, len(mbElements))
	i := 0

	for name, value := range mbElements {
		if value.IsNull() {
			mbNodes[i] = name
		} else {
			mbNodes[i] = fmt.Sprintf("%s:%d", name, value.(types.Int64).ValueInt64())
		}

		i++
	}

	return strings.Join(mbNodes, ",")
}
