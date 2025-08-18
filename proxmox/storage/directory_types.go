package storage

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// DirectoryStorageMutableFields defines the mutable attributes for 'dir' type storage.
type DirectoryStorageMutableFields struct {
	DataStoreCommonMutableFields
	Preallocation          *string          `json:"preallocation,omitempty" url:"preallocation,omitempty"`
	SnapshotsAsVolumeChain types.CustomBool `json:"snapshot-as-volume-chain,omitempty" url:"snapshot-as-volume-chain,omitempty"`
}

// DirectoryStorageImmutableFields defines the immutable attributes for 'dir' type storage.
type DirectoryStorageImmutableFields struct {
	Path *string `json:"path,omitempty" url:"path,omitempty"`
}

// DirectoryStorageCreateRequest defines options for 'dir' type storage.
type DirectoryStorageCreateRequest struct {
	DataStoreCommonImmutableFields
	DirectoryStorageMutableFields
	DirectoryStorageImmutableFields
}

type DirectoryStorageUpdateRequest struct {
	DirectoryStorageMutableFields
}
