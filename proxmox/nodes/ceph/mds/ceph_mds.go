/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package mds

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
	"github.com/bpg/terraform-provider-proxmox/proxmox/retry"
)

// confLock serializes MDS creation and removal. Both rewrite the cluster-wide ceph.conf without a lock, so two
// concurrent tasks can drop each other's [mds.<id>] section. Global variable by design.
//
//nolint:gochecknoglobals
var confLock sync.Mutex

// List returns every MDS in the Ceph cluster, not only those on this node.
func (c *Client) List(ctx context.Context) ([]*ListResponseData, error) {
	resBody := &ListResponseBody{}

	op := retry.NewAPICallOperation("Ceph MDS list", retry.WithRetryIf(retry.IsTransientAPIError))

	if err := op.Do(ctx, func() error {
		return c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), nil, resBody)
	}); err != nil {
		return nil, fmt.Errorf("error listing Ceph MDS: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// Get returns the MDS with the given ID, or api.ErrResourceDoesNotExist if the cluster doesn't list it.
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

	return nil, fmt.Errorf("error getting Ceph MDS %q: %w", name, api.ErrResourceDoesNotExist)
}

// Create creates an MDS on this node and waits for the task to complete.
func (c *Client) Create(ctx context.Context, name string) tasks.TaskResult {
	confLock.Lock()
	defer confLock.Unlock()

	op := retry.NewTaskOperation("Ceph MDS create", retry.WithRetryIf(retry.IsTransientAPIError))

	return c.Tasks().DoTask(ctx, op, func() (*string, error) {
		resBody := &CreateResponseBody{}

		if err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(url.PathEscape(name)), nil, resBody); err != nil {
			return nil, fmt.Errorf("error creating Ceph MDS %q: %w", name, err)
		}

		if resBody.Data == nil {
			return nil, api.ErrNoDataObjectInResponse
		}

		return resBody.Data, nil
	})
}

// Delete destroys an MDS on this node and waits for the task to complete.
func (c *Client) Delete(ctx context.Context, name string) tasks.TaskResult {
	confLock.Lock()
	defer confLock.Unlock()

	op := retry.NewTaskOperation("Ceph MDS delete",
		retry.WithRetryIf(func(err error) bool {
			return retry.IsTransientAPIError(err) && !errors.Is(err, api.ErrResourceDoesNotExist)
		}),
	)

	return c.Tasks().DoTask(ctx, op, func() (*string, error) {
		resBody := &DeleteResponseBody{}

		if err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(url.PathEscape(name)), nil, resBody); err != nil {
			return nil, fmt.Errorf("error deleting Ceph MDS %q: %w", name, err)
		}

		if resBody.Data == nil {
			return nil, api.ErrNoDataObjectInResponse
		}

		return resBody.Data, nil
	})
}
