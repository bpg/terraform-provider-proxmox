/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cdrom

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// Value represents the type for CD-ROM settings.
type Value = types.Map

// NewValue returns a new Value with the given CD-ROM settings from the PVE API.
func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value {
	// find storage devices with media=cdrom
	cdroms := vms.MapCustomStorageDevices(*config).Filter(func(device *vms.CustomStorageDevice) bool {
		return device.Media != nil && *device.Media == "cdrom"
	})

	elements := make(map[string]Model, len(cdroms))

	for iface, cdrom := range cdroms {
		m := Model{}
		m.importFromCustomStorageDevice(*cdrom)
		elements[iface] = m
	}

	obj, d := types.MapValueFrom(ctx, types.ObjectType{}.WithAttributeTypes(attributeTypes()), elements)
	diags.Append(d...)

	return obj
}

// FillCreateBody fills the CreateRequestBody with the CD-ROM settings from the Value.
//
// In the 'create' context, v is the plan.
func FillCreateBody(ctx context.Context, planValue Value, body *vms.CreateRequestBody, diags *diag.Diagnostics) {
	if planValue.IsNull() || planValue.IsUnknown() {
		return
	}

	var plan map[string]Model
	d := planValue.ElementsAs(ctx, &plan, false)
	diags.Append(d...)

	if d.HasError() {
		return
	}

	for iface, cdrom := range plan {
		err := body.AddCustomStorageDevice(cdrom.exportToCustomStorageDevice(iface))
		if err != nil {
			diags.AddError(err.Error(), "")
		}
	}
}

// FillUpdateBody fills the UpdateRequestBody with the CD-ROM settings from the Value.
//
// In the 'update' context, v is the plan and stateValue is the current state.
func FillUpdateBody(
	ctx context.Context,
	planValue, stateValue Value,
	updateBody *vms.UpdateRequestBody,
	_ bool,
	diags *diag.Diagnostics,
) {
	if planValue.IsNull() || planValue.IsUnknown() || planValue.Equal(stateValue) {
		return
	}

	var plan, state map[string]Model
	d := planValue.ElementsAs(ctx, &plan, false)
	diags.Append(d...)
	d = stateValue.ElementsAs(ctx, &state, false)
	diags.Append(d...)

	if diags.HasError() {
		return
	}

	toCreate, toUpdate, toDelete := utils.MapDiff(plan, state)

	for iface, dev := range toCreate {
		err := updateBody.AddCustomStorageDevice(dev.exportToCustomStorageDevice(iface))
		if err != nil {
			diags.AddError(err.Error(), "")
		}
	}

	for iface, dev := range toUpdate {
		// for CD-ROMs, the update fully override the existing device, we don't do per-attribute check
		err := updateBody.AddCustomStorageDevice(dev.exportToCustomStorageDevice(iface))
		if err != nil {
			diags.AddError(err.Error(), "")
		}
	}

	for iface := range toDelete {
		updateBody.Delete = append(updateBody.Delete, iface)
	}
}
