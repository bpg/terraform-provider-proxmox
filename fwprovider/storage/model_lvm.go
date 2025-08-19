package storage

import (
	"context"

	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	proxmox_types "github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LVMStorageModel maps the Terraform schema for LVM storage.
type LVMStorageModel struct {
	StorageModelBase
	VolumeGroup        types.String `tfsdk:"volume_group"`
	WipeRemovedVolumes types.Bool   `tfsdk:"wipe_removed_volumes"`
}

// GetStorageType returns the storage type identifier.
func (m *LVMStorageModel) GetStorageType() types.String {
	return types.StringValue("lvm")
}

// toCreateAPIRequest converts the Terraform model to a Proxmox API request body.
func (m *LVMStorageModel) toCreateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.LVMStorageCreateRequest{}
	request.Type = m.GetStorageType().ValueStringPointer()

	if err := m.populateCreateFields(ctx, &request.DataStoreCommonImmutableFields, &request.LVMStorageMutableFields.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.VolumeGroup = m.VolumeGroup.ValueStringPointer()
	request.WipeRemovedVolumes = proxmox_types.CustomBool(m.WipeRemovedVolumes.ValueBool())

	return request, nil
}

// toUpdateAPIRequest converts the Terraform model to a Proxmox API request body for updates.
func (m *LVMStorageModel) toUpdateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.LVMStorageUpdateRequest{}

	if err := m.populateUpdateFields(ctx, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.WipeRemovedVolumes = proxmox_types.CustomBool(m.WipeRemovedVolumes.ValueBool())

	return request, nil
}

// fromAPI populates the Terraform model from a Proxmox API response.
func (m *LVMStorageModel) fromAPI(ctx context.Context, datastore *storage.DatastoreGetResponseData) error {
	if err := m.populateBaseFromAPI(ctx, datastore); err != nil {
		return err
	}

	if datastore.VolumeGroup != nil {
		m.VolumeGroup = types.StringValue(*datastore.VolumeGroup)
	}
	if datastore.WipeRemovedVolumes != nil {
		m.WipeRemovedVolumes = types.BoolValue(*datastore.WipeRemovedVolumes.PointerBool())
	}

	return nil
}
