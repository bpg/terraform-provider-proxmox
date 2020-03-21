/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	getVMIDStep = 1
)

var (
	getVMIDCounter      = -1
	getVMIDCounterMutex = &sync.Mutex{}
)

// CloneVM clones a virtual machine.
func (c *VirtualEnvironmentClient) CloneVM(nodeName string, vmID int, d *VirtualEnvironmentVMCloneRequestBody) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu/%d/clone", url.PathEscape(nodeName), vmID), d, nil)
}

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
	getVMIDCounterMutex.Lock()
	defer getVMIDCounterMutex.Unlock()

	if getVMIDCounter < 0 {
		nextVMID, err := c.GetClusterNextID(nil)

		if err != nil {
			return nil, err
		}

		if nextVMID == nil {
			return nil, errors.New("Unable to retrieve the next available VM identifier")
		}

		getVMIDCounter = *nextVMID + getVMIDStep

		log.Printf("[DEBUG] Determined next available VM identifier to be %d", *nextVMID)

		return nextVMID, nil
	}

	vmID := getVMIDCounter

	for vmID <= 2147483637 {
		_, err := c.GetClusterNextID(&vmID)

		if err != nil {
			vmID += getVMIDStep

			continue
		}

		getVMIDCounter = vmID + getVMIDStep

		log.Printf("[DEBUG] Determined next available VM identifier to be %d", vmID)

		return &vmID, nil
	}

	return nil, errors.New("Unable to determine the next available VM identifier")
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

// MoveVMDisk moves a virtual machine disk.
func (c *VirtualEnvironmentClient) MoveVMDisk(nodeName string, vmID int, d *VirtualEnvironmentVMMoveDiskRequestBody) error {
	taskID, err := c.MoveVMDiskAsync(nodeName, vmID, d)

	if err != nil {
		return err
	}

	err = c.WaitForNodeTask(nodeName, *taskID, 86400, 5)

	if err != nil {
		return err
	}

	return nil
}

// MoveVMDiskAsync moves a virtual machine disk asynchronously.
func (c *VirtualEnvironmentClient) MoveVMDiskAsync(nodeName string, vmID int, d *VirtualEnvironmentVMMoveDiskRequestBody) (*string, error) {
	resBody := &VirtualEnvironmentVMMoveDiskResponseBody{}
	err := c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu/%d/move_disk", url.PathEscape(nodeName), vmID), d, resBody)

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
	taskID, err := c.RebootVMAsync(nodeName, vmID, d)

	if err != nil {
		return err
	}

	err = c.WaitForNodeTask(nodeName, *taskID, 1800, 5)

	if err != nil {
		return err
	}

	return nil
}

