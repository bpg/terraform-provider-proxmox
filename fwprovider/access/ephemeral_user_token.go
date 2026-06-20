/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	ephschema "github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
	proxmoxapi "github.com/bpg/terraform-provider-proxmox/proxmox/api"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_ ephemeral.EphemeralResource              = &ephemeralUserToken{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralUserToken{}
	_ ephemeral.EphemeralResourceWithClose     = &ephemeralUserToken{}
)

type ephemeralUserToken struct {
	client proxmox.Client
}

type ephemeralUserTokenModel struct {
	UserID         types.String `tfsdk:"user_id"`
	TokenName      types.String `tfsdk:"token_name"`
	Comment        types.String `tfsdk:"comment"`
	ExpirationDate types.String `tfsdk:"expiration_date"`
	PrivSeparation types.Bool   `tfsdk:"privileges_separation"`
	AutoRevoke     types.Bool   `tfsdk:"auto_revoke"`
	ID             types.String `tfsdk:"id"`
	Value          types.String `tfsdk:"value"`
}

type ephemeralTokenPrivate struct {
	UserID     string `json:"user_id"`
	TokenName  string `json:"token_name"`
	AutoRevoke bool   `json:"auto_revoke"`
}

const ephemeralTokenPrivateKey = "token"

// NewEphemeralUserToken creates the short-named proxmox_user_token ephemeral resource.
func NewEphemeralUserToken() ephemeral.EphemeralResource {
	return &ephemeralUserToken{}
}

func (r *ephemeralUserToken) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = "proxmox_user_token"
}

func (r *ephemeralUserToken) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = ephschema.Schema{
		Description: "Creates an ephemeral API token for a Proxmox VE user. " +
			"The token value is returned during apply but never stored in Terraform state. " +
			"By default the token is automatically revoked from Proxmox when Terraform " +
			"completes the apply operation.",
		Attributes: map[string]ephschema.Attribute{
			"user_id": ephschema.StringAttribute{
				Description: "The user ID to create the token for (e.g., 'myuser@pam').",
				Required:    true,
			},
			"token_name": ephschema.StringAttribute{
				Description: "The token name.",
				Required:    true,
			},
			"comment": ephschema.StringAttribute{
				Description: "Optional comment for the token.",
				Optional:    true,
			},
			"expiration_date": ephschema.StringAttribute{
				Description: "Expiration date for the token (RFC3339 format).",
				Optional:    true,
				Validators: []validator.String{
					validators.NewParseValidator(func(s string) (time.Time, error) {
						return time.Parse(time.RFC3339, s)
					}, "must be a valid RFC3339 date"),
				},
			},
			"privileges_separation": ephschema.BoolAttribute{
				Description: "Whether privilege separation is enabled for the token. " +
					"When enabled, the token only has access to the explicitly listed permissions.",
				Optional: true,
				Computed: true,
			},
			"auto_revoke": ephschema.BoolAttribute{
				Description: "When true (default), the token is automatically revoked in Proxmox " +
					"after Terraform completes the apply. Set to false to leave the token in Proxmox.",
				Optional: true,
				Computed: true,
			},
			"id": ephschema.StringAttribute{
				Description: "The full token ID in the format '<user_id>!<token_name>'.",
				Computed:    true,
			},
			"value": ephschema.StringAttribute{
				Description: "The token secret. Available only during apply; never stored in Terraform state.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (r *ephemeralUserToken) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.DataSource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected config.DataSource, got: %T.", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client
}

func (r *ephemeralUserToken) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var m ephemeralUserTokenModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &m)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError("Provider Not Configured",
			"Cannot create ephemeral user token: provider client is nil.")

		return
	}

	body := access.UserTokenCreateRequestBody{
		Comment:      attribute.StringPtrFromValue(m.Comment),
		PrivSeparate: proxmoxtypes.CustomBoolPtr(m.PrivSeparation.ValueBoolPointer()),
	}

	if !m.ExpirationDate.IsNull() && m.ExpirationDate.ValueString() != "" {
		t, err := time.Parse(time.RFC3339, m.ExpirationDate.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Unable to Create Ephemeral User Token", err.Error())
			return
		}

		v := t.Unix()
		body.ExpirationDate = &v
	}

	userID := m.UserID.ValueString()
	tokenName := m.TokenName.ValueString()
	tokenID := userID + "!" + tokenName

	value, err := r.client.Access().CreateUserToken(ctx, userID, tokenName, &body)
	if err != nil {
		if !errors.Is(err, proxmoxapi.ErrResourceAlreadyExists) {
			resp.Diagnostics.AddError(fmt.Sprintf("Unable to Create Ephemeral User Token %q", tokenID), err.Error())
			return
		}

		// Terraform opens ephemeral resources during both plan and apply. If a previous
		// plan open left the token behind (auto_revoke=false), the apply open would fail
		// with "already exists". Delete and recreate to ensure a fresh token value.
		if delErr := r.client.Access().DeleteUserToken(ctx, userID, tokenName); delErr != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Unable to Recreate Ephemeral User Token %q", tokenID), delErr.Error())
			return
		}

		value, err = r.client.Access().CreateUserToken(ctx, userID, tokenName, &body)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Unable to Create Ephemeral User Token %q", tokenID), err.Error())
			return
		}
	}

	autoRevoke := m.AutoRevoke.IsNull() || m.AutoRevoke.ValueBool()

	m.ID = types.StringValue(tokenID)
	m.Value = types.StringValue(value)
	m.AutoRevoke = types.BoolValue(autoRevoke)

	resp.Diagnostics.Append(resp.Result.Set(ctx, &m)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Store identity in private state for Close.
	privateData, jsonErr := json.Marshal(ephemeralTokenPrivate{
		UserID:     userID,
		TokenName:  tokenName,
		AutoRevoke: autoRevoke,
	})
	if jsonErr != nil {
		resp.Diagnostics.AddError("Unable to Save Ephemeral Token Private State", jsonErr.Error())
		return
	}

	resp.Diagnostics.Append(resp.Private.SetKey(ctx, ephemeralTokenPrivateKey, privateData)...)
}

