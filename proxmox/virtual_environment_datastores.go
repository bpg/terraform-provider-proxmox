/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"fmt"
	"sort"
)

// VirtualEnvironmentDatastoreListRequestBody contains the body for a datastore list request.
type VirtualEnvironmentDatastoreListRequestBody struct {
	ContentTypes CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Enabled      *CustomBool              `json:"enabled,omitempty" url:"enabled,omitempty,int"`
	Format       *CustomBool              `json:"format,omitempty" url:"format,omitempty,int"`
	ID           *string                  `json:"storage,omitempty" url:"storage,omitempty"`
	Target       *string                  `json:"target,omitempty" url:"target,omitempty"`
}

// VirtualEnvironmentDatastoreListResponseBody contains the body from a datastore list response.
type VirtualEnvironmentDatastoreListResponseBody struct {
	Data []*VirtualEnvironmentDatastoreListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentDatastoreListResponseData contains the data from a node list response.
type VirtualEnvironmentDatastoreListResponseData struct {
	Active              *CustomBool               `json:"active,omitempty"`
	ContentTypes        *CustomCommaSeparatedList `json:"content,omitempty"`
	Enabled             *CustomBool               `json:"enabled,omitempty"`
	ID                  string                    `json:"storage,omitempty"`
	Shared              *CustomBool               `json:"shared,omitempty"`
	SpaceAvailable      *int                      `json:"avail,omitempty"`
	SpaceTotal          *int                      `json:"total,omitempty"`
	SpaceUsed           *int                      `json:"used,omitempty"`
	SpaceUsedPercentage *float64                  `json:"used_fraction,omitempty"`
	Type                string                    `json:"type,omitempty"`
}

// ListDatastores retrieves a list of nodes.
func (c *VirtualEnvironmentClient) ListDatastores(nodeName string, d *VirtualEnvironmentDatastoreListRequestBody) ([]*VirtualEnvironmentDatastoreListResponseData, error) {
	resBody := &VirtualEnvironmentDatastoreListResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("nodes/%s/storage", nodeName), d, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	return resBody.Data, nil
}
