/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
)

// VirtualEnvironmentPoolCreateRequestBody contains the data for an pool create request.
type VirtualEnvironmentPoolCreateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	ID      string  `json:"groupid" url:"poolid"`
}

// VirtualEnvironmentPoolGetResponseBody contains the body from an pool get response.
type VirtualEnvironmentPoolGetResponseBody struct {
	Data *VirtualEnvironmentPoolGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentPoolGetResponseData contains the data from an pool get response.
type VirtualEnvironmentPoolGetResponseData struct {
	Comment *string                                    `json:"comment,omitempty"`
	Members []VirtualEnvironmentPoolGetResponseMembers `json:"members,omitempty"`
}

// VirtualEnvironmentPoolGetResponseMembers contains the members data from an pool get response.
type VirtualEnvironmentPoolGetResponseMembers struct {
	ID          string  `json:"id"`
	Node        string  `json:"node"`
	DatastoreID *string `json:"storage,omitempty"`
	Type        string  `json:"type"`
	VMID        *int    `json:"vmid"`
}

// VirtualEnvironmentPoolListResponseBody contains the body from an pool list response.
type VirtualEnvironmentPoolListResponseBody struct {
	Data []*VirtualEnvironmentPoolListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentPoolListResponseData contains the data from an pool list response.
type VirtualEnvironmentPoolListResponseData struct {
	Comment *string `json:"comment,omitempty"`
	ID      string  `json:"poolid"`
}

// VirtualEnvironmentPoolUpdateRequestBody contains the data for an pool update request.
type VirtualEnvironmentPoolUpdateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
}

// CreatePool creates an pool.
func (c *VirtualEnvironmentClient) CreatePool(d *VirtualEnvironmentPoolCreateRequestBody) error {
	return c.DoRequest(hmPOST, "pools", d, nil)
}

// DeletePool deletes an pool.
func (c *VirtualEnvironmentClient) DeletePool(id string) error {
	return c.DoRequest(hmDELETE, fmt.Sprintf("pools/%s", url.PathEscape(id)), nil, nil)
}

// GetPool retrieves an pool.
func (c *VirtualEnvironmentClient) GetPool(id string) (*VirtualEnvironmentPoolGetResponseData, error) {
	resBody := &VirtualEnvironmentPoolGetResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("pools/%s", url.PathEscape(id)), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	sort.Slice(resBody.Data.Members, func(i, j int) bool {
		return resBody.Data.Members[i].ID < resBody.Data.Members[j].ID
	})

	return resBody.Data, nil
}

// ListPools retrieves a list of pools.
func (c *VirtualEnvironmentClient) ListPools() ([]*VirtualEnvironmentPoolListResponseData, error) {
	resBody := &VirtualEnvironmentPoolListResponseBody{}
	err := c.DoRequest(hmGET, "pools", nil, resBody)

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

// UpdatePool updates an pool.
func (c *VirtualEnvironmentClient) UpdatePool(id string, d *VirtualEnvironmentPoolUpdateRequestBody) error {
	return c.DoRequest(hmPUT, fmt.Sprintf("pools/%s", url.PathEscape(id)), d, nil)
}
