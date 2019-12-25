/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"fmt"
	"net/url"
)

// CreateVM creates an virtual machine.
func (c *VirtualEnvironmentClient) CreateVM(nodeName string, d *VirtualEnvironmentVMCreateRequestBody) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu", url.PathEscape(nodeName)), d, nil)
}

// DeleteVM deletes an virtual machine.
func (c *VirtualEnvironmentClient) DeleteVM(nodeName string, vmID int) error {
	return c.DoRequest(hmDELETE, fmt.Sprintf("nodes/%s/qemu/%d", url.PathEscape(nodeName), vmID), nil, nil)
}

// GetVM retrieves an virtual machine.
func (c *VirtualEnvironmentClient) GetVM(nodeName string, vmID int) (*VirtualEnvironmentVMGetResponseData, error) {
	resBody := &VirtualEnvironmentVMGetResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("nodes/%s/qemu/%d/config", url.PathEscape(nodeName), vmID), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListVMs retrieves a list of virtual machines.
func (c *VirtualEnvironmentClient) ListVMs() ([]*VirtualEnvironmentVMListResponseData, error) {
	return nil, errors.New("Not implemented")
}

// UpdateVM updates a virtual machine.
func (c *VirtualEnvironmentClient) UpdateVM(nodeName string, vmID int, d *VirtualEnvironmentVMUpdateRequestBody) error {
	return c.DoRequest(hmPUT, fmt.Sprintf("nodes/%s/qemu/%d/config", url.PathEscape(nodeName), vmID), d, nil)
}

// UpdateVMAsync updates a virtual machine asynchronously.
func (c *VirtualEnvironmentClient) UpdateVMAsync(nodeName string, vmID int, d *VirtualEnvironmentVMUpdateRequestBody) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu/%d/config", url.PathEscape(nodeName), vmID), d, nil)
}
