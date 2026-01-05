/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"strings"

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
	CaPath              types.String `tfsdk:"capath"`
	Cert                types.String `tfsdk:"cert"`
	CertKey             types.String `tfsdk:"certkey"`
	Filter              types.String `tfsdk:"filter"`
	GroupDN             types.String `tfsdk:"group_dn"`
	GroupFilter         types.String `tfsdk:"group_filter"`
	GroupClasses        types.String `tfsdk:"group_classes"`
	GroupNameAttr       types.String `tfsdk:"group_name_attr"`
	Mode                types.String `tfsdk:"mode"`
	SSLVersion          types.String `tfsdk:"sslversion"`
	UserClasses         types.String `tfsdk:"user_classes"`
	SyncAttributes      types.String `tfsdk:"sync_attributes"`
	SyncDefaultsOptions types.String `tfsdk:"sync_defaults_options"`
	Comment             types.String `tfsdk:"comment"`
	Default             types.Bool   `tfsdk:"default"`
	CaseSensitive       types.Bool   `tfsdk:"case_sensitive"`
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

	if !m.Cert.IsNull() {
		req.CertPath = m.Cert.ValueStringPointer()
	}

	if !m.CertKey.IsNull() {
		req.CertKeyPath = m.CertKey.ValueStringPointer()
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
	if attribute.IsDefined(m.Server2) {
		if !m.Server2.Equal(state.Server2) {
			req.Server2 = m.Server2.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.Server2, state.Server2, &toDelete, "server2")
	}

	if attribute.IsDefined(m.BindDN) {
		if !m.BindDN.Equal(state.BindDN) {
			req.BindDN = m.BindDN.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.BindDN, state.BindDN, &toDelete, "bind_dn")
	}

	if attribute.IsDefined(m.BindPassword) {
		if !m.BindPassword.Equal(state.BindPassword) {
			req.BindPassword = m.BindPassword.ValueStringPointer()
		}
	} else {
		// The API field name is `password`.
		attribute.CheckDelete(m.BindPassword, state.BindPassword, &toDelete, "password")
	}

	if attribute.IsDefined(m.UserAttr) {
		if !m.UserAttr.Equal(state.UserAttr) {
			req.UserAttr = m.UserAttr.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.UserAttr, state.UserAttr, &toDelete, "user_attr")
	}

	if attribute.IsDefined(m.Port) {
		if !m.Port.Equal(state.Port) {
			port := int(m.Port.ValueInt64())
			req.Port = &port
		}
	} else {
		attribute.CheckDelete(m.Port, state.Port, &toDelete, "port")
	}

	// Booleans are sent on change (they are typically optional+computed).
	// Note: Secure is deprecated by Proxmox in favor of Mode, but is still
	// supported for backward compatibility.
	if !m.Secure.Equal(state.Secure) {
		req.Secure = proxmoxtypes.CustomBoolPtr(m.Secure.ValueBoolPointer())
	}

	if !m.Verify.Equal(state.Verify) {
		req.Verify = proxmoxtypes.CustomBoolPtr(m.Verify.ValueBoolPointer())
	}

	if attribute.IsDefined(m.CaPath) {
		if !m.CaPath.Equal(state.CaPath) {
			req.CaPath = m.CaPath.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.CaPath, state.CaPath, &toDelete, "capath")
	}

	if attribute.IsDefined(m.Cert) {
		if !m.Cert.Equal(state.Cert) {
			req.CertPath = m.Cert.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.Cert, state.Cert, &toDelete, "cert")
	}

	if attribute.IsDefined(m.CertKey) {
		if !m.CertKey.Equal(state.CertKey) {
			req.CertKeyPath = m.CertKey.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.CertKey, state.CertKey, &toDelete, "certkey")
	}

	if attribute.IsDefined(m.Filter) {
		if !m.Filter.Equal(state.Filter) {
			req.Filter = m.Filter.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.Filter, state.Filter, &toDelete, "filter")
	}

	if attribute.IsDefined(m.GroupDN) {
		if !m.GroupDN.Equal(state.GroupDN) {
			req.GroupDN = m.GroupDN.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.GroupDN, state.GroupDN, &toDelete, "group_dn")
	}

	if attribute.IsDefined(m.GroupFilter) {
		if !m.GroupFilter.Equal(state.GroupFilter) {
			req.GroupFilter = m.GroupFilter.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.GroupFilter, state.GroupFilter, &toDelete, "group_filter")
	}

	if attribute.IsDefined(m.GroupClasses) {
		if !m.GroupClasses.Equal(state.GroupClasses) {
			req.GroupClasses = m.GroupClasses.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.GroupClasses, state.GroupClasses, &toDelete, "group_classes")
	}

	if attribute.IsDefined(m.GroupNameAttr) {
		if !m.GroupNameAttr.Equal(state.GroupNameAttr) {
			req.GroupNameAttr = m.GroupNameAttr.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.GroupNameAttr, state.GroupNameAttr, &toDelete, "group_name_attr")
	}

	if attribute.IsDefined(m.Mode) {
		if !m.Mode.Equal(state.Mode) {
			req.Mode = m.Mode.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.Mode, state.Mode, &toDelete, "mode")
	}

	if attribute.IsDefined(m.SSLVersion) {
		if !m.SSLVersion.Equal(state.SSLVersion) {
			req.SSLVersion = m.SSLVersion.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.SSLVersion, state.SSLVersion, &toDelete, "sslversion")
	}

	if attribute.IsDefined(m.UserClasses) {
		if !m.UserClasses.Equal(state.UserClasses) {
			req.UserClasses = m.UserClasses.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.UserClasses, state.UserClasses, &toDelete, "user_classes")
	}

	if attribute.IsDefined(m.SyncAttributes) {
		if !m.SyncAttributes.Equal(state.SyncAttributes) {
			req.SyncAttributes = m.SyncAttributes.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.SyncAttributes, state.SyncAttributes, &toDelete, "sync_attributes")
	}

	if attribute.IsDefined(m.SyncDefaultsOptions) {
		if !m.SyncDefaultsOptions.Equal(state.SyncDefaultsOptions) {
			req.SyncDefaultsOpts = m.SyncDefaultsOptions.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.SyncDefaultsOptions, state.SyncDefaultsOptions, &toDelete, "sync-defaults-options")
	}

	if attribute.IsDefined(m.Comment) {
		if !m.Comment.Equal(state.Comment) {
			req.Comment = m.Comment.ValueStringPointer()
		}
	} else {
		attribute.CheckDelete(m.Comment, state.Comment, &toDelete, "comment")
	}

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

func (m *realmLDAPModel) fromAPIResponse(data *access.RealmGetResponseData) {
	if data.Server1 != nil {
		m.Server1 = types.StringPointerValue(data.Server1)
	} else {
		m.Server1 = types.StringNull()
	}

	if data.Server2 != nil {
		m.Server2 = types.StringPointerValue(data.Server2)
	} else {
		m.Server2 = types.StringNull()
	}

	if data.BaseDN != nil {
		m.BaseDN = types.StringPointerValue(data.BaseDN)
	} else {
		m.BaseDN = types.StringNull()
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
		m.Cert = types.StringPointerValue(data.CertPath)
	} else {
		m.Cert = types.StringNull()
	}

	if data.CertKeyPath != nil {
		m.CertKey = types.StringPointerValue(data.CertKeyPath)
	} else {
		m.CertKey = types.StringNull()
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
