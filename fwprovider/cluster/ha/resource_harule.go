/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ha

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	harules "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/ha/rules"
)

var (
	_ resource.Resource                   = &haruleResource{}
	_ resource.ResourceWithConfigure      = &haruleResource{}
	_ resource.ResourceWithImportState    = &haruleResource{}
	_ resource.ResourceWithValidateConfig = &haruleResource{}
)

// NewHARuleResource creates a new resource for managing HA rules.
func NewHARuleResource() resource.Resource {
	return &haruleResource{}
}

// haruleResource contains the resource's internal data.
type haruleResource struct {
	client *harules.Client
}

// Metadata defines the name of the resource.
func (r *haruleResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_harule"
}

// Schema defines the schema for the resource.
func (r *haruleResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a High Availability rule in a Proxmox VE cluster (PVE 9+). " +
			"HA rules replace the legacy HA groups and provide node affinity and resource affinity capabilities.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID(),
			"rule": schema.StringAttribute{
				Description: "The identifier of the High Availability rule to manage.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-_.]*[a-zA-Z0-9]$`),
						"must start with a letter, end with a letter or number, be composed of "+
							"letters, numbers, '-', '_' and '.', and must be at least 2 characters long",
					),
				},
			},
			"type": schema.StringAttribute{
				Description: "The HA rule type. Must be `node-affinity` or `resource-affinity`.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("node-affinity", "resource-affinity"),
				},
			},
			"comment": schema.StringAttribute{
				Description: "The comment associated with this rule.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
					stringvalidator.RegexMatches(regexp.MustCompile(`^\S|^$`), "must not start with whitespace"),
					stringvalidator.RegexMatches(regexp.MustCompile(`\S$|^$`), "must not end with whitespace"),
				},
			},
			"disable": schema.BoolAttribute{
				Description: "Whether the HA rule is disabled. Defaults to `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"resources": schema.SetAttribute{
				Description: "The set of HA resource IDs that this rule applies to (e.g. `vm:100`, `ct:101`). " +
					"The resources must already be managed by HA.",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^(vm|ct):\d+$`),
							"must be in format type:id (e.g. vm:100, ct:101)",
						),
					),
				},
			},
			"nodes": schema.MapAttribute{
				Description: "The member nodes for this rule (node-affinity only). They are provided as a map, " +
					"where the keys are the node names and the values represent their priority: integers for " +
					"known priorities or `null` for unset priorities.",
				Optional:    true,
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
			"strict": schema.BoolAttribute{
				Description: "Whether the node affinity rule is strict (node-affinity only). " +
					"When strict, resources cannot run on nodes not listed. Defaults to `false`.",
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
			},
			"affinity": schema.StringAttribute{
				Description: "The resource affinity type (resource-affinity only). " +
					"`positive` keeps resources on the same node, `negative` keeps them on separate nodes.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("positive", "negative"),
				},
			},
		},
	}
}

// ValidateConfig validates the resource configuration.
func (r *haruleResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var data RuleModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Skip validation if type is not yet known (e.g. during early plan phase).
	if data.Type.IsNull() || data.Type.IsUnknown() {
		return
	}

	ruleType := data.Type.ValueString()

	switch ruleType {
	case "node-affinity":
		if data.Nodes.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("nodes"),
				"Missing required attribute",
				"The `nodes` attribute is required when `type` is \"node-affinity\".",
			)
		}

		if !data.Affinity.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("affinity"),
				"Invalid attribute combination",
				"The `affinity` attribute is only valid when `type` is \"resource-affinity\".",
			)
		}

	case "resource-affinity":
		if data.Affinity.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("affinity"),
				"Missing required attribute",
				"The `affinity` attribute is required when `type` is \"resource-affinity\".",
			)
		}

		if !data.Nodes.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("nodes"),
				"Invalid attribute combination",
				"The `nodes` attribute is only valid when `type` is \"node-affinity\".",
			)
		}
	}
}

// Configure accesses the provider-configured Proxmox API client on behalf of the resource.
func (r *haruleResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
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

	r.client = cfg.Client.Cluster().HA().Rules()
}

// Create creates a new HA rule on the Proxmox cluster.
func (r *haruleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ruleID := data.Rule.ValueString()
	createRequest := &harules.HARuleCreateRequestBody{}
	createRequest.Rule = ruleID
	createRequest.Type = data.Type.ValueString()
	createRequest.Comment = data.Comment.ValueStringPointer()
	createRequest.Disable.FromValue(data.Disable)
	createRequest.Resources = r.resourcesToString(data.Resources)

	switch data.Type.ValueString() {
	case "node-affinity":
		nodesStr := r.nodesToString(data.Nodes)
		createRequest.Nodes = &nodesStr
		createRequest.Strict.FromValue(data.Strict)
	case "resource-affinity":
		createRequest.Affinity = data.Affinity.ValueStringPointer()
	}

	err := r.client.Create(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not create HA rule '%s'.", ruleID),
			err.Error(),
		)

		return
	}

	data.ID = types.StringValue(ruleID)

	r.readBack(ctx, &data, &resp.Diagnostics, &resp.State)
}

// Read reads a HA rule definition from the Proxmox cluster.
func (r *haruleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RuleModel

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

// Update updates a HA rule definition on the Proxmox cluster.
func (r *haruleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state RuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &harules.HARuleUpdateRequestBody{}
	updateRequest.Type = data.Type.ValueString()
	updateRequest.Comment = data.Comment.ValueStringPointer()
	updateRequest.Disable.FromValue(data.Disable)
	updateRequest.Resources = r.resourcesToString(data.Resources)

	var deleteFields []string

	if updateRequest.Comment == nil && !state.Comment.IsNull() {
		deleteFields = append(deleteFields, "comment")
	}

	switch data.Type.ValueString() {
	case "node-affinity":
		nodesStr := r.nodesToString(data.Nodes)
		updateRequest.Nodes = &nodesStr
		updateRequest.Strict.FromValue(data.Strict)
	case "resource-affinity":
		updateRequest.Affinity = data.Affinity.ValueStringPointer()
	}

	if len(deleteFields) > 0 {
		updateRequest.Delete = deleteFields
	}

	err := r.client.Update(ctx, state.Rule.ValueString(), updateRequest)
	if err == nil {
		r.readBack(ctx, &data, &resp.Diagnostics, &resp.State)
	} else {
		resp.Diagnostics.AddError(
			"Error updating HA rule",
			fmt.Sprintf("Could not update HA rule '%s', unexpected error: %s",
				state.Rule.ValueString(), err.Error()),
		)
	}
}

// Delete deletes a HA rule definition.
func (r *haruleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ruleID := data.Rule.ValueString()

	err := r.client.Delete(ctx, ruleID)
	if err != nil {
		if strings.Contains(err.Error(), "no such ha rule") {
			resp.Diagnostics.AddWarning(
				"HA rule does not exist",
				fmt.Sprintf(
					"Could not delete HA rule '%s', it does not exist or has been deleted outside of Terraform.",
					ruleID,
				),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error deleting HA rule",
				fmt.Sprintf("Could not delete HA rule '%s', unexpected error: %s",
					ruleID, err.Error()),
			)
		}
	}
}

// ImportState imports a HA rule from the Proxmox cluster.
func (r *haruleResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	reqID := req.ID
	data := RuleModel{
		ID:   types.StringValue(reqID),
		Rule: types.StringValue(reqID),
	}

	r.readBack(ctx, &data, &resp.Diagnostics, &resp.State)
}

// readBack reads information about a created or modified HA rule from the cluster then updates the response
// state accordingly.
func (r *haruleResource) readBack(
	ctx context.Context,
	data *RuleModel,
	respDiags *diag.Diagnostics,
	respState *tfsdk.State,
) {
	found, diags := r.read(ctx, data)

	respDiags.Append(diags...)

	if !found {
		respDiags.AddError(
			"HA rule not found after update",
			"Failed to find the rule when trying to read back the updated HA rule's data.",
		)
	}

	if !respDiags.HasError() {
		respDiags.Append(respState.Set(ctx, *data)...)
	}
}

// read reads information about a HA rule from the cluster.
func (r *haruleResource) read(ctx context.Context, data *RuleModel) (bool, diag.Diagnostics) {
	name := data.Rule.ValueString()

	rule, err := r.client.Get(ctx, name)
	if err != nil {
		diags := diag.Diagnostics{}

		if !strings.Contains(err.Error(), "no such ha rule") {
			diags.AddError("Could not read HA rule", err.Error())
		}

		return false, diags
	}

	return true, data.ImportFromAPI(*rule)
}

// resourcesToString converts the set of resource IDs into a comma-separated string.
func (r *haruleResource) resourcesToString(resources types.Set) string {
	elements := resources.Elements()
	parts := make([]string, len(elements))

	for i, elem := range elements {
		parts[i] = elem.(types.String).ValueString()
	}

	return strings.Join(parts, ",")
}

// nodesToString converts the map of node priorities into a comma-separated string.
func (r *haruleResource) nodesToString(nodes types.Map) string {
	elements := nodes.Elements()
	parts := make([]string, 0, len(elements))

	for name, value := range elements {
		if value.IsNull() {
			parts = append(parts, name)
		} else {
			parts = append(parts, fmt.Sprintf("%s:%d", name, value.(types.Int64).ValueInt64()))
		}
	}

	return strings.Join(parts, ",")
}
