package sdn

/*
VNET MODEL TERRAFORM
*/

import (
	"github.com/bpg/terraform-provider-proxmox/fwprovider/helpers/ptrConversion"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/vnets"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type sdnVnetModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Zone         types.String `tfsdk:"zone"`
	Alias        types.String `tfsdk:"alias"`
	IsolatePorts types.Bool   `tfsdk:"isolate_ports"`
	Tag          types.Int64  `tfsdk:"tag"`
	Type         types.String `tfsdk:"type"`
	VlanAware    types.Bool   `tfsdk:"vlanaware"`
	ZoneType     types.String `tfsdk:"zonetype"`
}

func (m *sdnVnetModel) importFromAPI(name string, data *vnets.VnetData) {
	m.ID = types.StringValue(name)
	m.Name = types.StringValue(name)

	m.Zone = types.StringPointerValue(data.Zone)
	m.Alias = types.StringPointerValue(data.Alias)
	m.IsolatePorts = types.BoolPointerValue(ptrConversion.Int64ToBoolPtr(data.IsolatePorts))
	m.Tag = types.Int64PointerValue(data.Tag)
	m.Type = types.StringPointerValue(data.Type)
	m.VlanAware = types.BoolPointerValue(ptrConversion.Int64ToBoolPtr(data.VlanAware))
}

func (m *sdnVnetModel) toAPIRequestBody() *vnets.VnetRequestData {
	data := &vnets.VnetRequestData{}

	data.ID = m.Name.ValueString()

	data.Zone = m.Zone.ValueStringPointer()
	data.Alias = m.Alias.ValueStringPointer()
	data.IsolatePorts = ptrConversion.BoolToInt64Ptr(m.IsolatePorts.ValueBoolPointer())
	data.Tag = m.Tag.ValueInt64Pointer()
	data.Type = m.Type.ValueStringPointer()
	data.VlanAware = ptrConversion.BoolToInt64Ptr(m.VlanAware.ValueBoolPointer())

	return data
}
