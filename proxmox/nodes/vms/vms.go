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
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
	"github.com/bpg/terraform-provider-proxmox/proxmox/retry"
	"github.com/bpg/terraform-provider-proxmox/utils/ip"
)

// CloneVM clones a virtual machine.
func (c *Client) CloneVM(ctx context.Context, retries int, d *CloneRequestBody) error {
	// just a guard in case someone sets retries to zero unknowingly
	if retries <= 0 {
		retries = 1
	}

	op := retry.NewTaskOperation("VM clone",
		retry.WithAttempts(uint(retries)),
		retry.WithBaseDelay(10*time.Second),
		retry.WithRetryIf(retry.IsTransientAPIError),
		retry.WithAlreadyDoneCheck(retry.ErrorContains("already exists")),
	)

	return op.DoTask(ctx,
		func() (*string, error) { return c.CloneVMAsync(ctx, d) },
		func(ctx context.Context, taskID string) error {
			return c.Tasks().WaitForTask(ctx, taskID, tasks.WithIgnoreWarnings())
		},
	)
}

// CloneVMAsync clones a virtual machine asynchronously.
func (c *Client) CloneVMAsync(ctx context.Context, d *CloneRequestBody) (*string, error) {
	resBody := &CloneResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("clone"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error cloning VM: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// ConvertToTemplate converts a virtual machine to a template using proper endpoint.
func (c *Client) ConvertToTemplate(ctx context.Context) error {
	resBody := &UpdateAsyncResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("template"), nil, resBody)
	if err != nil {
		return fmt.Errorf("error converting VM %d to template: %w", c.VMID, err)
	}

	if resBody.Data != nil {
		err = c.Tasks().WaitForTask(ctx, *resBody.Data)
		if err != nil {
			return fmt.Errorf("error waiting for VM %d template conversion: %w", c.VMID, err)
		}
	}

	return nil
}

// CreateVM creates a virtual machine.
func (c *Client) CreateVM(ctx context.Context, d *CreateRequestBody) error {
	op := retry.NewTaskOperation("VM create",
		retry.WithRetryIf(retry.ErrorContains("got no worker upid")),
		retry.WithAlreadyDoneCheck(retry.ErrorContains("already exists")),
	)

	return op.DoTask(ctx,
		func() (*string, error) { return c.CreateVMAsync(ctx, d) },
		func(ctx context.Context, taskID string) error { return c.Tasks().WaitForTask(ctx, taskID) },
	)
}

// CreateVMAsync creates a virtual machine asynchronously.
func (c *Client) CreateVMAsync(ctx context.Context, d *CreateRequestBody) (*string, error) {
	resBody := &CreateResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.basePath(), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error creating VM: %w", err)
	}

	if resBody.TaskID == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.TaskID, nil
}

// DeleteVM deletes a virtual machine.
func (c *Client) DeleteVM(ctx context.Context, purge bool, destroyUnreferencedDisks bool) error {
	purgeValue := 0
	if purge {
		purgeValue = 1
	}

	destroyUnreferencedDisksValue := 0
	if destroyUnreferencedDisks {
		destroyUnreferencedDisksValue = 1
	}

	op := retry.NewTaskOperation("VM delete",
		retry.WithRetryIf(func(err error) bool {
			return !errors.Is(err, api.ErrResourceDoesNotExist)
		}),
	)

	return op.DoTask(ctx,
		func() (*string, error) {
			resBody := &DeleteResponseBody{}
			path := fmt.Sprintf("?destroy-unreferenced-disks=%d&purge=%d", destroyUnreferencedDisksValue, purgeValue)

			err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(path), nil, resBody)
			if err != nil {
				return nil, fmt.Errorf("error deleting VM: %w", err)
			}

			if resBody.TaskID == nil {
				return nil, api.ErrNoDataObjectInResponse
			}

			return resBody.TaskID, nil
		},
		func(ctx context.Context, taskID string) error { return c.Tasks().WaitForTask(ctx, taskID) },
	)
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

