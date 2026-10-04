/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ha

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/retry"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// ManagerStatusResponseBody contains the body from a HA manager status response.
type ManagerStatusResponseBody struct {
	Data *ManagerStatusResponseData `json:"data,omitempty"`
}

// ManagerStatusResponseData contains the data from a HA manager status response.
type ManagerStatusResponseData struct {
	ManagerStatus struct {
		Disarm        map[string]any `json:"disarm,omitempty"`
		ServiceStatus map[string]any `json:"service_status"`
	} `json:"manager_status"`
}

// GetManagerStatus retrieves the HA manager status.
func (c *Client) GetManagerStatus(ctx context.Context) (*ManagerStatusResponseData, error) {
	resBody := &ManagerStatusResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("status/manager_status"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading HA manager status: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// serviceReleaseTimeout bounds WaitForServiceUnmanaged. The CRM drops an ignored service within one cycle (~10s), so
// running out of it means the CRM is not processing, e.g. no active master.
const serviceReleaseTimeout = 2 * time.Minute

// WaitForServiceUnmanaged waits until the HA manager has dropped a resource from its service status, which it does
// once the resource is set to "ignored" or removed. Returns right away while the HA stack is disarmed: the CRM keeps
// the service status as is then, and the LRMs don't act on services.
func (c *Client) WaitForServiceUnmanaged(ctx context.Context, id types.HAResourceID) error {
	errStillManaged := errors.New("still managed")

	waitCtx, cancel := context.WithTimeout(ctx, serviceReleaseTimeout)
	defer cancel()

	op := retry.NewPollOperation("HA service unmanaged",
		retry.WithRetryIf(func(err error) bool {
			return errors.Is(err, errStillManaged)
		}),
	)

	err := op.DoPoll(waitCtx, func() error {
		status, err := c.GetManagerStatus(waitCtx)
		if err != nil {
			return err
		}

		if status.ManagerStatus.Disarm != nil {
			return nil
		}

		if _, ok := status.ManagerStatus.ServiceStatus[id.String()]; ok {
			return errStillManaged
		}

		return nil
	})

	switch {
	case err == nil:
		return nil
	case ctx.Err() == nil && errors.Is(waitCtx.Err(), context.DeadlineExceeded):
		return fmt.Errorf("HA manager did not release %v within %s, check that the HA CRM has an active master", id,
			serviceReleaseTimeout)
	default:
		return fmt.Errorf("error waiting for HA manager to release %v: %w", id, err)
	}
}
