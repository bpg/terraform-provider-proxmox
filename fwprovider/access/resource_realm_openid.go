/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

var (
	_ resource.Resource                = (*realmOpenIDResource)(nil)
	_ resource.ResourceWithConfigure   = (*realmOpenIDResource)(nil)
	_ resource.ResourceWithImportState = (*realmOpenIDResource)(nil)
)

type realmOpenIDResource struct {
	client proxmox.Client
}

// NewRealmOpenIDResource creates a new OpenID realm resource.
func NewRealmOpenIDResource() resource.Resource {
	return &realmOpenIDResource{}
}

func (r *realmOpenIDResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_realm_openid"
}

func (r *realmOpenIDResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages an OpenID Connect authentication realm in Proxmox VE.",
		MarkdownDescription: "Manages an OpenID Connect authentication realm in Proxmox VE.\n\n" +
			"OpenID Connect realms allow Proxmox to authenticate users against an external OpenID Connect provider.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID("Realm identifier (same as realm)"),
			"realm": schema.StringAttribute{
				Description: "Realm identifier (e.g., 'my-oidc').",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(32),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z][A-Za-z0-9.\-_]+$`),
						"must be a valid realm identifier",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"issuer_url": schema.StringAttribute{
				Description: "OpenID Connect issuer URL. Proxmox uses OpenID Connect Discovery to configure the provider.",
				Required:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "OpenID Connect Client ID.",
				Required:    true,
			},
			"client_key": schema.StringAttribute{
				Description: "OpenID Connect Client Key (secret). Note: stored in Proxmox but not returned by API.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"autocreate": schema.BoolAttribute{
				Description: "Automatically create users on the Proxmox cluster if they do not exist.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"username_claim": schema.StringAttribute{
				Description: "OpenID claim used to generate the unique username (subject, username, or email).",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("subject", "username", "email"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"groups_claim": schema.StringAttribute{
				Description: "OpenID claim used to retrieve user group memberships.",
				Optional:    true,
			},
			"groups_autocreate": schema.BoolAttribute{
				Description: "Automatically create groups from claims rather than using existing Proxmox VE groups.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"groups_overwrite": schema.BoolAttribute{
				Description: "Replace assigned groups on login instead of appending to existing ones.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"scopes": schema.StringAttribute{
				Description: "Space-separated list of OpenID scopes to request.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("email profile"),
			},
			"prompt": schema.StringAttribute{
				Description: "Specifies whether the authorization server prompts for reauthentication and/or consent " +
					"(e.g., 'none', 'login', 'consent', 'select_account').",
				Optional: true,
			},
			"acr_values": schema.StringAttribute{
				Description: "Authentication Context Class Reference values for the OpenID provider.",
				Optional:    true,
			},
			"query_userinfo": schema.BoolAttribute{
				Description: "Query the OpenID userinfo endpoint for claims. " +
					"Required when the identity provider does not include claims in the ID token.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"comment": schema.StringAttribute{
				Description: "Description of the realm.",
				Optional:    true,
			},
			"default": schema.BoolAttribute{
				Description: "Use this realm as the default for login.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *realmOpenIDResource) Configure(
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

	r.client = cfg.Client
}

func (r *realmOpenIDResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan realmOpenIDModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createReq := plan.toCreateRequest()

	err := r.client.Access().CreateRealm(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating OpenID realm", err.Error())
		return
	}

	plan.ID = plan.Realm

	err = r.readOpenID(ctx, &plan, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError("Error reading OpenID realm after create", err.Error())
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *realmOpenIDResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state realmOpenIDModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.readOpenID(ctx, &state, &resp.Diagnostics)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading OpenID Realm",
			fmt.Sprintf("Could not read realm %q: %v", state.Realm.ValueString(), err),
		)

		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *realmOpenIDResource) readOpenID(
	ctx context.Context,
	model *realmOpenIDModel,
	diags *diag.Diagnostics,
) error {
	realmData, err := r.client.Access().GetRealm(ctx, model.Realm.ValueString())
	if err != nil {
		return err
	}

	model.fromAPIResponse(realmData, diags)

	return nil
}

//nolint:dupl
func (r *realmOpenIDResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan realmOpenIDModel
	var state realmOpenIDModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := plan.toUpdateRequest(&state)

	err := r.client.Access().UpdateRealm(ctx, plan.Realm.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating OpenID realm", err.Error())
		return
	}

	err = r.readOpenID(ctx, &plan, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError("Error reading OpenID realm after update", err.Error())
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *realmOpenIDResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state realmOpenIDModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Access().DeleteRealm(ctx, state.Realm.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			return
		}

		resp.Diagnostics.AddError("Error deleting OpenID realm", err.Error())

		return
	}
}

func (r *realmOpenIDResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("realm"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
