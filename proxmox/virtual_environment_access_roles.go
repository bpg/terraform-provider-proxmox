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

// VirtualEnvironmentAccessRoleCreateRequestBody contains the data for an access group create request.
type VirtualEnvironmentAccessRoleCreateRequestBody struct {
	ID         string           `json:"roleid" url:"roleid"`
	Privileges CustomPrivileges `json:"privs" url:"privs,comma"`
}

// VirtualEnvironmentAccessRoleGetResponseBody contains the body from an access group get response.
type VirtualEnvironmentAccessRoleGetResponseBody struct {
	Data *CustomPrivileges `json:"data,omitempty"`
}

// VirtualEnvironmentAccessRoleListResponseBody contains the body from an access group list response.
type VirtualEnvironmentAccessRoleListResponseBody struct {
	Data []*VirtualEnvironmentAccessRoleListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentAccessRoleListResponseData contains the data from an access group list response.
type VirtualEnvironmentAccessRoleListResponseData struct {
	ID         string            `json:"roleid"`
	Privileges *CustomPrivileges `json:"privs"`
	Special    CustomBool        `json:"special"`
}

// VirtualEnvironmentAccessRoleUpdateRequestBody contains the data for an access group update request.
type VirtualEnvironmentAccessRoleUpdateRequestBody struct {
	Privileges CustomPrivileges `json:"privs" url:"privs,comma"`
}

// CreateAccessRole creates an access role.
func (c *VirtualEnvironmentClient) CreateAccessRole(d *VirtualEnvironmentAccessRoleCreateRequestBody) error {
	return c.DoRequest(hmPOST, "access/roles", d, nil)
}

// DeleteAccessRole deletes an access role.
func (c *VirtualEnvironmentClient) DeleteAccessRole(id string) error {
	return c.DoRequest(hmDELETE, fmt.Sprintf("access/roles/%s", id), nil, nil)
}

// GetAccessRole retrieves an access role.
func (c *VirtualEnvironmentClient) GetAccessRole(id string) (*CustomPrivileges, error) {
	resBody := &VirtualEnvironmentAccessRoleGetResponseBody{}
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

// ListAccessRoles retrieves a list of access roles.
func (c *VirtualEnvironmentClient) ListAccessRoles() ([]*VirtualEnvironmentAccessRoleListResponseData, error) {
	resBody := &VirtualEnvironmentAccessRoleListResponseBody{}
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

	return resBody.Data, nil
}

// UpdateAccessRole updates an access role.
func (c *VirtualEnvironmentClient) UpdateAccessRole(id string, d *VirtualEnvironmentAccessRoleUpdateRequestBody) error {
	return c.DoRequest(hmPUT, fmt.Sprintf("access/roles/%s", id), d, nil)
}
