/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
)

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
