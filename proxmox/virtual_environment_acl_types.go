/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// VirtualEnvironmentACLGetResponseBody contains the body from an access control list response.
type VirtualEnvironmentACLGetResponseBody struct {
	Data []*VirtualEnvironmentACLGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentACLGetResponseData contains the data from an access control list response.
type VirtualEnvironmentACLGetResponseData struct {
	Path          string            `json:"path"`
	Propagate     *types.CustomBool `json:"propagate,omitempty"`
	RoleID        string            `json:"roleid"`
	Type          string            `json:"type"`
	UserOrGroupID string            `json:"ugid"`
}

// VirtualEnvironmentACLUpdateRequestBody contains the data for an access control list update request.
type VirtualEnvironmentACLUpdateRequestBody struct {
	Delete    *types.CustomBool `json:"delete,omitempty"    url:"delete,omitempty,int"`
	Groups    []string          `json:"groups,omitempty"    url:"groups,omitempty,comma"`
	Path      string            `json:"path"                url:"path"`
	Propagate *types.CustomBool `json:"propagate,omitempty" url:"propagate,omitempty,int"`
	Roles     []string          `json:"roles"               url:"roles,comma"`
	Users     []string          `json:"users,omitempty"     url:"users,omitempty,comma"`
}
