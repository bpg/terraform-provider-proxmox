/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_ resource.Resource                = &userTokenResource{}
	_ resource.ResourceWithConfigure   = &userTokenResource{}
	_ resource.ResourceWithImportState = &userTokenResource{}
)

type userTokenResource struct {
	client proxmox.Client
}

type userTokenModel struct {
	Comment        types.String `tfsdk:"comment"`
	ExpirationDate types.String `tfsdk:"expiration_date"`
	ID             types.String `tfsdk:"id"`
	PrivSeparation types.Bool   `tfsdk:"privileges_separation"`
	UserID         types.String `tfsdk:"user_id"`
	TokenName      types.String `tfsdk:"token_name"`
	Value          types.String `tfsdk:"value"`
}

// NewUserTokenResource creates a new user token resource.
func NewUserTokenResource() resource.Resource {
	return &userTokenResource{}
}

func (r *userTokenResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "User API tokens.",
		Attributes: map[string]schema.Attribute{
			"comment": schema.StringAttribute{
				Description: "Comment for the token.",
				Optional:    true,
			},
			"expiration_date": schema.StringAttribute{
				Description: "Expiration date for the token.",
				Optional:    true,
				Validators: []validator.String{
					validators.NewParseValidator(func(s string) (time.Time, error) {
						return time.Parse(time.RFC3339, s)
					}, "must be a valid RFC3339 date"),
				},
			},
			"id": attribute.ID("Unique token identifier with format `<user_id>!<token_name>`."),
			"privileges_separation": schema.BoolAttribute{
				Description: "Restrict API token privileges with separate ACLs (default)",
				MarkdownDescription: "Restrict API token privileges with separate ACLs (default), " +
					"or give full privileges of corresponding user.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"token_name": schema.StringAttribute{
				Description: "User-specific token identifier.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`[A-Za-z][A-Za-z0-9.\-_]+`), "must be a valid token identifier"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				Description: "User identifier.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				Description: "API token value used for authentication.",
				MarkdownDescription: "API token value used for authentication. It is populated only when creating a new token, " +
					"and can't be retrieved at import.",
				Computed:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					// the attribute can't be retrieved after token creation, so during update we have to use value
					// from state (i.e. populated at create) if available.
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *userTokenResource) Configure(
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

	r.client = cfg.Client
}

func (r *userTokenResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_user_token"
}

func (r *userTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userTokenModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	body := access.UserTokenCreateRequestBody{
		Comment:      plan.Comment.ValueStringPointer(),
		PrivSeparate: proxmoxtypes.CustomBoolPtr(plan.PrivSeparation.ValueBoolPointer()),
	}

	if !plan.ExpirationDate.IsNull() && plan.ExpirationDate.ValueString() != "" {
		expirationDate, err := time.Parse(
			time.RFC3339,
			plan.ExpirationDate.ValueString(),
		)
		if err != nil {
			resp.Diagnostics.AddError("Error parsing expiration date", err.Error())
			return
		}

		v := expirationDate.Unix()
		body.ExpirationDate = &v
	}

	value, err := r.client.Access().CreateUserToken(ctx, plan.UserID.ValueString(), plan.TokenName.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user token", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(plan.UserID.ValueString() + "!" + plan.TokenName.ValueString())
	plan.Value = types.StringValue(value)
	resp.State.Set(ctx, plan)
}

func (r *userTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userTokenModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.Access().GetUserToken(ctx, state.UserID.ValueString(), state.TokenName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading user token", err.Error())
		return
	}

	state.Comment = types.StringPointerValue(data.Comment)

	if data.ExpirationDate != nil && *data.ExpirationDate > 0 {
		dt := time.Unix(int64(*data.ExpirationDate), 0).UTC().Format(time.RFC3339)
		state.ExpirationDate = types.StringValue(dt)
	}

	state.PrivSeparation = types.BoolPointerValue(data.PrivSeparate.PointerBool())

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *userTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state userTokenModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := access.UserTokenUpdateRequestBody{
		// note: PVE API does not support resetting comment to empty string
		Comment:      plan.Comment.ValueStringPointer(),
		PrivSeparate: proxmoxtypes.CustomBoolPtr(plan.PrivSeparation.ValueBoolPointer()),
	}

	if !plan.ExpirationDate.IsNull() && plan.ExpirationDate.ValueString() != "" {
		// if planned value is not empty then set it
		expirationDate, err := time.Parse(
			time.RFC3339,
			plan.ExpirationDate.ValueString(),
		)
		if err != nil {
			resp.Diagnostics.AddError("Error parsing expiration date", err.Error())
			return
		}

		v := expirationDate.Unix()
		body.ExpirationDate = &v
	} else if !state.ExpirationDate.IsNull() {
		// if planned value is empty, but the current value is not then reset it
		body.ExpirationDate = new(int64)
	}

	err := r.client.Access().UpdateUserToken(ctx, plan.UserID.ValueString(), plan.TokenName.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user token", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.Set(ctx, plan)
}

func (r *userTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userTokenModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Access().DeleteUserToken(ctx, state.UserID.ValueString(), state.TokenName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user token", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *userTokenResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	idParts := strings.Split(req.ID, "!")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: 'user_id!token_name'. Got: %q", req.ID),
		)

		return
	}

	userID := idParts[0]
	tokenName := idParts[1]

	data, err := r.client.Access().GetUserToken(ctx, userID, tokenName)
	if err != nil {
		resp.Diagnostics.AddError("Error reading user token", err.Error())
		return
	}

	state := userTokenModel{
		Comment:        types.StringPointerValue(data.Comment),
		ID:             types.StringValue(req.ID),
		PrivSeparation: types.BoolPointerValue(data.PrivSeparate.PointerBool()),
		UserID:         types.StringValue(userID),
		TokenName:      types.StringValue(tokenName),
		Value:          types.StringNull(),
	}

	if data.ExpirationDate != nil && *data.ExpirationDate > 0 {
		state.ExpirationDate = types.StringValue(time.Unix(int64(*data.ExpirationDate), 0).UTC().Format(time.RFC3339))
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
