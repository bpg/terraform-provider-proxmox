/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cpu

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Model represents the CPU model. The sub-block mirrors the PVE web UI's "Processors" dialog:
// every knob exposed in that dialog — cores, sockets, vcpus, type+flags, limit, units,
// affinity, arch, numa — lives here.
type Model struct {
	Affinity     types.String  `tfsdk:"affinity"`
	Architecture types.String  `tfsdk:"architecture"`
	Cores        types.Int64   `tfsdk:"cores"`
	Flags        types.Set     `tfsdk:"flags"`
	Limit        types.Float64 `tfsdk:"limit"`
	Numa         types.Bool    `tfsdk:"numa"`
	Sockets      types.Int64   `tfsdk:"sockets"`
	Type         types.String  `tfsdk:"type"`
	Units        types.Int64   `tfsdk:"units"`
	Vcpus        types.Int64   `tfsdk:"vcpus"`
}

func attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"affinity":     types.StringType,
		"architecture": types.StringType,
		"cores":        types.Int64Type,
		"flags":        types.SetType{ElemType: types.StringType},
		"limit":        types.Float64Type,
		"numa":         types.BoolType,
		"sockets":      types.Int64Type,
		"type":         types.StringType,
		"units":        types.Int64Type,
		"vcpus":        types.Int64Type,
	}
}

// NullValue returns a properly typed null Value.
func NullValue() Value {
	return types.ObjectNull(attributeTypes())
}

// toAPI writes the CPU-related fields onto the shared create/update body. Like memory, cpu has
// no single nested API struct — `affinity`, `arch`, `cores`, `cpulimit`, `sockets`, `cpuunits`
// live directly on the top-level request body, while `cpu` (compound) maps to CPUEmulation.
// Signature therefore follows the ADR-004 §Model-API Conversion write-through variant.
//
// Null/unknown fields produce nil pointers so the request omits them entirely. The compound
// CPUEmulation is only populated when `Type` is defined in the plan — PVE rejects `cpu=...`
// without an explicit cputype.
func (m *Model) toAPI(ctx context.Context, body *vms.CreateRequestBody, diags *diag.Diagnostics) {
	body.CPUAffinity = attribute.StringPtrFromValue(m.Affinity)
	body.CPUArchitecture = attribute.StringPtrFromValue(m.Architecture)
	body.CPUCores = attribute.Int64PtrFromValue(m.Cores)
	body.CPULimit = attribute.Float64PtrFromValue(m.Limit)
	body.CPUSockets = attribute.Int64PtrFromValue(m.Sockets)
	body.CPUUnits = attribute.Int64PtrFromValue(m.Units)
	body.NUMAEnabled = attribute.CustomBoolPtrFromValue(m.Numa)
	body.VirtualCPUCount = attribute.Int64PtrFromValue(m.Vcpus)

	if !attribute.IsDefined(m.Type) {
		return
	}

	emulation := &vms.CustomCPUEmulation{Type: m.Type.ValueString()}

	if attribute.IsDefined(m.Flags) {
		var flags []string

		diags.Append(m.Flags.ElementsAs(ctx, &flags, false)...)
		emulation.Flags = &flags
	}

	body.CPUEmulation = emulation
}

// fromAPI populates the Model from a PVE API response. Caller guarantees `config` belongs to a
// "CPU block present" branch (i.e. at least one cpu-related key is set). PVE returns only the
// keys the user explicitly wrote, so fields absent from the response map to null in the Model —
// there is no implicit default to fill in.
func (m *Model) fromAPI(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) {
	m.Affinity = types.StringPointerValue(config.CPUAffinity)
	m.Architecture = types.StringPointerValue(config.CPUArchitecture)
	m.Cores = types.Int64PointerValue(config.CPUCores)
	m.Limit = types.Float64PointerValue(config.CPULimit.PointerFloat64())
	m.Numa = types.BoolPointerValue(config.NUMAEnabled.PointerBool())
	m.Sockets = types.Int64PointerValue(config.CPUSockets)
	m.Units = types.Int64PointerValue(config.CPUUnits)
	m.Vcpus = types.Int64PointerValue(config.VirtualCPUCount)

	if config.CPUEmulation == nil {
		m.Type = types.StringNull()
		m.Flags = types.SetNull(basetypes.StringType{})

		return
	}

	if config.CPUEmulation.Type == "" {
		m.Type = types.StringNull()
	} else {
		m.Type = types.StringValue(config.CPUEmulation.Type)
	}

	if config.CPUEmulation.Flags == nil || len(*config.CPUEmulation.Flags) == 0 {
		m.Flags = types.SetNull(basetypes.StringType{})

		return
	}

	flags, d := types.SetValueFrom(ctx, basetypes.StringType{}, *config.CPUEmulation.Flags)
	diags.Append(d...)

	m.Flags = flags
}
