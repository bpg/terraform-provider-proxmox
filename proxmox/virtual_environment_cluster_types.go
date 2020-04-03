/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

// VirtualEnvironmentClusterNextIDRequestBody contains the data for a cluster next id request.
type VirtualEnvironmentClusterNextIDRequestBody struct {
	VMID *int `json:"vmid,omitempty" url:"vmid,omitempty"`
}

// VirtualEnvironmentClusterNextIDResponseBody contains the body from a cluster next id response.
type VirtualEnvironmentClusterNextIDResponseBody struct {
	Data *CustomInt `json:"data,omitempty"`
}
