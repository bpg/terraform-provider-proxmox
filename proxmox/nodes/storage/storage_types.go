/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// DatastoreListRequestBody contains the body for a datastore list request.
type DatastoreListRequestBody struct {
	ContentTypes types.CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Enabled      *types.CustomBool              `json:"enabled,omitempty" url:"enabled,omitempty,int"`
	Format       *types.CustomBool              `json:"format,omitempty"  url:"format,omitempty,int"`
	ID           *string                        `json:"storage,omitempty" url:"storage,omitempty"`
	Target       *string                        `json:"target,omitempty"  url:"target,omitempty"`
}

// DatastoreListResponseBody contains the body from a datastore list response.
type DatastoreListResponseBody struct {
	Data []*DatastoreListResponseData `json:"data,omitempty"`
}

// DatastoreListResponseData contains the data from a datastore list response.
type DatastoreListResponseData struct {
	Active              *types.CustomBool               `json:"active,omitempty"`
	ContentTypes        *types.CustomCommaSeparatedList `json:"content,omitempty"`
	Enabled             *types.CustomBool               `json:"enabled,omitempty"`
	ID                  string                          `json:"storage,omitempty"`
	Shared              *types.CustomBool               `json:"shared,omitempty"`
	SpaceAvailable      *types.CustomInt64              `json:"avail,omitempty"`
	SpaceTotal          *types.CustomInt64              `json:"total,omitempty"`
	SpaceUsed           *types.CustomInt64              `json:"used,omitempty"`
	SpaceUsedPercentage *types.CustomFloat64            `json:"used_fraction,omitempty"`
	Type                string                          `json:"type,omitempty"`
}
