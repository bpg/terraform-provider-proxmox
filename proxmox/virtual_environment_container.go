/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// CloneContainer clones a container.
func (c *VirtualEnvironmentClient) CloneContainer(nodeName string, vmID int, d *VirtualEnvironmentContainerCloneRequestBody) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/lxc/%d/clone", url.PathEscape(nodeName), vmID), d, nil)
}

// CreateContainer creates a container.
func (c *VirtualEnvironmentClient) CreateContainer(nodeName string, d *VirtualEnvironmentContainerCreateRequestBody) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/lxc", url.PathEscape(nodeName)), d, nil)
}

// DeleteContainer deletes a container.
func (c *VirtualEnvironmentClient) DeleteContainer(nodeName string, vmID int) error {
	return c.DoRequest(hmDELETE, fmt.Sprintf("nodes/%s/lxc/%d", url.PathEscape(nodeName), vmID), nil, nil)
}

// GetContainer retrieves a container.
func (c *VirtualEnvironmentClient) GetContainer(nodeName string, vmID int) (*VirtualEnvironmentContainerGetResponseData, error) {
	resBody := &VirtualEnvironmentContainerGetResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("nodes/%s/lxc/%d/config", url.PathEscape(nodeName), vmID), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// GetContainerStatus retrieves the status for a container.
func (c *VirtualEnvironmentClient) GetContainerStatus(nodeName string, vmID int) (*VirtualEnvironmentContainerGetStatusResponseData, error) {
	resBody := &VirtualEnvironmentContainerGetStatusResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("nodes/%s/lxc/%d/status/current", url.PathEscape(nodeName), vmID), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// RebootContainer reboots a container.
func (c *VirtualEnvironmentClient) RebootContainer(nodeName string, vmID int, d *VirtualEnvironmentContainerRebootRequestBody) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/lxc/%d/status/reboot", url.PathEscape(nodeName), vmID), d, nil)
}

// ShutdownContainer shuts down a container.
func (c *VirtualEnvironmentClient) ShutdownContainer(nodeName string, vmID int, d *VirtualEnvironmentContainerShutdownRequestBody) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/lxc/%d/status/shutdown", url.PathEscape(nodeName), vmID), d, nil)
}

// StartContainer starts a container.
func (c *VirtualEnvironmentClient) StartContainer(nodeName string, vmID int) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/lxc/%d/status/start", url.PathEscape(nodeName), vmID), nil, nil)
}

// StopContainer stops a container immediately.
func (c *VirtualEnvironmentClient) StopContainer(nodeName string, vmID int) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/lxc/%d/status/stop", url.PathEscape(nodeName), vmID), nil, nil)
}

// UpdateContainer updates a container.
func (c *VirtualEnvironmentClient) UpdateContainer(nodeName string, vmID int, d *VirtualEnvironmentContainerUpdateRequestBody) error {
	return c.DoRequest(hmPUT, fmt.Sprintf("nodes/%s/lxc/%d/config", url.PathEscape(nodeName), vmID), d, nil)
}

// WaitForContainerState waits for a container to reach a specific state.
func (c *VirtualEnvironmentClient) WaitForContainerState(ctx context.Context, nodeName string, vmID int, state string, timeout int, delay int) error {
	state = strings.ToLower(state)

	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			data, err := c.GetContainerStatus(nodeName, vmID)

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

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	return fmt.Errorf("timeout while waiting for container \"%d\" to enter the state \"%s\"", vmID, state)
}

// WaitForContainerLock waits for a container lock to be released.
func (c *VirtualEnvironmentClient) WaitForContainerLock(ctx context.Context, nodeName string, vmID int, timeout int, delay int, ignoreErrorResponse bool) error {
	timeDelay := int64(delay)
	timeMax := float64(timeout)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			data, err := c.GetContainerStatus(nodeName, vmID)

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

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	return fmt.Errorf("timeout while waiting for container \"%d\" to become unlocked", vmID)
}
