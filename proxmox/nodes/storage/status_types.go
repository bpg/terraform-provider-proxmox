/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// DatastoreGetStatusResponseBody contains the body from a datastore status get request.
type DatastoreGetStatusResponseBody struct {
	Data *DatastoreGetStatusResponseData `json:"data,omitempty"`
}

// DatastoreGetStatusResponseData contains the data from a datastore status get request.
type DatastoreGetStatusResponseData struct {
	Active         *types.CustomBool               `json:"active,omitempty"`
	AvailableBytes *int64                          `json:"avail,omitempty"`
	Content        *types.CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Enabled        *types.CustomBool               `json:"enabled,omitempty"`
	Shared         *types.CustomBool               `json:"shared,omitempty"`
	TotalBytes     *int64                          `json:"total,omitempty"`
	Type           *string                         `json:"type,omitempty"`
	UsedBytes      *int64                          `json:"used,omitempty"`
}
