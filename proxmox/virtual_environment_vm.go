/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// CreateVM creates a virtual machine.
func (c *VirtualEnvironmentClient) CreateVM(nodeName string, d *VirtualEnvironmentVMCreateRequestBody) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu", url.PathEscape(nodeName)), d, nil)
}

// DeleteVM deletes a virtual machine.
func (c *VirtualEnvironmentClient) DeleteVM(nodeName string, vmID int) error {
	return c.DoRequest(hmDELETE, fmt.Sprintf("nodes/%s/qemu/%d", url.PathEscape(nodeName), vmID), nil, nil)
}

// GetVM retrieves a virtual machine.
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

// GetVMID retrieves the next available VM identifier.
func (c *VirtualEnvironmentClient) GetVMID() (*int, error) {
	nodes, err := c.ListNodes()

	if err != nil {
		return nil, err
	}

	vmID := 100

VMID:
	for vmID <= 2147483637 {
		for _, n := range nodes {
			err := c.DoRequest(hmGET, fmt.Sprintf("nodes/%s/qemu/%d/status/current", url.PathEscape(n.Name), vmID), nil, nil)

			if err == nil {
				vmID += 5

				continue VMID
			}
		}

		return &vmID, nil
	}

	return nil, errors.New("Unable to retrieve the next available VM identifier")
}

// GetVMNetworkInterfacesFromAgent retrieves the network interfaces reported by the QEMU agent.
func (c *VirtualEnvironmentClient) GetVMNetworkInterfacesFromAgent(nodeName string, vmID int) (*VirtualEnvironmentVMGetQEMUNetworkInterfacesResponseData, error) {
	resBody := &VirtualEnvironmentVMGetQEMUNetworkInterfacesResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("nodes/%s/qemu/%d/agent/network-get-interfaces", url.PathEscape(nodeName), vmID), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// GetVMStatus retrieves the status for a virtual machine.
func (c *VirtualEnvironmentClient) GetVMStatus(nodeName string, vmID int) (*VirtualEnvironmentVMGetStatusResponseData, error) {
	resBody := &VirtualEnvironmentVMGetStatusResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("nodes/%s/qemu/%d/status/current", url.PathEscape(nodeName), vmID), nil, resBody)

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

// RebootVM reboots a virtual machine.
func (c *VirtualEnvironmentClient) RebootVM(nodeName string, vmID int, d *VirtualEnvironmentVMRebootRequestBody) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu/%d/status/reboot", url.PathEscape(nodeName), vmID), d, nil)
}

// ShutdownVM shuts down a virtual machine.
func (c *VirtualEnvironmentClient) ShutdownVM(nodeName string, vmID int, d *VirtualEnvironmentVMShutdownRequestBody) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu/%d/status/shutdown", url.PathEscape(nodeName), vmID), d, nil)
}

// StartVM starts a virtual machine.
func (c *VirtualEnvironmentClient) StartVM(nodeName string, vmID int) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu/%d/status/start", url.PathEscape(nodeName), vmID), nil, nil)
}

// StopVM stops a virtual machine immediately.
func (c *VirtualEnvironmentClient) StopVM(nodeName string, vmID int) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu/%d/status/stop", url.PathEscape(nodeName), vmID), nil, nil)
}

// UpdateVM updates a virtual machine.
func (c *VirtualEnvironmentClient) UpdateVM(nodeName string, vmID int, d *VirtualEnvironmentVMUpdateRequestBody) error {
	return c.DoRequest(hmPUT, fmt.Sprintf("nodes/%s/qemu/%d/config", url.PathEscape(nodeName), vmID), d, nil)
}

// UpdateVMAsync updates a virtual machine asynchronously.
func (c *VirtualEnvironmentClient) UpdateVMAsync(nodeName string, vmID int, d *VirtualEnvironmentVMUpdateRequestBody) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu/%d/config", url.PathEscape(nodeName), vmID), d, nil)
}

// WaitForNetworkInterfacesFromAgent waits for a virtual machine's QEMU agent to publish the network interfaces.
func (c *VirtualEnvironmentClient) WaitForNetworkInterfacesFromAgent(nodeName string, vmID int, timeout int, delay int) (*VirtualEnvironmentVMGetQEMUNetworkInterfacesResponseData, error) {
	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			data, err := c.GetVMNetworkInterfacesFromAgent(nodeName, vmID)

			if err == nil && data != nil {
				return data, err
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(200 * time.Millisecond)

		timeElapsed = time.Now().Sub(timeStart)
	}

	return nil, fmt.Errorf("Timeout while waiting for the QEMU agent on VM \"%d\" to publish the network interfaces", vmID)
}

// WaitForNoNetworkInterfacesFromAgent waits for a virtual machine's QEMU agent to unpublish the network interfaces.
func (c *VirtualEnvironmentClient) WaitForNoNetworkInterfacesFromAgent(nodeName string, vmID int, timeout int, delay int) error {
	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			_, err := c.GetVMNetworkInterfacesFromAgent(nodeName, vmID)

			if err != nil {
				return nil
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(200 * time.Millisecond)

		timeElapsed = time.Now().Sub(timeStart)
	}

	return fmt.Errorf("Timeout while waiting for the QEMU agent on VM \"%d\" to unpublish the network interfaces", vmID)
}

// WaitForState waits for a virtual machine to reach a specific state.
func (c *VirtualEnvironmentClient) WaitForState(nodeName string, vmID int, state string, timeout int, delay int) error {
	state = strings.ToLower(state)

	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			data, err := c.GetVMStatus(nodeName, vmID)

			if err != nil {
				return err
			}

			if data.Status == state {
				return nil
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(200 * time.Millisecond)

		timeElapsed = time.Now().Sub(timeStart)
	}

	return fmt.Errorf("Timeout while waiting for VM \"%d\" to enter the state \"%s\"", vmID, state)
}
