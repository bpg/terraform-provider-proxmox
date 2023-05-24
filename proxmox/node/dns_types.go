/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package node

// DNSGetResponseBody contains the body from a DNS get response.
type DNSGetResponseBody struct {
	Data *DNSGetResponseData `json:"data,omitempty"`
}

// DNSGetResponseData contains the data from a DNS get response.
type DNSGetResponseData struct {
	Server1      *string `json:"dns1,omitempty"   url:"dns1,omitempty"`
	Server2      *string `json:"dns2,omitempty"   url:"dns2,omitempty"`
	Server3      *string `json:"dns3,omitempty"   url:"dns3,omitempty"`
	SearchDomain *string `json:"search,omitempty" url:"search,omitempty"`
}

// DNSUpdateRequestBody contains the body for a DNS update request.
type DNSUpdateRequestBody struct {
	Server1      *string `json:"dns1,omitempty"   url:"dns1,omitempty"`
	Server2      *string `json:"dns2,omitempty"   url:"dns2,omitempty"`
	Server3      *string `json:"dns3,omitempty"   url:"dns3,omitempty"`
	SearchDomain *string `json:"search,omitempty" url:"search,omitempty"`
}
