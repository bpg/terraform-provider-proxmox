package storage

// LVMThinStorageMutableFields defines options for 'lvmthin' type storage.
type LVMThinStorageMutableFields struct {
	DataStoreCommonMutableFields
}

// LVMThinStorageImmutableFields defines options for 'lvmthin' type storage.
type LVMThinStorageImmutableFields struct {
	VolumeGroup *string `json:"vgname"             url:"vgname"`
	ThinPool    *string `json:"thinpool,omitempty" url:"thinpool,omitempty"`
}

// LVMThinStorageCreateRequest defines the request body for creating a new LVM thin storage.
type LVMThinStorageCreateRequest struct {
	DataStoreCommonImmutableFields
	LVMThinStorageMutableFields
	LVMThinStorageImmutableFields
}

// LVMThinStorageUpdateRequest defines the request body for updating an existing LVM thin storage.
type LVMThinStorageUpdateRequest struct {
	LVMThinStorageMutableFields
}
