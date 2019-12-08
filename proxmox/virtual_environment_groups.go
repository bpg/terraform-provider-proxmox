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

// VirtualEnvironmentGroupCreateRequestBody contains the data for an access group create request.
type VirtualEnvironmentGroupCreateRequestBody struct {
	Comment string `json:"comment" url:"comment"`
	ID      string `json:"groupid" url:"groupid"`
}

// VirtualEnvironmentGroupGetResponseBody contains the body from an access group get response.
type VirtualEnvironmentGroupGetResponseBody struct {
	Data *VirtualEnvironmentGroupGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentGroupGetResponseData contains the data from an access group get response.
type VirtualEnvironmentGroupGetResponseData struct {
	Comment string   `json:"comment"`
	Members []string `json:"members"`
}

// VirtualEnvironmentGroupListResponseBody contains the body from an access group list response.
type VirtualEnvironmentGroupListResponseBody struct {
	Data []*VirtualEnvironmentGroupListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentGroupListResponseData contains the data from an access group list response.
type VirtualEnvironmentGroupListResponseData struct {
	Comment string `json:"comment"`
	ID      string `json:"groupid"`
}

// VirtualEnvironmentGroupUpdateRequestBody contains the data for an access group update request.
type VirtualEnvironmentGroupUpdateRequestBody struct {
	Comment string `json:"comment" url:"comment"`
}

// CreateGroup creates an access group.
func (c *VirtualEnvironmentClient) CreateGroup(d *VirtualEnvironmentGroupCreateRequestBody) error {
	return c.DoRequest(hmPOST, "access/groups", d, nil)
}

// DeleteGroup deletes an access group.
func (c *VirtualEnvironmentClient) DeleteGroup(id string) error {
	return c.DoRequest(hmDELETE, fmt.Sprintf("access/groups/%s", id), nil, nil)
}

// GetGroup retrieves an access group.
func (c *VirtualEnvironmentClient) GetGroup(id string) (*VirtualEnvironmentGroupGetResponseData, error) {
	resBody := &VirtualEnvironmentGroupGetResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("access/groups/%s", url.PathEscape(id)), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	sort.Strings(resBody.Data.Members)

	return resBody.Data, nil
}

// ListGroups retrieves a list of access groups.
func (c *VirtualEnvironmentClient) ListGroups() ([]*VirtualEnvironmentGroupListResponseData, error) {
	resBody := &VirtualEnvironmentGroupListResponseBody{}
	err := c.DoRequest(hmGET, "access/groups", nil, resBody)

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

// UpdateGroup updates an access group.
func (c *VirtualEnvironmentClient) UpdateGroup(id string, d *VirtualEnvironmentGroupUpdateRequestBody) error {
	return c.DoRequest(hmPUT, fmt.Sprintf("access/groups/%s", id), d, nil)
}
