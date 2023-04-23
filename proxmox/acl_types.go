/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmox

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// ACLGetResponseBody contains the body from an access control list response.
type ACLGetResponseBody struct {
	Data []*ACLGetResponseData `json:"data,omitempty"`
}

// ACLGetResponseData contains the data from an access control list response.
type ACLGetResponseData struct {
	Path          string            `json:"path"`
	Propagate     *types.CustomBool `json:"propagate,omitempty"`
	RoleID        string            `json:"roleid"`
	Type          string            `json:"type"`
	UserOrGroupID string            `json:"ugid"`
}

// ACLUpdateRequestBody contains the data for an access control list update request.
type ACLUpdateRequestBody struct {
	Delete    *types.CustomBool `json:"delete,omitempty"    url:"delete,omitempty,int"`
	Groups    []string          `json:"groups,omitempty"    url:"groups,omitempty,comma"`
	Path      string            `json:"path"                url:"path"`
	Propagate *types.CustomBool `json:"propagate,omitempty" url:"propagate,omitempty,int"`
	Roles     []string          `json:"roles"               url:"roles,comma"`
	Users     []string          `json:"users,omitempty"     url:"users,omitempty,comma"`
}
