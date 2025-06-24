package sdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
)

var _ datasource.DataSource = &sdnZoneDataSource{}

var _ datasource.DataSourceWithConfigure = &sdnZoneDataSource{}

type sdnZoneDataSource struct {
	client *zones.Client
}

func NewSDNZoneDataSource() datasource.DataSource {
	return &sdnZoneDataSource{}
}

func (d *sdnZoneDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_sdn_zone"
}

func (d *sdnZoneDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.DataSource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Configuration",
			fmt.Sprintf("Expected config.DataSource but got: %T", req.ProviderData),
		)

		return
	}

	d.client = cfg.Client.Cluster().SDNZones()
}

func (d *sdnZoneDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch a Proxmox SDN Zone by name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the SDN zone.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name (ID) of the SDN zone.",
			},
			"type":                       schema.StringAttribute{Computed: true},
			"ipam":                       schema.StringAttribute{Computed: true},
			"dns":                        schema.StringAttribute{Computed: true},
			"reversedns":                 schema.StringAttribute{Computed: true},
			"dns_zone":                   schema.StringAttribute{Computed: true},
			"nodes":                      schema.StringAttribute{Computed: true},
			"mtu":                        schema.Int64Attribute{Computed: true},
			"bridge":                     schema.StringAttribute{Computed: true},
			"tag":                        schema.Int64Attribute{Computed: true},
			"vlan_protocol":              schema.StringAttribute{Computed: true},
			"peers":                      schema.StringAttribute{Computed: true},
			"controller":                 schema.StringAttribute{Computed: true},
			"vrf_vxlan":                  schema.Int64Attribute{Computed: true},
			"exit_nodes":                 schema.StringAttribute{Computed: true},
			"primary_exit_node":          schema.StringAttribute{Computed: true},
			"exit_nodes_local_routing":   schema.BoolAttribute{Computed: true},
			"advertise_subnets":          schema.BoolAttribute{Computed: true},
			"disable_arp_nd_suppression": schema.BoolAttribute{Computed: true},
			"rt_import":                  schema.StringAttribute{Computed: true},
		},
	}
}

func (d *sdnZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data sdnZoneModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	zone, err := d.client.GetZone(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to fetch SDN Zone", err.Error())
		return
	}

	readModel := &sdnZoneModel{}
	readModel.importFromAPI(zone.ID, zone)
	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}
