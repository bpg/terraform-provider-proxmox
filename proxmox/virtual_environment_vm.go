/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
)

// VirtualEnvironmentVMCreateRequestBody contains the data for an virtual machine create request.
type VirtualEnvironmentVMCreateRequestBody struct {
	ACPI            *CustomBool            `json:"acpi,omitempty" url:"acpi,omitempty,int"`
	Agent           *CustomAgent           `json:"agent,omitempty" url:"agent,omitempty"`
	AudioDevice     *CustomAudioDevice     `json:"audio0,omitempty" url:"audio0,omitempty"`
	Autostart       *CustomBool            `json:"autostart,omitempty" url:"autostart,omitempty,int"`
	BackupFile      *string                `json:"archive,omitempty" url:"archive,omitempty"`
	BandwidthLimit  *int                   `json:"bwlimit,omitempty" url:"bwlimit,omitempty"`
	BIOS            *string                `json:"bios,omitempty" url:"bios,omitempty"`
	BootDiskID      *string                `json:"bootdisk,omitempty" url:"bootdisk,omitempty"`
	BootOrder       *string                `json:"boot,omitempty" url:"boot,omitempty"`
	CDROM           *string                `json:"cdrom,omitempty" url:"cdrom,omitempty"`
	CloudInitConfig *CustomCloudInitConfig `json:"cloudinit,omitempty" url:"cloudinit,omitempty"`
	CPUArchitecture *string                `json:"arch,omitempty" url:"arch,omitempty"`
	CPUCores        *int                   `json:"cores,omitempty" url:"cores,omitempty"`
	CPULimit        *int                   `json:"cpulimit,omitempty" url:"cpulimit,omitempty"`
	CPUUnits        *int                   `json:"cpuunits,omitempty" url:"cpuunits,omitempty"`
	DedicatedMemory *int                   `json:"memory,omitempty" url:"memory,omitempty"`
	Description     *string                `json:"description,omitempty" url:"description,omitempty"`
	EFIDisk         *CustomEFIDisk         `json:"efidisk0,omitempty" url:"efidisk0,omitempty"`
	FloatingMemory  *int                   `json:"balloon,omitempty" url:"balloon,omitempty"`
	Freeze          *CustomBool            `json:"freeze,omitempty" url:"freeze,omitempty,int"`
	HookScript      *string                `json:"hookscript,omitempty" url:"hookscript,omitempty"`
	Hotplug         []string               `json:"hotplug,omitempty" url:"hotplug,omitempty,comma"`
	Hugepages       *string                `json:"hugepages,omitempty" url:"hugepages,omitempty"`
	IDEDevices      CustomIDEDevices       `json:"ide,omitempty" url:"ide,omitempty"`
	KeyboardLayout  *string                `json:"keyboard,omitempty" url:"keyboard,omitempty"`
	KVMArguments    []string               `json:"args,omitempty" url:"args,omitempty,space"`
	KVMEnabled      *CustomBool            `json:"kvm,omitempty" url:"kvm,omitempty,int"`
	LocalTime       *CustomBool            `json:"localtime,omitempty" url:"localtime,omitempty,int"`
	Lock            *string                `json:"lock,omitempty" url:"lock,omitempty"`
	MachineType     *string                `json:"machine,omitempty" url:"machine,omitempty"`
	MigrateDowntime *float64               `json:"migrate_downtime,omitempty" url:"migrate_downtime,omitempty"`
	MigrateSpeed    *int                   `json:"migrate_speed,omitempty" url:"migrate_speed,omitempty"`
	Name            *string                `json:"name,omitempty" url:"name,omitempty"`
	NetworkDevices  CustomNetworkDevices   `json:"net,omitempty" url:"net,omitempty"`
	NodeName        string                 `json:"node" url:"node"`
	NUMADevices     CustomNUMADevices      `json:"numa_devices,omitempty" url:"numa,omitempty"`
	NUMAEnabled     *CustomBool            `json:"numa,omitempty" url:"numa,omitempty,int"`
	OSType          *string                `json:"ostype,omitempty" url:"ostype,omitempty"`
	Overwrite       *CustomBool            `json:"force,omitempty" url:"force,omitempty,int"`
	PCIDevices      CustomPCIDevices       `json:"hostpci,omitempty" url:"hostpci,omitempty"`
	SharedMemory    *CustomSharedMemory    `json:"ivshmem,omitempty" url:"ivshmem,omitempty"`
	StartOnBoot     *CustomBool            `json:"onboot,omitempty" url:"onboot,omitempty,int"`
	VMID            int                    `json:"vmid" url:"vmid"`
}

// VirtualEnvironmentVMGetResponseBody contains the body from an virtual machine get response.
type VirtualEnvironmentVMGetResponseBody struct {
	Data *VirtualEnvironmentVMGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentVMGetResponseData contains the data from an virtual machine get response.
type VirtualEnvironmentVMGetResponseData struct {
	ACPI *CustomBool `json:"acpi,omitempty" url:"acpi,omitempty,int"`
}

// VirtualEnvironmentVMListResponseBody contains the body from an virtual machine list response.
type VirtualEnvironmentVMListResponseBody struct {
	Data []*VirtualEnvironmentVMListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentVMListResponseData contains the data from an virtual machine list response.
type VirtualEnvironmentVMListResponseData struct {
	ACPI *CustomBool `json:"acpi,omitempty" url:"acpi,omitempty,int"`
}

// VirtualEnvironmentVMUpdateRequestBody contains the data for an virtual machine update request.
type VirtualEnvironmentVMUpdateRequestBody struct {
	ACPI *CustomBool `json:"acpi,omitempty" url:"acpi,omitempty,int"`
}

// CreateVM creates an virtual machine.
func (c *VirtualEnvironmentClient) CreateVM(d *VirtualEnvironmentVMCreateRequestBody) error {
	return c.DoRequest(hmPOST, "nodes/%s/qemu", d, nil)
}

// DeleteVM deletes an virtual machine.
func (c *VirtualEnvironmentClient) DeleteVM(id string) error {
	return errors.New("Not implemented")
}

// GetVM retrieves an virtual machine.
func (c *VirtualEnvironmentClient) GetVM(nodeName string, vmID int) (*VirtualEnvironmentVMGetResponseData, error) {
	return nil, errors.New("Not implemented")
}

// ListVMs retrieves a list of virtual machines.
func (c *VirtualEnvironmentClient) ListVMs() ([]*VirtualEnvironmentVMListResponseData, error) {
	return nil, errors.New("Not implemented")
}

// UpdateVM updates an virtual machine.
func (c *VirtualEnvironmentClient) UpdateVM(id string, d *VirtualEnvironmentVMUpdateRequestBody) error {
	return errors.New("Not implemented")
}
