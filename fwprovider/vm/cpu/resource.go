/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cpu

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// Value represents the type for CPU settings.
type Value = types.Object

// NewValue returns a new Value with the given CPU settings from the PVE API.
func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value {
	cpu := Model{}

	cpu.Affinity = types.StringPointerValue(config.CPUAffinity)
	cpu.Architecture = types.StringPointerValue(config.CPUArchitecture)
	cpu.Hotplugged = types.Int64PointerValue(config.VirtualCPUCount)
	cpu.Limit = types.Int64PointerValue(config.CPULimit.PointerInt64())
	cpu.Numa = types.BoolPointerValue(config.NUMAEnabled.PointerBool())
	cpu.Units = types.Int64PointerValue(config.CPUUnits)

	// special cases: PVE does not return actual value for cores VM, etc is using default (i.e. a value is not specified)

	if config.CPUCores != nil {
		cpu.Cores = types.Int64PointerValue(config.CPUCores)
	} else {
		cpu.Cores = types.Int64Value(1)
	}

	if config.CPUSockets != nil {
		cpu.Sockets = types.Int64PointerValue(config.CPUSockets)
	} else {
		cpu.Sockets = types.Int64Value(1)
	}

	if config.CPUEmulation != nil {
		cpu.Type = types.StringValue(config.CPUEmulation.Type)

		flags, d := types.SetValueFrom(ctx, basetypes.StringType{}, config.CPUEmulation.Flags)
		diags.Append(d...)

		cpu.Flags = flags
	} else {
		cpu.Type = types.StringValue("kvm64")
		cpu.Flags = types.SetNull(basetypes.StringType{})
	}

	obj, d := types.ObjectValueFrom(ctx, attributeTypes(), cpu)
	diags.Append(d...)

	return obj
}

