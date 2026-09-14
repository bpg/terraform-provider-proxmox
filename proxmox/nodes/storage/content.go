/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/avast/retry-go/v5"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/tasks"
	retryv2 "github.com/bpg/terraform-provider-proxmox/proxmox/retry"
)

// DeleteDatastoreFile deletes a file in a datastore, waiting for the async
// deletion task to complete.
func (c *Client) DeleteDatastoreFile(
	ctx context.Context,
	volumeID string,
) tasks.TaskResult {
	op := retryv2.NewTaskOperation("storage delete file",
		retryv2.WithRetryIf(func(err error) bool {
			return retryv2.IsTransientAPIError(err) && !errors.Is(err, api.ErrResourceDoesNotExist)
		}),
	)

	return c.Tasks().DoTask(ctx, op,
		func() (*string, error) { return c.DeleteDatastoreFileAsync(ctx, volumeID) },
	)
}

// DeleteDatastoreFileAsync deletes a file in a datastore asynchronously.
func (c *Client) DeleteDatastoreFileAsync(
	ctx context.Context,
	volumeID string,
) (*string, error) {
	resBody := &DeleteDatastoreFileResponseBody{}

	err := c.DoRequest(
		ctx,
		http.MethodDelete,
		c.ExpandPath(fmt.Sprintf("content/%s", url.PathEscape(volumeID))),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting file %s from datastore %s: %w", volumeID, c.StorageName, err)
	}

	// nil data means the delete completed synchronously (no task to wait for).
	return resBody.TaskID, nil
}

// ListDatastoreFiles retrieves a list of the files in a datastore.
// contentType is optional and filters the results by content type (e.g., "iso", "vztmpl", "backup").
// If contentType is nil, all files in the datastore are returned (unfiltered).
func (c *Client) ListDatastoreFiles(
	ctx context.Context,
	contentType *string,
) ([]*DatastoreFileListResponseData, error) {
	resBody := &DatastoreFileListResponseBody{}
	reqBody := &DatastoreFileListRequestBody{
		ContentType: contentType,
	}

	err := retry.New(
		retry.Context(ctx),
		retry.RetryIf(func(err error) bool {
			var httpError *api.HTTPError
			if errors.As(err, &httpError) && httpError.Code == http.StatusForbidden {
				return false
			}

			return !errors.Is(err, api.ErrResourceDoesNotExist)
		}),
		retry.LastErrorOnly(true),
	).Do(
		func() error {
			return c.DoRequest(ctx, http.MethodGet, c.ExpandPath("content"), reqBody, resBody)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error listing files from datastore %s: %w", c.StorageName, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].VolumeID < resBody.Data[j].VolumeID
	})

	return resBody.Data, nil
}

// GetDatastoreFile get a file details in a datastore.
func (c *Client) GetDatastoreFile(
	ctx context.Context,
	volumeID string,
) (*DatastoreFileGetResponseData, error) {
	resBody := &DatastoreFileGetResponseBody{}

	err := c.DoRequest(
		ctx,
		http.MethodGet,
		c.ExpandPath(
			fmt.Sprintf(
				"content/%s",
				url.PathEscape(volumeID),
			),
		),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error get file %s from datastore %s: %w", volumeID, c.StorageName, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}
