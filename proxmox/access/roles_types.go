/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// RoleCreateRequestBody contains the data for an access group create request.
type RoleCreateRequestBody struct {
	ID         string                 `json:"roleid" url:"roleid"`
	Privileges types.CustomPrivileges `json:"privs"  url:"privs,comma"`
}

// RoleGetResponseBody contains the body from an access group get response.
type RoleGetResponseBody struct {
	Data *types.CustomPrivileges `json:"data,omitempty"`
}

// RoleListResponseBody contains the body from an access group list response.
type RoleListResponseBody struct {
	Data []*RoleListResponseData `json:"data,omitempty"`
}

// RoleListResponseData contains the data from an access group list response.
type RoleListResponseData struct {
	ID         string                  `json:"roleid"`
	Privileges *types.CustomPrivileges `json:"privs,omitempty"`
	Special    *types.CustomBool       `json:"special,omitempty"`
}

// RoleUpdateRequestBody contains the data for an access group update request.
type RoleUpdateRequestBody struct {
	Privileges types.CustomPrivileges `json:"privs" url:"privs,comma"`
}
