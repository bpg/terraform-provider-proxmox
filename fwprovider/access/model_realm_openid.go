/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type realmOpenIDModel struct {
	ID        types.String `tfsdk:"id"`
	Realm     types.String `tfsdk:"realm"`
	IssuerURL types.String `tfsdk:"issuer_url"`
	ClientID  types.String `tfsdk:"client_id"`
	ClientKey types.String `tfsdk:"client_key"`
	// Write-only client secret; never persisted to state. Read via req.Config.
	ClientKeyWO types.String `tfsdk:"client_key_wo"`
	// Change-detection trigger for the write-only ClientKeyWO; bump to force a resend.
	ClientKeyWOVersion types.Int64  `tfsdk:"client_key_wo_version"`
	AutoCreate         types.Bool   `tfsdk:"autocreate"`
	UsernameClaim      types.String `tfsdk:"username_claim"`
	GroupsClaim        types.String `tfsdk:"groups_claim"`
	GroupsAutocreate   types.Bool   `tfsdk:"groups_autocreate"`
	GroupsOverwrite    types.Bool   `tfsdk:"groups_overwrite"`
	Scopes             types.String `tfsdk:"scopes"`
	Prompt             types.String `tfsdk:"prompt"`
	ACRValues          types.String `tfsdk:"acr_values"`
	Audiences          types.String `tfsdk:"audiences"`
	QueryUserinfo      types.Bool   `tfsdk:"query_userinfo"`
	Comment            types.String `tfsdk:"comment"`
	Default            types.Bool   `tfsdk:"default"`
}

func (m *realmOpenIDModel) toCreateRequest(clientKeyWO types.String) *access.RealmCreateRequestBody {
	req := &access.RealmCreateRequestBody{
		Realm:     m.Realm.ValueString(),
		Type:      "openid",
		IssuerURL: m.IssuerURL.ValueStringPointer(),
		ClientID:  m.ClientID.ValueStringPointer(),
	}

	// client_key (state-persisted) and client_key_wo (write-only) are mutually
	// exclusive (enforced by ConfigValidators); prefer the write-only value.
	switch {
	case !clientKeyWO.IsNull():
		req.ClientKey = clientKeyWO.ValueStringPointer()
	case !m.ClientKey.IsNull():
		req.ClientKey = m.ClientKey.ValueStringPointer()
	}

	if !m.AutoCreate.IsNull() {
		req.AutoCreate = proxmoxtypes.CustomBoolPtr(m.AutoCreate.ValueBoolPointer())
	}

	if !m.UsernameClaim.IsNull() {
		req.UsernameClaim = m.UsernameClaim.ValueStringPointer()
	}

	if !m.GroupsClaim.IsNull() {
		req.GroupsClaim = m.GroupsClaim.ValueStringPointer()
	}

	if !m.GroupsAutocreate.IsNull() {
		req.GroupsAutocreate = proxmoxtypes.CustomBoolPtr(m.GroupsAutocreate.ValueBoolPointer())
	}

	if !m.GroupsOverwrite.IsNull() {
		req.GroupsOverwrite = proxmoxtypes.CustomBoolPtr(m.GroupsOverwrite.ValueBoolPointer())
	}

	if !m.Scopes.IsNull() {
		req.Scopes = m.Scopes.ValueStringPointer()
	}

	if !m.Prompt.IsNull() {
		req.Prompt = m.Prompt.ValueStringPointer()
	}

	if !m.ACRValues.IsNull() {
		req.ACRValues = m.ACRValues.ValueStringPointer()
	}

	if !m.Audiences.IsNull() {
		req.Audiences = m.Audiences.ValueStringPointer()
	}

	if !m.QueryUserinfo.IsNull() {
		req.QueryUserinfo = proxmoxtypes.CustomBoolPtr(m.QueryUserinfo.ValueBoolPointer())
	}

	if !m.Comment.IsNull() {
		req.Comment = m.Comment.ValueStringPointer()
	}

	if !m.Default.IsNull() {
		req.Default = proxmoxtypes.CustomBoolPtr(m.Default.ValueBoolPointer())
	}

	return req
}

