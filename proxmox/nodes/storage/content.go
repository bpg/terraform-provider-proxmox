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

	"github.com/avast/retry-go/v4"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// DeleteDatastoreFile deletes a file in a datastore.
func (c *Client) DeleteDatastoreFile(
	ctx context.Context,
	volumeID string,
) error {
	path := c.ExpandPath(fmt.Sprintf("content/%s", url.PathEscape(volumeID)))

	err := retry.Do(
		func() error {
			return c.DoRequest(ctx, http.MethodDelete, path, nil, nil)
		},
		retry.Context(ctx),
		retry.RetryIf(func(err error) bool {
			return !errors.Is(err, api.ErrResourceDoesNotExist)
		}),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		return fmt.Errorf("error deleting file %s from datastore %s: %w", volumeID, c.StorageName, err)
	}

	return nil
}

// ListDatastoreFiles retrieves a list of the files in a datastore.
func (c *Client) ListDatastoreFiles(
	ctx context.Context,
) ([]*DatastoreFileListResponseData, error) {
	resBody := &DatastoreFileListResponseBody{}

	err := retry.Do(
		func() error {
			return c.DoRequest(ctx, http.MethodGet, c.ExpandPath("content"), nil, resBody)
		},
		retry.Context(ctx),
		retry.RetryIf(func(err error) bool {
			return !errors.Is(err, api.ErrResourceDoesNotExist)
		}),
		retry.LastErrorOnly(true),
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
