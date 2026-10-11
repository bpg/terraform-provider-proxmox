/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zfs

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

// Create provisions a new ZFS pool and waits for the task to complete.
func (c *Client) Create(ctx context.Context, body *CreateRequestBody) tasks.TaskResult {
	op := retry.NewTaskOperation("ZFS pool create",
		retry.WithRetryIf(retry.IsTransientAPIError),
	)

	return c.Tasks().DoTask(ctx, op, func() (*string, error) {
		resBody := &CreateResponseBody{}

		if err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), body, resBody); err != nil {
			return nil, fmt.Errorf("error creating ZFS pool: %w", err)
		}

		if resBody.Data == nil {
			return nil, api.ErrNoDataObjectInResponse
		}

		return resBody.Data, nil
	})
}

// Get returns the details of a named ZFS pool, or api.ErrResourceDoesNotExist if absent.
func (c *Client) Get(ctx context.Context, name string) (*GetResponseData, error) {
	resBody := &GetResponseBody{}

	if err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(url.PathEscape(name)), nil, resBody); err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) || isPoolNotFoundErr(err) {
			return nil, errors.Join(api.ErrResourceDoesNotExist, err)
		}

		return nil, fmt.Errorf("error getting ZFS pool %q: %w", name, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// List returns all ZFS pools present on the node.
func (c *Client) List(ctx context.Context) ([]*ListResponseData, error) {
	resBody := &ListResponseBody{}

	if err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), nil, resBody); err != nil {
		return nil, fmt.Errorf("error listing ZFS pools: %w", err)
	}

	return resBody.Data, nil
}

// Delete destroys a ZFS pool and waits for the task to complete.
func (c *Client) Delete(ctx context.Context, name string, params *DeleteRequestParams) tasks.TaskResult {
	op := retry.NewTaskOperation("ZFS pool delete",
		retry.WithRetryIf(func(err error) bool {
			return retry.IsTransientAPIError(err) && !errors.Is(err, api.ErrResourceDoesNotExist)
		}),
	)

	return c.Tasks().DoTask(ctx, op, func() (*string, error) {
		path := c.ExpandPath(url.PathEscape(name))

		if params != nil {
			values, err := query.Values(params)
			if err != nil {
				return nil, fmt.Errorf("error encoding ZFS pool delete params: %w", err)
			}

			if encoded := values.Encode(); encoded != "" {
				path = path + "?" + encoded
			}
		}

		resBody := &DeleteResponseBody{}

		if err := c.DoRequest(ctx, http.MethodDelete, path, nil, resBody); err != nil {
			return nil, fmt.Errorf("error deleting ZFS pool %q: %w", name, err)
		}

		if resBody.Data == nil {
			return nil, api.ErrNoDataObjectInResponse
		}

		return resBody.Data, nil
	})
}

// isPoolNotFoundErr checks for PVE's ZFS "no such pool" HTTP 500 error, which the
// generic API client does not map to ErrResourceDoesNotExist on its own.
func isPoolNotFoundErr(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "no such pool")
}