// RebuildCloudInitDisk regenerates and changes cloud-init config drive.
func (c *Client) RebuildCloudInitDisk(ctx context.Context) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath("cloudinit"), nil, nil)
	if err != nil {
		return fmt.Errorf("error rebuilding cloud-init drive: %w", err)
	}

	return nil
}

// RebootVMAndWaitForRunning reboots a virtual machine and waits for it to be running.
func (c *Client) RebootVMAndWaitForRunning(ctx context.Context, rebootTimeoutSec int) error {
	// We add 3 seconds padding to the timeout to account for retries and delays down the callstack.
	ctx, cancel := context.WithTimeout(ctx, time.Duration(rebootTimeoutSec+3)*time.Second)
	defer cancel()

	err := c.RebootVM(
		ctx,
		&RebootRequestBody{
			Timeout: &rebootTimeoutSec,
		},
	)
	if err != nil {
		return err
	}

	return c.WaitForVMStatus(ctx, "running")
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
	// Retry wraps the entire operation (dispatch + wait) because "does not exist"
	// errors can come from WaitForTask on NFS storage when a disk was just moved
	// and the storage needs time to sync.
	op := retry.NewAPICallOperation("VM disk resize",
		retry.WithAttempts(5),
		retry.WithRetryIf(func(err error) bool {
			return strings.Contains(err.Error(), "got timeout") ||
				strings.Contains(err.Error(), "does not exist")
		}),
	)

	return op.Do(ctx, func() error {
		taskID, err := c.ResizeVMDiskAsync(ctx, d)
		if err != nil {
			return err
		}

		return c.Tasks().WaitForTask(ctx, *taskID)
	})
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
			tflog.Error(ctx, "error retrieving task log", map[string]any{
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

	op := retry.NewAPICallOperation("VM start",
		retry.WithRetryIf(retry.ErrorContains("got no worker upid")),
	)

	err := op.Do(ctx, func() error {
		err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("status/start"), reqBody, resBody)
		if err != nil && strings.Contains(err.Error(), "already running") {
			return nil
		}

		return err
	})
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
	op := retry.NewAPICallOperation("VM config update",
		retry.WithRetryIf(retry.ErrorContains("got timeout")),
	)

	return op.Do(ctx, func() error {
		return c.DoRequest(ctx, http.MethodPut, c.ExpandPath("config"), d, nil)
	})
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

// isAgentNotReadyError checks if an HTTP error indicates the agent is not ready yet.
// This includes HTTP 500 errors with messages like "QEMU guest agent is not running"
// which can occur with certain SCSI controller types (e.g., virtio-scsi-single).
func isAgentNotReadyError(err error) bool {
	var httpError *api.HTTPError
	if !errors.As(err, &httpError) {
		return false
	}

	if httpError.Code == http.StatusBadRequest {
		return true
	}

	if httpError.Code == http.StatusInternalServerError {
		msg := strings.ToLower(httpError.Message)

		return strings.Contains(msg, "qemu guest agent") &&
			(strings.Contains(msg, "not running") ||
				strings.Contains(msg, "not available") ||
				strings.Contains(msg, "not ready"))
	}

	return false
}

// WaitForNetworkInterfacesFromVMAgent waits for a virtual machine's QEMU agent to publish the network interfaces.
func (c *Client) WaitForNetworkInterfacesFromVMAgent(
	ctx context.Context,
	timeout time.Duration,
	waitForIPConfig *WaitForIPConfig, // configuration for which IP types to wait for (nil = wait for any global unicast)
) (*GetQEMUNetworkInterfacesResponseData, error) {
	errNoIPsYet := errors.New("no ips yet")

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt)
	defer signal.Stop(ch)

	go func() {
		select {
		case <-ch:
			const msg = "interrupted by signal"
			tflog.Warn(ctx, msg)
			cancel()
		case <-ctxWithTimeout.Done():
		}
	}()

	op := retry.NewPollOperation("VM network interfaces",
		retry.WithRetryIf(func(err error) bool {
			return isAgentNotReadyError(err) || errors.Is(err, errNoIPsYet)
		}),
	)

	var result *GetQEMUNetworkInterfacesResponseData

	err := op.DoPoll(ctxWithTimeout, func() error {
		data, err := c.GetVMNetworkInterfacesFromAgent(ctx)
		if err != nil {
			var httpError *api.HTTPError
			if errors.As(err, &httpError) {
				if httpError.Code == http.StatusForbidden {
					return err
				}

				if isAgentNotReadyError(err) {
					return errNoIPsYet
				}
			}

			return errNoIPsYet
		}

		if data == nil || data.Result == nil {
			return errNoIPsYet
		}

		hasIPv4, hasIPv6 := c.checkIPAddresses(*data.Result)

		if waitForIPConfig == nil {
			if !hasIPv4 && !hasIPv6 {
				return errNoIPsYet
			}

			result = data

			return nil
		}

		requiredIPv4 := waitForIPConfig.IPv4
		requiredIPv6 := waitForIPConfig.IPv6

		if !requiredIPv4 && !requiredIPv6 {
			if !hasIPv4 && !hasIPv6 {
				return errNoIPsYet
			}

			result = data

			return nil
		}

		if (requiredIPv4 && !hasIPv4) || (requiredIPv6 && !hasIPv6) {
			return errNoIPsYet
		}

		result = data

		return nil
	})

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return nil, fmt.Errorf(
			"timeout while waiting for the QEMU agent on VM \"%d\" to publish the network interfaces",
			c.VMID,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("error waiting for VM network interfaces: %w", err)
	}

	return result, nil
}

// checkIPAddresses checks network interfaces for valid IP addresses and returns whether IPv4 and IPv6 are present.
func (c *Client) checkIPAddresses(
	nics []GetQEMUNetworkInterfacesResponseResult,
) (bool, bool) {
	hasIPv4 := false
	hasIPv6 := false

	for _, nic := range nics {
		if nic.Name == "lo" {
			continue
		}

		if nic.IPAddresses == nil || len(*nic.IPAddresses) == 0 {
			continue
		}

		for _, addr := range *nic.IPAddresses {
			if !ip.IsValidGlobalUnicast(addr.Address) {
				continue
			}

			if ip.IsIPv4(addr.Address) {
				hasIPv4 = true
			} else if ip.IsIPv6(addr.Address) {
				hasIPv6 = true
			}
		}
	}

	return hasIPv4, hasIPv6
}

// WaitForVMConfigUnlock waits for a virtual machine configuration to become unlocked.
func (c *Client) WaitForVMConfigUnlock(ctx context.Context, ignoreErrorResponse bool) error {
	stillLocked := errors.New("still locked")

	op := retry.NewPollOperation("VM config unlock",
		retry.WithRetryIf(func(err error) bool {
			return errors.Is(err, stillLocked) || ignoreErrorResponse
		}),
	)

	err := op.DoPoll(ctx, func() error {
		data, err := c.GetVMStatus(ctx)
		if err != nil {
			return err
		}

		if data.Lock != nil && *data.Lock != "" {
			return stillLocked
		}

		return nil
	})

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

	op := retry.NewPollOperation("VM status",
		retry.WithRetryIf(func(err error) bool {
			return errors.Is(err, unexpectedStatus)
		}),
	)

	err := op.DoPoll(ctx, func() error {
		data, err := c.GetVMStatus(ctx)
		if err != nil {
			return err
		}

		if data.Status != status {
			return unexpectedStatus
		}

		return nil
	})

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("timeout while waiting for VM %d to enter the status %q", c.VMID, status)
	}

	if err != nil {
		return fmt.Errorf("error waiting for VM %d to enter the status %q: %w", c.VMID, status, err)
	}

	return nil
}
