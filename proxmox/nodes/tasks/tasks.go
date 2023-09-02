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

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// GetTaskStatus retrieves the status of a task.
func (c *Client) GetTaskStatus(ctx context.Context, upid string) (*GetTaskStatusResponseData, error) {
	resBody := &GetTaskStatusResponseBody{}

	path, err := c.BuildPath(upid, "status")
	if err != nil {
		return nil, fmt.Errorf("error building path for task status: %w", err)
	}

	err = c.DoRequest(
		ctx,
		http.MethodGet,
		path,
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving task status: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// WaitForTask waits for a specific task to complete.
func (c *Client) WaitForTask(ctx context.Context, upid string, timeoutSec, delaySec int) error {
	timeDelay := int64(delaySec)
	timeMax := float64(timeoutSec)
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	isCriticalError := func(err error) bool {
		var target *api.HTTPError
		if errors.As(err, &target) {
			if target.Code != http.StatusBadRequest {
				// this is a special case to account for eventual consistency
				// when creating a task -- the task may not be available via status API
				// immediately after creation
				return true
			}
		}

		return err != nil
	}

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			status, err := c.GetTaskStatus(ctx, upid)
			if isCriticalError(err) {
				return err
			}

			if status.Status != "running" {
				if status.ExitCode != "OK" {
					return fmt.Errorf(
						"task \"%s\" failed to complete with exit code: %s",
						upid,
						status.ExitCode,
					)
				}

				return nil
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(200 * time.Millisecond)

		timeElapsed = time.Since(timeStart)

		if ctx.Err() != nil {
			return fmt.Errorf(
				"context error while waiting for task \"%s\" to complete: %w",
				upid, ctx.Err(),
			)
		}
	}

	return fmt.Errorf(
		"timeout while waiting for task \"%s\" to complete",
		upid,
	)
}
