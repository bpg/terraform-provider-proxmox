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

// ListDatastores retrieves a list of nodes.
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
