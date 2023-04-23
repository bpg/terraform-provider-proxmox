/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	getVMIDStep = 1
)

var (
	getVMIDCounter      = -1
	getVMIDCounterMutex = &sync.Mutex{}
)

// CloneVM clones a virtual machine.
func (c *VirtualEnvironmentClient) CloneVM(
	ctx context.Context,
	nodeName string,
	vmID int,
	retries int,
	d *VirtualEnvironmentVMCloneRequestBody,
	timeout int,
) error {
	resBody := &VirtualEnvironmentVMMoveDiskResponseBody{}
	var err error

	// just a guard in case someone sets retries to 0 unknowingly
	if retries <= 0 {
		retries = 1
	}

	for i := 0; i < retries; i++ {
		err = c.DoRequest(
			ctx,
			http.MethodPost,
			fmt.Sprintf("nodes/%s/qemu/%d/clone", url.PathEscape(nodeName), vmID),
			d,
			resBody,
		)

		if err != nil {
			return err
		}

		if resBody.Data == nil {
			return errors.New("the server did not include a data object in the response")
		}

		err = c.WaitForNodeTask(ctx, nodeName, *resBody.Data, timeout, 5)

		if err == nil {
			return nil
		}
		time.Sleep(10 * time.Second)
	}

	return err
}

// CreateVM creates a virtual machine.
func (c *VirtualEnvironmentClient) CreateVM(
	ctx context.Context,
	nodeName string,
	d *VirtualEnvironmentVMCreateRequestBody,
	timeout int,
) error {
	taskID, err := c.CreateVMAsync(ctx, nodeName, d)
	if err != nil {
		return err
	}

	err = c.WaitForNodeTask(ctx, nodeName, *taskID, timeout, 1)

	if err != nil {
		return fmt.Errorf("error waiting for VM creation: %w", err)
	}

	return nil
}

// CreateVMAsync creates a virtual machine asynchronously.
func (c *VirtualEnvironmentClient) CreateVMAsync(
	ctx context.Context,
	nodeName string,
	d *VirtualEnvironmentVMCreateRequestBody,
) (*string, error) {
	resBody := &VirtualEnvironmentVMCreateResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("nodes/%s/qemu", url.PathEscape(nodeName)),
		d,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// DeleteVM deletes a virtual machine.
func (c *VirtualEnvironmentClient) DeleteVM(ctx context.Context, nodeName string, vmID int) error {
	return c.DoRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf(
			"nodes/%s/qemu/%d?destroy-unreferenced-disks=1&purge=1",
			url.PathEscape(nodeName),
			vmID,
		),
		nil,
		nil,
	)
}

