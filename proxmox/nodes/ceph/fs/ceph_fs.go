/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
	"github.com/bpg/terraform-provider-proxmox/proxmox/retry"
)

// List returns every CephFS in the cluster.
func (c *Client) List(ctx context.Context) ([]*ListResponseData, error) {
	resBody := &ListResponseBody{}

	op := retry.NewAPICallOperation("CephFS list", retry.WithRetryIf(retry.IsTransientAPIError))

	if err := op.Do(ctx, func() error {
		return c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), nil, resBody)
	}); err != nil {
		return nil, fmt.Errorf("error listing CephFS: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// Get returns the CephFS with the given name, or api.ErrResourceDoesNotExist if the cluster doesn't list it.
func (c *Client) Get(ctx context.Context, name string) (*ListResponseData, error) {
	list, err := c.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, d := range list {
		if d != nil && d.Name == name {
			return d, nil
		}
	}

	return nil, fmt.Errorf("error getting CephFS %q: %w", name, api.ErrResourceDoesNotExist)
}

// Create creates a CephFS and its two pools, and waits for the task to complete. Proxmox refuses while no MDS is in
// standby, which a newly created MDS takes a moment to reach, so that refusal is retried for about a minute. Nothing
// else is: the request isn't idempotent.
func (c *Client) Create(ctx context.Context, name string, body *CreateRequestBody) tasks.TaskResult {
	op := retry.NewTaskOperation("CephFS create",
		retry.WithRetryIf(retry.ErrorContains("Metadata Server (MDS) found")),
		retry.WithAttempts(7),
		retry.WithBaseDelay(time.Second),
	)

	return c.Tasks().DoTask(ctx, op, func() (*string, error) {
		resBody := &CreateResponseBody{}

		if err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(url.PathEscape(name)), body, resBody); err != nil {
			return nil, fmt.Errorf("error creating CephFS %q: %w", name, err)
		}

		if resBody.Data == nil {
			return nil, api.ErrNoDataObjectInResponse
		}

		return resBody.Data, nil
	})
}
