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

	"github.com/bpg/terraform-provider-proxmox/fwprovider/structure"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_ resource.Resource                     = (*aclResource)(nil)
	_ resource.ResourceWithConfigure        = (*aclResource)(nil)
	_ resource.ResourceWithImportState      = (*aclResource)(nil)
	_ resource.ResourceWithConfigValidators = (*aclResource)(nil)
)

type aclResource struct {
	client proxmox.Client
}

// NewACLResource creates a new ACL resource.
func NewACLResource() resource.Resource {
	return &aclResource{}
}

func (r *aclResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"id": structure.IDAttribute(),
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

func (r *aclResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("group_id"),
			path.MatchRoot("token_id"),
			path.MatchRoot("user_id"),
		),
	}
}

func (r *aclResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *aclResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_acl"
}

func (r *aclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aclResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.intoUpdateBody()

	err := r.client.Access().UpdateACL(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create ACL", apiCallFailed+err.Error())
		return
	}

	err = plan.generateID()
	if err != nil {
		resp.Diagnostics.AddError("Unable to create ACL", "failed to generate ID: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *aclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aclResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	acls, err := r.client.Access().GetACL(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable read ACL", apiCallFailed+err.Error())
		return
	}

	for _, acl := range acls {
		switch acl.Type {
		case "group":
			if acl.UserOrGroupID != state.GroupID.ValueString() {
				continue
			}
		case "token":
			if acl.UserOrGroupID != state.TokenID.ValueString() {
				continue
			}
		case "user":
			if acl.UserOrGroupID != state.UserID.ValueString() {
				continue
			}
		default:
			// ignore unknown values
			continue
		}

		if acl.Path != state.Path {
			continue
		}

		if acl.RoleID != state.RoleID {
			continue
		}

		state.Propagate = ptr.Or(acl.Propagate.PointerBool(), true)

		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *aclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		state aclResourceModel
		plan  aclResourceModel
	)

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	stateBody := state.intoUpdateBody()
	stateBody.Delete = proxmoxtypes.CustomBool(true).Pointer()

	err := r.client.Access().UpdateACL(ctx, stateBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete old ACL", apiCallFailed+err.Error())
		return
	}

	planBody := plan.intoUpdateBody()

	err = r.client.Access().UpdateACL(ctx, planBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create ACL", apiCallFailed+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *aclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aclResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	stateBody := state.intoUpdateBody()
	stateBody.Delete = proxmoxtypes.CustomBool(true).Pointer()

	err := r.client.Access().UpdateACL(ctx, stateBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete old ACL", apiCallFailed+err.Error())
		return
	}
}

func (r *aclResource) ImportState(
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

const apiCallFailed = "API call failed: "
