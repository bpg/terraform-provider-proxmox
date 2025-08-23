package storage

import (
	"context"

	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DirectoryStorageModel maps the Terraform schema for directory storage.
type DirectoryStorageModel struct {
	StorageModelBase
	Path          types.String `tfsdk:"path"`
	Preallocation types.String `tfsdk:"preallocation"`
	Backups       *BackupModel `tfsdk:"backups"`
}

func (m *DirectoryStorageModel) GetStorageType() types.String {
	return types.StringValue("dir")
}

func (m *DirectoryStorageModel) toCreateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.DirectoryStorageCreateRequest{}
	request.Type = m.GetStorageType().ValueStringPointer()

	if err := m.populateCreateFields(ctx, &request.DataStoreCommonImmutableFields, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.Path = m.Path.ValueStringPointer()
	request.Preallocation = m.Preallocation.ValueStringPointer()

	return request, nil
}

func (m *DirectoryStorageModel) toUpdateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.DirectoryStorageUpdateRequest{}

	if err := m.populateUpdateFields(ctx, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.Preallocation = m.Preallocation.ValueStringPointer()

	return request, nil
}

func (m *DirectoryStorageModel) fromAPI(ctx context.Context, datastore *storage.DatastoreGetResponseData) error {
	if err := m.populateBaseFromAPI(ctx, datastore); err != nil {
		return err
	}

	if datastore.Path != nil {
		m.Path = types.StringValue(*datastore.Path)
	}
	if datastore.Preallocation != nil {
		m.Preallocation = types.StringValue(*datastore.Preallocation)
	}

	return nil
}
