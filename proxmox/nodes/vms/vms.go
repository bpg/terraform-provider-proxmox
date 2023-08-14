/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// CloneVM clones a virtual machine.
func (c *Client) CloneVM(ctx context.Context, retries int, d *CloneRequestBody, timeout int) error {
	var err error

	resBody := &MoveDiskResponseBody{}

	// just a guard in case someone sets retries to 0 unknowingly
	if retries <= 0 {
		retries = 1
	}

	for i := 0; i < retries; i++ {
		err = c.DoRequest(ctx, http.MethodPost, c.ExpandPath("clone"), d, resBody)

		if err != nil {
			return fmt.Errorf("error cloning VM: %w", err)
		}

		if resBody.Data == nil {
			return api.ErrNoDataObjectInResponse
		}

		err = c.Tasks().WaitForTask(ctx, *resBody.Data, timeout, 5)
		if err == nil {
			return nil
		}

		time.Sleep(10 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("error waiting for VM clone: %w", err)
	}

	return nil
}

// CreateVM creates a virtual machine.
func (c *Client) CreateVM(ctx context.Context, d *CreateRequestBody, timeout int) error {
	taskID, err := c.CreateVMAsync(ctx, d)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID, timeout, 1)

	if err != nil {
		return fmt.Errorf("error waiting for VM creation: %w", err)
	}

	return nil
}

// CreateVMAsync creates a virtual machine asynchronously.
func (c *Client) CreateVMAsync(ctx context.Context, d *CreateRequestBody) (*string, error) {
	resBody := &CreateResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.basePath(), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error creating VM: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// DeleteVM deletes a virtual machine.
func (c *Client) DeleteVM(ctx context.Context) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath("?destroy-unreferenced-disks=1&purge=1"), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting VM: %w", err)
	}

	return nil
}

