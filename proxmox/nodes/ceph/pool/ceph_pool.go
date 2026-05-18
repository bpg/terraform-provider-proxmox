/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pool

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
	"github.com/bpg/terraform-provider-proxmox/proxmox/retry"
)

// List returns all Ceph pools on the node along with their full settings.
func (c *Client) List(ctx context.Context) ([]*ListResponseData, error) {
	resBody := &ListResponseBody{}

	if err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), nil, resBody); err != nil {
		return nil, fmt.Errorf("error listing Ceph pools: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// Get returns a single Ceph pool by name, or api.ErrResourceDoesNotExist if absent.
// The Proxmox GET-by-name endpoint only returns a sub-index, so we filter the list response.
func (c *Client) Get(ctx context.Context, name string) (*ListResponseData, error) {
	pools, err := c.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range pools {
		if p.PoolName == name {
			return p, nil
		}
	}

	return nil, api.ErrResourceDoesNotExist
}

// Create creates a new Ceph pool and waits for the dispatch task to complete.
func (c *Client) Create(ctx context.Context, body *CreateRequestBody) tasks.TaskResult {
	op := retry.NewTaskOperation("Ceph pool create",
		retry.WithRetryIf(retry.IsTransientAPIError),
		retry.WithAlreadyDoneCheck(retry.ErrorContains("already exists")),
	)

	return c.Tasks().DoTask(ctx, op, func() (*string, error) {
		resBody := &CreateResponseBody{}

		if err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), body, resBody); err != nil {
			return nil, fmt.Errorf("error creating Ceph pool: %w", err)
		}

		if resBody.Data == nil {
			return nil, api.ErrNoDataObjectInResponse
		}

		return resBody.Data, nil
	})
}

// Update modifies an existing Ceph pool and waits for the dispatch task to complete.
func (c *Client) Update(ctx context.Context, name string, body *UpdateRequestBody) tasks.TaskResult {
	op := retry.NewTaskOperation("Ceph pool update",
		retry.WithRetryIf(retry.IsTransientAPIError),
	)

	return c.Tasks().DoTask(ctx, op, func() (*string, error) {
		resBody := &UpdateResponseBody{}

		if err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(url.PathEscape(name)), body, resBody); err != nil {
			return nil, fmt.Errorf("error updating Ceph pool %q: %w", name, err)
		}

		if resBody.Data == nil {
			return nil, api.ErrNoDataObjectInResponse
		}

		return resBody.Data, nil
	})
}

// Delete destroys a Ceph pool and waits for the dispatch task to complete.
// params may be nil to use Proxmox defaults.
func (c *Client) Delete(ctx context.Context, name string, params *DeleteRequestParams) tasks.TaskResult {
	op := retry.NewTaskOperation("Ceph pool delete",
		retry.WithRetryIf(func(err error) bool {
			return retry.IsTransientAPIError(err) && !errors.Is(err, api.ErrResourceDoesNotExist)
		}),
	)

	return c.Tasks().DoTask(ctx, op, func() (*string, error) {
		path := c.ExpandPath(url.PathEscape(name))

		if params != nil {
			values, err := query.Values(params)
			if err != nil {
				return nil, fmt.Errorf("error encoding Ceph pool delete params: %w", err)
			}

			if encoded := values.Encode(); encoded != "" {
				path = path + "?" + encoded
			}
		}

		resBody := &DeleteResponseBody{}

		if err := c.DoRequest(ctx, http.MethodDelete, path, nil, resBody); err != nil {
			return nil, fmt.Errorf("error deleting Ceph pool %q: %w", name, err)
		}

		if resBody.Data == nil {
			return nil, api.ErrNoDataObjectInResponse
		}

		return resBody.Data, nil
	})
}
