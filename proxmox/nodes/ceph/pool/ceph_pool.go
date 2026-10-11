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
	"strings"

	"github.com/google/go-querystring/query"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
	"github.com/bpg/terraform-provider-proxmox/proxmox/retry"
)

// Get returns the full settings of a single Ceph pool via /status?verbose=1, or
// api.ErrResourceDoesNotExist if absent. /status returns all settable settings
// (including pg_num_min) plus application_list (verbose=1 only) and crush_rule
// as a name string — fields the list endpoint omits or returns in less usable forms.
//
// PVE surfaces a missing pool as HTTP 500 with body
// `"error with 'osd pool get': mon_cmd failed - unrecognized pool 'X'"`
// (the underlying Ceph mon error), which the shared client does not map to
// api.ErrResourceDoesNotExist on its own — so we translate it here.
func (c *Client) Get(ctx context.Context, name string) (*StatusResponseData, error) {
	resBody := &StatusResponseBody{}
	path := c.ExpandPath(url.PathEscape(name)) + "/status?verbose=1"

	op := retry.NewAPICallOperation("Ceph pool status",
		retry.WithRetryIf(func(err error) bool {
			return retry.IsTransientAPIError(err) && !isUnrecognizedPoolErr(err)
		}),
	)

	if err := op.Do(ctx, func() error {
		return c.DoRequest(ctx, http.MethodGet, path, nil, resBody)
	}); err != nil {
		if isUnrecognizedPoolErr(err) {
			return nil, errors.Join(api.ErrResourceDoesNotExist, err)
		}

		return nil, fmt.Errorf("error getting Ceph pool %q: %w", name, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// isUnrecognizedPoolErr matches the Ceph mon "unrecognized pool" message that
// PVE surfaces (via HTTP 500) when a pool does not exist.
func isUnrecognizedPoolErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "unrecognized pool")
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
