/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vnet

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/vnets"
)

// Ensure the implementation satisfies the required interfaces.
var (
	_ datasource.DataSource              = &vnetsDataSource{}
	_ datasource.DataSourceWithConfigure = &vnetsDataSource{}
)

// vnetsDataSource is the data source implementation for SDN VNets.
type vnetsDataSource struct {
	client *vnets.Client
}

// vnetsDataSourceModel represents the data source model for listing VNets.
type vnetsDataSourceModel struct {
	VNets types.List `tfsdk:"vnets"`
}

// vnetDataModel represents individual VNet data in the list.
type vnetDataModel struct {
	ID           types.String `tfsdk:"id"`
	Zone         types.String `tfsdk:"zone"`
	Alias        types.String `tfsdk:"alias"`
	IsolatePorts types.Bool   `tfsdk:"isolate_ports"`
	Tag          types.Int64  `tfsdk:"tag"`
	VlanAware    types.Bool   `tfsdk:"vlan_aware"`
}

// Configure adds the provider-configured client to the data source.
func (d *vnetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.DataSource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected config.DataSource, got: %T", req.ProviderData),
		)

		return
	}

	d.client = &vnets.Client{Client: cfg.Client.Cluster()}
}

// Metadata returns the data source type name.
func (d *vnetsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdn_vnets"
}

// Schema defines the schema for the data source.
func (d *vnetsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about all SDN VNets in Proxmox.",
		MarkdownDescription: "Retrieves information about all SDN VNets in Proxmox. " +
			"This data source lists all virtual networks configured in the Software-Defined Networking setup.",
		Attributes: map[string]schema.Attribute{
			"vnets": schema.ListAttribute{
				Description: "List of SDN VNets.",
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":            types.StringType,
						"zone":          types.StringType,
						"alias":         types.StringType,
						"isolate_ports": types.BoolType,
						"tag":           types.Int64Type,
						"vlan_aware":    types.BoolType,
					},
				},
			},
		},
	}
}

// Read fetches all SDN VNets from the Proxmox VE API.
func (d *vnetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data vnetsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vnetsList, err := d.client.GetVnetsWithParams(ctx, &sdn.QueryParams{})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SDN VNets",
			err.Error(),
		)

		return
	}

	// Convert VNets to list elements
	vnetElements := make([]attr.Value, len(vnetsList))
	for i, vnet := range vnetsList {
		vnetData := vnetDataModel{
			ID:           types.StringValue(vnet.ID),
			Zone:         types.StringPointerValue(vnet.Zone),
			Alias:        types.StringPointerValue(vnet.Alias),
			IsolatePorts: types.BoolPointerValue(vnet.IsolatePorts.PointerBool()),
			Tag:          types.Int64PointerValue(vnet.Tag),
			VlanAware:    types.BoolPointerValue(vnet.VlanAware.PointerBool()),
		}

		objValue, objDiag := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"id":            types.StringType,
			"zone":          types.StringType,
			"alias":         types.StringType,
			"isolate_ports": types.BoolType,
			"tag":           types.Int64Type,
			"vlan_aware":    types.BoolType,
		}, vnetData)
		resp.Diagnostics.Append(objDiag...)

		if resp.Diagnostics.HasError() {
			return
		}

		vnetElements[i] = objValue
	}

	listValue, listDiag := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":            types.StringType,
			"zone":          types.StringType,
			"alias":         types.StringType,
			"isolate_ports": types.BoolType,
			"tag":           types.Int64Type,
			"vlan_aware":    types.BoolType,
		},
	}, vnetElements)
	resp.Diagnostics.Append(listDiag...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.VNets = listValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// NewVNetsDataSource returns a new data source for SDN VNets.
func NewVNetsDataSource() datasource.DataSource {
	return &vnetsDataSource{}
}
