package storage

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// CIFSStorageMutableFields defines specific options for 'smb'/'cifs' type storage.
type CIFSStorageMutableFields struct {
	DataStoreCommonMutableFields
	DataStoreWithBackups

	Preallocation *string `json:"preallocation,omitempty" url:"preallocation,omitempty"`
}

type CIFSStorageImmutableFields struct {
	Server                 *string          `json:"server"                             url:"server"`
	Username               *string          `json:"username"                           url:"username"`
	Password               *string          `json:"password"                           url:"password"`
	Share                  *string          `json:"share"                              url:"share"`
	Domain                 *string          `json:"domain,omitempty"                   url:"domain,omitempty"`
	Subdirectory           *string          `json:"subdir,omitempty"                   url:"subdir,omitempty"`
	SnapshotsAsVolumeChain types.CustomBool `json:"snapshot-as-volume-chain,omitempty" url:"snapshot-as-volume-chain,omitempty"`
}

// CIFSStorageCreateRequest defines the request body for creating a new SMB/CIFS storage.
type CIFSStorageCreateRequest struct {
	DataStoreCommonImmutableFields
	CIFSStorageMutableFields
	CIFSStorageImmutableFields
}

// CIFSStorageUpdateRequest defines the request body for updating an existing SMB/CIFS storage.
type CIFSStorageUpdateRequest struct {
	CIFSStorageMutableFields
}
