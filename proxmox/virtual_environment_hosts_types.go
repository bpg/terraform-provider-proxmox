/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

// VirtualEnvironmentHostsGetResponseBody contains the body from a hosts get response.
type VirtualEnvironmentHostsGetResponseBody struct {
	Data *VirtualEnvironmentHostsGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentHostsGetResponseData contains the data from a hosts get response.
type VirtualEnvironmentHostsGetResponseData struct {
	Data   string  `json:"data"`
	Digest *string `json:"digest,omitempty"`
}

// VirtualEnvironmentHostsUpdateRequestBody contains the body for a hosts update request.
type VirtualEnvironmentHostsUpdateRequestBody struct {
	Data   string  `json:"data"             url:"data"`
	Digest *string `json:"digest,omitempty" url:"digest,omitempty"`
}
