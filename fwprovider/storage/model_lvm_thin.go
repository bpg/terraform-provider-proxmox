package storage

import (
	"context"

	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LVMThinStorageModel maps the Terraform schema for LVM storage.
type LVMThinStorageModel struct {
	StorageModelBase
	VolumeGroup types.String `tfsdk:"volume_group"`
	ThinPool    types.String `tfsdk:"thin_pool"`
}

// GetStorageType returns the storage type identifier.
func (m *LVMThinStorageModel) GetStorageType() types.String {
	return types.StringValue("lvmthin")
}

// toCreateAPIRequest converts the Terraform model to a Proxmox API request body.
func (m *LVMThinStorageModel) toCreateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.LVMThinStorageCreateRequest{}
	request.Type = m.GetStorageType().ValueStringPointer()

	if err := m.populateCreateFields(ctx, &request.DataStoreCommonImmutableFields, &request.LVMThinStorageMutableFields.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.VolumeGroup = m.VolumeGroup.ValueStringPointer()
	request.ThinPool = m.ThinPool.ValueStringPointer()

	return request, nil
}

// toUpdateAPIRequest converts the Terraform model to a Proxmox API request body for updates.
func (m *LVMThinStorageModel) toUpdateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.LVMThinStorageUpdateRequest{}

	if err := m.populateUpdateFields(ctx, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	return request, nil
}

// fromAPI populates the Terraform model from a Proxmox API response.
func (m *LVMThinStorageModel) fromAPI(ctx context.Context, datastore *storage.DatastoreGetResponseData) error {
	if err := m.populateBaseFromAPI(ctx, datastore); err != nil {
		return err
	}

	if datastore.VolumeGroup != nil {
		m.VolumeGroup = types.StringValue(*datastore.VolumeGroup)
	}
	if datastore.ThinPool != nil {
		m.ThinPool = types.StringValue(*datastore.ThinPool)
	}

	return nil
}
