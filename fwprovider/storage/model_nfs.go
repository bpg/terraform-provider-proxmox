/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"

	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NFSStorageModel maps the Terraform schema for NFS storage.
type NFSStorageModel struct {
	modelBase

	Server                 types.String `tfsdk:"server"`
	Export                 types.String `tfsdk:"export"`
	Options                types.String `tfsdk:"options"`
	Preallocation          types.String `tfsdk:"preallocation"`
	SnapshotsAsVolumeChain types.Bool   `tfsdk:"snapshot_as_volume_chain"`
	Backups                *BackupModel `tfsdk:"backups"`
}

func (m *NFSStorageModel) GetStorageType() types.String {
	return types.StringValue("nfs")
}

func (m *NFSStorageModel) toCreateAPIRequest(ctx context.Context) (any, error) {
	request := storage.NFSStorageCreateRequest{}
	request.Type = m.GetStorageType().ValueStringPointer()

	if err := m.populateCreateFields(ctx, &request.DataStoreCommonImmutableFields, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.Server = m.Server.ValueStringPointer()
	request.Export = m.Export.ValueStringPointer()
	request.Options = m.Options.ValueStringPointer()
	request.Preallocation = m.Preallocation.ValueStringPointer()
	request.SnapshotsAsVolumeChain = proxmoxtypes.CustomBool(m.SnapshotsAsVolumeChain.ValueBool())

	if m.Backups != nil {
		backups, err := m.Backups.toAPI()
		if err != nil {
			return nil, err
		}

		request.Backups = backups
	}

	return request, nil
}

func (m *NFSStorageModel) toUpdateAPIRequest(ctx context.Context) (any, error) {
	request := storage.NFSStorageUpdateRequest{}

	if err := m.populateUpdateFields(ctx, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.Options = m.Options.ValueStringPointer()

	if m.Backups != nil {
		backups, err := m.Backups.toAPI()
		if err != nil {
			return nil, err
		}

		request.Backups = backups
	}

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

	if datastore.SnapshotsAsVolumeChain != nil {
		m.SnapshotsAsVolumeChain = types.BoolValue(*datastore.SnapshotsAsVolumeChain.PointerBool())
	}

	if datastore.MaxProtectedBackups != nil || (datastore.PruneBackups != nil && *datastore.PruneBackups != "") {
		if m.Backups == nil {
			m.Backups = &BackupModel{}
		}

		if err := m.Backups.fromAPI(datastore.MaxProtectedBackups, datastore.PruneBackups); err != nil {
			return err
		}
	} else {
		m.Backups = nil
	}

	return nil
}
