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

	"github.com/avast/retry-go/v4"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// CloneContainer clones a container.
func (c *Client) CloneContainer(ctx context.Context, d *CloneRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("/clone"), d, nil)
	if err != nil {
		return fmt.Errorf("error cloning container: %w", err)
	}

	return nil
}

// CreateContainer creates a container.
func (c *Client) CreateContainer(ctx context.Context, d *CreateRequestBody) error {
	taskID, err := c.CreateContainerAsync(ctx, d)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID)
	if err != nil {
		return fmt.Errorf("error waiting for container created: %w", err)
	}

	return nil
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
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(""), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting container: %w", err)
	}

	return nil
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

// WaitForContainerNetworkInterfaces waits for a container to publish its network interfaces.
func (c *Client) WaitForContainerNetworkInterfaces(
	ctx context.Context,
	timeout time.Duration,
) ([]GetNetworkInterfacesData, error) {
	errNoIPsYet := errors.New("no ips yet")

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ifaces, err := retry.DoWithData(
		func() ([]GetNetworkInterfacesData, error) {
			ifaces, err := c.GetContainerNetworkInterfaces(ctx)
			if err != nil {
				return nil, err
			}

			for _, iface := range ifaces {
				if iface.Name != "lo" && iface.IPAddresses != nil && len(*iface.IPAddresses) > 0 {
					// we have at least one non-loopback interface with an IP address
					return ifaces, nil
				}
			}

			return nil, errNoIPsYet
		},
		retry.Context(ctxWithTimeout),
		retry.RetryIf(func(err error) bool {
			var target *api.HTTPError
			if errors.As(err, &target) {
				if target.Code == http.StatusBadRequest {
					// this is a special case to account for eventual consistency
					// when creating a task -- the task may not be available via status API
					// immediately after creation
					return true
				}
			}

			return errors.Is(err, api.ErrNoDataObjectInResponse) || errors.Is(err, errNoIPsYet)
		}),
		retry.LastErrorOnly(true),
		retry.UntilSucceeded(),
		retry.DelayType(retry.FixedDelay),
		retry.Delay(time.Second),
	)
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, errors.New("timeout while waiting for container IP addresses")
	}

	if err != nil {
		return nil, fmt.Errorf("error while waiting for container IP addresses: %w", err)
	}

	return ifaces, nil
}

// RebootContainer reboots a container.
func (c *Client) RebootContainer(ctx context.Context, d *RebootRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("status/reboot"), d, nil)
	if err != nil {
		return fmt.Errorf("error rebooting container: %w", err)
	}

	return nil
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

	taskID, err := c.StartContainerAsync(ctx)
	if err != nil {
		return fmt.Errorf("error starting container: %w", err)
	}

	err = c.Tasks().WaitForTask(ctx, *taskID)
	if err != nil {
		return fmt.Errorf("error waiting for container start: %w", err)
	}

	// the timeout here should probably be configurable
	err = c.WaitForContainerStatus(ctx, "running")
	if err != nil {
		return fmt.Errorf("error waiting for container start: %w", err)
	}

	return nil
}

// StartContainerAsync starts a container asynchronously.
func (c *Client) StartContainerAsync(ctx context.Context) (*string, error) {
	resBody := &StartResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("status/start"), nil, resBody)
	if err != nil {
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
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath("config"), d, nil)
	if err != nil {
		return fmt.Errorf("error updating container: %w", err)
	}

	return nil
}

// WaitForContainerStatus waits for a container to reach a specific state.
func (c *Client) WaitForContainerStatus(ctx context.Context, status string) error {
	status = strings.ToLower(status)

	unexpectedStatus := fmt.Errorf("unexpected status %q", status)

	err := retry.Do(
		func() error {
			data, err := c.GetContainerStatus(ctx)
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

	err := retry.Do(
		func() error {
			data, err := c.GetContainerStatus(ctx)
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
		return fmt.Errorf("timeout while waiting for container %d configuration to become unlocked", c.VMID)
	}

	if err != nil && !ignoreErrorResponse {
		return fmt.Errorf("error waiting for container %d configuration to become unlocked: %w", c.VMID, err)
	}

	return nil
}
