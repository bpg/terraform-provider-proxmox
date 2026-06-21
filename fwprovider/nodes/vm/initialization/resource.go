/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package initialization

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Value is the resource-side type for the initialization block (includes password field).
type Value = types.Object

// DataSourceValue is the datasource-side type for the initialization block (no password field).
type DataSourceValue = types.Object

// NewValue returns a Value populated from the PVE API response, or NullValue() when
// no cloud-init configuration is present.
func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value {
	if !hasCloudInitData(config) {
		return NullValue()
	}

	var m Model

	m.fromAPI(ctx, config, diags)

	obj, d := types.ObjectValueFrom(ctx, attributeTypes(), m)
	diags.Append(d...)

	return obj
}

// NewDataSourceValue returns a DataSourceValue populated from the PVE API response,
// or NullDataSourceValue() when no cloud-init configuration is present.
func NewDataSourceValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) DataSourceValue {
	return fromAPIForDatasource(ctx, config, diags)
}

// FillCreateBody writes the initialization block into the VM create request body.
func FillCreateBody(ctx context.Context, planValue Value, body *vms.CreateRequestBody, diags *diag.Diagnostics) {
	if planValue.IsNull() || planValue.IsUnknown() {
		return
	}

	var plan Model

	diags.Append(planValue.As(ctx, &plan, basetypes.ObjectAsOptions{})...)

	if diags.HasError() {
		return
	}

	plan.toAPI(ctx, body, diags)
}

// FillUpdateBody writes the initialization block diff into the VM update request body.
//
// When the whole block is removed from the plan, all cloud-init API keys are queued for
// deletion. Otherwise, the current plan values are sent and per-field deletions are added
// for fields removed since the last apply.
//
// Note on password: because password is write-only it is always null in state, so state
// comparison cannot detect whether a password was previously set. The password is sent
// whenever the plan includes a value; removal from the plan does not automatically delete
// it from Proxmox (to clear the password, replace the VM or set it explicitly).
func FillUpdateBody(
	ctx context.Context,
	planValue, stateValue Value,
	updateBody *vms.UpdateRequestBody,
	diags *diag.Diagnostics,
) {
	if planValue.IsUnknown() || planValue.Equal(stateValue) {
		return
	}

	// Whole block removed — queue deletion of all cloud-init API keys.
	if planValue.IsNull() {
		deleteAllCloudInit(ctx, stateValue, updateBody, diags)
		return
	}

	plan := unpackOrEmpty(ctx, planValue, diags)
	state := unpackOrEmpty(ctx, stateValue, diags)

	if diags.HasError() {
		return
	}

	// Apply plan values to the request body.
	plan.toAPI(ctx, updateBody, diags)

	if diags.HasError() {
		return
	}

	// Queue deletions for individual fields that were removed from the plan.
	planDNS := unpackDNS(ctx, plan.DNS, diags)
	stateDNS := unpackDNS(ctx, state.DNS, diags)

	attribute.CheckDeleteBody(planDNS.Domain, stateDNS.Domain, updateBody, "searchdomain")

	if planDNS.Servers.IsNull() && !stateDNS.Servers.IsNull() {
		updateBody.AppendDelete("nameserver")
	}

	planUA := unpackUserAccount(ctx, plan.UserAccount, diags)
	stateUA := unpackUserAccount(ctx, state.UserAccount, diags)

	attribute.CheckDeleteBody(planUA.Username, stateUA.Username, updateBody, "ciuser")

	if planUA.Keys.IsNull() && !stateUA.Keys.IsNull() {
		updateBody.AppendDelete("sshkeys")
	}

	attribute.CheckDeleteBody(plan.Type, state.Type, updateBody, "citype")
	attribute.CheckDeleteBody(plan.Upgrade, state.Upgrade, updateBody, "ciupgrade")

	// File ID deletions (RequiresReplace, but handle defensively).
	attribute.CheckDeleteBody(plan.UserDataFileID, state.UserDataFileID, updateBody, "cicustom")
	attribute.CheckDeleteBody(plan.VendorDataFileID, state.VendorDataFileID, updateBody, "cicustom")
	attribute.CheckDeleteBody(plan.NetworkDataFileID, state.NetworkDataFileID, updateBody, "cicustom")
	attribute.CheckDeleteBody(plan.MetaDataFileID, state.MetaDataFileID, updateBody, "cicustom")

	// IP config: delete slots that exceeded the new plan length.
	planLen := ipConfigLen(ctx, plan.IPConfig, diags)
	stateLen := ipConfigLen(ctx, state.IPConfig, diags)

	for i := planLen; i < stateLen; i++ {
		updateBody.AppendDelete(fmt.Sprintf("ipconfig%d", i))
	}
}