// RebootVMAsync reboots a virtual machine asynchronously.
func (c *VirtualEnvironmentClient) RebootVMAsync(nodeName string, vmID int, d *VirtualEnvironmentVMRebootRequestBody) (*string, error) {
	resBody := &VirtualEnvironmentVMRebootResponseBody{}
	err := c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu/%d/status/reboot", url.PathEscape(nodeName), vmID), d, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ResizeVMDisk resizes a virtual machine disk.
func (c *VirtualEnvironmentClient) ResizeVMDisk(nodeName string, vmID int, d *VirtualEnvironmentVMResizeDiskRequestBody) error {
	return c.DoRequest(hmPUT, fmt.Sprintf("nodes/%s/qemu/%d/resize", url.PathEscape(nodeName), vmID), d, nil)
}

// ShutdownVM shuts down a virtual machine.
func (c *VirtualEnvironmentClient) ShutdownVM(nodeName string, vmID int, d *VirtualEnvironmentVMShutdownRequestBody) error {
	taskID, err := c.ShutdownVMAsync(nodeName, vmID, d)

	if err != nil {
		return err
	}

	err = c.WaitForNodeTask(nodeName, *taskID, 1800, 5)

	if err != nil {
		return err
	}

	return nil
}

// ShutdownVMAsync shuts down a virtual machine asynchronously.
func (c *VirtualEnvironmentClient) ShutdownVMAsync(nodeName string, vmID int, d *VirtualEnvironmentVMShutdownRequestBody) (*string, error) {
	resBody := &VirtualEnvironmentVMShutdownResponseBody{}
	err := c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu/%d/status/shutdown", url.PathEscape(nodeName), vmID), d, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// StartVM starts a virtual machine.
func (c *VirtualEnvironmentClient) StartVM(nodeName string, vmID int) error {
	taskID, err := c.StartVMAsync(nodeName, vmID)

	if err != nil {
		return err
	}

	err = c.WaitForNodeTask(nodeName, *taskID, 1800, 5)

	if err != nil {
		return err
	}

	return nil
}

// StartVMAsync starts a virtual machine asynchronously.
func (c *VirtualEnvironmentClient) StartVMAsync(nodeName string, vmID int) (*string, error) {
	resBody := &VirtualEnvironmentVMStartResponseBody{}
	err := c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu/%d/status/start", url.PathEscape(nodeName), vmID), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// StopVM stops a virtual machine.
func (c *VirtualEnvironmentClient) StopVM(nodeName string, vmID int) error {
	taskID, err := c.StopVMAsync(nodeName, vmID)

	if err != nil {
		return err
	}

	err = c.WaitForNodeTask(nodeName, *taskID, 300, 5)

	if err != nil {
		return err
	}

	return nil
}

// StopVMAsync stops a virtual machine asynchronously.
func (c *VirtualEnvironmentClient) StopVMAsync(nodeName string, vmID int) (*string, error) {
	resBody := &VirtualEnvironmentVMStopResponseBody{}
	err := c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu/%d/status/stop", url.PathEscape(nodeName), vmID), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// UpdateVM updates a virtual machine.
func (c *VirtualEnvironmentClient) UpdateVM(nodeName string, vmID int, d *VirtualEnvironmentVMUpdateRequestBody) error {
	return c.DoRequest(hmPUT, fmt.Sprintf("nodes/%s/qemu/%d/config", url.PathEscape(nodeName), vmID), d, nil)
}

// UpdateVMAsync updates a virtual machine asynchronously.
func (c *VirtualEnvironmentClient) UpdateVMAsync(nodeName string, vmID int, d *VirtualEnvironmentVMUpdateRequestBody) (*string, error) {
	resBody := &VirtualEnvironmentVMUpdateAsyncResponseBody{}
	err := c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/qemu/%d/config", url.PathEscape(nodeName), vmID), d, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// WaitForNetworkInterfacesFromVMAgent waits for a virtual machine's QEMU agent to publish the network interfaces.
func (c *VirtualEnvironmentClient) WaitForNetworkInterfacesFromVMAgent(nodeName string, vmID int, timeout int, delay int, waitForIP bool) (*VirtualEnvironmentVMGetQEMUNetworkInterfacesResponseData, error) {
	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			data, err := c.GetVMNetworkInterfacesFromAgent(nodeName, vmID)

			if err == nil && data != nil && data.Result != nil {
				missingIP := false

				if waitForIP {
					for _, nic := range *data.Result {
						if nic.Name == "lo" {
							continue
						}

						if nic.IPAddresses == nil || (nic.IPAddresses != nil && len(*nic.IPAddresses) == 0) {
							missingIP = true
							break
						}
					}
				}

				if !missingIP {
					return data, err
				}
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(200 * time.Millisecond)

		timeElapsed = time.Now().Sub(timeStart)
	}

	return nil, fmt.Errorf("Timeout while waiting for the QEMU agent on VM \"%d\" to publish the network interfaces", vmID)
}

// WaitForNoNetworkInterfacesFromVMAgent waits for a virtual machine's QEMU agent to unpublish the network interfaces.
func (c *VirtualEnvironmentClient) WaitForNoNetworkInterfacesFromVMAgent(nodeName string, vmID int, timeout int, delay int) error {
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

// WaitForVMConfigUnlock waits for a virtual machine configuration to become unlocked.
func (c *VirtualEnvironmentClient) WaitForVMConfigUnlock(nodeName string, vmID int, timeout int, delay int, ignoreErrorResponse bool) error {
	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			data, err := c.GetVMStatus(nodeName, vmID)

			if err != nil {
				if !ignoreErrorResponse {
					return err
				}
			} else if data.Lock == nil || *data.Lock == "" {
				return nil
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(200 * time.Millisecond)

		timeElapsed = time.Now().Sub(timeStart)
	}

	return fmt.Errorf("Timeout while waiting for VM \"%d\" configuration to become unlocked", vmID)
}

// WaitForVMState waits for a virtual machine to reach a specific state.
func (c *VirtualEnvironmentClient) WaitForVMState(nodeName string, vmID int, state string, timeout int, delay int) error {
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