// GetVM retrieves a virtual machine.
func (c *VirtualEnvironmentClient) GetVM(
	ctx context.Context,
	nodeName string,
	vmID int,
) (*VirtualEnvironmentVMGetResponseData, error) {
	resBody := &VirtualEnvironmentVMGetResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("nodes/%s/qemu/%d/config", url.PathEscape(nodeName), vmID),
		nil,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// GetVMID retrieves the next available VM identifier.
func (c *VirtualEnvironmentClient) GetVMID(ctx context.Context) (*int, error) {
	getVMIDCounterMutex.Lock()
	defer getVMIDCounterMutex.Unlock()

	if getVMIDCounter < 0 {
		nextVMID, err := c.API().Cluster().GetNextID(ctx, nil)
		if err != nil {
			return nil, err
		}

		if nextVMID == nil {
			return nil, errors.New("unable to retrieve the next available VM identifier")
		}

		getVMIDCounter = *nextVMID + getVMIDStep

		tflog.Debug(ctx, "next VM identifier", map[string]interface{}{
			"id": *nextVMID,
		})

		return nextVMID, nil
	}

	vmID := getVMIDCounter

	for vmID <= 2147483637 {
		_, err := c.API().Cluster().GetNextID(ctx, &vmID)
		if err != nil {
			vmID += getVMIDStep

			continue
		}

		getVMIDCounter = vmID + getVMIDStep

		tflog.Debug(ctx, "next VM identifier", map[string]interface{}{
			"id": vmID,
		})

		return &vmID, nil
	}

	return nil, errors.New("unable to determine the next available VM identifier")
}

// GetVMNetworkInterfacesFromAgent retrieves the network interfaces reported by the QEMU agent.
func (c *VirtualEnvironmentClient) GetVMNetworkInterfacesFromAgent(
	ctx context.Context,
	nodeName string,
	vmID int,
) (*VirtualEnvironmentVMGetQEMUNetworkInterfacesResponseData, error) {
	resBody := &VirtualEnvironmentVMGetQEMUNetworkInterfacesResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"nodes/%s/qemu/%d/agent/network-get-interfaces",
			url.PathEscape(nodeName),
			vmID,
		),
		nil,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// GetVMStatus retrieves the status for a virtual machine.
func (c *VirtualEnvironmentClient) GetVMStatus(
	ctx context.Context,
	nodeName string,
	vmID int,
) (*VirtualEnvironmentVMGetStatusResponseData, error) {
	resBody := &VirtualEnvironmentVMGetStatusResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("nodes/%s/qemu/%d/status/current", url.PathEscape(nodeName), vmID),
		nil,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// MigrateVM migrates a virtual machine.
func (c *VirtualEnvironmentClient) MigrateVM(
	ctx context.Context,
	nodeName string,
	vmID int,
	d *VirtualEnvironmentVMMigrateRequestBody,
	timeout int,
) error {
	taskID, err := c.MigrateVMAsync(ctx, nodeName, vmID, d)
	if err != nil {
		return err
	}

	err = c.WaitForNodeTask(ctx, nodeName, *taskID, timeout, 5)

	if err != nil {
		return err
	}

	return nil
}

// MigrateVMAsync migrates a virtual machine asynchronously.
func (c *VirtualEnvironmentClient) MigrateVMAsync(
	ctx context.Context,
	nodeName string,
	vmID int,
	d *VirtualEnvironmentVMMigrateRequestBody,
) (*string, error) {
	resBody := &VirtualEnvironmentVMMigrateResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("nodes/%s/qemu/%d/migrate", url.PathEscape(nodeName), vmID),
		d,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// MoveVMDisk moves a virtual machine disk.
func (c *VirtualEnvironmentClient) MoveVMDisk(
	ctx context.Context,
	nodeName string,
	vmID int,
	d *VirtualEnvironmentVMMoveDiskRequestBody,
	timeout int,
) error {
	taskID, err := c.MoveVMDiskAsync(ctx, nodeName, vmID, d)
	if err != nil {
		if strings.Contains(err.Error(), "you can't move to the same storage with same format") {
			// if someone tries to move to the same storage, the move is considered to be successful
			return nil
		}

		return err
	}

	err = c.WaitForNodeTask(ctx, nodeName, *taskID, timeout, 5)

	if err != nil {
		return err
	}

	return nil
}

// MoveVMDiskAsync moves a virtual machine disk asynchronously.
func (c *VirtualEnvironmentClient) MoveVMDiskAsync(
	ctx context.Context,
	nodeName string,
	vmID int,
	d *VirtualEnvironmentVMMoveDiskRequestBody,
) (*string, error) {
	resBody := &VirtualEnvironmentVMMoveDiskResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("nodes/%s/qemu/%d/move_disk", url.PathEscape(nodeName), vmID),
		d,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListVMs retrieves a list of virtual machines.
func (c *VirtualEnvironmentClient) ListVMs(
	ctx context.Context,
	nodeName string,
) ([]*VirtualEnvironmentVMListResponseData, error) {
	resBody := &VirtualEnvironmentVMListResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("nodes/%s/qemu", url.PathEscape(nodeName)),
		nil,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// RebootVM reboots a virtual machine.
func (c *VirtualEnvironmentClient) RebootVM(
	ctx context.Context,
	nodeName string,
	vmID int,
	d *VirtualEnvironmentVMRebootRequestBody,
	timeout int,
) error {
	taskID, err := c.RebootVMAsync(ctx, nodeName, vmID, d)
	if err != nil {
		return err
	}

	err = c.WaitForNodeTask(ctx, nodeName, *taskID, timeout, 5)

	if err != nil {
		return err
	}

	return nil
}

// RebootVMAsync reboots a virtual machine asynchronously.
func (c *VirtualEnvironmentClient) RebootVMAsync(
	ctx context.Context,
	nodeName string,
	vmID int,
	d *VirtualEnvironmentVMRebootRequestBody,
) (*string, error) {
	resBody := &VirtualEnvironmentVMRebootResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("nodes/%s/qemu/%d/status/reboot", url.PathEscape(nodeName), vmID),
		d,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ResizeVMDisk resizes a virtual machine disk.
func (c *VirtualEnvironmentClient) ResizeVMDisk(
	ctx context.Context,
	nodeName string,
	vmID int,
	d *VirtualEnvironmentVMResizeDiskRequestBody,
) error {
	var err error
	tflog.Debug(ctx, "resize disk", map[string]interface{}{
		"disk": d.Disk,
		"size": d.Size,
	})
	for i := 0; i < 5; i++ {
		err = c.DoRequest(
			ctx,
			http.MethodPut,
			fmt.Sprintf("nodes/%s/qemu/%d/resize", url.PathEscape(nodeName), vmID),
			d,
			nil,
		)
		if err == nil {
			return nil
		}
		tflog.Debug(ctx, "resize disk failed", map[string]interface{}{
			"retry": i,
		})
		time.Sleep(5 * time.Second)

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return err
}

// ShutdownVM shuts down a virtual machine.
func (c *VirtualEnvironmentClient) ShutdownVM(
	ctx context.Context,
	nodeName string,
	vmID int,
	d *VirtualEnvironmentVMShutdownRequestBody,
	timeout int,
) error {
	taskID, err := c.ShutdownVMAsync(ctx, nodeName, vmID, d)
	if err != nil {
		return err
	}

	err = c.WaitForNodeTask(ctx, nodeName, *taskID, timeout, 5)

	if err != nil {
		return err
	}

	return nil
}

// ShutdownVMAsync shuts down a virtual machine asynchronously.
func (c *VirtualEnvironmentClient) ShutdownVMAsync(
	ctx context.Context,
	nodeName string,
	vmID int,
	d *VirtualEnvironmentVMShutdownRequestBody,
) (*string, error) {
	resBody := &VirtualEnvironmentVMShutdownResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("nodes/%s/qemu/%d/status/shutdown", url.PathEscape(nodeName), vmID),
		d,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// StartVM starts a virtual machine.
func (c *VirtualEnvironmentClient) StartVM(
	ctx context.Context,
	nodeName string,
	vmID int,
	timeout int,
) error {
	taskID, err := c.StartVMAsync(ctx, nodeName, vmID)
	if err != nil {
		return err
	}

	err = c.WaitForNodeTask(ctx, nodeName, *taskID, timeout, 5)

	if err != nil {
		return err
	}

	return nil
}

// StartVMAsync starts a virtual machine asynchronously.
func (c *VirtualEnvironmentClient) StartVMAsync(
	ctx context.Context,
	nodeName string,
	vmID int,
) (*string, error) {
	resBody := &VirtualEnvironmentVMStartResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("nodes/%s/qemu/%d/status/start", url.PathEscape(nodeName), vmID),
		nil,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// StopVM stops a virtual machine.
func (c *VirtualEnvironmentClient) StopVM(
	ctx context.Context,
	nodeName string,
	vmID int,
	timeout int,
) error {
	taskID, err := c.StopVMAsync(ctx, nodeName, vmID)
	if err != nil {
		return err
	}

	err = c.WaitForNodeTask(ctx, nodeName, *taskID, timeout, 5)

	if err != nil {
		return err
	}

	return nil
}

// StopVMAsync stops a virtual machine asynchronously.
func (c *VirtualEnvironmentClient) StopVMAsync(
	ctx context.Context,
	nodeName string,
	vmID int,
) (*string, error) {
	resBody := &VirtualEnvironmentVMStopResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("nodes/%s/qemu/%d/status/stop", url.PathEscape(nodeName), vmID),
		nil,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// UpdateVM updates a virtual machine.
func (c *VirtualEnvironmentClient) UpdateVM(
	ctx context.Context,
	nodeName string,
	vmID int,
	d *VirtualEnvironmentVMUpdateRequestBody,
) error {
	return c.DoRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("nodes/%s/qemu/%d/config", url.PathEscape(nodeName), vmID),
		d,
		nil,
	)
}

// UpdateVMAsync updates a virtual machine asynchronously.
func (c *VirtualEnvironmentClient) UpdateVMAsync(
	ctx context.Context,
	nodeName string,
	vmID int,
	d *VirtualEnvironmentVMUpdateRequestBody,
) (*string, error) {
	resBody := &VirtualEnvironmentVMUpdateAsyncResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("nodes/%s/qemu/%d/config", url.PathEscape(nodeName), vmID),
		d,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// WaitForNetworkInterfacesFromVMAgent waits for a virtual machine's QEMU agent to publish the network interfaces.
func (c *VirtualEnvironmentClient) WaitForNetworkInterfacesFromVMAgent(
	ctx context.Context,
	nodeName string,
	vmID int,
	timeout int,
	delay int,
	waitForIP bool,
) (*VirtualEnvironmentVMGetQEMUNetworkInterfacesResponseData, error) {
	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			data, err := c.GetVMNetworkInterfacesFromAgent(ctx, nodeName, vmID)

			if err == nil && data != nil && data.Result != nil {
				hasAnyGlobalUnicast := false

				if waitForIP {
					for _, nic := range *data.Result {
						if nic.Name == "lo" {
							continue
						}

						if nic.IPAddresses == nil ||
							(nic.IPAddresses != nil && len(*nic.IPAddresses) == 0) {
							break
						}

						for _, addr := range *nic.IPAddresses {
							if ip := net.ParseIP(addr.Address); ip != nil && ip.IsGlobalUnicast() {
								hasAnyGlobalUnicast = true
							}
						}
					}
				}

				if hasAnyGlobalUnicast {
					return data, err
				}
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(200 * time.Millisecond)

		timeElapsed = time.Since(timeStart)

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf(
		"timeout while waiting for the QEMU agent on VM \"%d\" to publish the network interfaces",
		vmID,
	)
}

// WaitForNoNetworkInterfacesFromVMAgent waits for a virtual machine's QEMU agent to unpublish the network interfaces.
func (c *VirtualEnvironmentClient) WaitForNoNetworkInterfacesFromVMAgent(
	ctx context.Context,
	nodeName string,
	vmID int,
	timeout int,
	delay int,
) error {
	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			_, err := c.GetVMNetworkInterfacesFromAgent(ctx, nodeName, vmID)
			if err != nil {
				return nil
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(200 * time.Millisecond)

		timeElapsed = time.Since(timeStart)

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	return fmt.Errorf(
		"timeout while waiting for the QEMU agent on VM \"%d\" to unpublish the network interfaces",
		vmID,
	)
}

// WaitForVMConfigUnlock waits for a virtual machine configuration to become unlocked.
func (c *VirtualEnvironmentClient) WaitForVMConfigUnlock(
	ctx context.Context,
	nodeName string,
	vmID int,
	timeout int,
	delay int,
	ignoreErrorResponse bool,
) error {
	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			data, err := c.GetVMStatus(ctx, nodeName, vmID)

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

		timeElapsed = time.Since(timeStart)

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	return fmt.Errorf("timeout while waiting for VM \"%d\" configuration to become unlocked", vmID)
}

// WaitForVMState waits for a virtual machine to reach a specific state.
func (c *VirtualEnvironmentClient) WaitForVMState(
	ctx context.Context,
	nodeName string,
	vmID int,
	state string,
	timeout int,
	delay int,
) error {
	state = strings.ToLower(state)

	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			data, err := c.GetVMStatus(ctx, nodeName, vmID)
			if err != nil {
				return err
			}

			if data.Status == state {
				return nil
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(200 * time.Millisecond)

		timeElapsed = time.Since(timeStart)

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	return fmt.Errorf("timeout while waiting for VM \"%d\" to enter the state \"%s\"", vmID, state)
}
