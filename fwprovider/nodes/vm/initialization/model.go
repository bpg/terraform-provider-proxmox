/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package initialization

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// ---- nested model types ----

type dnsModel struct {
	Domain  types.String `tfsdk:"domain"`
	Servers types.List   `tfsdk:"servers"`
}

type ipConfigModel struct {
	IPv4Address types.String `tfsdk:"ipv4_address"`
	IPv4Gateway types.String `tfsdk:"ipv4_gateway"`
	IPv6Address types.String `tfsdk:"ipv6_address"`
	IPv6Gateway types.String `tfsdk:"ipv6_gateway"`
}

// userAccountModel is the resource-side user_account model (includes write-only password).
type userAccountModel struct {
	Keys     types.List   `tfsdk:"keys"`
	Password types.String `tfsdk:"password"`
	Username types.String `tfsdk:"username"`
}

// dsUserAccountModel is the datasource-side user_account model (no password — never returned by API).
type dsUserAccountModel struct {
	Keys     types.List   `tfsdk:"keys"`
	Username types.String `tfsdk:"username"`
}

// ---- attribute type maps ----

func dnsAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"domain":  types.StringType,
		"servers": types.ListType{ElemType: types.StringType},
	}
}

func ipConfigAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ipv4_address": types.StringType,
		"ipv4_gateway": types.StringType,
		"ipv6_address": types.StringType,
		"ipv6_gateway": types.StringType,
	}
}

func userAccountAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"keys":     types.ListType{ElemType: types.StringType},
		"password": types.StringType,
		"username": types.StringType,
	}
}

func dsUserAccountAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"keys":     types.ListType{ElemType: types.StringType},
		"username": types.StringType,
	}
}

// attributeTypes returns the attribute types for the resource-side initialization Value.
func attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"dns":                  types.ObjectType{AttrTypes: dnsAttributeTypes()},
		"ip_config":            types.ListType{ElemType: types.ObjectType{AttrTypes: ipConfigAttributeTypes()}},
		"meta_data_file_id":    types.StringType,
		"network_data_file_id": types.StringType,
		"type":                 types.StringType,
		"upgrade":              types.BoolType,
		"user_account":         types.ObjectType{AttrTypes: userAccountAttributeTypes()},
		"user_data_file_id":    types.StringType,
		"vendor_data_file_id":  types.StringType,
	}
}

// dsAttributeTypes returns the attribute types for the datasource-side initialization Value.
func dsAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"dns":                  types.ObjectType{AttrTypes: dnsAttributeTypes()},
		"ip_config":            types.ListType{ElemType: types.ObjectType{AttrTypes: ipConfigAttributeTypes()}},
		"meta_data_file_id":    types.StringType,
		"network_data_file_id": types.StringType,
		"type":                 types.StringType,
		"upgrade":              types.BoolType,
		"user_account":         types.ObjectType{AttrTypes: dsUserAccountAttributeTypes()},
		"user_data_file_id":    types.StringType,
		"vendor_data_file_id":  types.StringType,
	}
}

// ---- resource-side Model ----

// Model represents the initialization block for a resource.
type Model struct {
	DNS               types.Object `tfsdk:"dns"`
	IPConfig          types.List   `tfsdk:"ip_config"`
	MetaDataFileID    types.String `tfsdk:"meta_data_file_id"`
	NetworkDataFileID types.String `tfsdk:"network_data_file_id"`
	Type              types.String `tfsdk:"type"`
	Upgrade           types.Bool   `tfsdk:"upgrade"`
	UserAccount       types.Object `tfsdk:"user_account"`
	UserDataFileID    types.String `tfsdk:"user_data_file_id"`
	VendorDataFileID  types.String `tfsdk:"vendor_data_file_id"`
}

// NullValue returns a properly typed null Value for use in the resource model.
func NullValue() Value {
	return types.ObjectNull(attributeTypes())
}

// NullDataSourceValue returns a properly typed null Value for use in the datasource model.
func NullDataSourceValue() DataSourceValue {
	return types.ObjectNull(dsAttributeTypes())
}

