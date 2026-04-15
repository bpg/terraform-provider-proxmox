/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package hardware

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// ListPCIDevicesRequestBody contains the request body for listing PCI devices.
type ListPCIDevicesRequestBody struct {
	ClassBlacklist *string `url:"pci-class-blacklist,omitempty"`
}

// ListPCIDevicesResponseBody contains the response body from listing PCI devices.
type ListPCIDevicesResponseBody struct {
	Data []*PCIDeviceData `json:"data,omitempty"`
}

// PCIDeviceData contains data for a single PCI device.
type PCIDeviceData struct {
	ID                  string            `json:"id"`
	Class               string            `json:"class"`
	Device              string            `json:"device"`
	DeviceName          *string           `json:"device_name,omitempty"`
	IOMMUGroup          int64             `json:"iommugroup"`
	MediatedDevices     *types.CustomBool `json:"mdev,omitempty"`
	SubsystemDevice     *string           `json:"subsystem_device,omitempty"`
	SubsystemDeviceName *string           `json:"subsystem_device_name,omitempty"`
	SubsystemVendor     *string           `json:"subsystem_vendor,omitempty"`
	SubsystemVendorName *string           `json:"subsystem_vendor_name,omitempty"`
	Vendor              string            `json:"vendor"`
	VendorName          *string           `json:"vendor_name,omitempty"`
}
