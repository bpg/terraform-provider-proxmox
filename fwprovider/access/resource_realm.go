/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
)

var (
	_ resource.Resource                     = (*realmResource)(nil)
	_ resource.ResourceWithConfigure        = (*realmResource)(nil)
	_ resource.ResourceWithImportState      = (*realmResource)(nil)
	_ resource.ResourceWithConfigValidators = (*realmResource)(nil)
)

type realmResource struct {
	client proxmox.Client
}

// NewACLResource creates a new ACL resource.
func NewRealmResource() resource.Resource {
	return &realmResource{}
}

func (r *realmResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages ACLs on the Proxmox cluster",
		MarkdownDescription: "Manages ACLs on the Proxmox cluster.\n\n" +
			"ACLs are used to control access to resources in the Proxmox cluster.\n" +
			"Each ACL consists of a path, a user, group or token, a role, and a flag to allow propagation of permissions.",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Description: "The group the ACL should apply to (mutually exclusive with `token_id` and `user_id`)",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": attribute.ID(),
			"path": schema.StringAttribute{
				Description: "Access control path",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"propagate": schema.BoolAttribute{
				Description: "Allow to propagate (inherit) permissions.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"role_id": schema.StringAttribute{
				Description: "The role to apply",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"token_id": schema.StringAttribute{
				Description: "The token the ACL should apply to (mutually exclusive with `group_id` and `user_id`)",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "The user the ACL should apply to (mutually exclusive with `group_id` and `token_id`)",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *realmResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("group_id"),
			path.MatchRoot("token_id"),
			path.MatchRoot("user_id"),
		),
	}
}

func (r *realmResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = cfg.Client
}

func (r *realmResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_realm"
}

func (r *realmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan realmResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.intoUpdateBody()

	err := r.client.Access().UpdateRealm(ctx, "TODO", body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create ACL", apiCallFailed+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *realmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state realmResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	realms, err := r.client.Access().GetRealm(ctx, "TODO")
	if err != nil {
		resp.Diagnostics.AddError("Unable read ACL", apiCallFailed+err.Error())
		return
	}

	_ = realms //TODO

	resp.State.RemoveResource(ctx)
}

func (r *realmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		state realmResourceModel
		plan  realmResourceModel
	)

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	stateBody := state.intoUpdateBody()

	err := r.client.Access().UpdateRealm(ctx, "TODO", stateBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete old ACL", apiCallFailed+err.Error())
		return
	}

	planBody := plan.intoUpdateBody()

	err = r.client.Access().UpdateRealm(ctx, "TODO", planBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create ACL", apiCallFailed+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *realmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state realmResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	stateBody := state.intoUpdateBody()

	err := r.client.Access().UpdateRealm(ctx, "TODO", stateBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete old ACL", apiCallFailed+err.Error())
		return
	}
}

func (r *realmResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	model, err := parseACLResourceModelFromID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to import ACL", "failed to parse ID: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
