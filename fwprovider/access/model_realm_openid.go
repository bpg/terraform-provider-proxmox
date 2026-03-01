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
	ID               types.String `tfsdk:"id"`
	Realm            types.String `tfsdk:"realm"`
	IssuerURL        types.String `tfsdk:"issuer_url"`
	ClientID         types.String `tfsdk:"client_id"`
	ClientKey        types.String `tfsdk:"client_key"`
	AutoCreate       types.Bool   `tfsdk:"autocreate"`
	UsernameClaim    types.String `tfsdk:"username_claim"`
	GroupsClaim      types.String `tfsdk:"groups_claim"`
	GroupsAutocreate types.Bool   `tfsdk:"groups_autocreate"`
	GroupsOverwrite  types.Bool   `tfsdk:"groups_overwrite"`
	Scopes           types.String `tfsdk:"scopes"`
	Prompt           types.String `tfsdk:"prompt"`
	ACRValues        types.String `tfsdk:"acr_values"`
	QueryUserinfo    types.Bool   `tfsdk:"query_userinfo"`
	Comment          types.String `tfsdk:"comment"`
	Default          types.Bool   `tfsdk:"default"`
}

func (m *realmOpenIDModel) toCreateRequest() *access.RealmCreateRequestBody {
	req := &access.RealmCreateRequestBody{
		Realm:     m.Realm.ValueString(),
		Type:      "openid",
		IssuerURL: m.IssuerURL.ValueStringPointer(),
		ClientID:  m.ClientID.ValueStringPointer(),
	}

	if !m.ClientKey.IsNull() {
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

func (m *realmOpenIDModel) toUpdateRequest(state *realmOpenIDModel) *access.RealmUpdateRequestBody {
	req := &access.RealmUpdateRequestBody{}
	var toDelete []string

	// Required fields: update directly.
	if !m.IssuerURL.Equal(state.IssuerURL) {
		req.IssuerURL = m.IssuerURL.ValueStringPointer()
	}

	if !m.ClientID.Equal(state.ClientID) {
		req.ClientID = m.ClientID.ValueStringPointer()
	}

	// Optional fields: support unsetting using the API's `delete` parameter.
	updateStringAttribute(&req.ClientKey, m.ClientKey, state.ClientKey, &toDelete, "client-key")
	updateStringAttribute(&req.Scopes, m.Scopes, state.Scopes, &toDelete, "scopes")
	updateStringAttribute(&req.Prompt, m.Prompt, state.Prompt, &toDelete, "prompt")
	updateStringAttribute(&req.ACRValues, m.ACRValues, state.ACRValues, &toDelete, "acr-values")
	updateStringAttribute(&req.GroupsClaim, m.GroupsClaim, state.GroupsClaim, &toDelete, "groups-claim")
	updateStringAttribute(&req.Comment, m.Comment, state.Comment, &toDelete, "comment")

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

	// Note: client_key is never returned by the API, preserve from state

	// Set optional string fields
	m.UsernameClaim = types.StringPointerValue(data.UsernameClaim)
	m.GroupsClaim = types.StringPointerValue(data.GroupsClaim)
	m.Scopes = types.StringPointerValue(data.Scopes)
	m.Prompt = types.StringPointerValue(data.Prompt)
	m.ACRValues = types.StringPointerValue(data.ACRValues)
	m.Comment = types.StringPointerValue(data.Comment)

	// Set optional boolean fields with defaults
	m.AutoCreate = types.BoolPointerValue(data.AutoCreate.PointerBool())
	if m.AutoCreate.IsNull() {
		m.AutoCreate = types.BoolValue(false)
	}

	m.GroupsAutocreate = types.BoolPointerValue(data.GroupsAutocreate.PointerBool())
	if m.GroupsAutocreate.IsNull() {
		m.GroupsAutocreate = types.BoolValue(false)
	}

	m.GroupsOverwrite = types.BoolPointerValue(data.GroupsOverwrite.PointerBool())
	if m.GroupsOverwrite.IsNull() {
		m.GroupsOverwrite = types.BoolValue(false)
	}

	m.QueryUserinfo = types.BoolPointerValue(data.QueryUserinfo.PointerBool())
	if m.QueryUserinfo.IsNull() {
		m.QueryUserinfo = types.BoolValue(true)
	}

	m.Default = types.BoolPointerValue(data.Default.PointerBool())
	if m.Default.IsNull() {
		m.Default = types.BoolValue(false)
	}
}
