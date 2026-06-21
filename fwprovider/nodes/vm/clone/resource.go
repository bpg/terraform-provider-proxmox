/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package clone

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// BuildCloneBody extracts the clone model and returns a CloneRequestBody ready to send to PVE.
// newVMID must be the destination VM ID. Returns nil if the clone block is null.
func BuildCloneBody(ctx context.Context, cloneValue Value, newVMID int, diags *diag.Diagnostics) *vms.CloneRequestBody {
	if cloneValue.IsNull() || cloneValue.IsUnknown() {
		return nil
	}

	var m Model

	diags.Append(cloneValue.As(ctx, &m, basetypes.ObjectAsOptions{})...)

	if diags.HasError() {
		return nil
	}

	body := &vms.CloneRequestBody{
		VMIDNew: newVMID,
	}

	if !m.DatastoreID.IsNull() && !m.DatastoreID.IsUnknown() {
		body.TargetStorage = m.DatastoreID.ValueStringPointer()
	}

	if !m.Full.IsNull() && !m.Full.IsUnknown() {
		v := proxmoxtypes.CustomBool(m.Full.ValueBool())
		body.FullCopy = &v
	}

	return body
}

// Retries extracts the retries count from the clone block (defaults to 3 if not set).
func Retries(ctx context.Context, cloneValue Value, diags *diag.Diagnostics) int {
	if cloneValue.IsNull() || cloneValue.IsUnknown() {
		return 3
	}

	var m Model

	diags.Append(cloneValue.As(ctx, &m, basetypes.ObjectAsOptions{})...)

	if diags.HasError() || m.Retries.IsNull() || m.Retries.IsUnknown() {
		return 3
	}

	return int(m.Retries.ValueInt64())
}

// SourceNodeName returns the source node for the clone, falling back to the provided default.
func SourceNodeName(ctx context.Context, cloneValue Value, targetNodeName string, diags *diag.Diagnostics) string {
	if cloneValue.IsNull() || cloneValue.IsUnknown() {
		return targetNodeName
	}

	var m Model

	diags.Append(cloneValue.As(ctx, &m, basetypes.ObjectAsOptions{})...)

	if diags.HasError() || m.NodeName.IsNull() || m.NodeName.IsUnknown() {
		return targetNodeName
	}

	return m.NodeName.ValueString()
}

// SourceVMID returns the source VM ID from the clone block.
func SourceVMID(ctx context.Context, cloneValue Value, diags *diag.Diagnostics) int {
	if cloneValue.IsNull() || cloneValue.IsUnknown() {
		return 0
	}

	var m Model

	diags.Append(cloneValue.As(ctx, &m, basetypes.ObjectAsOptions{})...)

	if diags.HasError() {
		return 0
	}

	return int(m.VMID.ValueInt64())
}
