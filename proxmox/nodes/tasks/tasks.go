/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tasks

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// GetTaskStatus retrieves the status of a task.
func (c *Client) GetTaskStatus(ctx context.Context, upid string) (*GetTaskStatusResponseData, error) {
	resBody := &GetTaskStatusResponseBody{}

	path, err := c.BuildPath(upid, "status")
	if err != nil {
		return nil, fmt.Errorf("error building path for task status: %w", err)
	}

	err = c.DoRequest(ctx, http.MethodGet, path, nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving task status: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetTaskLog retrieves the log of a task. The log is returned as an array of
// lines. Each line is an object with a line number and the text of the line.
// Reads first 50 lines by default.
func (c *Client) GetTaskLog(ctx context.Context, upid string) ([]string, error) {
	resBody := &GetTaskLogResponseBody{}
	lines := []string{}

	path, err := c.BuildPath(upid, "log")
	if err != nil {
		return lines, fmt.Errorf("error building path for task status: %w", err)
	}

	err = c.DoRequest(ctx, http.MethodGet, path, nil, resBody)
	if err != nil {
		return lines, fmt.Errorf("error retrieving task status: %w", err)
	}

	if resBody.Data == nil {
		return lines, api.ErrNoDataObjectInResponse
	}

	for _, line := range resBody.Data {
		lines = append(lines, line.LineText)
	}

	return lines, nil
}

// DeleteTask deletes specific task.
func (c *Client) DeleteTask(ctx context.Context, upid string) error {
	path, err := c.baseTaskPath(upid)
	if err != nil {
		return fmt.Errorf("error creating task path: %w", err)
	}

	err = c.DoRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			return nil
		}

		return fmt.Errorf("error deleting task: %w", err)
	}

	return nil
}

// WaitForTask waits for a specific task to complete.
func (c *Client) WaitForTask(ctx context.Context, upid string) error {
	errStillRunning := errors.New("still running")

	status, err := retry.DoWithData(
		func() (*GetTaskStatusResponseData, error) {
			status, err := c.GetTaskStatus(ctx, upid)
			if err != nil {
				return nil, err
			}

			if status.Status == "running" {
				return nil, errStillRunning
			}

			return status, err
		},
		retry.Context(ctx),
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

			return errors.Is(err, errStillRunning)
		}),
		retry.LastErrorOnly(true),
		retry.Attempts(0), // retry until context deadline
		retry.DelayType(retry.FixedDelay),
		retry.Delay(time.Second),
	)

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("timeout while waiting for task %q to complete", upid)
	}

	if err != nil {
		return fmt.Errorf("error while waiting for task %q to complete: %w", upid, err)
	}

	if status.ExitCode != "OK" {
		return fmt.Errorf("task %q failed to complete with exit code: %s", upid, status.ExitCode)
	}

	return nil
}
