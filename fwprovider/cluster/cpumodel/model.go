/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cpumodel

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/cpumodels"
)

// customPrefix is the prefix PVE adds to custom CPU model names when they are referenced as a CPU type.
const customPrefix = "custom-"

type customCPUModel struct {
	ID            types.String `tfsdk:"id"`
	CPUType       types.String `tfsdk:"cputype"`
	ReportedModel types.String `tfsdk:"reported_model"`
	Flags         types.List   `tfsdk:"flags"`
	PhysBits      types.String `tfsdk:"phys_bits"`
	GuestPhysBits types.Int64  `tfsdk:"guest_phys_bits"`
	Hidden        types.Bool   `tfsdk:"hidden"`
	HVVendorID    types.String `tfsdk:"hv_vendor_id"`
	Level         types.Int64  `tfsdk:"level"`
}

func (m *customCPUModel) toAPICommon(ctx context.Context, diags *diag.Diagnostics) cpumodels.RequestBodyCommon {
	body := cpumodels.RequestBodyCommon{
		GuestPhysBits: attribute.Int64PtrFromValue(m.GuestPhysBits),
		Hidden:        attribute.CustomBoolPtrFromValue(m.Hidden),
		HVVendorID:    attribute.StringPtrFromValue(m.HVVendorID),
		Level:         attribute.Int64PtrFromValue(m.Level),
		PhysBits:      attribute.StringPtrFromValue(m.PhysBits),
		ReportedModel: attribute.StringPtrFromValue(m.ReportedModel),
	}

	if attribute.IsDefined(m.Flags) {
		var flags []string

		diags.Append(m.Flags.ElementsAs(ctx, &flags, false)...)

		if len(flags) > 0 {
			joined := strings.Join(flags, ";")
			body.Flags = &joined
		}
	}

	return body
}

func (m *customCPUModel) toAPICreate(ctx context.Context, diags *diag.Diagnostics) *cpumodels.CreateRequestBody {
	return &cpumodels.CreateRequestBody{
		RequestBodyCommon: m.toAPICommon(ctx, diags),
		CPUType:           m.ID.ValueString(),
	}
}

func (m *customCPUModel) toAPIUpdate(
	ctx context.Context,
	state *customCPUModel,
	diags *diag.Diagnostics,
) *cpumodels.UpdateRequestBody {
	body := &cpumodels.UpdateRequestBody{
		RequestBodyCommon: m.toAPICommon(ctx, diags),
	}

	var toDelete []string

	attribute.CheckDelete(m.Flags, state.Flags, &toDelete, "flags")
	attribute.CheckDelete(m.PhysBits, state.PhysBits, &toDelete, "phys-bits")
	attribute.CheckDelete(m.GuestPhysBits, state.GuestPhysBits, &toDelete, "guest-phys-bits")
	attribute.CheckDelete(m.HVVendorID, state.HVVendorID, &toDelete, "hv-vendor-id")
	attribute.CheckDelete(m.Level, state.Level, &toDelete, "level")

	if len(toDelete) > 0 {
		body.Delete = toDelete
	}

	return body
}

func (m *customCPUModel) fromAPI(data *cpumodels.GetResponseData) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(strings.TrimPrefix(data.CPUType, customPrefix))
	m.CPUType = types.StringValue(data.CPUType)
	m.ReportedModel = types.StringPointerValue(data.ReportedModel)
	m.PhysBits = types.StringPointerValue(data.PhysBits)
	m.HVVendorID = types.StringPointerValue(data.HVVendorID)

	// PVE omits `hidden` unless it is set; omission means false.
	m.Hidden = types.BoolValue(data.Hidden != nil && bool(*data.Hidden))

	if data.GuestPhysBits != nil {
		m.GuestPhysBits = types.Int64Value(int64(*data.GuestPhysBits))
	} else {
		m.GuestPhysBits = types.Int64Null()
	}

	if data.Level != nil {
		m.Level = types.Int64Value(int64(*data.Level))
	} else {
		m.Level = types.Int64Null()
	}

	if data.Flags != nil && *data.Flags != "" {
		parts := strings.Split(*data.Flags, ";")
		values := make([]attr.Value, len(parts))

		for i, flag := range parts {
			values[i] = types.StringValue(strings.TrimSpace(flag))
		}

		list, d := types.ListValue(types.StringType, values)
		diags.Append(d...)

		m.Flags = list
	} else {
		m.Flags = types.ListNull(types.StringType)
	}

	return diags
}