// toAPI writes cloud-init configuration onto the create/update request body.
// The password field is included when set; since it is write-only it is never
// read back from state, so the caller is responsible for passing the plan value.
func (m *Model) toAPI(ctx context.Context, body *vms.CreateRequestBody, diags *diag.Diagnostics) {
	ci := &vms.CustomCloudInitConfig{}

	// DNS
	if attribute.IsDefined(m.DNS) {
		var dns dnsModel

		diags.Append(m.DNS.As(ctx, &dns, basetypes.ObjectAsOptions{})...)

		if !diags.HasError() {
			ci.SearchDomain = attribute.StringPtrFromValue(dns.Domain)

			if attribute.IsDefined(dns.Servers) {
				var servers []string

				diags.Append(dns.Servers.ElementsAs(ctx, &servers, false)...)

				if len(servers) > 0 {
					ns := strings.Join(servers, " ")
					ci.Nameserver = &ns
				}
			}
		}
	}

	// IP configuration
	if attribute.IsDefined(m.IPConfig) {
		var ipConfigs []ipConfigModel

		diags.Append(m.IPConfig.ElementsAs(ctx, &ipConfigs, false)...)

		for _, c := range ipConfigs {
			ci.IPConfig = append(ci.IPConfig, vms.CustomCloudInitIPConfig{
				IPv4:        attribute.StringPtrFromValue(c.IPv4Address),
				GatewayIPv4: attribute.StringPtrFromValue(c.IPv4Gateway),
				IPv6:        attribute.StringPtrFromValue(c.IPv6Address),
				GatewayIPv6: attribute.StringPtrFromValue(c.IPv6Gateway),
			})
		}
	}

	// User account
	if attribute.IsDefined(m.UserAccount) {
		var ua userAccountModel

		diags.Append(m.UserAccount.As(ctx, &ua, basetypes.ObjectAsOptions{})...)

		if !diags.HasError() {
			ci.Username = attribute.StringPtrFromValue(ua.Username)
			ci.Password = attribute.StringPtrFromValue(ua.Password)

			if attribute.IsDefined(ua.Keys) {
				var keys []string

				diags.Append(ua.Keys.ElementsAs(ctx, &keys, false)...)

				if len(keys) > 0 {
					sshKeys := vms.CustomCloudInitSSHKeys(keys)
					ci.SSHKeys = &sshKeys
				}
			}
		}
	}

	// Custom data files
	if attribute.IsDefined(m.UserDataFileID) || attribute.IsDefined(m.VendorDataFileID) ||
		attribute.IsDefined(m.NetworkDataFileID) || attribute.IsDefined(m.MetaDataFileID) {
		ci.Files = &vms.CustomCloudInitFiles{
			UserVolume:    attribute.StringPtrFromValue(m.UserDataFileID),
			VendorVolume:  attribute.StringPtrFromValue(m.VendorDataFileID),
			NetworkVolume: attribute.StringPtrFromValue(m.NetworkDataFileID),
			MetaVolume:    attribute.StringPtrFromValue(m.MetaDataFileID),
		}
	}

	// Scalar fields
	ci.Type = attribute.StringPtrFromValue(m.Type)

	if attribute.IsDefined(m.Upgrade) {
		v := proxmoxtypes.CustomBool(m.Upgrade.ValueBool())
		ci.Upgrade = &v
	}

	body.CloudInitConfig = ci
}

// fromAPI populates the resource Model from the PVE API response.
// Password is always left null — the API returns a masked value, not the plaintext.
func (m *Model) fromAPI(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) {
	// DNS
	if config.CloudInitDNSDomain != nil || config.CloudInitDNSServer != nil {
		dns := dnsModel{
			Domain: types.StringPointerValue(config.CloudInitDNSDomain),
		}

		if config.CloudInitDNSServer != nil {
			parts := strings.Fields(*config.CloudInitDNSServer)

			servers, d := types.ListValueFrom(ctx, types.StringType, parts)
			diags.Append(d...)

			dns.Servers = servers
		} else {
			dns.Servers = types.ListNull(types.StringType)
		}

		obj, d := types.ObjectValueFrom(ctx, dnsAttributeTypes(), dns)
		diags.Append(d...)

		m.DNS = obj
	} else {
		m.DNS = types.ObjectNull(dnsAttributeTypes())
	}

	// IP configuration — convert the map (ipconfig0..7) to an ordered list.
	if len(config.IPConfigs) > 0 {
		maxIdx := -1

		for key := range config.IPConfigs {
			var idx int

			fmt.Sscanf(key, "ipconfig%d", &idx) //nolint:errcheck

			if idx > maxIdx {
				maxIdx = idx
			}
		}

		ipConfigs := make([]ipConfigModel, maxIdx+1)

		for i := range ipConfigs {
			key := fmt.Sprintf("ipconfig%d", i)

			if c, ok := config.IPConfigs[key]; ok && c != nil {
				ipConfigs[i] = ipConfigModel{
					IPv4Address: types.StringPointerValue(c.IPv4),
					IPv4Gateway: types.StringPointerValue(c.GatewayIPv4),
					IPv6Address: types.StringPointerValue(c.IPv6),
					IPv6Gateway: types.StringPointerValue(c.GatewayIPv6),
				}
			} else {
				ipConfigs[i] = ipConfigModel{
					IPv4Address: types.StringNull(),
					IPv4Gateway: types.StringNull(),
					IPv6Address: types.StringNull(),
					IPv6Gateway: types.StringNull(),
				}
			}
		}

		list, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: ipConfigAttributeTypes()}, ipConfigs)
		diags.Append(d...)

		m.IPConfig = list
	} else {
		m.IPConfig = types.ListNull(types.ObjectType{AttrTypes: ipConfigAttributeTypes()})
	}

	// User account — password is always null (API returns "**********", never plaintext)
	if config.CloudInitUsername != nil || config.CloudInitSSHKeys != nil {
		ua := userAccountModel{
			Username: types.StringPointerValue(config.CloudInitUsername),
			Password: types.StringNull(),
		}

		if config.CloudInitSSHKeys != nil && len(*config.CloudInitSSHKeys) > 0 {
			keys, d := types.ListValueFrom(ctx, types.StringType, []string(*config.CloudInitSSHKeys))
			diags.Append(d...)

			ua.Keys = keys
		} else {
			ua.Keys = types.ListNull(types.StringType)
		}

		obj, d := types.ObjectValueFrom(ctx, userAccountAttributeTypes(), ua)
		diags.Append(d...)

		m.UserAccount = obj
	} else {
		m.UserAccount = types.ObjectNull(userAccountAttributeTypes())
	}

	// Custom data file IDs
	if config.CloudInitFiles != nil {
		m.UserDataFileID = types.StringPointerValue(config.CloudInitFiles.UserVolume)
		m.VendorDataFileID = types.StringPointerValue(config.CloudInitFiles.VendorVolume)
		m.NetworkDataFileID = types.StringPointerValue(config.CloudInitFiles.NetworkVolume)
		m.MetaDataFileID = types.StringPointerValue(config.CloudInitFiles.MetaVolume)
	} else {
		m.UserDataFileID = types.StringNull()
		m.VendorDataFileID = types.StringNull()
		m.NetworkDataFileID = types.StringNull()
		m.MetaDataFileID = types.StringNull()
	}

	// Scalar fields
	m.Type = types.StringPointerValue(config.CloudInitType)

	if config.CloudInitUpgrade != nil {
		m.Upgrade = types.BoolValue(bool(*config.CloudInitUpgrade))
	} else {
		m.Upgrade = types.BoolNull()
	}
}

