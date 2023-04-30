/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// IPSetListResponseBody contains the data from an IPSet get response.
type IPSetListResponseBody struct {
	Data []*IPSetListResponseData `json:"data,omitempty"`
}

// IPSetCreateRequestBody contains the data for an IPSet create request
type IPSetCreateRequestBody struct {
	Comment string `json:"comment,omitempty" url:"comment,omitempty"`
	Name    string `json:"name"              url:"name"`
}

// IPSetGetResponseBody contains the body from an IPSet get response.
type IPSetGetResponseBody struct {
	Data []*IPSetGetResponseData `json:"data,omitempty"`
}

// IPSetGetResponseData contains the data from an IPSet get response.
type IPSetGetResponseData struct {
	CIDR    string            `json:"cidr"              url:"cidr"`
	NoMatch *types.CustomBool `json:"nomatch,omitempty" url:"nomatch,omitempty,int"`
	Comment *string           `json:"comment,omitempty" url:"comment,omitempty"`
}

// IPSetUpdateRequestBody contains the data for an IPSet update request.
type IPSetUpdateRequestBody struct {
	ReName  string  `json:"rename,omitempty"  url:"rename,omitempty"`
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Name    string  `json:"name"              url:"name"`
}

// IPSetListResponseData contains list of IPSets from
type IPSetListResponseData struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Name    string  `json:"name"              url:"name"`
}

// IPSetContent is an array of IPSetGetResponseData.
type IPSetContent []IPSetGetResponseData
