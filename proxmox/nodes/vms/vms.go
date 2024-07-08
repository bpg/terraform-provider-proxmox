/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
)

// CloneVM clones a virtual machine.
func (c *Client) CloneVM(ctx context.Context, retries int, d *CloneRequestBody) error {
	var err error

	resBody := &MoveDiskResponseBody{}

	// just a guard in case someone sets retries to zero unknowingly
	if retries <= 0 {
		retries = 1
	}

	err = retry.Do(func() error {
		err = c.DoRequest(ctx, http.MethodPost, c.ExpandPath("clone"), d, resBody)
		if err != nil {
			return fmt.Errorf("error cloning VM: %w", err)
		}

		if resBody.Data == nil {
			return api.ErrNoDataObjectInResponse
		}

		// ignoring warnings as per https://www.mail-archive.com/pve-devel@lists.proxmox.com/msg17724.html
		return c.Tasks().WaitForTask(ctx, *resBody.Data, tasks.WithIgnoreWarnings())
	}, retry.Attempts(uint(retries)), retry.Delay(10*time.Second))
	if err != nil {
		return fmt.Errorf("error waiting for VM clone: %w", err)
	}

	return nil
}

// CreateVM creates a virtual machine.
func (c *Client) CreateVM(ctx context.Context, d *CreateRequestBody) error {
	taskID, err := c.CreateVMAsync(ctx, d)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID)
	if err != nil {
		return fmt.Errorf("error waiting for VM creation: %w", err)
	}

	return nil
}

// CreateVMAsync creates a virtual machine asynchronously. Returns ID of the started task.
func (c *Client) CreateVMAsync(ctx context.Context, d *CreateRequestBody) (*string, error) {
	resBody := &CreateResponseBody{}

	err := retry.Do(
		func() error {
			return c.DoRequest(ctx, http.MethodPost, c.basePath(), d, resBody)
		},
		retry.Context(ctx),
		retry.OnRetry(func(n uint, err error) {
			tflog.Warn(ctx, "retrying VM creation", map[string]interface{}{
				"attempt": n,
				"error":   err.Error(),
			})

			e := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath("?destroy-unreferenced-disks=1&purge=1"), nil, nil)
			if e != nil {
				tflog.Warn(ctx, "deleting VM after failed creation", map[string]interface{}{
					"error": e,
				})
			}
		}),
		retry.LastErrorOnly(false),
		retry.Attempts(3),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating VM: %w", err)
	}

	if resBody.TaskID == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.TaskID, nil
}

// DeleteVM creates a virtual machine.
func (c *Client) DeleteVM(ctx context.Context) error {
	taskID, err := c.DeleteVMAsync(ctx)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID)
	if err != nil {
		return fmt.Errorf("error waiting for VM deletion: %w", err)
	}

	return nil
}

// DeleteVMAsync deletes a virtual machine asynchronously. Returns ID of the started task.
func (c *Client) DeleteVMAsync(ctx context.Context) (*string, error) {
	// PVE may return a 500 error "got no worker upid - start worker failed", so we retry few times.
	resBody := &DeleteResponseBody{}

	err := retry.Do(
		func() error {
			return c.DoRequest(ctx, http.MethodDelete, c.ExpandPath("?destroy-unreferenced-disks=1&purge=1"), nil, resBody)
		},
		retry.Context(ctx),
		retry.RetryIf(func(err error) bool {
			return !errors.Is(err, api.ErrResourceDoesNotExist)
		}),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting VM: %w", err)
	}

	return resBody.TaskID, nil
}

