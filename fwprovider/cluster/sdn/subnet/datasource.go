/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package subnet

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
)

var _ datasource.DataSource = &DataSource{}

var _ datasource.DataSourceWithConfigure = &DataSource{}

type DataSource struct {
	client *cluster.Client
}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

func (d *DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdn_subnet"
}

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.DataSource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Configuration",
			fmt.Sprintf("Expected config.DataSource, got: %T", req.ProviderData),
		)

		return
	}

	d.client = cfg.Client.Cluster()
}

func (d *DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve details about a specific SDN Subnet in Proxmox VE.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The full ID in the format 'vnet-id/subnet-id'.",
				Computed:    true,
			},
			"cidr": schema.StringAttribute{
				Description: "A CIDR network address, for example 10.0.0.0/8",
				Required:    true,
				CustomType:  customtypes.IPCIDRType{},
			},
			"vnet": schema.StringAttribute{
				Description: "The VNet this subnet belongs to.",
				Required:    true,
			},
			"dhcp_dns_server": schema.StringAttribute{
				Description: "The DNS server used for DHCP.",
				CustomType:  customtypes.IPAddrType{},
				Computed:    true,
			},
			"dhcp_range": schema.SingleNestedAttribute{
				Description: "DHCP range (start and end IPs).",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"start_address": schema.StringAttribute{
						CustomType:  customtypes.IPAddrType{},
						Computed:    true,
						Description: "Start of the DHCP range.",
					},
					"end_address": schema.StringAttribute{
						CustomType:  customtypes.IPAddrType{},
						Computed:    true,
						Description: "End of the DHCP range.",
					},
				},
			},
			"dns_zone_prefix": schema.StringAttribute{
				Description: "Prefix used for DNS zone delegation.",
				Computed:    true,
			},
			"gateway": schema.StringAttribute{
				Description: "The gateway address for the subnet.",
				CustomType:  customtypes.IPAddrType{},
				Computed:    true,
			},
			"snat": schema.BoolAttribute{
				Description: "Whether SNAT is enabled for the subnet.",
				Computed:    true,
			},
		},
	}
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := d.client.SDNVnets(data.VNet.ValueString()).Subnets()

	canonicalID, err := resolveCanonicalSubnetID(ctx, client, data.CIDR.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Resolve SDN Subnet ID", err.Error())
		return
	}

	subnet, err := client.GetSubnet(ctx, canonicalID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(
				"SDN Subnet Not Found",
				fmt.Sprintf("SDN Subnet with ID '%s' in VNet '%s' was not found",
					data.CIDR.ValueString(), data.VNet.ValueString()),
			)

			return
		}

		resp.Diagnostics.AddError("Unable to Read SDN Subnet", err.Error())

		return
	}

	state := &model{}
	state.CIDR = data.CIDR
	state.VNet = data.VNet

	if err := state.fromAPI(subnet); err != nil {
		resp.Diagnostics.AddError("Invalid Subnet Data", err.Error())
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s/%s", data.VNet.ValueString(), data.CIDR.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
