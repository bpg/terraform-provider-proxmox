/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// CephFSStorageMutableFields defines the mutable attributes for 'cephfs' type storage.
type CephFSStorageMutableFields struct {
	DataStoreCommonMutableFields

	Backups DataStoreWithBackups `json:"-" url:"backups,omitempty"`

	Options *string `json:"options,omitempty" url:"options,omitempty"`
}

// CephFSStorageImmutableFields defines the attributes of a 'cephfs' storage that are only sent on creation. Proxmox
// would accept them in an update, but a changed source doesn't replace an existing mount.
type CephFSStorageImmutableFields struct {
	FSName *string           `json:"fs-name,omitempty" url:"fs-name,omitempty"`
	Subdir *string           `json:"subdir,omitempty"  url:"subdir,omitempty"`
	FUSE   *types.CustomBool `json:"fuse,omitempty"    url:"fuse,omitempty,int"`
}

// CephFSStorageCreateRequest defines the request body for creating a new CephFS storage.
type CephFSStorageCreateRequest struct {
	DataStoreCommonImmutableFields
	CephFSStorageMutableFields
	CephFSStorageImmutableFields
}

// CephFSStorageUpdateRequest defines the request body for updating an existing CephFS storage.
type CephFSStorageUpdateRequest struct {
	CephFSStorageMutableFields
}
