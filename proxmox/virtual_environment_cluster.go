/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"context"
	"errors"
)

// GetClusterNextID retrieves the next free VM identifier for the cluster.
func (c *VirtualEnvironmentClient) GetClusterNextID(ctx context.Context, vmID *int) (*int, error) {
	reqBody := &VirtualEnvironmentClusterNextIDRequestBody{
		VMID: vmID,
	}

	resBody := &VirtualEnvironmentClusterNextIDResponseBody{}
	err := c.DoRequest(ctx, hmGET, "cluster/nextid", reqBody, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return (*int)(resBody.Data), nil
}
