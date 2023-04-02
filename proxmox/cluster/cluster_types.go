/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cluster

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// NextIDRequestBody contains the data for a cluster next id request.
type NextIDRequestBody struct {
	VMID *int `json:"vmid,omitempty" url:"vmid,omitempty"`
}

// NextIDResponseBody contains the body from a cluster next id response.
type NextIDResponseBody struct {
	Data *types.CustomInt `json:"data,omitempty"`
}
