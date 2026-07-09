package disk

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type Model struct {
	Aio         types.String `tfsdk:"aio"`
	Backup      types.Bool   `tfsdk:"backup"`
	Cache       types.String `tfsdk:"cache"`
	DatastoreId types.String `tfsdk:"datastore_id"`
	Discard     types.String `tfsdk:"discard"`
	FileFormat  types.String `tfsdk:"file_format"`
	ImportFrom  types.String `tfsdk:"import_from"`
	IOThread    types.Bool   `tfsdk:"iothread"`
	Size        types.Int64  `tfsdk:"size"`
}

func attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"aio":          types.StringType,
		"backup":       types.BoolType,
		"cache":        types.StringType,
		"datastore_id": types.StringType,
		"discard":      types.StringType,
		"file_format":  types.StringType,
		"import_from":  types.StringType,
		"iothread":     types.BoolType,
		"size":         types.Int64Type,
	}
}

// NullValue returns a properly typed null Value.
func NullValue() Value { return types.MapNull(types.ObjectType{}.WithAttributeTypes(attributeTypes())) }

// toAPI writes the Disk-related fields onto the shared create/update body.
func (m *Model) toAPI() vms.CustomStorageDevice {
	return vms.CustomStorageDevice{
		AIO:         attribute.StringPtrFromValue(m.Aio),
		Backup:      attribute.CustomBoolPtrFromValue(m.Backup),
		Cache:       attribute.StringPtrFromValue(m.Cache),
		DatastoreID: attribute.StringPtrFromValue(m.DatastoreId),
		Discard:     attribute.StringPtrFromValue(m.Discard),
		Format:      attribute.StringPtrFromValue(m.FileFormat),
		ImportFrom:  attribute.StringPtrFromValue(m.ImportFrom),
		IOThread:    attribute.CustomBoolPtrFromValue(m.IOThread),
		Media:       new("disk"),
		Size:        proxmoxtypes.DiskSizeFromGigabytes(m.Size.ValueInt64()),
	}
}

// fromAPI populates the Model from the PVE API response.
func (m *Model) fromAPI(d vms.CustomStorageDevice) {
	m.Aio = types.StringPointerValue(d.AIO)
	m.Backup = types.BoolPointerValue(d.Backup.PointerBool())
	m.Cache = types.StringPointerValue(d.Cache)
	m.DatastoreId = types.StringPointerValue(d.DatastoreID)
	m.Discard = types.StringPointerValue(d.Discard)
	m.FileFormat = types.StringPointerValue(d.Format)
	m.ImportFrom = types.StringPointerValue(d.ImportFrom)
	m.IOThread = types.BoolPointerValue(d.IOThread.PointerBool())
	m.Size = types.Int64Value(d.Size.InGigabytes())
}
