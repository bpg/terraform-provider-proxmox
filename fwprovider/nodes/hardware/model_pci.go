/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package hardware

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/hardware"
)

// pciDataSourceModel is the top-level model for the proxmox_hardware_pci data source.
type pciDataSourceModel struct {
	NodeName          types.String `tfsdk:"node_name"`
	PCIClassBlacklist types.List   `tfsdk:"pci_class_blacklist"`
	Filters           *pciFilters  `tfsdk:"filters"`
	Devices           []pciDevice  `tfsdk:"devices"`
}

// pciFilters holds the client-side filter parameters.
type pciFilters struct {
	ID       types.String `tfsdk:"id"`
	Class    types.String `tfsdk:"class"`
	VendorID types.String `tfsdk:"vendor_id"`
	DeviceID types.String `tfsdk:"device_id"`
}

// pciDevice is the model for a single PCI device in the output list.
type pciDevice struct {
	ID                  types.String `tfsdk:"id"`
	Class               types.String `tfsdk:"class"`
	Device              types.String `tfsdk:"device"`
	DeviceName          types.String `tfsdk:"device_name"`
	IOMMUGroup          types.Int64  `tfsdk:"iommu_group"`
	MediatedDevices     types.Bool   `tfsdk:"mdev"`
	SubsystemDevice     types.String `tfsdk:"subsystem_device"`
	SubsystemDeviceName types.String `tfsdk:"subsystem_device_name"`
	SubsystemVendor     types.String `tfsdk:"subsystem_vendor"`
	SubsystemVendorName types.String `tfsdk:"subsystem_vendor_name"`
	Vendor              types.String `tfsdk:"vendor"`
	VendorName          types.String `tfsdk:"vendor_name"`
}

// pciDeviceFromAPI converts an API PCI device to the Terraform model.
// All Computed fields get known values (never null).
func pciDeviceFromAPI(d *hardware.PCIDeviceData) pciDevice {
	m := pciDevice{
		ID:         types.StringValue(d.ID),
		Class:      types.StringValue(d.Class),
		Device:     types.StringValue(d.Device),
		IOMMUGroup: types.Int64Value(d.IOMMUGroup),
		Vendor:     types.StringValue(d.Vendor),
	}

	m.DeviceName = attribute.StringValueFromPtr(d.DeviceName)

	m.VendorName = attribute.StringValueFromPtr(d.VendorName)

	m.SubsystemDevice = attribute.StringValueFromPtr(d.SubsystemDevice)

	m.SubsystemDeviceName = attribute.StringValueFromPtr(d.SubsystemDeviceName)

	m.SubsystemVendor = attribute.StringValueFromPtr(d.SubsystemVendor)

	m.SubsystemVendorName = attribute.StringValueFromPtr(d.SubsystemVendorName)

	m.MediatedDevices = attribute.BoolValueFromCustomBoolPtr(d.MediatedDevices)

	return m
}