func (r *ephemeralUserToken) Close(ctx context.Context, req ephemeral.CloseRequest, resp *ephemeral.CloseResponse) {
	data, diags := req.Private.GetKey(ctx, ephemeralTokenPrivateKey)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if len(data) == 0 {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError("Provider Not Configured",
			"Cannot revoke ephemeral user token: provider client is nil.")

		return
	}

	var private ephemeralTokenPrivate

	if err := json.Unmarshal(data, &private); err != nil {
		resp.Diagnostics.AddError("Unable to Read Ephemeral Token Private State", err.Error())
		return
	}

	if !private.AutoRevoke {
		return
	}

	if err := r.client.Access().DeleteUserToken(ctx, private.UserID, private.TokenName); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Revoke Ephemeral User Token %q", private.UserID+"!"+private.TokenName),
			err.Error(),
		)
	}
}

// Long-name alias: proxmox_virtual_environment_user_token (deprecated).

type ephemeralUserTokenLong struct {
	ephemeralUserToken
}

var (
	_ ephemeral.EphemeralResource              = &ephemeralUserTokenLong{}
	_ ephemeral.EphemeralResourceWithConfigure = &ephemeralUserTokenLong{}
	_ ephemeral.EphemeralResourceWithClose     = &ephemeralUserTokenLong{}
)

// NewEphemeralUserTokenLong creates the long-named proxmox_virtual_environment_user_token ephemeral resource.
func NewEphemeralUserTokenLong() ephemeral.EphemeralResource {
	return &ephemeralUserTokenLong{}
}

func (r *ephemeralUserTokenLong) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = "proxmox_virtual_environment_user_token"
}

func (r *ephemeralUserTokenLong) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	r.ephemeralUserToken.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use proxmox_user_token instead. This resource will be removed in a future major release."
}
