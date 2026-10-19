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

// rbdDefaultPool is the pool Proxmox uses when an RBD storage names none.
const rbdDefaultPool = "rbd"

// RBDStorageModel maps the Terraform schema for RBD storage.
type RBDStorageModel struct {
	modelBase

	Pool      types.String `tfsdk:"pool"`
	DataPool  types.String `tfsdk:"data_pool"`
	Namespace types.String `tfsdk:"namespace"`
	KRBD      types.Bool   `tfsdk:"krbd"`
}

// GetStorageType returns the storage type identifier.
func (m *RBDStorageModel) GetStorageType() types.String {
	return types.StringValue("rbd")
}

// toCreateAPIRequest converts the Terraform model to a Proxmox API request body.
func (m *RBDStorageModel) toCreateAPIRequest(ctx context.Context) (any, error) {
	request := storage.RBDStorageCreateRequest{}
	request.Type = m.GetStorageType().ValueStringPointer()

	if err := m.populateCreateFields(ctx,
		&request.DataStoreCommonImmutableFields,
		&request.DataStoreCommonMutableFields,
	); err != nil {
		return nil, err
	}

	request.Pool = attribute.StringPtrFromValue(m.Pool)
	request.DataPool = attribute.StringPtrFromValue(m.DataPool)
	request.Namespace = attribute.StringPtrFromValue(m.Namespace)
	request.KRBD = attribute.CustomBoolPtrFromValue(m.KRBD)

	return request, nil
}

// toUpdateAPIRequest converts the Terraform model to a Proxmox API request body for updates.
func (m *RBDStorageModel) toUpdateAPIRequest(ctx context.Context) (any, error) {
	request := storage.RBDStorageUpdateRequest{}

	if err := m.populateUpdateFields(ctx, &request.DataStoreCommonMutableFields); err != nil {
		return nil, err
	}

	request.KRBD = attribute.CustomBoolPtrFromValue(m.KRBD)

	return request, nil
}

// fromAPI populates the Terraform model from a Proxmox API response.
func (m *RBDStorageModel) fromAPI(ctx context.Context, datastore *storage.DatastoreGetResponseData) error {
	if err := m.populateBaseFromAPI(ctx, datastore); err != nil {
		return err
	}

	if datastore.Pool != nil {
		m.Pool = types.StringValue(*datastore.Pool)
	} else {
		m.Pool = types.StringValue(rbdDefaultPool)
	}

	m.DataPool = types.StringPointerValue(datastore.DataPool)
	m.Namespace = types.StringPointerValue(datastore.Namespace)
	m.KRBD = attribute.BoolValueFromCustomBoolPtr(datastore.KRBD)

	return nil
}
