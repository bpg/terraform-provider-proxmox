/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tasks

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// GetTaskStatus retrieves the status of a task.
func (c *Client) GetTaskStatus(ctx context.Context, upid string) (*GetTaskStatusResponseData, error) {
	resBody := &GetTaskStatusResponseBody{}

	err := c.DoRequest(
		ctx,
		http.MethodGet,
		c.ExpandPath(fmt.Sprintf("%s/status", url.PathEscape(upid))),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrievinf task status: %w", err)
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

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			status, err := c.GetTaskStatus(ctx, upid)
			if err != nil {
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
