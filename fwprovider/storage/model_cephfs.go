/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
)

// CephFSStorageModel maps the Terraform schema for CephFS storage.
type CephFSStorageModel struct {
	modelBase
	modelDirOptions

	FSName  types.String `tfsdk:"fs_name"`
	Subdir  types.String `tfsdk:"subdir"`
	FUSE    types.Bool   `tfsdk:"fuse"`
	Options types.String `tfsdk:"options"`
	Backups *BackupModel `tfsdk:"backups"`
}

// GetStorageType returns the storage type identifier.
func (m *CephFSStorageModel) GetStorageType() types.String {
	return types.StringValue("cephfs")
}

// toCreateAPIRequest converts the Terraform model to a Proxmox API request body.
func (m *CephFSStorageModel) toCreateAPIRequest(ctx context.Context) (any, error) {
	request := storage.CephFSStorageCreateRequest{}
	request.Type = m.GetStorageType().ValueStringPointer()

	if err := m.populateCreateFields(ctx,
		&request.DataStoreCommonImmutableFields,
		&request.DataStoreCommonMutableFields,
	); err != nil {
		return nil, err
	}

	request.FSName = attribute.StringPtrFromValue(m.FSName)
	request.Subdir = attribute.StringPtrFromValue(m.Subdir)
	request.FUSE = attribute.CustomBoolPtrFromValue(m.FUSE)
	request.Options = attribute.StringPtrFromValue(m.Options)
	request.CreateBasePath = attribute.CustomBoolPtrFromValue(m.CreateBasePath)
	request.CreateSubdirs = attribute.CustomBoolPtrFromValue(m.CreateSubdirs)

	if m.Backups != nil {
		backups, err := m.Backups.toAPI()
		if err != nil {
			return nil, err
		}

		request.Backups = backups
	}

	return request, nil
}

// toUpdateAPIRequest converts the Terraform model to a Proxmox API request body for updates.
func (m *CephFSStorageModel) toUpdateAPIRequest(ctx context.Context) (any, error) {
	request := storage.CephFSStorageUpdateRequest{}

	if err := m.populateUpdateFields(ctx, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.Options = attribute.StringPtrFromValue(m.Options)
	request.CreateBasePath = attribute.CustomBoolPtrFromValue(m.CreateBasePath)
	request.CreateSubdirs = attribute.CustomBoolPtrFromValue(m.CreateSubdirs)

	if m.Backups != nil {
		backups, err := m.Backups.toAPI()
		if err != nil {
			return nil, err
		}

		request.Backups = backups
	}

	return request, nil
}

// fromAPI populates the Terraform model from a Proxmox API response.
func (m *CephFSStorageModel) fromAPI(ctx context.Context, datastore *storage.DatastoreGetResponseData) error {
	if err := m.populateBaseFromAPI(ctx, datastore); err != nil {
		return err
	}

	m.populateFromAPI(datastore)

	m.FSName = types.StringPointerValue(datastore.FSName)
	m.Subdir = types.StringPointerValue(datastore.SubDirectory)
	m.FUSE = attribute.BoolValueFromCustomBoolPtr(datastore.FUSE)
	m.Options = types.StringPointerValue(datastore.Options)

	// only populate backups if user has configured it to avoid "was absent, but now present" error
	if m.Backups != nil {
		if err := m.Backups.fromAPI(datastore.MaxProtectedBackups, datastore.PruneBackups); err != nil {
			return err
		}
	}

	return nil
}
