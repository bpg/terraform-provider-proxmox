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

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// GetDatastoreStatus gets status information for a given datastore.
func (c *Client) GetDatastoreStatus(
	ctx context.Context,
) (*DatastoreGetStatusResponseData, error) {
	resBody := &DatastoreGetStatusResponseBody{}

	err := c.DoRequest(
		ctx,
		http.MethodGet,
		c.ExpandPath("status"),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving status for datastore %s: %w", c.StorageName, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}
