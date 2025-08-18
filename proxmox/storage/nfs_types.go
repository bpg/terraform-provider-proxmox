package storage

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// NFSStorageMutableFields defines the mutable attributes for 'nfs' type storage.
type NFSStorageMutableFields struct {
	DataStoreCommonMutableFields
	Preallocation          *string          `json:"preallocation,omitempty" url:"preallocation,omitempty"`
	SnapshotsAsVolumeChain types.CustomBool `json:"snapshot-as-volume-chain,omitempty" url:"snapshot-as-volume-chain,omitempty"`
	Options                *string          `json:"options,omitempty" url:"options,omitempty"`
}

// NFSStorageImmutableFields defines the immutable attributes for 'nfs' type storage.
type NFSStorageImmutableFields struct {
	Server *string `json:"server,omitempty" url:"server,omitempty"`
	Export *string `json:"export,omitempty" url:"export,omitempty"`
}

// NFSStorageCreateRequest defines the request body for creating a new NFS storage.
type NFSStorageCreateRequest struct {
	DataStoreCommonImmutableFields
	NFSStorageMutableFields
	NFSStorageImmutableFields
}

// NFSStorageUpdateRequest defines the request body for updating an existing NFS storage.
type NFSStorageUpdateRequest struct {
	NFSStorageMutableFields
}