// GetVM retrieves a virtual machine.
func (c *Client) GetVM(ctx context.Context) (*GetResponseData, error) {
	resBody := &GetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("config"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving VM %d: %w", c.VMID, err)
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
func (c *Client) MigrateVM(ctx context.Context, d *MigrateRequestBody) error {
	taskID, err := c.MigrateVMAsync(ctx, d)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID)
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
func (c *Client) MoveVMDisk(ctx context.Context, d *MoveDiskRequestBody) error {
	taskID, err := c.MoveVMDiskAsync(ctx, d)
	if err != nil {
		if strings.Contains(err.Error(), "you can't move to the same storage with same format") {
			// if someone tries to move to the same storage, the move is considered to be successful
			return nil
		}

		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID)
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
func (c *Client) RebootVM(ctx context.Context, d *RebootRequestBody) error {
	taskID, err := c.RebootVMAsync(ctx, d)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID)
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
	err := retry.Do(
		func() error {
			taskID, err := c.ResizeVMDiskAsync(ctx, d)
			if err != nil {
				return err
			}

			return c.Tasks().WaitForTask(ctx, *taskID)
		},
		retry.Context(ctx),
		retry.Attempts(3),
		retry.Delay(1*time.Second),
		retry.LastErrorOnly(false),
		retry.RetryIf(func(err error) bool {
			return strings.Contains(err.Error(), "got timeout")
		}),
	)
	if err != nil {
		return fmt.Errorf("error waiting for VM disk resize: %w", err)
	}

	return nil
}

// ResizeVMDiskAsync resizes a virtual machine disk asynchronously.
func (c *Client) ResizeVMDiskAsync(ctx context.Context, d *ResizeDiskRequestBody) (*string, error) {
	resBody := &MoveDiskResponseBody{}

	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath("resize"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error moving VM disk: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// ShutdownVM shuts down a virtual machine.
func (c *Client) ShutdownVM(ctx context.Context, d *ShutdownRequestBody) error {
	taskID, err := c.ShutdownVMAsync(ctx, d)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID)
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
// Returns the task log if the VM had warnings at startup, or fails to start.
func (c *Client) StartVM(ctx context.Context, timeoutSec int) ([]string, error) {
	taskID, err := c.StartVMAsync(ctx, timeoutSec)
	if err != nil {
		return nil, err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID, tasks.WithIgnoreStatus(599))
	if err != nil {
		log, e := c.Tasks().GetTaskLog(ctx, *taskID)
		if e != nil {
			tflog.Error(ctx, "error retrieving task log", map[string]interface{}{
				"task_id": *taskID,
				"error":   e.Error(),
			})

			log = []string{}
		}

		if strings.Contains(err.Error(), "WARNING") && len(log) > 0 {
			return log, nil
		}

		return log, fmt.Errorf("error waiting for VM start: %w", err)
	}

	return nil, nil
}

// StartVMAsync starts a virtual machine asynchronously.
func (c *Client) StartVMAsync(ctx context.Context, timeoutSec int) (*string, error) {
	reqBody := &StartRequestBody{
		TimeoutSeconds: &timeoutSec,
	}
	resBody := &StartResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("status/start"), reqBody, resBody)
	if err != nil {
		return nil, fmt.Errorf("error starting VM: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// StopVM stops a virtual machine.
func (c *Client) StopVM(ctx context.Context) error {
	taskID, err := c.StopVMAsync(ctx)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID)
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
	timeout int, // time in seconds to wait until giving up
	delay int, // the delay in seconds between requests to the agent
	waitForIP bool, // whether to block until an IP is found, or just block until the interfaces are published
) (*GetQEMUNetworkInterfacesResponseData, error) {
	delaySeconds := int64(delay)
	timeMaxSeconds := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	for timeElapsed.Seconds() < timeMaxSeconds {
		timeElapsed = time.Since(timeStart)

		// check if terraform wants to shut us down (we try to poll the ctx every 200ms)
		if ctx.Err() != nil {
			return nil, fmt.Errorf("error waiting for VM network interfaces: %w", ctx.Err())
		}

		select {
		case <-ch:
			{
				// the returned error will be eaten by the terraform runtime, so we log it here as well
				const msg = "interrupted by signal"

				tflog.Warn(ctx, msg)

				return nil, fmt.Errorf(msg)
			}
		default:
		}

		// sleep another 200 milliseconds if we haven't delayed enough since our last call
		if int64(timeElapsed.Seconds())%delaySeconds != 0 {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		// request the network interfaces from the agent
		data, err := c.GetVMNetworkInterfacesFromAgent(ctx)

		// tick ahead and continue if we got an error from the api
		if err != nil || data == nil || data.Result == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// If we're waiting for an IP, check if we have one yet; if not then keep looping
		if waitForIP {
			for _, nic := range *data.Result {
				// skip the loopback interface
				if nic.Name == "lo" {
					continue
				}

				// skip the interface if it has no IP addresses
				if nic.IPAddresses == nil ||
					(nic.IPAddresses != nil && len(*nic.IPAddresses) == 0) {
					continue
				}

				// return if the interface has any global unicast addresses
				for _, addr := range *nic.IPAddresses {
					if ip := net.ParseIP(addr.Address); ip != nil && ip.IsGlobalUnicast() {
						return data, err
					}
				}
			}

			// no IP address has come through the agent yet
			time.Sleep(1 * time.Second)
			continue //nolint
		}

		// if not waiting for an IP, and the agent sent us an interface, return
		if data.Result != nil && len(*data.Result) > 0 {
			return data, err
		}

		// we didn't get any interfaces so tick ahead to keep looping
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf(
		"timeout while waiting for the QEMU agent on VM \"%d\" to publish the network interfaces",
		c.VMID,
	)
}

// WaitForVMConfigUnlock waits for a virtual machine configuration to become unlocked.
func (c *Client) WaitForVMConfigUnlock(ctx context.Context, ignoreErrorResponse bool) error {
	stillLocked := errors.New("still locked")

	err := retry.Do(
		func() error {
			data, err := c.GetVMStatus(ctx)
			if err != nil {
				return err
			}

			if data.Lock != nil && *data.Lock != "" {
				return stillLocked
			}

			return nil
		},
		retry.Context(ctx),
		retry.RetryIf(func(err error) bool {
			return errors.Is(err, stillLocked) || ignoreErrorResponse
		}),
		retry.UntilSucceeded(),
		retry.Delay(1*time.Second),
		retry.LastErrorOnly(true),
	)

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("timeout while waiting for VM %d configuration to become unlocked", c.VMID)
	}

	if err != nil && !ignoreErrorResponse {
		return fmt.Errorf("error waiting for VM %d configuration to become unlocked: %w", c.VMID, err)
	}

	return nil
}

// WaitForVMStatus waits for a virtual machine to reach a specific status.
func (c *Client) WaitForVMStatus(ctx context.Context, status string) error {
	status = strings.ToLower(status)

	unexpectedStatus := fmt.Errorf("unexpected status %q", status)

	err := retry.Do(
		func() error {
			data, err := c.GetVMStatus(ctx)
			if err != nil {
				return err
			}

			if data.Status != status {
				return unexpectedStatus
			}

			return nil
		},
		retry.Context(ctx),
		retry.RetryIf(func(err error) bool {
			return errors.Is(err, unexpectedStatus)
		}),
		retry.UntilSucceeded(),
		retry.Delay(1*time.Second),
		retry.LastErrorOnly(true),
	)

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("timeout while waiting for VM %d to enter the status %q", c.VMID, status)
	}

	if err != nil {
		return fmt.Errorf("error waiting for VM %d to enter the status %q: %w", c.VMID, status, err)
	}

	return nil
}
