package sdn

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/vnets"
)

var (
	_ datasource.DataSource              = &sdnVnetDataSource{}
	_ datasource.DataSourceWithConfigure = &sdnVnetDataSource{}
)

type sdnVnetDataSource struct {
	client *vnets.Client
}

func NewSDNVnetDataSource() datasource.DataSource {
	return &sdnVnetDataSource{}
}

func (d *sdnVnetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdn_vnet"
}

func (d *sdnVnetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.DataSource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data",
			fmt.Sprintf("Expected config.DataSource, got: %T", req.ProviderData),
		)
		return
	}

	d.client = cfg.Client.Cluster().SDNVnets()
}

func (d *sdnVnetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about an existing SDN Vnet in Proxmox VE.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the vnet (usually the name).",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the vnet.",
			},
			"zone": schema.StringAttribute{
				Computed:    true,
				Description: "The zone associated with the vnet.",
			},
			"zonetype": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the zone associated with this vnet.",
			},
			"alias": schema.StringAttribute{
				Computed:    true,
				Description: "An alias for this vnet.",
			},
			"isolate_ports": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether ports are isolated.",
			},
			"tag": schema.Int64Attribute{
				Computed:    true,
				Description: "VLAN/VXLAN tag.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of the vnet.",
			},
			"vlanaware": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this vnet is VLAN aware.",
			},
		},
	}
}

func (d *sdnVnetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config sdnVnetModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vnetID := config.Name.ValueString()
	vnet, err := d.client.GetVnet(ctx, vnetID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError("Vnet not found", fmt.Sprintf("No vnet with ID %q exists", vnetID))
			return
		}
		resp.Diagnostics.AddError("Error retrieving vnet", err.Error())
		return
	}

	state := sdnVnetModel{}
	state.importFromAPI(vnetID, vnet)
	state.ID = types.StringValue(vnetID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
