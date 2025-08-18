package storage

import (
	"context"
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/storage"
	proxmox_types "github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NFSStorageModel maps the Terraform schema for NFS storage.
type NFSStorageModel struct {
	ID                     types.String `tfsdk:"id" json:"storage"`
	Type                   types.String `tfsdk:"type" json:"type"`
	Nodes                  types.Set    `tfsdk:"nodes" json:"nodes"`
	ContentTypes           types.Set    `tfsdk:"content" json:"content"`
	Disable                types.Bool   `tfsdk:"disable" json:"disable"`
	Server                 types.String `tfsdk:"server" json:"server"`
	Export                 types.String `tfsdk:"export"  json:"export"`
	Options                types.String `tfsdk:"options" json:"options"`
	Preallocation          types.String `tfsdk:"preallocation" json:"preallocation"`
	SnapshotsAsVolumeChain types.Bool   `tfsdk:"snapshot_as_volume_chain" json:"snapshot-as-volume-chain"`
}

func (m *NFSStorageModel) GetID() types.String {
	return m.ID
}

func (m *NFSStorageModel) GetStorageType() types.String {
	return types.StringValue("nfs")
}

func (m *NFSStorageModel) toCreateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.NFSStorageCreateRequest{}

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
	request.Type = m.GetStorageType().ValueStringPointer()
	request.Nodes = &nodes
	request.ContentTypes = &contentTypes
	request.Disable = proxmox_types.CustomBoolPtr(m.Disable.ValueBoolPointer())
	request.Server = m.Server.ValueStringPointer()
	request.Export = m.Export.ValueStringPointer()
	request.Options = m.Options.ValueStringPointer()

	return request, nil
}

func (m *NFSStorageModel) toUpdateAPIRequest(ctx context.Context) (interface{}, error) {
	request := storage.NFSStorageUpdateRequest{}

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

	request.Nodes = &nodes
	request.ContentTypes = &contentTypes
	request.Disable = proxmox_types.CustomBoolPtr(m.Disable.ValueBoolPointer())
	request.Options = m.Options.ValueStringPointer()

	return request, nil
}

func (m *NFSStorageModel) fromAPI(ctx context.Context, datastore *storage.DatastoreGetResponseData) error {
	m.ID = types.StringValue(*datastore.ID)
	m.Type = m.GetStorageType()
	if datastore.ContentTypes != nil {
		contentTypes, diags := types.SetValueFrom(ctx, types.StringType, *datastore.ContentTypes)
		if diags.HasError() {
			return fmt.Errorf("cannot parse content from datastore: %s", diags)
		}
		m.ContentTypes = contentTypes
	}
	if datastore.Nodes != nil {
		nodes, diags := types.SetValueFrom(ctx, types.StringType, *datastore.Nodes)
		if diags.HasError() {
			return fmt.Errorf("cannot parse nodes from datastore: %s", diags)
		}
		m.Nodes = nodes
	} else {
		m.Nodes = types.SetValueMust(types.StringType, []attr.Value{})
	}
	if datastore.Disable != nil {
		m.Disable = datastore.Disable.ToValue()
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

	return nil
}
