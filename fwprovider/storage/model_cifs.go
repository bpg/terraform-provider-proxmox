/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// CIFSStorageModel maps the Terraform schema for CIFS storage.
type CIFSStorageModel struct {
	modelBase

	Server                 types.String `tfsdk:"server"`
	Username               types.String `tfsdk:"username"`
	Password               types.String `tfsdk:"password"`
	Share                  types.String `tfsdk:"share"`
	Domain                 types.String `tfsdk:"domain"`
	SubDirectory           types.String `tfsdk:"subdirectory"`
	Preallocation          types.String `tfsdk:"preallocation"`
	SnapshotsAsVolumeChain types.Bool   `tfsdk:"snapshot_as_volume_chain"`
	Backups                *BackupModel `tfsdk:"backups"`
}

func (m *CIFSStorageModel) GetStorageType() types.String {
	return types.StringValue("cifs")
}

func (m *CIFSStorageModel) toCreateAPIRequest(ctx context.Context) (any, error) {
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

func (m *CIFSStorageModel) toUpdateAPIRequest(ctx context.Context) (any, error) {
	request := storage.CIFSStorageUpdateRequest{}

	if err := m.populateUpdateFields(ctx, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.Preallocation = m.Preallocation.ValueStringPointer()

	if m.Backups != nil {
		backups, err := m.Backups.toAPI()
		if err != nil {
			return nil, err
		}

		request.Backups = backups
	}

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

	if datastore.SnapshotsAsVolumeChain != nil {
		m.SnapshotsAsVolumeChain = types.BoolValue(*datastore.SnapshotsAsVolumeChain.PointerBool())
	}

	// only populate backups if user has configured it to avoid "was absent, but now present" error
	if m.Backups != nil {
		if err := m.Backups.fromAPI(datastore.MaxProtectedBackups, datastore.PruneBackups); err != nil {
			return err
		}
	}

	return nil
}
