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

// alreadyExistsMessage is part of Proxmox's refusal to create a CephFS whose filesystem or pools already exist.
const alreadyExistsMessage = "already exist"

// isRetryableCreateError reports whether a failed create may be sent again. Proxmox raises its refusals as HTTP 500,
// which counts as transient, so an existing filesystem has to be excluded by name.
func isRetryableCreateError(err error) bool {
	return (retry.IsTransientAPIError(err) || retry.ErrorContains("Metadata Server (MDS) found")(err)) &&
		!retry.ErrorContains(alreadyExistsMessage)(err)
}

// Create creates a CephFS and its two pools, and waits for the task to complete. Proxmox refuses while no MDS is in
// standby, which a newly created MDS takes a moment to reach, so errors are retried for about a minute. Proxmox checks
// for an existing filesystem or pools before it starts the task: on a retry that means an earlier attempt went
// through, but on the first attempt it's a real conflict, so that error is never retried.
func (c *Client) Create(ctx context.Context, name string, body *CreateRequestBody) tasks.TaskResult {
	op := retry.NewTaskOperation("CephFS create",
		retry.WithRetryIf(isRetryableCreateError),
		retry.WithAlreadyDoneCheck(retry.ErrorContains(alreadyExistsMessage)),
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
