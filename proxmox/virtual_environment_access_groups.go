/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"fmt"
	"net/url"
)

// VirtualEnvironmentAccessGroupCreateRequestData contains the data for an access group create request.
type VirtualEnvironmentAccessGroupCreateRequestData struct {
	Comment string `json:"comment" url:"comment"`
	ID      string `json:"groupid" url:"groupid"`
}

// VirtualEnvironmentAccessGroupGetResponseBody contains the body from an access group get response.
type VirtualEnvironmentAccessGroupGetResponseBody struct {
	Data *VirtualEnvironmentAccessGroupGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentAccessGroupGetResponseData contains the data from an access group get response.
type VirtualEnvironmentAccessGroupGetResponseData struct {
	Comment string   `json:"comment"`
	Members []string `json:"members"`
}

// VirtualEnvironmentAccessGroupListResponseBody contains the body from an access group list response.
type VirtualEnvironmentAccessGroupListResponseBody struct {
	Data []*VirtualEnvironmentAccessGroupListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentAccessGroupListResponseData contains the data from an access group list response.
type VirtualEnvironmentAccessGroupListResponseData struct {
	Comment string `json:"comment"`
	ID      string `json:"groupid"`
}

// CreateAccessGroup creates an access group.
func (c *VirtualEnvironmentClient) CreateAccessGroup(id, comment string) error {
	reqBody := &VirtualEnvironmentAccessGroupCreateRequestData{
		Comment: comment,
		ID:      id,
	}

	err := c.DoRequest(hmPOST, "access/groups", reqBody, nil)

	if err != nil {
		return err
	}

	return nil
}

// GetAccessGroup retrieves an access group.
func (c *VirtualEnvironmentClient) GetAccessGroup(id string) (*VirtualEnvironmentAccessGroupGetResponseData, error) {
	resBody := &VirtualEnvironmentAccessGroupGetResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("access/groups/%s", url.PathEscape(id)), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListAccessGroups retrieves a list of access groups.
func (c *VirtualEnvironmentClient) ListAccessGroups() ([]*VirtualEnvironmentAccessGroupListResponseData, error) {
	resBody := &VirtualEnvironmentAccessGroupListResponseBody{}
	err := c.DoRequest(hmGET, "access/groups", nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	return resBody.Data, nil
}
