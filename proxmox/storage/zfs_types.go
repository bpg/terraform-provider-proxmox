package storage

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// ZFSStorageMutableFields defines options for 'zfspool' type storage.
type ZFSStorageMutableFields struct {
	DataStoreCommonMutableFields
	ThinProvision types.CustomBool `json:"sparse,omitempty" url:"sparse,omitempty,int"`
	Blocksize     *string          `json:"blocksize,omitempty" url:"blocksize,omitempty"`
}

type ZFSStorageImmutableFields struct {
	DataStoreCommonMutableFields
	ZFSPool *string `json:"pool" url:"pool"`
}

// ZFSStorageCreateRequest defines the request body for creating a new ZFS storage.
type ZFSStorageCreateRequest struct {
	DataStoreCommonImmutableFields
	ZFSStorageMutableFields
	ZFSStorageImmutableFields
}

// ZFSStorageUpdateRequest defines the request body for updating an existing ZFS storage.
type ZFSStorageUpdateRequest struct {
	ZFSStorageMutableFields
}
