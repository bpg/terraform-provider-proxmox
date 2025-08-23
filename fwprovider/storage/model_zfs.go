package storage

import (
	"context"

	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	proxmox_types "github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ZFSStorageModel maps the Terraform schema for ZFS storage.
type ZFSStorageModel struct {
	StorageModelBase
	ZFSPool       types.String `tfsdk:"zfs_pool"`
	ThinProvision types.Bool   `tfsdk:"thin_provision"`
	Blocksize     types.String `tfsdk:"blocksize"`
}

// GetStorageType returns the storage type identifier.
func (m *ZFSStorageModel) GetStorageType() types.String {
	return types.StringValue("zfspool")
}

// toCreateAPIRequest converts the Terraform model to a Proxmox API request body.
func (m *ZFSStorageModel) toCreateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.ZFSStorageCreateRequest{}
	request.Type = m.GetStorageType().ValueStringPointer()

	if err := m.populateCreateFields(ctx, &request.DataStoreCommonImmutableFields, &request.ZFSStorageMutableFields.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.ZFSPool = m.ZFSPool.ValueStringPointer()
	request.ThinProvision = proxmox_types.CustomBool(m.ThinProvision.ValueBool())
	request.Blocksize = m.Blocksize.ValueStringPointer()

	return request, nil
}

// toUpdateAPIRequest converts the Terraform model to a Proxmox API request body for updates.
func (m *ZFSStorageModel) toUpdateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.ZFSStorageUpdateRequest{}

	if err := m.populateUpdateFields(ctx, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.ThinProvision = proxmox_types.CustomBool(m.ThinProvision.ValueBool())
	request.Blocksize = m.Blocksize.ValueStringPointer()

	return request, nil
}

// fromAPI populates the Terraform model from a Proxmox API response.
func (m *ZFSStorageModel) fromAPI(ctx context.Context, datastore *storage.DatastoreGetResponseData) error {
	if err := m.populateBaseFromAPI(ctx, datastore); err != nil {
		return err
	}

	if datastore.ZFSPool != nil {
		m.ZFSPool = types.StringValue(*datastore.ZFSPool)
	}
	if datastore.ThinProvision != nil {
		m.ThinProvision = types.BoolValue(*datastore.ThinProvision.PointerBool())
	}
	if datastore.Blocksize != nil {
		m.Blocksize = types.StringValue(*datastore.Blocksize)
	}

	return nil
}
