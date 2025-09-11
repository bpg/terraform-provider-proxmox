/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// ListDatastore retrieves a list of the cluster.
func (c *Client) ListDatastore(ctx context.Context, d *DatastoreListRequest) ([]*DatastoreGetResponseData, error) {
	resBody := &DatastoreListResponse{}

	err := c.DoRequest(
		ctx,
		http.MethodGet,
		c.basePath(),
		d,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving datastores: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return *(resBody.Data[i]).ID < *(resBody.Data[j]).ID
	})

	return resBody.Data, nil
}

func (c *Client) GetDatastore(ctx context.Context, d *DatastoreGetRequest) (*DatastoreGetResponseData, error) {
	resBody := &DatastoreGetResponse{}

	err := c.DoRequest(
		ctx,
		http.MethodGet,
		c.ExpandPath(*d.ID),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error reading datastore: %w", err)
	}

	return resBody.Data, nil
}

func (c *Client) CreateDatastore(ctx context.Context, d interface{}) (*DatastoreCreateResponseData, error) {
	resBody := &DatastoreCreateResponse{}

	err := c.DoRequest(
		ctx,
		http.MethodPost,
		c.basePath(),
		d,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating datastore: %w", err)
	}

	return resBody.Data, nil
}

func (c *Client) UpdateDatastore(ctx context.Context, storeID string, d interface{}) error {
	err := c.DoRequest(
		ctx,
		http.MethodPut,
		c.ExpandPath(storeID),
		d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error updating datastore: %w", err)
	}

	return nil
}

func (c *Client) DeleteDatastore(ctx context.Context, storeID string) error {
	err := c.DoRequest(
		ctx,
		http.MethodDelete,
		c.ExpandPath(storeID),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error deleting datastore: %w", err)
	}

	return nil
}
