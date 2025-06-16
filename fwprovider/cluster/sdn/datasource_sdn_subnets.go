package sdn

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/subnets"
)

var (
	_ datasource.DataSource              = &sdnSubnetDataSource{}
	_ datasource.DataSourceWithConfigure = &sdnSubnetDataSource{}
)

type sdnSubnetDataSource struct {
	client *subnets.Client
}

func NewSDNSubnetDataSource() datasource.DataSource {
	return &sdnSubnetDataSource{}
}

func (d *sdnSubnetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdn_subnet"
}

func (d *sdnSubnetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = cfg.Client.Cluster().SDNSubnets()
}

func (d *sdnSubnetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve details about a specific SDN Subnet in Proxmox VE.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The full ID in the format 'vnet-id/subnet-id'.",
			},
			"subnet": schema.StringAttribute{
				Required: true,
			},
			"canonical_name": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
			"vnet": schema.StringAttribute{
				Required:    true,
				Description: "The VNet this subnet belongs to.",
			},
			"dhcp_dns_server": schema.StringAttribute{
				Computed:    true,
				Description: "The DNS server used for DHCP.",
			},
			"dhcp_range": schema.ListNestedAttribute{
				Optional:    false,
				Computed:    true,
				Description: "List of DHCP ranges (start and end IPs).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start_address": schema.StringAttribute{
							Computed:    true,
							Description: "Start of the DHCP range.",
						},
						"end_address": schema.StringAttribute{
							Computed:    true,
							Description: "End of the DHCP range.",
						},
					},
				},
			},
			"dnszoneprefix": schema.StringAttribute{
				Computed:    true,
				Description: "Prefix used for DNS zone delegation.",
			},
			"gateway": schema.StringAttribute{
				Computed:    true,
				Description: "The gateway address for the subnet.",
			},
			"snat": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether SNAT is enabled for the subnet.",
			},
		},
	}
}

func (d *sdnSubnetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config sdnSubnetModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnet, err := d.client.GetSubnet(ctx, config.Vnet.ValueString(), config.Subnet.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError("Subnet not found", err.Error())
			return
		}
		resp.Diagnostics.AddError("Failed to retrieve subnet", err.Error())
		return
	}

	// Set the state
	state := &sdnSubnetModel{}
	state.Subnet = config.Subnet
	state.Vnet = config.Vnet
	state.importFromAPI(config.Subnet.ValueString(), subnet)

	// Set canonical name and ID (both = user-supplied subnet)
	state.ID = config.Subnet
	state.CanonicalName = config.Subnet

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
