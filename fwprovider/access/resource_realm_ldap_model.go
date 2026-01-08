/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type realmLDAPModel struct {
	ID                  types.String `tfsdk:"id"`
	Realm               types.String `tfsdk:"realm"`
	Server1             types.String `tfsdk:"server1"`
	Server2             types.String `tfsdk:"server2"`
	BaseDN              types.String `tfsdk:"base_dn"`
	BindDN              types.String `tfsdk:"bind_dn"`
	BindPassword        types.String `tfsdk:"bind_password"`
	UserAttr            types.String `tfsdk:"user_attr"`
	Port                types.Int64  `tfsdk:"port"`
	Secure              types.Bool   `tfsdk:"secure"`
	Verify              types.Bool   `tfsdk:"verify"`
	CaPath              types.String `tfsdk:"ca_path"`
	CertPath            types.String `tfsdk:"cert_path"`
	CertKeyPath         types.String `tfsdk:"cert_key_path"`
	Filter              types.String `tfsdk:"filter"`
	GroupDN             types.String `tfsdk:"group_dn"`
	GroupFilter         types.String `tfsdk:"group_filter"`
	GroupClasses        types.String `tfsdk:"group_classes"`
	GroupNameAttr       types.String `tfsdk:"group_name_attr"`
	Mode                types.String `tfsdk:"mode"`
	SSLVersion          types.String `tfsdk:"ssl_version"`
	UserClasses         types.String `tfsdk:"user_classes"`
	SyncAttributes      types.String `tfsdk:"sync_attributes"`
	SyncDefaultsOptions types.String `tfsdk:"sync_defaults_options"`
	Comment             types.String `tfsdk:"comment"`
	Default             types.Bool   `tfsdk:"default"`
	CaseSensitive       types.Bool   `tfsdk:"case_sensitive"`
}

