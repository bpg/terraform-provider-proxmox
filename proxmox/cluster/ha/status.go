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

// WaitForServiceUnmanaged waits until the HA manager has dropped a resource from its service status, which it does
// once the resource is set to "ignored" or removed.
func (c *Client) WaitForServiceUnmanaged(ctx context.Context, id types.HAResourceID) error {
	errStillManaged := errors.New("still managed")

	op := retry.NewPollOperation("HA service unmanaged",
		retry.WithRetryIf(func(err error) bool {
			return errors.Is(err, errStillManaged)
		}),
	)

	err := op.DoPoll(ctx, func() error {
		status, err := c.GetManagerStatus(ctx)
		if err != nil {
			return err
		}

		if _, ok := status.ManagerStatus.ServiceStatus[id.String()]; ok {
			return errStillManaged
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error waiting for HA manager to release %v: %w", id, err)
	}

	return nil
}
