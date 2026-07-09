package disk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// Value represents the type for Disk settings.
type Value = types.Map

// NewValue returns a new Value with the given Disk settings from the PVE API.
//
// Returns NullValue() when the VM has no disks attached - PVE does not auto-attach.
func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value {
	// find storage devices with media=disk
	disks := config.StorageDevices.Filter(func(device *vms.CustomStorageDevice) bool {
		return device.Media != nil && *device.Media == "disk"
	})

	if len(disks) == 0 {
		return NullValue()
	}

	elements := make(map[string]Model, len(disks))

	for iface, disk := range disks {
		m := Model{}
		m.fromAPI(*disk)
		elements[iface] = m
	}

	obj, d := types.MapValueFrom(ctx, types.ObjectType{}.WithAttributeTypes(attributeTypes()), elements)
	diags.Append(d...)

	return obj
}

// FillCreateBody fills the CreateRequestBody with the Disk settings from the Value
//
// In the 'create' context, planValue is the plan.
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

	for iface, disk := range plan {
		body.AddCustomStorageDevice(iface, disk.toAPI())
	}
}

// FillUpdateBody fills the UpdateRequestBody with the Disk settings from the Value.
//
// In the 'update' context, planValue is the plan and stateValue is the current state. Either
// side may be null (e.g. plan is null when the user removes the whole `disk` block; state is
// null on a refresh of a VM that never had a disk). Null is treated as an empty map so
// MapDiff can still produce slot-level creates/deletes.
//
// Note: disk size changes are NOT handled here — they require a separate ResizeVMDisk API call.
// Use ResizeDisks for that.
func FillUpdateBody(
	ctx context.Context,
	planValue, stateValue Value,
	updateBody *vms.UpdateRequestBody,
	diags *diag.Diagnostics,
) {
	if planValue.IsUnknown() || planValue.Equal(stateValue) {
		return
	}

	var plan, state map[string]Model

	if !planValue.IsNull() {
		d := planValue.ElementsAs(ctx, &plan, false)
		diags.Append(d...)
	}

	if !stateValue.IsNull() {
		d := stateValue.ElementsAs(ctx, &state, false)
		diags.Append(d...)
	}

	if diags.HasError() {
		return
	}

	toCreate, toUpdate, toDelete := utils.MapDiff(plan, state)

	for iface, dev := range toCreate {
		updateBody.AddCustomStorageDevice(iface, dev.toAPI())
	}

	for iface, dev := range toUpdate {
		updateBody.AddCustomStorageDevice(iface, dev.toAPI())
	}

	for iface := range toDelete {
		updateBody.Delete = append(updateBody.Delete, iface)
	}
}

// ResizeDisks issues ResizeVMDisk API calls for any disks whose size increased between state
// and plan. PVE does not allow shrinking disks, so only growth is applied. Size changes are
// not part of the regular UpdateVM call — they require a dedicated resize endpoint.
func ResizeDisks(
	ctx context.Context,
	planValue, stateValue Value,
	vmAPI *vms.Client,
	diags *diag.Diagnostics,
) {
	if planValue.IsNull() || planValue.IsUnknown() {
		return
	}

	var plan, state map[string]Model

	d := planValue.ElementsAs(ctx, &plan, false)
	diags.Append(d...)

	if !stateValue.IsNull() {
		d = stateValue.ElementsAs(ctx, &state, false)
		diags.Append(d...)
	}

	if diags.HasError() {
		return
	}

	for iface, planDisk := range plan {
		stateDisk, exists := state[iface]
		if !exists {
			continue // new disk, size is set via create/update body
		}

		planSize := planDisk.Size.ValueInt64()
		stateSize := stateDisk.Size.ValueInt64()

		if planSize <= stateSize {
			if planSize < stateSize {
				diags.AddError(
					fmt.Sprintf("Unable to Shrink Disk %s", iface),
					fmt.Sprintf(
						"Disk size can only be increased. Current size: %d GiB, requested: %d GiB",
						stateSize, planSize,
					),
				)
			}

			continue
		}

		tflog.Info(ctx, fmt.Sprintf("Resizing disk %s from %d GiB to %d GiB", iface, stateSize, planSize))

		result := vmAPI.ResizeVMDisk(ctx, &vms.ResizeDiskRequestBody{
			Disk: iface,
			Size: *proxmoxtypes.DiskSizeFromGigabytes(planSize),
		})
		if result.Err() != nil {
			diags.AddError(
				fmt.Sprintf("Unable to Resize Disk %s", iface),
				result.Err().Error(),
			)

			return
		}

		for _, w := range result.Warnings() {
			diags.AddWarning(fmt.Sprintf("Disk %s Resize", iface), w)
		}
	}
}
