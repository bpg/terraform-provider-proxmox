/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

// VirtualEnvironmentClusterAliasCreateRequestBody contains the data for an alias create request.
type VirtualEnvironmentClusterAliasCreateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Name    string  `json:"name"              url:"name"`
	CIDR    string  `json:"cidr"              url:"cidr"`
}

// VirtualEnvironmentClusterAliasGetResponseBody contains the body from an alias get response.
type VirtualEnvironmentClusterAliasGetResponseBody struct {
	Data *VirtualEnvironmentClusterAliasGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentClusterAliasGetResponseData contains the data from an alias get response.
type VirtualEnvironmentClusterAliasGetResponseData struct {
	Comment   *string `json:"comment,omitempty" url:"comment,omitempty"`
	Name      string  `json:"name"              url:"name"`
	CIDR      string  `json:"cidr"              url:"cidr"`
	Digest    *string `json:"digest"            url:"digest"`
	IPVersion int     `json:"ipversion"         url:"ipversion"`
}

// VirtualEnvironmentClusterAliasListResponseBody contains the data from an alias get response.
type VirtualEnvironmentClusterAliasListResponseBody struct {
	Data []*VirtualEnvironmentClusterAliasGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentClusterAliasUpdateRequestBody contains the data for an alias update request.
type VirtualEnvironmentClusterAliasUpdateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	ReName  string  `json:"rename"            url:"rename"`
	CIDR    string  `json:"cidr"              url:"cidr"`
}
