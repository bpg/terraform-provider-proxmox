/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fs

// ListResponseBody contains the body from a CephFS list response.
type ListResponseBody struct {
	Data []*ListResponseData `json:"data,omitempty"`
}

// ListResponseData describes one filesystem. data_pools and the pool IDs arrived in pve-manager 9.1.18; older
// releases return only data_pool, the first data pool.
type ListResponseData struct {
	Name           string   `json:"name"`
	MetadataPool   string   `json:"metadata_pool"`
	MetadataPoolID *int64   `json:"metadata_pool_id,omitempty"`
	DataPool       string   `json:"data_pool,omitempty"`
	DataPools      []string `json:"data_pools,omitempty"`
	DataPoolIDs    []int64  `json:"data_pool_ids,omitempty"`
}

// FirstDataPool returns the filesystem's first data pool, the one Proxmox created with it, and its ID when known.
func (d *ListResponseData) FirstDataPool() (string, *int64) {
	name := d.DataPool
	if len(d.DataPools) > 0 {
		name = d.DataPools[0]
	}

	if len(d.DataPoolIDs) > 0 {
		id := d.DataPoolIDs[0]

		return name, &id
	}

	return name, nil
}

// CreateRequestBody contains the body for creating a CephFS.
type CreateRequestBody struct {
	PGNum *int64 `url:"pg_num,omitempty"`
}

// CreateResponseBody contains the body from a CephFS create response.
type CreateResponseBody struct {
	Data *string `json:"data,omitempty"`
}
