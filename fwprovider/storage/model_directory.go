package storage

import (
	"context"
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	proxmox_types "github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DirectoryStorageModel maps the Terraform schema for directory storage.
type DirectoryStorageModel struct {
	ID            types.String `tfsdk:"id" json:"storage"`
	Type          types.String `tfsdk:"type" json:"type"`
	Path          types.String `tfsdk:"path" json:"path"`
	Nodes         types.Set    `tfsdk:"nodes" json:"nodes"`
	ContentTypes  types.Set    `tfsdk:"content" json:"content"`
	Disable       types.Bool   `tfsdk:"disable" json:"disable"`
	Shared        types.Bool   `tfsdk:"shared" json:"shared"`
	Preallocation types.String `tfsdk:"preallocation" json:"preallocation"`
}

// toCreateAPIRequest converts the Terraform model to a Proxmox API request body.
func (m *DirectoryStorageModel) toCreateAPIRequest(ctx context.Context) (storage.DirectoryStorageCreateRequest, error) {
	storageType := "dir"
	request := storage.DirectoryStorageCreateRequest{}

	nodes := proxmox_types.CustomCommaSeparatedList{}
	diags := m.Nodes.ElementsAs(ctx, &nodes, false)
	if diags.HasError() {
		return request, fmt.Errorf("cannot convert nodes to directory storage: %s", diags)
	}
	contentTypes := proxmox_types.CustomCommaSeparatedList{}
	diags = m.ContentTypes.ElementsAs(ctx, &contentTypes, false)
	if diags.HasError() {
		return request, fmt.Errorf("cannot convert content-types to directory storage: %s", diags)
	}

	request.ID = m.ID.ValueStringPointer()
	request.Type = &storageType
	request.Nodes = &nodes
	request.ContentTypes = &contentTypes
	request.Disable = proxmox_types.CustomBoolPtr(m.Disable.ValueBoolPointer())
	request.Shared = proxmox_types.CustomBoolPtr(m.Shared.ValueBoolPointer())
	request.Path = m.Path.ValueStringPointer()

	return request, nil
}

func (m *DirectoryStorageModel) toUpdateAPIRequest(ctx context.Context) (storage.DirectoryStorageUpdateRequest, error) {
	request := storage.DirectoryStorageUpdateRequest{}

	nodes := proxmox_types.CustomCommaSeparatedList{}
	diags := m.Nodes.ElementsAs(ctx, &nodes, false)
	if diags.HasError() {
		return request, fmt.Errorf("cannot convert nodes to directory storage: %s", diags)
	}

	contentTypes := proxmox_types.CustomCommaSeparatedList{}
	diags = m.ContentTypes.ElementsAs(ctx, &contentTypes, false)
	if diags.HasError() {
		return request, fmt.Errorf("cannot convert content-types to directory storage: %s", diags)
	}

	request.ContentTypes = &contentTypes
	request.Nodes = &nodes
	request.Disable = proxmox_types.CustomBoolPtr(m.Disable.ValueBoolPointer())
	request.Shared = proxmox_types.CustomBoolPtr(m.Shared.ValueBoolPointer())

	return request, nil
}

func (m *DirectoryStorageModel) importFromAPI(ctx context.Context, datastore storage.DatastoreGetResponseData) error {
	m.ID = types.StringValue(*datastore.ID)
	m.Type = types.StringValue(*datastore.Type)
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
