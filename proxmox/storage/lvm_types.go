package storage

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// LVMStorageMutableFields defines options for 'lvm' type storage.
type LVMStorageMutableFields struct {
	DataStoreCommonMutableFields
	WipeRemovedVolumes types.CustomBool `json:"saferemove" url:"saferemove,int"`
}

// LVMStorageImmutableFields defines options for 'lvm' type storage.
type LVMStorageImmutableFields struct {
	VolumeGroup *string `json:"vgname" url:"vgname"`
}

// LVMStorageCreateRequest defines the request body for creating a new LVM storage.
type LVMStorageCreateRequest struct {
	DataStoreCommonImmutableFields
	LVMStorageMutableFields
	LVMStorageImmutableFields
}

// LVMStorageUpdateRequest defines the request body for updating an existing LVM storage.
type LVMStorageUpdateRequest struct {
	LVMStorageMutableFields
}
