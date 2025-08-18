package storage

import (
	"context"

	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	proxmox_types "github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NFSStorageModel maps the Terraform schema for NFS storage.
type NFSStorageModel struct {
	StorageModelBase
	Server                 types.String `tfsdk:"server"`
	Export                 types.String `tfsdk:"export"`
	Options                types.String `tfsdk:"options"`
	Preallocation          types.String `tfsdk:"preallocation"`
	SnapshotsAsVolumeChain types.Bool   `tfsdk:"snapshot_as_volume_chain"`
}

func (m *NFSStorageModel) GetStorageType() types.String {
	return types.StringValue("nfs")
}

func (m *NFSStorageModel) toCreateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.NFSStorageCreateRequest{}
	request.Type = m.GetStorageType().ValueStringPointer()

	if err := m.populateCreateFields(ctx, &request.DataStoreCommonImmutableFields, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.Server = m.Server.ValueStringPointer()
	request.Export = m.Export.ValueStringPointer()
	request.Options = m.Options.ValueStringPointer()
	request.Preallocation = m.Preallocation.ValueStringPointer()
	request.SnapshotsAsVolumeChain = proxmox_types.CustomBool(m.SnapshotsAsVolumeChain.ValueBool())

	return request, nil
}

func (m *NFSStorageModel) toUpdateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.NFSStorageUpdateRequest{}

	if err := m.populateUpdateFields(ctx, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.Options = m.Options.ValueStringPointer()

	return request, nil
}

func (m *NFSStorageModel) fromAPI(ctx context.Context, datastore *storage.DatastoreGetResponseData) error {
	if err := m.populateBaseFromAPI(ctx, datastore); err != nil {
		return err
	}

	if datastore.Server != nil {
		m.Server = types.StringValue(*datastore.Server)
	}
	if datastore.Export != nil {
		m.Export = types.StringValue(*datastore.Export)
	}
	if datastore.Options != nil {
		m.Options = types.StringValue(*datastore.Options)
	}
	if datastore.Preallocation != nil {
		m.Preallocation = types.StringValue(*datastore.Preallocation)
	}

	return nil
}
