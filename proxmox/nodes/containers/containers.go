/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package containers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/retry"
	"github.com/bpg/terraform-provider-proxmox/utils/ip"
)

var errContainerAlreadyRunning = errors.New("container is already running")

// CloneContainer clones a container.
func (c *Client) CloneContainer(ctx context.Context, d *CloneRequestBody) error {
	op := retry.NewTaskOperation("container clone",
		retry.WithBaseDelay(10*time.Second),
		retry.WithRetryIf(retry.IsTransientAPIError),
		retry.WithAlreadyDoneCheck(retry.ErrorContains("already exists")),
	)

	return op.DoTask(ctx,
		func() (*string, error) { return c.CloneContainerAsync(ctx, d) },
		func(ctx context.Context, taskID string) error { return c.Tasks().WaitForTask(ctx, taskID) },
	)
}

// CloneContainerAsync clones a container asynchronously.
func (c *Client) CloneContainerAsync(ctx context.Context, d *CloneRequestBody) (*string, error) {
	resBody := &CloneResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("clone"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error cloning container: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// CreateContainer creates a container.
func (c *Client) CreateContainer(ctx context.Context, d *CreateRequestBody) error {
	op := retry.NewTaskOperation("container create",
		retry.WithRetryIf(retry.IsTransientAPIError),
		retry.WithAlreadyDoneCheck(retry.ErrorContains("already exists")),
	)

	return op.DoTask(ctx,
		func() (*string, error) { return c.CreateContainerAsync(ctx, d) },
		func(ctx context.Context, taskID string) error { return c.Tasks().WaitForTask(ctx, taskID) },
	)
}

// CreateContainerAsync creates a container asynchronously.
func (c *Client) CreateContainerAsync(ctx context.Context, d *CreateRequestBody) (*string, error) {
	resBody := &CreateResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.basePath(), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error creating container: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// DeleteContainer deletes a container.
func (c *Client) DeleteContainer(ctx context.Context) error {
	op := retry.NewTaskOperation("container delete",
		retry.WithRetryIf(func(err error) bool {
			return retry.IsTransientAPIError(err) && !errors.Is(err, api.ErrResourceDoesNotExist)
		}),
	)

	return op.DoTask(ctx,
		func() (*string, error) { return c.DeleteContainerAsync(ctx) },
		func(ctx context.Context, taskID string) error { return c.Tasks().WaitForTask(ctx, taskID) },
	)
}

// DeleteContainerAsync deletes a container asynchronously.
func (c *Client) DeleteContainerAsync(ctx context.Context) (*string, error) {
	resBody := &DeleteResponseBody{}

	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(""), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error deleting container: %w", err)
	}

	// nil data means the delete completed synchronously (no task to wait for).
	return resBody.Data, nil
}

// GetContainer retrieves a container.
func (c *Client) GetContainer(ctx context.Context) (*GetResponseData, error) {
	resBody := &GetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("config"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving container: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetContainerStatus retrieves the status for a container.
func (c *Client) GetContainerStatus(ctx context.Context) (*GetStatusResponseData, error) {
	resBody := &GetStatusResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("status/current"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving container status: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetContainerNetworkInterfaces retrieves details about the container network interfaces.
func (c *Client) GetContainerNetworkInterfaces(ctx context.Context) ([]GetNetworkInterfacesData, error) {
	resBody := &GetNetworkInterfaceResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("interfaces"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving container network interfaces: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// ListContainers retrieves a list of containers.
func (c *Client) ListContainers(ctx context.Context) ([]*ListResponseData, error) {
	resBody := &ListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.basePath(), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Containers: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// WaitForContainerNetworkInterfaces waits for a container to publish its network interfaces.
func (c *Client) WaitForContainerNetworkInterfaces(
	ctx context.Context,
	timeout time.Duration,
	waitForIPConfig *WaitForIPConfig, // configuration for which IP types to wait for (nil = wait for any global unicast)
) ([]GetNetworkInterfacesData, error) {
	errNoIPsYet := errors.New("no ips yet")

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	op := retry.NewPollOperation("container network interfaces",
		retry.WithRetryIf(func(err error) bool {
			var target *api.HTTPError
			if errors.As(err, &target) {
				if target.Code == http.StatusBadRequest {
					return true
				}
			}

			return errors.Is(err, api.ErrNoDataObjectInResponse) || errors.Is(err, errNoIPsYet)
		}),
	)

	var ifaces []GetNetworkInterfacesData

	err := op.DoPoll(ctxWithTimeout, func() error {
		var err error

		ifaces, err = c.GetContainerNetworkInterfaces(ctxWithTimeout)
		if err != nil {
			return err
		}

		hasIPv4, hasIPv6 := c.checkIPAddresses(ifaces)

		if waitForIPConfig == nil {
			if !hasIPv4 && !hasIPv6 {
				return errNoIPsYet
			}

			return nil
		}

		requiredIPv4 := waitForIPConfig.IPv4
		requiredIPv6 := waitForIPConfig.IPv6

		if !requiredIPv4 && !requiredIPv6 {
			if !hasIPv4 && !hasIPv6 {
				return errNoIPsYet
			}

			return nil
		}

		if (requiredIPv4 && !hasIPv4) || (requiredIPv6 && !hasIPv6) {
			return errNoIPsYet
		}

		return nil
	})

	if errors.Is(err, context.DeadlineExceeded) {
		return nil, errors.New("timeout while waiting for container IP addresses")
	}

	if err != nil {
		return nil, fmt.Errorf("error while waiting for container IP addresses: %w", err)
	}

	return ifaces, nil
}

// checkIPAddresses checks network interfaces for valid IP addresses and returns whether IPv4 and IPv6 are present.
func (c *Client) checkIPAddresses(
	ifaces []GetNetworkInterfacesData,
) (bool, bool) {
	hasIPv4 := false
	hasIPv6 := false

	for _, iface := range ifaces {
		if iface.Name == "lo" || iface.IPAddresses == nil || len(*iface.IPAddresses) == 0 {
			continue
		}

		for _, ipAddr := range *iface.IPAddresses {
			if !ip.IsValidGlobalUnicast(ipAddr.Address) {
				continue
			}

			if ip.IsIPv4(ipAddr.Address) {
				hasIPv4 = true
			} else if ip.IsIPv6(ipAddr.Address) {
				hasIPv6 = true
			}
		}
	}

	return hasIPv4, hasIPv6
}

// RebootContainer reboots a container.
func (c *Client) RebootContainer(ctx context.Context, d *RebootRequestBody) error {
	taskID, err := c.RebootContainerAsync(ctx, d)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID)
	if err != nil {
		return fmt.Errorf("error waiting for container reboot: %w", err)
	}

	return nil
}

// RebootContainerAsync reboots a container asynchronously.
func (c *Client) RebootContainerAsync(ctx context.Context, d *RebootRequestBody) (*string, error) {
	resBody := &RebootResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("status/reboot"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error rebooting container: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// ShutdownContainer shuts down a container.
func (c *Client) ShutdownContainer(ctx context.Context, d *ShutdownRequestBody) error {
	taskID, err := c.ShutdownContainerAsync(ctx, d)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID)
	if err != nil {
		return fmt.Errorf("error waiting for container shut down: %w", err)
	}

	return nil
}

// ShutdownContainerAsync shuts down a container asynchronously.
func (c *Client) ShutdownContainerAsync(ctx context.Context, d *ShutdownRequestBody) (*string, error) {
	resBody := &ShutdownResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("status/shutdown"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error shutting down container: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// StartContainer starts a container if is not already running.
func (c *Client) StartContainer(ctx context.Context) error {
	status, err := c.GetContainerStatus(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving container status: %w", err)
	}

	if status.Status == "running" {
		return nil
	}

	op := retry.NewTaskOperation("container start",
		retry.WithRetryIf(retry.ErrorContains("got no worker upid")),
	)

	if err := op.DoTask(ctx,
		func() (*string, error) { return c.StartContainerAsync(ctx) },
		func(ctx context.Context, taskID string) error { return c.Tasks().WaitForTask(ctx, taskID) },
	); err != nil {
		if errors.Is(err, errContainerAlreadyRunning) {
			return nil
		}

		return err
	}

	return c.WaitForContainerStatus(ctx, "running")
}

// StartContainerAsync starts a container asynchronously.
// Returns errContainerAlreadyRunning if the container is already running.
func (c *Client) StartContainerAsync(ctx context.Context) (*string, error) {
	resBody := &StartResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("status/start"), nil, resBody)
	if err != nil {
		if strings.Contains(err.Error(), "already running") {
			return nil, errContainerAlreadyRunning
		}

		return nil, fmt.Errorf("error starting container: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// StopContainer stops a container immediately.
func (c *Client) StopContainer(ctx context.Context) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("status/stop"), nil, nil)
	if err != nil {
		return fmt.Errorf("error stopping container: %w", err)
	}

	return nil
}

// UpdateContainer updates a container.
func (c *Client) UpdateContainer(ctx context.Context, d *UpdateRequestBody) error {
	op := retry.NewAPICallOperation("container config update",
		retry.WithRetryIf(retry.ErrorContains("got timeout")),
	)

	return op.Do(ctx, func() error {
		return c.DoRequest(ctx, http.MethodPut, c.ExpandPath("config"), d, nil)
	})
}

// WaitForContainerStatus waits for a container to reach a specific state.
func (c *Client) WaitForContainerStatus(ctx context.Context, status string) error {
	status = strings.ToLower(status)
	unexpectedStatus := fmt.Errorf("unexpected status %q", status)

	op := retry.NewPollOperation("container status",
		retry.WithRetryIf(func(err error) bool {
			return errors.Is(err, unexpectedStatus)
		}),
	)

	err := op.DoPoll(ctx, func() error {
		data, err := c.GetContainerStatus(ctx)
		if err != nil {
			return err
		}

		if data.Status != status {
			return unexpectedStatus
		}

		return nil
	})

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("timeout while waiting for container %d to enter the status %q", c.VMID, status)
	}

	if err != nil {
		return fmt.Errorf("error waiting for container %d to enter the status %q: %w", c.VMID, status, err)
	}

	return nil
}

// WaitForContainerConfigUnlock waits for a container lock to be released.
func (c *Client) WaitForContainerConfigUnlock(ctx context.Context, ignoreErrorResponse bool) error {
	stillLocked := errors.New("still locked")

	op := retry.NewPollOperation("container config unlock",
		retry.WithRetryIf(func(err error) bool {
			return errors.Is(err, stillLocked) || ignoreErrorResponse
		}),
	)

	err := op.DoPoll(ctx, func() error {
		data, err := c.GetContainerStatus(ctx)
		if err != nil {
			return err
		}

		if data.Lock != nil && *data.Lock != "" {
			return stillLocked
		}

		return nil
	})

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("timeout while waiting for container %d configuration to become unlocked", c.VMID)
	}

	if err != nil && !ignoreErrorResponse {
		return fmt.Errorf("error waiting for container %d configuration to become unlocked: %w", c.VMID, err)
	}

	return nil
}

// ResizeContainerDisk resizes a container disk.
func (c *Client) ResizeContainerDisk(ctx context.Context, d *ResizeRequestBody) error {
	op := retry.NewTaskOperation("container disk resize",
		retry.WithRetryIf(retry.IsTransientAPIError),
	)

	return op.DoTask(ctx,
		func() (*string, error) { return c.ResizeContainerDiskAsync(ctx, d) },
		func(ctx context.Context, taskID string) error { return c.Tasks().WaitForTask(ctx, taskID) },
	)
}

// ResizeContainerDiskAsync resizes a container disk asynchronously.
func (c *Client) ResizeContainerDiskAsync(ctx context.Context, d *ResizeRequestBody) (*string, error) {
	resBody := &ResizeResponseBody{}

	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath("resize"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error resizing container disk: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}
