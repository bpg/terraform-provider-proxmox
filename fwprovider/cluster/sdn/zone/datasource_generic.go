/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone

import (
	"context"
	"errors"
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
)

type zoneDataSourceConfig struct {
	typeNameSuffix string
	zoneType       string
	modelFunc      func() zoneModel
}

type genericZoneDataSource struct {
	client *zones.Client
	config zoneDataSourceConfig
}

func newGenericZoneDataSource(cfg zoneDataSourceConfig) *genericZoneDataSource {
	return &genericZoneDataSource{config: cfg}
}

func (d *genericZoneDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.config.typeNameSuffix
}

func (d *genericZoneDataSource) Configure(
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
			"Unexpected DataSource Configure Type",
			fmt.Sprintf(
				"Expected config.DataSource, got: %T",
				req.ProviderData,
			),
		)

		return
	}

	d.client = cfg.Client.Cluster().SDNZones()
}

func genericDataSourceAttributesWith(extraAttributes map[string]schema.Attribute) map[string]schema.Attribute {
	// Start with generic attributes as the base
	result := map[string]schema.Attribute{
		"dns": schema.StringAttribute{
			Computed:    true,
			Description: "DNS API server address.",
		},
		"dns_zone": schema.StringAttribute{
			Computed:    true,
			Description: "DNS domain name. The DNS zone must already exist on the DNS server.",
			MarkdownDescription: "DNS domain name. Used to register hostnames, such as `<hostname>.<domain>`. " +
				"The DNS zone must already exist on the DNS server.",
		},
		"id": schema.StringAttribute{
			Description: "The unique identifier of the SDN zone.",
			Required:    true,
		},
		"ipam": schema.StringAttribute{
			Computed:    true,
			Description: "IP Address Management system.",
		},
		"mtu": schema.Int64Attribute{
			Computed:    true,
			Description: "MTU value for the zone.",
		},
		"nodes": schema.SetAttribute{
			CustomType: stringset.Type{
				SetType: types.SetType{
					ElemType: types.StringType,
				},
			},
			Description: "The Proxmox nodes which the zone and associated VNets are deployed on",
			ElementType: types.StringType,
			Computed:    true,
		},
		"reverse_dns": schema.StringAttribute{
			Computed:    true,
			Description: "Reverse DNS API server address.",
		},
	}

	// Add extra attributes, allowing them to override generic ones if needed
	if extraAttributes != nil {
		maps.Copy(result, extraAttributes)
	}

	return result
}

func (d *genericZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	state := d.config.modelFunc()
	resp.Diagnostics.Append(req.Config.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	zone, err := d.client.GetZone(ctx, state.getID())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(
				"SDN Zone Not Found",
				fmt.Sprintf("SDN zone with ID '%s' was not found", state.getID()),
			)

			return
		}

		resp.Diagnostics.AddError(
			"Unable to Read SDN Zone",
			err.Error(),
		)

		return
	}

	// Verify the zone type matches what this datasource expects
	if zone.Type != nil && *zone.Type != d.config.zoneType {
		resp.Diagnostics.AddError(
			"SDN Zone Type Mismatch",
			fmt.Sprintf(
				"Expected zone type '%s' but found '%s' for zone '%s'",
				d.config.zoneType,
				*zone.Type,
				zone.ID,
			),
		)

		return
	}

	readModel := d.config.modelFunc()
	diags := &diag.Diagnostics{}
	readModel.importFromAPI(zone.ID, zone, diags)
	resp.Diagnostics.Append(*diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}
