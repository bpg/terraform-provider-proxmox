/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type optionsdatabase struct {
	BandwidthLimit string  `json:"bwlimit,omitempty"    url:"bwlimit,omitempty"`
	EmailFrom      *string `json:"email_from,omitempty" url:"email_from,omitempty"`
	Keyboard       *string `json:"keyboard,omitempty"   url:"keyboard,omitempty"`
	Language       *string `json:"language,omitempty"   url:"language,omitempty"`
}

// OptionsResponseBody contains the body from a cluster options response.
type OptionsResponseBody struct {
	Data *OptionsResponseData `json:"data,omitempty"`
}

// OptionsResponseData contains the data from a cluster options response.
type OptionsResponseData struct {
	optionsdatabase
	MaxWorkers *types.CustomInt `json:"max_workers,omitempty"`
}

// OptionsRequestData contains the body for cluster options request.
type OptionsRequestData struct {
	optionsdatabase
	MaxWorkers *int64 `json:"max_workers,omitempty" url:"max_workers,omitempty"`
	Delete     string `json:"delete,omitempty"      url:"delete,omitempty"`
}
