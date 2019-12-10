/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"sort"
)

// VirtualEnvironmentACLGetResponseBody contains the body from an access control list response.
type VirtualEnvironmentACLGetResponseBody struct {
	Data []*VirtualEnvironmentACLGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentACLGetResponseData contains the data from an access control list response.
type VirtualEnvironmentACLGetResponseData struct {
	Path          string      `json:"path"`
	Propagate     *CustomBool `json:"propagate,omitempty"`
	RoleID        string      `json:"roleid"`
	Type          string      `json:"type"`
	UserOrGroupID string      `json:"ugid"`
}

// VirtualEnvironmentACLUpdateRequestBody contains the data for an access control list update request.
type VirtualEnvironmentACLUpdateRequestBody struct {
	Delete    *CustomBool `json:"delete,omitempty" url:"delete,omitempty,int"`
	Groups    []string    `json:"groups,omitempty" url:"groups,omitempty,comma"`
	Path      string      `json:"path" url:"path"`
	Propagate *CustomBool `json:"propagate,omitempty" url:"propagate,omitempty,int"`
	Roles     []string    `json:"roles" url:"roles,comma"`
	Users     []string    `json:"users,omitempty" url:"users,omitempty,comma"`
}

// GetACL retrieves the access control list.
func (c *VirtualEnvironmentClient) GetACL() ([]*VirtualEnvironmentACLGetResponseData, error) {
	resBody := &VirtualEnvironmentACLGetResponseBody{}
	err := c.DoRequest(hmGET, "access/acl", nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Path < resBody.Data[j].Path
	})

	return resBody.Data, nil
}

// UpdateACL updates the access control list.
func (c *VirtualEnvironmentClient) UpdateACL(d *VirtualEnvironmentACLUpdateRequestBody) error {
	return c.DoRequest(hmPUT, "access/acl", d, nil)
}
