/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zfs

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// CreateRequestBody contains the body for creating a ZFS pool.
type CreateRequestBody struct {
	Name      string `url:"name"`
	Devices   string `url:"devices"`
	RaidLevel string `url:"raidlevel"`

	AddStorage  *types.CustomBool `url:"add_storage,omitempty,int"`
	AShift      *int64            `url:"ashift,omitempty"`
	Compression *string           `url:"compression,omitempty"`
	DraidConfig *string           `url:"draid-config,omitempty"`
}

// DeleteRequestParams contains the query parameters for destroying a ZFS pool.
type DeleteRequestParams struct {
	CleanupConfig *types.CustomBool `url:"cleanup-config,omitempty,int"`
	CleanupDisks  *types.CustomBool `url:"cleanup-disks,omitempty,int"`
}

// GetResponseBody wraps the GET /nodes/{node}/disks/zfs/{name} response.
type GetResponseBody struct {
	Data *GetResponseData `json:"data,omitempty"`
}

// GetResponseData describes a single ZFS pool from the detail endpoint.
type GetResponseData struct {
	Name   string `json:"name"`
	State  string `json:"state"`
	Errors string `json:"errors"`
}

// ListResponseBody wraps the GET /nodes/{node}/disks/zfs response.
type ListResponseBody struct {
	Data []*ListResponseData `json:"data,omitempty"`
}

// ListResponseData describes one ZFS pool entry from the list endpoint.
type ListResponseData struct {
	Name   string  `json:"name"`
	Health string  `json:"health"`
	Alloc  int64   `json:"alloc"`
	Free   int64   `json:"free"`
	Size   int64   `json:"size"`
	Frag   int64   `json:"frag"`
	Dedup  float64 `json:"dedup"`
}

// CreateResponseBody wraps the create response (a UPID task identifier).
type CreateResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// DeleteResponseBody wraps the delete response (a UPID task identifier).
type DeleteResponseBody struct {
	Data *string `json:"data,omitempty"`
}
