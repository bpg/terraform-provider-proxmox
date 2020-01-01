/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

// VirtualEnvironmentDNSGetResponseBody contains the body from an pool get response.
type VirtualEnvironmentDNSGetResponseBody struct {
	Data *VirtualEnvironmentDNSGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentDNSGetResponseData contains the data from an pool get response.
type VirtualEnvironmentDNSGetResponseData struct {
	Server1      *string `json:"dns1,omitempty" url:"dns1,omitempty"`
	Server2      *string `json:"dns2,omitempty" url:"dns2,omitempty"`
	Server3      *string `json:"dns3,omitempty" url:"dns3,omitempty"`
	SearchDomain *string `json:"search,omitempty" url:"search,omitempty"`
}

// VirtualEnvironmentDNSUpdateRequestBody contains the data for an pool create request.
type VirtualEnvironmentDNSUpdateRequestBody struct {
	Server1      *string `json:"dns1,omitempty" url:"dns1,omitempty"`
	Server2      *string `json:"dns2,omitempty" url:"dns2,omitempty"`
	Server3      *string `json:"dns3,omitempty" url:"dns3,omitempty"`
	SearchDomain *string `json:"search,omitempty" url:"search,omitempty"`
}
