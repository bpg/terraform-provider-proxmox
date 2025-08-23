package storage

import (
	"context"
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	proxmox_types "github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StorageModelBase contains the common fields for all storage models.
type StorageModelBase struct {
	ID           types.String `tfsdk:"id"`
	Nodes        types.Set    `tfsdk:"nodes"`
	ContentTypes types.Set    `tfsdk:"content"`
	Disable      types.Bool   `tfsdk:"disable"`
	Shared       types.Bool   `tfsdk:"shared"`
}

// GetID returns the storage identifier from the base model.
func (m *StorageModelBase) GetID() types.String {
	return m.ID
}

// populateBaseFromAPI is a helper to populate the common fields from an API response.
func (m *StorageModelBase) populateBaseFromAPI(ctx context.Context, datastore *storage.DatastoreGetResponseData) error {
	m.ID = types.StringValue(*datastore.ID)

	if datastore.Nodes != nil {
		nodes, diags := types.SetValueFrom(ctx, types.StringType, *datastore.Nodes)
		if diags.HasError() {
			return fmt.Errorf("cannot parse nodes from datastore: %s", diags)
		}
		m.Nodes = nodes
	} else {
		m.Nodes = types.SetValueMust(types.StringType, []attr.Value{})
	}

	if datastore.ContentTypes != nil {
		contentTypes, diags := types.SetValueFrom(ctx, types.StringType, *datastore.ContentTypes)
		if diags.HasError() {
			return fmt.Errorf("cannot parse content from datastore: %s", diags)
		}
		m.ContentTypes = contentTypes
	}

	if datastore.Disable != nil {
		m.Disable = datastore.Disable.ToValue()
	}
	if datastore.Shared != nil {
		m.Shared = datastore.Shared.ToValue()
	}

	return nil
}

// populateCreateFields is a helper to populate the common fields for a create request.
func (m *StorageModelBase) populateCreateFields(ctx context.Context, immutableReq *storage.DataStoreCommonImmutableFields, mutableReq *storage.DataStoreCommonMutableFields) error {
	var nodes proxmox_types.CustomCommaSeparatedList
	if diags := m.Nodes.ElementsAs(ctx, &nodes, false); diags.HasError() {
		return fmt.Errorf("cannot convert nodes: %s", diags)
	}

	var contentTypes proxmox_types.CustomCommaSeparatedList
	if diags := m.ContentTypes.ElementsAs(ctx, &contentTypes, false); diags.HasError() {
		return fmt.Errorf("cannot convert content-types: %s", diags)
	}

	immutableReq.ID = m.ID.ValueStringPointer()
	mutableReq.Nodes = &nodes
	mutableReq.ContentTypes = &contentTypes
	mutableReq.Disable = proxmox_types.CustomBoolPtr(m.Disable.ValueBoolPointer())

	return nil
}

// populateUpdateFields is a helper to populate the common fields for an update request.
func (m *StorageModelBase) populateUpdateFields(ctx context.Context, mutableReq *storage.DataStoreCommonMutableFields) error {
	var nodes proxmox_types.CustomCommaSeparatedList
	if diags := m.Nodes.ElementsAs(ctx, &nodes, false); diags.HasError() {
		return fmt.Errorf("cannot convert nodes: %s", diags)
	}

	var contentTypes proxmox_types.CustomCommaSeparatedList
	if diags := m.ContentTypes.ElementsAs(ctx, &contentTypes, false); diags.HasError() {
		return fmt.Errorf("cannot convert content-types: %s", diags)
	}

	mutableReq.Nodes = &nodes
	mutableReq.ContentTypes = &contentTypes
	mutableReq.Disable = proxmox_types.CustomBoolPtr(m.Disable.ValueBoolPointer())

	return nil
}