// updateStringAttribute is a helper function for handling optional string attributes in updates.
// It sets the request field if the value changed, or adds to the delete list if being unset.
func updateStringAttribute(reqField **string, planVal, stateVal types.String, toDelete *[]string, apiName string) {
	if attribute.IsDefined(planVal) {
		if !planVal.Equal(stateVal) {
			*reqField = planVal.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(planVal, stateVal, toDelete, apiName)
	}
}

// updateInt64Attribute is a helper function for handling optional int64 attributes in updates.
// It sets the request field if the value changed, or adds to the delete list if being unset.
func updateInt64Attribute(reqField **int, planVal, stateVal types.Int64, toDelete *[]string, apiName string) {
	if attribute.IsDefined(planVal) {
		if !planVal.Equal(stateVal) {
			val := int(planVal.ValueInt64())
			*reqField = &val
		}
	} else {
		attribute.CheckDelete(planVal, stateVal, toDelete, apiName)
	}
}

func (m *realmLDAPModel) toCreateRequest() *access.RealmCreateRequestBody {
	req := &access.RealmCreateRequestBody{
		Realm:   m.Realm.ValueString(),
		Type:    "ldap",
		Server1: m.Server1.ValueStringPointer(),
		BaseDN:  m.BaseDN.ValueStringPointer(),
	}

	if !m.Server2.IsNull() {
		req.Server2 = m.Server2.ValueStringPointer()
	}

	if !m.BindDN.IsNull() {
		req.BindDN = m.BindDN.ValueStringPointer()
	}

	if !m.BindPassword.IsNull() {
		req.BindPassword = m.BindPassword.ValueStringPointer()
	}

	if !m.UserAttr.IsNull() {
		req.UserAttr = m.UserAttr.ValueStringPointer()
	}

	if !m.Port.IsNull() {
		port := int(m.Port.ValueInt64())
		req.Port = &port
	}

	if !m.Secure.IsNull() {
		req.Secure = proxmoxtypes.CustomBoolPtr(m.Secure.ValueBoolPointer())
	}

	if !m.Verify.IsNull() {
		req.Verify = proxmoxtypes.CustomBoolPtr(m.Verify.ValueBoolPointer())
	}

	if !m.CaPath.IsNull() {
		req.CaPath = m.CaPath.ValueStringPointer()
	}

	if !m.CertPath.IsNull() {
		req.CertPath = m.CertPath.ValueStringPointer()
	}

	if !m.CertKeyPath.IsNull() {
		req.CertKeyPath = m.CertKeyPath.ValueStringPointer()
	}

	if !m.Filter.IsNull() {
		req.Filter = m.Filter.ValueStringPointer()
	}

	if !m.GroupDN.IsNull() {
		req.GroupDN = m.GroupDN.ValueStringPointer()
	}

	if !m.GroupFilter.IsNull() {
		req.GroupFilter = m.GroupFilter.ValueStringPointer()
	}

	if !m.GroupClasses.IsNull() {
		req.GroupClasses = m.GroupClasses.ValueStringPointer()
	}

	if !m.GroupNameAttr.IsNull() {
		req.GroupNameAttr = m.GroupNameAttr.ValueStringPointer()
	}

	if !m.Mode.IsNull() {
		req.Mode = m.Mode.ValueStringPointer()
	}

	if !m.SSLVersion.IsNull() {
		req.SSLVersion = m.SSLVersion.ValueStringPointer()
	}

	if !m.UserClasses.IsNull() {
		req.UserClasses = m.UserClasses.ValueStringPointer()
	}

	if !m.SyncAttributes.IsNull() {
		req.SyncAttributes = m.SyncAttributes.ValueStringPointer()
	}

	if !m.SyncDefaultsOptions.IsNull() {
		req.SyncDefaultsOpts = m.SyncDefaultsOptions.ValueStringPointer()
	}

	if !m.Comment.IsNull() {
		req.Comment = m.Comment.ValueStringPointer()
	}

	if !m.Default.IsNull() {
		req.Default = proxmoxtypes.CustomBoolPtr(m.Default.ValueBoolPointer())
	}

	if !m.CaseSensitive.IsNull() {
		req.CaseSensitive = proxmoxtypes.CustomBoolPtr(m.CaseSensitive.ValueBoolPointer())
	}

	return req
}

func (m *realmLDAPModel) toUpdateRequest(state *realmLDAPModel) *access.RealmUpdateRequestBody {
	req := &access.RealmUpdateRequestBody{}
	var toDelete []string

	// Required fields: update directly.
	if !m.Server1.Equal(state.Server1) {
		req.Server1 = m.Server1.ValueStringPointer()
	}

	if !m.BaseDN.Equal(state.BaseDN) {
		req.BaseDN = m.BaseDN.ValueStringPointer()
	}

	// Optional fields: support unsetting using the API's `delete` parameter.
	updateStringAttribute(&req.Server2, m.Server2, state.Server2, &toDelete, "server2")
	updateStringAttribute(&req.BindDN, m.BindDN, state.BindDN, &toDelete, "bind_dn")
	// The API field name for BindPassword is `password`.
	updateStringAttribute(&req.BindPassword, m.BindPassword, state.BindPassword, &toDelete, "password")
	updateStringAttribute(&req.UserAttr, m.UserAttr, state.UserAttr, &toDelete, "user_attr")
	updateInt64Attribute(&req.Port, m.Port, state.Port, &toDelete, "port")

	// Booleans are sent on change (they are typically optional+computed).
	// Note: Secure is deprecated by Proxmox in favor of Mode, but is still
	// supported for backward compatibility.
	if !m.Secure.Equal(state.Secure) {
		req.Secure = proxmoxtypes.CustomBoolPtr(m.Secure.ValueBoolPointer())
	}

	if !m.Verify.Equal(state.Verify) {
		req.Verify = proxmoxtypes.CustomBoolPtr(m.Verify.ValueBoolPointer())
	}

	updateStringAttribute(&req.CaPath, m.CaPath, state.CaPath, &toDelete, "capath")
	updateStringAttribute(&req.CertPath, m.CertPath, state.CertPath, &toDelete, "cert")
	updateStringAttribute(&req.CertKeyPath, m.CertKeyPath, state.CertKeyPath, &toDelete, "certkey")
	updateStringAttribute(&req.Filter, m.Filter, state.Filter, &toDelete, "filter")
	updateStringAttribute(&req.GroupDN, m.GroupDN, state.GroupDN, &toDelete, "group_dn")
	updateStringAttribute(&req.GroupFilter, m.GroupFilter, state.GroupFilter, &toDelete, "group_filter")
	updateStringAttribute(&req.GroupClasses, m.GroupClasses, state.GroupClasses, &toDelete, "group_classes")
	updateStringAttribute(&req.GroupNameAttr, m.GroupNameAttr, state.GroupNameAttr, &toDelete, "group_name_attr")
	updateStringAttribute(&req.Mode, m.Mode, state.Mode, &toDelete, "mode")
	updateStringAttribute(&req.SSLVersion, m.SSLVersion, state.SSLVersion, &toDelete, "sslversion")
	updateStringAttribute(&req.UserClasses, m.UserClasses, state.UserClasses, &toDelete, "user_classes")
	updateStringAttribute(&req.SyncAttributes, m.SyncAttributes, state.SyncAttributes, &toDelete, "sync_attributes")
	updateStringAttribute(&req.SyncDefaultsOpts, m.SyncDefaultsOptions, state.SyncDefaultsOptions, &toDelete, "sync-defaults-options")
	updateStringAttribute(&req.Comment, m.Comment, state.Comment, &toDelete, "comment")

	if !m.Default.Equal(state.Default) {
		req.Default = proxmoxtypes.CustomBoolPtr(m.Default.ValueBoolPointer())
	}

	if !m.CaseSensitive.Equal(state.CaseSensitive) {
		req.CaseSensitive = proxmoxtypes.CustomBoolPtr(m.CaseSensitive.ValueBoolPointer())
	}

	if len(toDelete) > 0 {
		deleteStr := strings.Join(toDelete, ",")
		req.Delete = &deleteStr
	}

	return req
}

func (m *realmLDAPModel) fromAPIResponse(data *access.RealmGetResponseData, diags *diag.Diagnostics) {
	// Validate required fields
	if data.Server1 == nil {
		diags.AddError(
			"Missing Required Field",
			"API response is missing required field 'server1' for LDAP realm",
		)

		return
	}

	if data.BaseDN == nil {
		diags.AddError(
			"Missing Required Field",
			"API response is missing required field 'base_dn' for LDAP realm",
		)

		return
	}

	// Set required fields
	m.Server1 = types.StringPointerValue(data.Server1)
	m.BaseDN = types.StringPointerValue(data.BaseDN)

	// Set optional fields
	if data.Server2 != nil {
		m.Server2 = types.StringPointerValue(data.Server2)
	} else {
		m.Server2 = types.StringNull()
	}

	if data.BindDN != nil {
		m.BindDN = types.StringPointerValue(data.BindDN)
	} else {
		m.BindDN = types.StringNull()
	}

	// Note: bind_password is never returned by the API, preserve from state

	if data.UserAttr != nil {
		m.UserAttr = types.StringPointerValue(data.UserAttr)
	} else {
		m.UserAttr = types.StringValue("uid") // default value
	}

	if data.Port != nil {
		m.Port = types.Int64Value(int64(*data.Port))
	} else {
		m.Port = types.Int64Null()
	}

	if data.Secure != nil {
		m.Secure = types.BoolValue(bool(*data.Secure))
	} else {
		m.Secure = types.BoolValue(false)
	}

	if data.Verify != nil {
		m.Verify = types.BoolValue(bool(*data.Verify))
	} else {
		m.Verify = types.BoolValue(false)
	}

	if data.CaPath != nil {
		m.CaPath = types.StringPointerValue(data.CaPath)
	} else {
		m.CaPath = types.StringNull()
	}

	if data.CertPath != nil {
		m.CertPath = types.StringPointerValue(data.CertPath)
	} else {
		m.CertPath = types.StringNull()
	}

	if data.CertKeyPath != nil {
		m.CertKeyPath = types.StringPointerValue(data.CertKeyPath)
	} else {
		m.CertKeyPath = types.StringNull()
	}

	if data.Filter != nil {
		m.Filter = types.StringPointerValue(data.Filter)
	} else {
		m.Filter = types.StringNull()
	}

	if data.GroupDN != nil {
		m.GroupDN = types.StringPointerValue(data.GroupDN)
	} else {
		m.GroupDN = types.StringNull()
	}

	if data.GroupFilter != nil {
		m.GroupFilter = types.StringPointerValue(data.GroupFilter)
	} else {
		m.GroupFilter = types.StringNull()
	}

	if data.GroupClasses != nil {
		m.GroupClasses = types.StringPointerValue(data.GroupClasses)
	} else {
		m.GroupClasses = types.StringNull()
	}

	if data.GroupNameAttr != nil {
		m.GroupNameAttr = types.StringPointerValue(data.GroupNameAttr)
	} else {
		m.GroupNameAttr = types.StringNull()
	}

	if data.Mode != nil {
		m.Mode = types.StringPointerValue(data.Mode)
	} else {
		m.Mode = types.StringNull()
	}

	if data.SSLVersion != nil {
		m.SSLVersion = types.StringPointerValue(data.SSLVersion)
	} else {
		m.SSLVersion = types.StringNull()
	}

	if data.UserClasses != nil {
		m.UserClasses = types.StringPointerValue(data.UserClasses)
	} else {
		m.UserClasses = types.StringNull()
	}

	if data.SyncAttributes != nil {
		m.SyncAttributes = types.StringPointerValue(data.SyncAttributes)
	} else {
		m.SyncAttributes = types.StringNull()
	}

	if data.SyncDefaultsOpts != nil {
		m.SyncDefaultsOptions = types.StringPointerValue(data.SyncDefaultsOpts)
	} else {
		m.SyncDefaultsOptions = types.StringNull()
	}

	if data.Comment != nil {
		m.Comment = types.StringPointerValue(data.Comment)
	} else {
		m.Comment = types.StringNull()
	}

	if data.Default != nil {
		m.Default = types.BoolValue(bool(*data.Default))
	} else {
		m.Default = types.BoolValue(false)
	}

	if data.CaseSensitive != nil {
		m.CaseSensitive = types.BoolValue(bool(*data.CaseSensitive))
	} else {
		m.CaseSensitive = types.BoolValue(true) // default value
	}
}