// fromAPIForDatasource populates a datasource-compatible Value from the API response.
// Identical to fromAPI but produces an object using dsUserAccountAttributeTypes (no password).
func fromAPIForDatasource(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) DataSourceValue {
	if !hasCloudInitData(config) {
		return NullDataSourceValue()
	}

	var m struct {
		DNS               types.Object `tfsdk:"dns"`
		IPConfig          types.List   `tfsdk:"ip_config"`
		MetaDataFileID    types.String `tfsdk:"meta_data_file_id"`
		NetworkDataFileID types.String `tfsdk:"network_data_file_id"`
		Type              types.String `tfsdk:"type"`
		Upgrade           types.Bool   `tfsdk:"upgrade"`
		UserAccount       types.Object `tfsdk:"user_account"`
		UserDataFileID    types.String `tfsdk:"user_data_file_id"`
		VendorDataFileID  types.String `tfsdk:"vendor_data_file_id"`
	}

	// Reuse the resource Model for all fields except user_account.
	var rm Model
	rm.fromAPI(ctx, config, diags)

	m.DNS = rm.DNS
	m.IPConfig = rm.IPConfig
	m.MetaDataFileID = rm.MetaDataFileID
	m.NetworkDataFileID = rm.NetworkDataFileID
	m.Type = rm.Type
	m.Upgrade = rm.Upgrade
	m.UserDataFileID = rm.UserDataFileID
	m.VendorDataFileID = rm.VendorDataFileID

	// Build datasource user_account — no password field.
	if config.CloudInitUsername != nil || config.CloudInitSSHKeys != nil {
		dsUA := dsUserAccountModel{
			Username: types.StringPointerValue(config.CloudInitUsername),
		}

		if config.CloudInitSSHKeys != nil && len(*config.CloudInitSSHKeys) > 0 {
			keys, d := types.ListValueFrom(ctx, types.StringType, []string(*config.CloudInitSSHKeys))
			diags.Append(d...)

			dsUA.Keys = keys
		} else {
			dsUA.Keys = types.ListNull(types.StringType)
		}

		obj, d := types.ObjectValueFrom(ctx, dsUserAccountAttributeTypes(), dsUA)
		diags.Append(d...)

		m.UserAccount = obj
	} else {
		m.UserAccount = types.ObjectNull(dsUserAccountAttributeTypes())
	}

	obj, d := types.ObjectValueFrom(ctx, dsAttributeTypes(), m)
	diags.Append(d...)

	return obj
}

// hasCloudInitData reports whether any cloud-init field is set in the API response.
func hasCloudInitData(config *vms.GetResponseData) bool {
	return config.CloudInitUsername != nil ||
		config.CloudInitDNSDomain != nil ||
		config.CloudInitDNSServer != nil ||
		config.CloudInitFiles != nil ||
		config.CloudInitSSHKeys != nil ||
		config.CloudInitType != nil ||
		config.CloudInitUpgrade != nil ||
		len(config.IPConfigs) > 0
}
