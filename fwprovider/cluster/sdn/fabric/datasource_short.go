/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// --- OpenFabric ---

var (
	_ datasource.DataSource              = &openFabricDSShort{}
	_ datasource.DataSourceWithConfigure = &openFabricDSShort{}
)

type openFabricDSShort struct{ *OpenFabricDataSource }

// NewOpenFabricShortDataSource creates the short-name alias proxmox_sdn_fabric_openfabric (data source).
func NewOpenFabricShortDataSource() datasource.DataSource {
	inner := NewOpenFabricDataSource().(*OpenFabricDataSource)
	return &openFabricDSShort{OpenFabricDataSource: inner}
}

func (d *openFabricDSShort) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_fabric_openfabric"
}

func (d *openFabricDSShort) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.OpenFabricDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

// --- OSPF ---

var (
	_ datasource.DataSource              = &ospfFabricDSShort{}
	_ datasource.DataSourceWithConfigure = &ospfFabricDSShort{}
)

type ospfFabricDSShort struct{ *OSPFDataSource }

// NewOSPFShortDataSource creates the short-name alias proxmox_sdn_fabric_ospf (data source).
func NewOSPFShortDataSource() datasource.DataSource {
	inner := NewOSPFDataSource().(*OSPFDataSource)
	return &ospfFabricDSShort{OSPFDataSource: inner}
}

func (d *ospfFabricDSShort) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_fabric_ospf"
}

func (d *ospfFabricDSShort) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.OSPFDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}