// deleteAllCloudInit queues deletion for every cloud-init API key that was active in state.
func deleteAllCloudInit(ctx context.Context, stateValue Value, body *vms.UpdateRequestBody, diags *diag.Diagnostics) {
	state := unpackOrEmpty(ctx, stateValue, diags)

	stateDNS := unpackDNS(ctx, state.DNS, diags)
	stateUA := unpackUserAccount(ctx, state.UserAccount, diags)

	if !stateDNS.Domain.IsNull() {
		body.AppendDelete("searchdomain")
	}

	if !stateDNS.Servers.IsNull() {
		body.AppendDelete("nameserver")
	}

	if !stateUA.Username.IsNull() {
		body.AppendDelete("ciuser")
	}

	if !stateUA.Keys.IsNull() {
		body.AppendDelete("sshkeys")
	}

	if !state.Type.IsNull() {
		body.AppendDelete("citype")
	}

	if !state.Upgrade.IsNull() {
		body.AppendDelete("ciupgrade")
	}

	if !state.UserDataFileID.IsNull() || !state.VendorDataFileID.IsNull() ||
		!state.NetworkDataFileID.IsNull() || !state.MetaDataFileID.IsNull() {
		body.AppendDelete("cicustom")
	}

	n := ipConfigLen(ctx, state.IPConfig, diags)

	for i := range n {
		body.AppendDelete(fmt.Sprintf("ipconfig%d", i))
	}
}

// unpackOrEmpty returns a Model decoded from the Value, or a zero Model when the Value is
// null or unknown — mirrors the pattern in cpu.unpackOrEmpty.
func unpackOrEmpty(ctx context.Context, value Value, diags *diag.Diagnostics) Model {
	if value.IsNull() || value.IsUnknown() {
		return Model{
			DNS:               types.ObjectNull(dnsAttributeTypes()),
			IPConfig:          types.ListNull(types.ObjectType{AttrTypes: ipConfigAttributeTypes()}),
			MetaDataFileID:    types.StringNull(),
			NetworkDataFileID: types.StringNull(),
			Type:              types.StringNull(),
			Upgrade:           types.BoolNull(),
			UserAccount:       types.ObjectNull(userAccountAttributeTypes()),
			UserDataFileID:    types.StringNull(),
			VendorDataFileID:  types.StringNull(),
		}
	}

	var m Model

	diags.Append(value.As(ctx, &m, basetypes.ObjectAsOptions{})...)

	return m
}

func unpackDNS(ctx context.Context, obj types.Object, diags *diag.Diagnostics) dnsModel {
	if obj.IsNull() || obj.IsUnknown() {
		return dnsModel{Domain: types.StringNull(), Servers: types.ListNull(types.StringType)}
	}

	var dns dnsModel

	diags.Append(obj.As(ctx, &dns, basetypes.ObjectAsOptions{})...)

	return dns
}

func unpackUserAccount(ctx context.Context, obj types.Object, diags *diag.Diagnostics) userAccountModel {
	if obj.IsNull() || obj.IsUnknown() {
		return userAccountModel{
			Keys:     types.ListNull(types.StringType),
			Password: types.StringNull(),
			Username: types.StringNull(),
		}
	}

	var ua userAccountModel

	diags.Append(obj.As(ctx, &ua, basetypes.ObjectAsOptions{})...)

	return ua
}

func ipConfigLen(ctx context.Context, list types.List, diags *diag.Diagnostics) int {
	if list.IsNull() || list.IsUnknown() {
		return 0
	}

	var items []ipConfigModel

	diags.Append(list.ElementsAs(ctx, &items, false)...)

	return len(items)
}
