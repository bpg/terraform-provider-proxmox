/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// VirtualEnvironmentRoleCreateRequestBody contains the data for an access group create request.
type VirtualEnvironmentRoleCreateRequestBody struct {
	ID         string                 `json:"roleid" url:"roleid"`
	Privileges types.CustomPrivileges `json:"privs"  url:"privs,comma"`
}

// VirtualEnvironmentRoleGetResponseBody contains the body from an access group get response.
type VirtualEnvironmentRoleGetResponseBody struct {
	Data *types.CustomPrivileges `json:"data,omitempty"`
}

// VirtualEnvironmentRoleListResponseBody contains the body from an access group list response.
type VirtualEnvironmentRoleListResponseBody struct {
	Data []*VirtualEnvironmentRoleListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentRoleListResponseData contains the data from an access group list response.
type VirtualEnvironmentRoleListResponseData struct {
	ID         string                  `json:"roleid"`
	Privileges *types.CustomPrivileges `json:"privs,omitempty"`
	Special    *types.CustomBool       `json:"special,omitempty"`
}

// VirtualEnvironmentRoleUpdateRequestBody contains the data for an access group update request.
type VirtualEnvironmentRoleUpdateRequestBody struct {
	Privileges types.CustomPrivileges `json:"privs" url:"privs,comma"`
}
