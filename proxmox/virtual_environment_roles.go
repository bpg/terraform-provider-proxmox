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

// VirtualEnvironmentRoleCreateRequestBody contains the data for an access group create request.
type VirtualEnvironmentRoleCreateRequestBody struct {
	ID         string           `json:"roleid" url:"roleid"`
	Privileges CustomPrivileges `json:"privs" url:"privs,comma"`
}

// VirtualEnvironmentRoleGetResponseBody contains the body from an access group get response.
type VirtualEnvironmentRoleGetResponseBody struct {
	Data *CustomPrivileges `json:"data,omitempty"`
}

// VirtualEnvironmentRoleListResponseBody contains the body from an access group list response.
type VirtualEnvironmentRoleListResponseBody struct {
	Data []*VirtualEnvironmentRoleListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentRoleListResponseData contains the data from an access group list response.
type VirtualEnvironmentRoleListResponseData struct {
	ID         string            `json:"roleid"`
	Privileges *CustomPrivileges `json:"privs,omitempty"`
	Special    *CustomBool       `json:"special,omitempty"`
}

// VirtualEnvironmentRoleUpdateRequestBody contains the data for an access group update request.
type VirtualEnvironmentRoleUpdateRequestBody struct {
	Privileges CustomPrivileges `json:"privs" url:"privs,comma"`
}

// CreateRole creates an access role.
func (c *VirtualEnvironmentClient) CreateRole(d *VirtualEnvironmentRoleCreateRequestBody) error {
	return c.DoRequest(hmPOST, "access/roles", d, nil)
}

// DeleteRole deletes an access role.
func (c *VirtualEnvironmentClient) DeleteRole(id string) error {
	return c.DoRequest(hmDELETE, fmt.Sprintf("access/roles/%s", id), nil, nil)
}

// GetRole retrieves an access role.
func (c *VirtualEnvironmentClient) GetRole(id string) (*CustomPrivileges, error) {
	resBody := &VirtualEnvironmentRoleGetResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("access/roles/%s", url.PathEscape(id)), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	sort.Strings(*resBody.Data)

	return resBody.Data, nil
}

// ListRoles retrieves a list of access roles.
func (c *VirtualEnvironmentClient) ListRoles() ([]*VirtualEnvironmentRoleListResponseData, error) {
	resBody := &VirtualEnvironmentRoleListResponseBody{}
	err := c.DoRequest(hmGET, "access/roles", nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	for i := range resBody.Data {
		if resBody.Data[i].Privileges != nil {
			sort.Strings(*resBody.Data[i].Privileges)
		}
	}

	return resBody.Data, nil
}

// UpdateRole updates an access role.
func (c *VirtualEnvironmentClient) UpdateRole(id string, d *VirtualEnvironmentRoleUpdateRequestBody) error {
	return c.DoRequest(hmPUT, fmt.Sprintf("access/roles/%s", id), d, nil)
}