// GetVM retrieves a virtual machine.
func (c *Client) GetVM(ctx context.Context) (*GetResponseData, error) {
	resBody := &GetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("config"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving VM: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetVMNetworkInterfacesFromAgent retrieves the network interfaces reported by the QEMU agent.
func (c *Client) GetVMNetworkInterfacesFromAgent(ctx context.Context) (*GetQEMUNetworkInterfacesResponseData, error) {
	resBody := &GetQEMUNetworkInterfacesResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("agent/network-get-interfaces"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving VM network interfaces from agent: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetVMStatus retrieves the status for a virtual machine.
func (c *Client) GetVMStatus(ctx context.Context) (*GetStatusResponseData, error) {
	resBody := &GetStatusResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("status/current"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving VM status: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// MigrateVM migrates a virtual machine.
func (c *Client) MigrateVM(ctx context.Context, d *MigrateRequestBody, timeout int) error {
	taskID, err := c.MigrateVMAsync(ctx, d)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID, timeout, 5)
	if err != nil {
		return fmt.Errorf("error waiting for VM migration: %w", err)
	}

	return nil
}

// MigrateVMAsync migrates a virtual machine asynchronously.
func (c *Client) MigrateVMAsync(ctx context.Context, d *MigrateRequestBody) (*string, error) {
	resBody := &MigrateResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("migrate"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error migrating VM: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// MoveVMDisk moves a virtual machine disk.
func (c *Client) MoveVMDisk(ctx context.Context, d *MoveDiskRequestBody, timeout int) error {
	taskID, err := c.MoveVMDiskAsync(ctx, d)
	if err != nil {
		if strings.Contains(err.Error(), "you can't move to the same storage with same format") {
			// if someone tries to move to the same storage, the move is considered to be successful
			return nil
		}

		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID, timeout, 5)
	if err != nil {
		return fmt.Errorf("error waiting for VM disk move: %w", err)
	}

	return nil
}

// MoveVMDiskAsync moves a virtual machine disk asynchronously.
func (c *Client) MoveVMDiskAsync(ctx context.Context, d *MoveDiskRequestBody) (*string, error) {
	resBody := &MoveDiskResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("move_disk"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error moving VM disk: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// ListVMs retrieves a list of virtual machines.
func (c *Client) ListVMs(ctx context.Context) ([]*ListResponseData, error) {
	resBody := &ListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.basePath(), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving VMs: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// RebootVM reboots a virtual machine.
func (c *Client) RebootVM(ctx context.Context, d *RebootRequestBody, timeout int) error {
	taskID, err := c.RebootVMAsync(ctx, d)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID, timeout, 5)
	if err != nil {
		return fmt.Errorf("error waiting for VM reboot: %w", err)
	}

	return nil
}

// RebootVMAsync reboots a virtual machine asynchronously.
func (c *Client) RebootVMAsync(ctx context.Context, d *RebootRequestBody) (*string, error) {
	resBody := &RebootResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("status/reboot"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error rebooting VM: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// ResizeVMDisk resizes a virtual machine disk.
func (c *Client) ResizeVMDisk(ctx context.Context, d *ResizeDiskRequestBody) error {
	var err error

	tflog.Debug(ctx, "resize disk", map[string]interface{}{
		"disk": d.Disk,
		"size": d.Size,
	})

	for i := 0; i < 5; i++ {
		err = c.DoRequest(
			ctx,
			http.MethodPut,
			c.ExpandPath("resize"),
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
			return fmt.Errorf("error resizing VM disk: %w", ctx.Err())
		}
	}

	if err != nil {
		return fmt.Errorf("error resizing VM disk: %w", err)
	}

	return nil
}

// ShutdownVM shuts down a virtual machine.
func (c *Client) ShutdownVM(ctx context.Context, d *ShutdownRequestBody, timeout int) error {
	taskID, err := c.ShutdownVMAsync(ctx, d)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID, timeout, 5)
	if err != nil {
		return fmt.Errorf("error waiting for VM shutdown: %w", err)
	}

	return nil
}

// ShutdownVMAsync shuts down a virtual machine asynchronously.
func (c *Client) ShutdownVMAsync(ctx context.Context, d *ShutdownRequestBody) (*string, error) {
	resBody := &ShutdownResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("status/shutdown"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error shutting down VM: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// StartVM starts a virtual machine.
func (c *Client) StartVM(ctx context.Context, timeout int) error {
	taskID, err := c.StartVMAsync(ctx)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID, timeout, 5)
	if err != nil {
		return fmt.Errorf("error waiting for VM start: %w", err)
	}

	return nil
}

// StartVMAsync starts a virtual machine asynchronously.
func (c *Client) StartVMAsync(ctx context.Context) (*string, error) {
	resBody := &StartResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("status/start"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error starting VM: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// StopVM stops a virtual machine.
func (c *Client) StopVM(ctx context.Context, timeout int) error {
	taskID, err := c.StopVMAsync(ctx)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID, timeout, 5)
	if err != nil {
		return fmt.Errorf("error waiting for VM stop: %w", err)
	}

	return nil
}

// StopVMAsync stops a virtual machine asynchronously.
func (c *Client) StopVMAsync(ctx context.Context) (*string, error) {
	resBody := &StopResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("status/stop"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error stopping VM: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// UpdateVM updates a virtual machine.
func (c *Client) UpdateVM(ctx context.Context, d *UpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath("config"), d, nil)
	if err != nil {
		return fmt.Errorf("error updating VM: %w", err)
	}

	return nil
}

// UpdateVMAsync updates a virtual machine asynchronously.
func (c *Client) UpdateVMAsync(ctx context.Context, d *UpdateRequestBody) (*string, error) {
	resBody := &UpdateAsyncResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("config"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error updating VM: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// WaitForNetworkInterfacesFromVMAgent waits for a virtual machine's QEMU agent to publish the network interfaces.
func (c *Client) WaitForNetworkInterfacesFromVMAgent(
	ctx context.Context,
	timeout int,
	delay int,
	waitForIP bool,
) (*GetQEMUNetworkInterfacesResponseData, error) {
	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		//nolint:nestif
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			data, err := c.GetVMNetworkInterfacesFromAgent(ctx)

			if err == nil && data != nil && data.Result != nil {
				hasAnyGlobalUnicast := false

				if waitForIP {
					for _, nic := range *data.Result {
						if nic.Name == "lo" {
							continue
						}

						if nic.IPAddresses == nil ||
							(nic.IPAddresses != nil && len(*nic.IPAddresses) == 0) {
							continue
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
			return nil, fmt.Errorf("error waiting for VM network interfaces: %w", ctx.Err())
		}
	}

	return nil, fmt.Errorf(
		"timeout while waiting for the QEMU agent on VM \"%d\" to publish the network interfaces",
		c.VMID,
	)
}

// WaitForNoNetworkInterfacesFromVMAgent waits for a virtual machine's QEMU agent to unpublish the network interfaces.
func (c *Client) WaitForNoNetworkInterfacesFromVMAgent(ctx context.Context, timeout int, delay int) error {
	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			_, err := c.GetVMNetworkInterfacesFromAgent(ctx)
			if err == nil {
				return nil
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(200 * time.Millisecond)

		timeElapsed = time.Since(timeStart)

		if ctx.Err() != nil {
			return fmt.Errorf("error waiting for VM network interfaces: %w", ctx.Err())
		}
	}

	return fmt.Errorf(
		"timeout while waiting for the QEMU agent on VM \"%d\" to unpublish the network interfaces",
		c.VMID,
	)
}

// WaitForVMConfigUnlock waits for a virtual machine configuration to become unlocked.
func (c *Client) WaitForVMConfigUnlock(ctx context.Context, timeout int, delay int, ignoreErrorResponse bool) error {
	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			data, err := c.GetVMStatus(ctx)

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
			return fmt.Errorf("error waiting for VM configuration to become unlocked: %w", ctx.Err())
		}
	}

	return fmt.Errorf("timeout while waiting for VM \"%d\" configuration to become unlocked", c.VMID)
}

// WaitForVMState waits for a virtual machine to reach a specific state.
func (c *Client) WaitForVMState(ctx context.Context, state string, timeout int, delay int) error {
	state = strings.ToLower(state)

	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			data, err := c.GetVMStatus(ctx)
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
			return fmt.Errorf("error waiting for VM state: %w", ctx.Err())
		}
	}

	return fmt.Errorf("timeout while waiting for VM \"%d\" to enter the state \"%s\"", c.VMID, state)
}
