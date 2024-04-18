/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package containers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

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
func (c *Client) CreateContainer(ctx context.Context, d *CreateRequestBody, timeout time.Duration) error {
	taskID, err := c.CreateContainerAsync(ctx, d)
	if err != nil {
		return err
	}

	err = c.Tasks().WaitForTask(ctx, *taskID, timeout)
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
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("status/shutdown"), d, nil)
	if err != nil {
		return fmt.Errorf("error shutting down container: %w", err)
	}

	return nil
}

// StartContainer starts a container if is not already running.
func (c *Client) StartContainer(ctx context.Context, timeout time.Duration) error {
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

	err = c.Tasks().WaitForTask(ctx, *taskID, timeout)
	if err != nil {
		return fmt.Errorf("error waiting for container start: %w", err)
	}

	// the timeout here should probably be configurable
	err = c.WaitForContainerStatus(ctx, "running", int(timeout.Seconds())*2, 5)
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
func (c *Client) WaitForContainerStatus(ctx context.Context, status string, timeout int, delay int) error {
	status = strings.ToLower(status)

	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			data, err := c.GetContainerStatus(ctx)
			if err != nil {
				return fmt.Errorf("error retrieving container status: %w", err)
			}

			if data.Status == status {
				return nil
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(200 * time.Millisecond)

		timeElapsed = time.Since(timeStart)

		if ctx.Err() != nil {
			return fmt.Errorf("context error: %w", ctx.Err())
		}
	}

	return fmt.Errorf(
		"timeout while waiting for container \"%d\" to enter the status \"%s\"",
		c.VMID,
		status,
	)
}

// WaitForContainerLock waits for a container lock to be released.
func (c *Client) WaitForContainerLock(ctx context.Context, timeout int, delay int, ignoreErrorResponse bool) error {
	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			data, err := c.GetContainerStatus(ctx)

			if err != nil {
				if !ignoreErrorResponse {
					return fmt.Errorf("error retrieving container status: %w", err)
				}
			} else if data.Lock == nil || *data.Lock == "" {
				return nil
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(200 * time.Millisecond)

		timeElapsed = time.Since(timeStart)

		if ctx.Err() != nil {
			return fmt.Errorf("context error: %w", ctx.Err())
		}
	}

	return fmt.Errorf("timeout while waiting for container \"%d\" to become unlocked", c.VMID)
}
