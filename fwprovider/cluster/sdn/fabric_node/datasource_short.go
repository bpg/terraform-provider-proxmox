/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric_node

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// --- OpenFabric Node ---

var (
	_ datasource.DataSource              = &openFabricNodeDSShort{}
	_ datasource.DataSourceWithConfigure = &openFabricNodeDSShort{}
)

type openFabricNodeDSShort struct{ *OpenFabricDataSource }

// NewOpenFabricShortDataSource creates the short-name alias proxmox_sdn_fabric_node_openfabric (data source).
func NewOpenFabricShortDataSource() datasource.DataSource {
	inner := NewOpenFabricDataSource().(*OpenFabricDataSource)
	return &openFabricNodeDSShort{OpenFabricDataSource: inner}
}

func (d *openFabricNodeDSShort) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "proxmox_sdn_fabric_node_openfabric"
}

func (d *openFabricNodeDSShort) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	d.OpenFabricDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

// --- OSPF Node ---

var (
	_ datasource.DataSource              = &ospfFabricNodeDSShort{}
	_ datasource.DataSourceWithConfigure = &ospfFabricNodeDSShort{}
)

type ospfFabricNodeDSShort struct{ *OSPFDataSource }

// NewOSPFShortDataSource creates the short-name alias proxmox_sdn_fabric_node_ospf (data source).
func NewOSPFShortDataSource() datasource.DataSource {
	inner := NewOSPFDataSource().(*OSPFDataSource)
	return &ospfFabricNodeDSShort{OSPFDataSource: inner}
}

func (d *ospfFabricNodeDSShort) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "proxmox_sdn_fabric_node_ospf"
}

func (d *ospfFabricNodeDSShort) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	d.OSPFDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}
