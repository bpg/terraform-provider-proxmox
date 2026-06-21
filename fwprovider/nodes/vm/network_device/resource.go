/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network_device

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// NewValue converts a PVE CustomNetworkDeviceMap to a typed list Value.
func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value {
	if len(config.NetworkDevices) == 0 {
		return NullValue()
	}

	keys := make([]string, 0, len(config.NetworkDevices))

	for k := range config.NetworkDevices {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		ni, err := strconv.Atoi(keys[i][3:])
		if err != nil {
			ni = 0
		}

		nj, err := strconv.Atoi(keys[j][3:])
		if err != nil {
			nj = 0
		}

		return ni < nj
	})

	elements := make([]attr.Value, 0, len(keys))

	for _, k := range keys {
		dev := config.NetworkDevices[k]
		if dev == nil {
			continue
		}

		m := fromAPI(dev)

		obj, d := types.ObjectValueFrom(ctx, attributeTypes(), m)
		diags.Append(d...)

		elements = append(elements, obj)
	}

	list, d := types.ListValue(elementType(), elements)
	diags.Append(d...)

	return list
}

// FillCreateBody populates the create request body with the network devices from the plan.
func FillCreateBody(ctx context.Context, planValue Value, body *vms.CreateRequestBody, diags *diag.Diagnostics) {
	if planValue.IsNull() || planValue.IsUnknown() {
		return
	}

	models := extractModels(ctx, planValue, diags)
	if diags.HasError() {
		return
	}

	devices := make(vms.CustomNetworkDevices, len(models))
	for i, m := range models {
		devices[i] = m.toAPI()
	}

	body.NetworkDevices = devices
}

// FillUpdateBody computes the network device changes between plan and state and applies them to the update body.
func FillUpdateBody(ctx context.Context, planValue Value, stateValue Value, body *vms.CreateRequestBody, diags *diag.Diagnostics) {
	planModels := extractModels(ctx, planValue, diags)
	stateModels := extractModels(ctx, stateValue, diags)

	if diags.HasError() {
		return
	}

	// Build update slice from plan devices
	devices := make(vms.CustomNetworkDevices, len(planModels))
	for i, m := range planModels {
		devices[i] = m.toAPI()
	}

	body.NetworkDevices = devices

	// Delete slots that exist in state but are not in plan
	for i := len(planModels); i < len(stateModels); i++ {
		body.Delete = append(body.Delete, fmt.Sprintf("net%d", i))
	}
}

// extractModels decodes the list Value into individual Model instances.
func extractModels(ctx context.Context, listValue Value, diags *diag.Diagnostics) []Model {
	if listValue.IsNull() || listValue.IsUnknown() {
		return nil
	}

	elements := listValue.Elements()
	models := make([]Model, 0, len(elements))

	for _, elem := range elements {
		obj, ok := elem.(basetypes.ObjectValue)
		if !ok {
			diags.AddError("network_device: unexpected element type", "expected ObjectValue")
			return nil
		}

		var m Model

		diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)

		if diags.HasError() {
			return nil
		}

		models = append(models, m)
	}

	return models
}
