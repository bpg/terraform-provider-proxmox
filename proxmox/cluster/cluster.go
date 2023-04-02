/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// GetNextID retrieves the next free VM identifier for the cluster.
func (c *Client) GetNextID(ctx context.Context, vmID *int) (*int, error) {
	reqBody := &NextIDRequestBody{
		VMID: vmID,
	}

	resBody := &NextIDResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, "cluster/nextid", reqBody, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving next VM ID: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return (*int)(resBody.Data), nil
}
