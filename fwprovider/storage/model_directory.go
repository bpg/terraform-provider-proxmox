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

func (m *DirectoryStorageModel) toCreateAPIRequest(ctx context.Context) (any, error) {
	request := storage.DirectoryStorageCreateRequest{}
	request.Type = m.GetStorageType().ValueStringPointer()

	if err := m.populateCreateFields(ctx, &request.DataStoreCommonImmutableFields, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.Path = m.Path.ValueStringPointer()
	request.Preallocation = m.Preallocation.ValueStringPointer()
	request.Shared = proxmoxtypes.CustomBoolPtr(m.Shared.ValueBoolPointer())

	if m.Backups != nil {
		backups, err := m.Backups.toAPI()
		if err != nil {
			return nil, err
		}

		request.Backups = backups
	}

	return request, nil
}

func (m *DirectoryStorageModel) toUpdateAPIRequest(ctx context.Context) (any, error) {
	request := storage.DirectoryStorageUpdateRequest{}

	if err := m.populateUpdateFields(ctx, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.Preallocation = m.Preallocation.ValueStringPointer()
	request.Shared = proxmoxtypes.CustomBoolPtr(m.Shared.ValueBoolPointer())

	if m.Backups != nil {
		backups, err := m.Backups.toAPI()
		if err != nil {
			return nil, err
		}

		request.Backups = backups
	}

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