func (m *realmOpenIDModel) toUpdateRequest(state *realmOpenIDModel, clientKeyWO types.String) *access.RealmUpdateRequestBody {
	req := &access.RealmUpdateRequestBody{}
	var toDelete []string

	// Required fields: update directly.
	if !m.IssuerURL.Equal(state.IssuerURL) {
		req.IssuerURL = m.IssuerURL.ValueStringPointer()
	}

	if !m.ClientID.Equal(state.ClientID) {
		req.ClientID = m.ClientID.ValueStringPointer()
	}

	// Client key: when supplied via the write-only client_key_wo, resend it on every
	// update (write-only values are invisible to diffs, so rotation is driven by the
	// client_key_wo_version bump). Otherwise diff the state-tracked client_key.
	if !clientKeyWO.IsNull() {
		req.ClientKey = clientKeyWO.ValueStringPointer()
	} else {
		updateStringAttribute(&req.ClientKey, m.ClientKey, state.ClientKey, &toDelete, "client-key")

		// Write-only client_key_wo is never mirrored into state, so the version counter
		// leaving state is the removal signal.
		if m.ClientKey.IsNull() && m.ClientKeyWOVersion.IsNull() && !state.ClientKeyWOVersion.IsNull() {
			toDelete = append(toDelete, "client-key")
		}
	}

	// Optional fields: support unsetting using the API's `delete` parameter.
	updateStringAttribute(&req.Scopes, m.Scopes, state.Scopes, &toDelete, "scopes")
	updateStringAttribute(&req.Prompt, m.Prompt, state.Prompt, &toDelete, "prompt")
	updateStringAttribute(&req.ACRValues, m.ACRValues, state.ACRValues, &toDelete, "acr-values")
	updateStringAttribute(&req.Audiences, m.Audiences, state.Audiences, &toDelete, "audiences")
	updateStringAttribute(&req.GroupsClaim, m.GroupsClaim, state.GroupsClaim, &toDelete, "groups-claim")
	updateStringAttribute(&req.Comment, m.Comment, state.Comment, &toDelete, "comment")
	updateStringAttribute(&req.UsernameClaim, m.UsernameClaim, state.UsernameClaim, &toDelete, "username-claim")

	// Booleans are sent on change.
	if !m.AutoCreate.Equal(state.AutoCreate) {
		req.AutoCreate = proxmoxtypes.CustomBoolPtr(m.AutoCreate.ValueBoolPointer())
	}

	if !m.GroupsAutocreate.Equal(state.GroupsAutocreate) {
		req.GroupsAutocreate = proxmoxtypes.CustomBoolPtr(m.GroupsAutocreate.ValueBoolPointer())
	}

	if !m.GroupsOverwrite.Equal(state.GroupsOverwrite) {
		req.GroupsOverwrite = proxmoxtypes.CustomBoolPtr(m.GroupsOverwrite.ValueBoolPointer())
	}

	if !m.QueryUserinfo.Equal(state.QueryUserinfo) {
		req.QueryUserinfo = proxmoxtypes.CustomBoolPtr(m.QueryUserinfo.ValueBoolPointer())
	}

	if !m.Default.Equal(state.Default) {
		req.Default = proxmoxtypes.CustomBoolPtr(m.Default.ValueBoolPointer())
	}

	if len(toDelete) > 0 {
		req.Delete = toDelete
	}

	return req
}

func (m *realmOpenIDModel) fromAPIResponse(data *access.RealmGetResponseData, diags *diag.Diagnostics) {
	// Validate required fields
	if data.IssuerURL == nil {
		diags.AddError(
			"Missing Required Field",
			"API response is missing required field 'issuer-url' for OpenID realm",
		)

		return
	}

	if data.ClientID == nil {
		diags.AddError(
			"Missing Required Field",
			"API response is missing required field 'client-id' for OpenID realm",
		)

		return
	}

	// Set required fields
	m.IssuerURL = types.StringPointerValue(data.IssuerURL)
	m.ClientID = types.StringPointerValue(data.ClientID)

	// client_key is deliberately not mirrored from the API response (which does return
	// it): when the key is supplied via the write-only client_key_wo it must stay null
	// in state, so preserve whatever state already holds.

	// Set optional string fields
	m.UsernameClaim = types.StringPointerValue(data.UsernameClaim)
	m.GroupsClaim = types.StringPointerValue(data.GroupsClaim)
	m.Prompt = types.StringPointerValue(data.Prompt)
	m.ACRValues = types.StringPointerValue(data.ACRValues)
	m.Audiences = types.StringPointerValue(data.Audiences)
	m.Comment = types.StringPointerValue(data.Comment)

	m.Scopes = types.StringPointerValue(data.Scopes)

	// Set optional boolean fields
	m.AutoCreate = types.BoolPointerValue(data.AutoCreate.PointerBool())
	m.GroupsAutocreate = types.BoolPointerValue(data.GroupsAutocreate.PointerBool())
	m.GroupsOverwrite = types.BoolPointerValue(data.GroupsOverwrite.PointerBool())
	m.QueryUserinfo = types.BoolPointerValue(data.QueryUserinfo.PointerBool())
	m.Default = types.BoolPointerValue(data.Default.PointerBool())
}
