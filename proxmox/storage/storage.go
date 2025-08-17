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

// ListDatastores retrieves a list of the cluster.
func (c *Client) ListDatastores(
	ctx context.Context,
	d *DatastoreListRequestBody,
) ([]*DatastoreListResponseData, error) {
	resBody := &DatastoreListResponseBody{}

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
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	return resBody.Data, nil
}

func (c *Client) CreateDatastore(
	ctx context.Context,
	d interface{},
) error {
	err := c.DoRequest(
		ctx,
		http.MethodPost,
		c.basePath(),
		d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error creating datastore: %w", err)
	}

	return nil
}

func (c *Client) UpdateDatastore(
	ctx context.Context,
	d interface{},
) error {

	err := c.DoRequest(
		ctx,
		http.MethodPost,
		c.ExpandPath(d.(DataStoreBase).Storage),
		d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error updating datastore: %w", err)
	}

	return nil
}

func (c *Client) DeleteDatastore(
	ctx context.Context,
	d interface{},
) error {

	err := c.DoRequest(
		ctx,
		http.MethodDelete,
		c.ExpandPath(d.(DataStoreBase).Storage),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error deleting datastore: %w", err)
	}

	return nil
}
