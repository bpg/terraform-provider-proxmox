package storage

import (
	"context"

	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	proxmox_types "github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CIFSStorageModel maps the Terraform schema for CIFS storage.
type CIFSStorageModel struct {
	StorageModelBase
	Server                 types.String `tfsdk:"server"`
	Username               types.String `tfsdk:"username"`
	Password               types.String `tfsdk:"password"`
	Share                  types.String `tfsdk:"share"`
	Domain                 types.String `tfsdk:"domain"`
	SubDirectory           types.String `tfsdk:"subdirectory"`
	Preallocation          types.String `tfsdk:"preallocation"`
	SnapshotsAsVolumeChain types.Bool   `tfsdk:"snapshot_as_volume_chain"`
}

func (m *CIFSStorageModel) GetStorageType() types.String {
	return types.StringValue("cifs")
}

func (m *CIFSStorageModel) toCreateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.CIFSStorageCreateRequest{}
	request.Type = m.GetStorageType().ValueStringPointer()

	if err := m.populateCreateFields(ctx, &request.DataStoreCommonImmutableFields, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.Server = m.Server.ValueStringPointer()
	request.Username = m.Username.ValueStringPointer()
	request.Password = m.Password.ValueStringPointer()
	request.Share = m.Share.ValueStringPointer()
	request.Domain = m.Domain.ValueStringPointer()
	request.Subdirectory = m.SubDirectory.ValueStringPointer()
	request.Preallocation = m.Preallocation.ValueStringPointer()
	request.SnapshotsAsVolumeChain = proxmox_types.CustomBool(m.SnapshotsAsVolumeChain.ValueBool())

	return request, nil
}

func (m *CIFSStorageModel) toUpdateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.CIFSStorageUpdateRequest{}

	if err := m.populateUpdateFields(ctx, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.Preallocation = m.Preallocation.ValueStringPointer()

	return request, nil
}

func (m *CIFSStorageModel) fromAPI(ctx context.Context, datastore *storage.DatastoreGetResponseData) error {
	if err := m.populateBaseFromAPI(ctx, datastore); err != nil {
		return err
	}

	if datastore.Server != nil {
		m.Server = types.StringValue(*datastore.Server)
	}
	if datastore.Username != nil {
		m.Username = types.StringValue(*datastore.Username)
	}
	if datastore.Share != nil {
		m.Share = types.StringValue(*datastore.Share)
	}
	if datastore.Domain != nil {
		m.Domain = types.StringValue(*datastore.Domain)
	}
	if datastore.SubDirectory != nil {
		m.SubDirectory = types.StringValue(*datastore.SubDirectory)
	}
	if datastore.Preallocation != nil {
		m.Preallocation = types.StringValue(*datastore.Preallocation)
	}

	return nil
}