// FillCreateBody fills the CreateRequestBody with the CPU settings from the Value.
//
// In the 'create' context, v is the plan.
func FillCreateBody(ctx context.Context, planValue Value, body *vms.CreateRequestBody, diags *diag.Diagnostics) {
	var plan Model

	if planValue.IsNull() || planValue.IsUnknown() {
		return
	}

	d := planValue.As(ctx, &plan, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	if d.HasError() {
		return
	}

	// for computed fields, we need to check if they are unknown
	if !plan.Affinity.IsUnknown() {
		body.CPUAffinity = plan.Affinity.ValueStringPointer()
	}

	if !plan.Architecture.IsUnknown() {
		body.CPUArchitecture = plan.Architecture.ValueStringPointer()
	}

	if !plan.Cores.IsUnknown() {
		body.CPUCores = plan.Cores.ValueInt64Pointer()
	}

	if !plan.Limit.IsUnknown() {
		body.CPULimit = plan.Limit.ValueInt64Pointer()
	}

	if !plan.Sockets.IsUnknown() {
		body.CPUSockets = plan.Sockets.ValueInt64Pointer()
	}

	if !plan.Units.IsUnknown() {
		body.CPUUnits = plan.Units.ValueInt64Pointer()
	}

	if !plan.Numa.IsUnknown() {
		body.NUMAEnabled = proxmoxtypes.CustomBoolPtr(plan.Numa.ValueBoolPointer())
	}

	if !plan.Hotplugged.IsUnknown() {
		body.VirtualCPUCount = plan.Hotplugged.ValueInt64Pointer()
	}

	body.CPUEmulation = &vms.CustomCPUEmulation{}

	if !plan.Type.IsUnknown() {
		body.CPUEmulation.Type = plan.Type.ValueString()
	}

	if !plan.Flags.IsUnknown() {
		d = plan.Flags.ElementsAs(ctx, &body.CPUEmulation.Flags, false)
		diags.Append(d...)
	}
}

// FillUpdateBody fills the UpdateRequestBody with the CPU settings from the Value.
//
// In the 'update' context, v is the plan and stateValue is the current state.
func FillUpdateBody(
	ctx context.Context,
	planValue, stateValue Value,
	updateBody *vms.UpdateRequestBody,
	isClone bool,
	diags *diag.Diagnostics,
) {
	var plan, state Model

	if planValue.IsNull() || planValue.IsUnknown() || planValue.Equal(stateValue) {
		return
	}

	d := planValue.As(ctx, &plan, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	d = stateValue.As(ctx, &state, basetypes.ObjectAsOptions{})
	diags.Append(d...)

	if diags.HasError() {
		return
	}

	var errs []error

	del := func(field string) {
		errs = append(errs, updateBody.ToDelete(field))
	}

	if !plan.Affinity.Equal(state.Affinity) {
		if attribute.ShouldBeRemoved(plan.Affinity, state.Affinity, isClone) {
			del("CPUAffinity")
		} else if attribute.IsDefined(plan.Affinity) {
			updateBody.CPUAffinity = plan.Affinity.ValueStringPointer()
		}
	}

	if !plan.Architecture.Equal(state.Architecture) {
		if attribute.ShouldBeRemoved(plan.Architecture, state.Architecture, isClone) {
			del("CPUArchitecture")
		} else if attribute.IsDefined(plan.Architecture) {
			updateBody.CPUArchitecture = plan.Architecture.ValueStringPointer()
		}
	}

	if !plan.Cores.Equal(state.Cores) {
		if attribute.ShouldBeRemoved(plan.Cores, state.Cores, isClone) {
			del("CPUCores")
		} else if attribute.IsDefined(plan.Cores) {
			updateBody.CPUCores = plan.Cores.ValueInt64Pointer()
		}
	}

	if !plan.Limit.Equal(state.Limit) {
		if attribute.ShouldBeRemoved(plan.Limit, state.Limit, isClone) {
			del("CPULimit")
		} else if attribute.IsDefined(plan.Sockets) {
			updateBody.CPULimit = plan.Limit.ValueInt64Pointer()
		}
	}

	if !plan.Sockets.Equal(state.Sockets) {
		if attribute.ShouldBeRemoved(plan.Sockets, state.Sockets, isClone) {
			del("CPUSockets")
		} else if attribute.IsDefined(plan.Sockets) {
			updateBody.CPUSockets = plan.Sockets.ValueInt64Pointer()
		}
	}

	if !plan.Units.Equal(state.Units) {
		if attribute.ShouldBeRemoved(plan.Units, state.Units, isClone) {
			del("CPUUnits")
		} else if attribute.IsDefined(plan.Units) {
			updateBody.CPUUnits = plan.Units.ValueInt64Pointer()
		}
	}

	if !plan.Numa.Equal(state.Numa) {
		if attribute.ShouldBeRemoved(plan.Numa, state.Numa, isClone) {
			del("NUMAEnabled")
		} else if attribute.IsDefined(plan.Numa) {
			updateBody.NUMAEnabled = proxmoxtypes.CustomBoolPtr(plan.Numa.ValueBoolPointer())
		}
	}

	if !plan.Hotplugged.Equal(state.Hotplugged) {
		if attribute.ShouldBeRemoved(plan.Hotplugged, state.Hotplugged, isClone) {
			del("VirtualCPUCount")
		} else if attribute.IsDefined(plan.Hotplugged) {
			updateBody.VirtualCPUCount = plan.Hotplugged.ValueInt64Pointer()
		}
	}

	var delType, delFlags bool

	cpuEmulation := &vms.CustomCPUEmulation{}

	if !plan.Type.Equal(state.Type) {
		if attribute.ShouldBeRemoved(plan.Type, state.Type, isClone) {
			delType = true
		} else if attribute.IsDefined(plan.Type) {
			cpuEmulation.Type = plan.Type.ValueString()
		}
	}

	if !plan.Flags.Equal(state.Flags) {
		if attribute.ShouldBeRemoved(plan.Flags, state.Flags, isClone) {
			delFlags = true
		} else if attribute.IsDefined(plan.Flags) {
			d = plan.Flags.ElementsAs(ctx, &cpuEmulation.Flags, false)
			diags.Append(d...)
		}
	}

	switch {
	case delType && !delFlags:
		diags.AddError("Cannot have CPU flags without explicit definition of CPU type", "")
	case delType:
		del("CPUEmulation")
	case !reflect.DeepEqual(cpuEmulation, &vms.CustomCPUEmulation{}):
		updateBody.CPUEmulation = cpuEmulation
	}
}
