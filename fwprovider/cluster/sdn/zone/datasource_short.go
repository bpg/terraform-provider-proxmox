/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// --- EVPN ---

var (
	_ datasource.DataSource              = &evpnZoneDSShort{}
	_ datasource.DataSourceWithConfigure = &evpnZoneDSShort{}
)

type evpnZoneDSShort struct{ *EVPNDataSource }

// NewEVPNShortDataSource creates the short-name alias proxmox_sdn_zone_evpn (data source).
func NewEVPNShortDataSource() datasource.DataSource {
	inner := NewEVPNDataSource().(*EVPNDataSource)
	return &evpnZoneDSShort{EVPNDataSource: inner}
}

func (d *evpnZoneDSShort) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_zone_evpn"
}

func (d *evpnZoneDSShort) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.EVPNDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

// --- QinQ ---

var (
	_ datasource.DataSource              = &qinqZoneDSShort{}
	_ datasource.DataSourceWithConfigure = &qinqZoneDSShort{}
)

type qinqZoneDSShort struct{ *QinQDataSource }

// NewQinQShortDataSource creates the short-name alias proxmox_sdn_zone_qinq (data source).
func NewQinQShortDataSource() datasource.DataSource {
	inner := NewQinQDataSource().(*QinQDataSource)
	return &qinqZoneDSShort{QinQDataSource: inner}
}

func (d *qinqZoneDSShort) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_zone_qinq"
}

func (d *qinqZoneDSShort) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.QinQDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

// --- Simple ---

var (
	_ datasource.DataSource              = &simpleZoneDSShort{}
	_ datasource.DataSourceWithConfigure = &simpleZoneDSShort{}
)

type simpleZoneDSShort struct{ *SimpleDataSource }

// NewSimpleShortDataSource creates the short-name alias proxmox_sdn_zone_simple (data source).
func NewSimpleShortDataSource() datasource.DataSource {
	inner := NewSimpleDataSource().(*SimpleDataSource)
	return &simpleZoneDSShort{SimpleDataSource: inner}
}

func (d *simpleZoneDSShort) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_zone_simple"
}

func (d *simpleZoneDSShort) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.SimpleDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

// --- VLAN ---

var (
	_ datasource.DataSource              = &vlanZoneDSShort{}
	_ datasource.DataSourceWithConfigure = &vlanZoneDSShort{}
)

type vlanZoneDSShort struct{ *VLANDataSource }

// NewVLANShortDataSource creates the short-name alias proxmox_sdn_zone_vlan (data source).
func NewVLANShortDataSource() datasource.DataSource {
	inner := NewVLANDataSource().(*VLANDataSource)
	return &vlanZoneDSShort{VLANDataSource: inner}
}

func (d *vlanZoneDSShort) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_zone_vlan"
}

func (d *vlanZoneDSShort) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.VLANDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

// --- VXLAN ---

var (
	_ datasource.DataSource              = &vxlanZoneDSShort{}
	_ datasource.DataSourceWithConfigure = &vxlanZoneDSShort{}
)

type vxlanZoneDSShort struct{ *VXLANDataSource }

// NewVXLANShortDataSource creates the short-name alias proxmox_sdn_zone_vxlan (data source).
func NewVXLANShortDataSource() datasource.DataSource {
	inner := NewVXLANDataSource().(*VXLANDataSource)
	return &vxlanZoneDSShort{VXLANDataSource: inner}
}

func (d *vxlanZoneDSShort) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_zone_vxlan"
}

func (d *vxlanZoneDSShort) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.VXLANDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

// --- Zones (list) ---

var (
	_ datasource.DataSource              = &zonesDSShort{}
	_ datasource.DataSourceWithConfigure = &zonesDSShort{}
)

type zonesDSShort struct{ zonesDataSource }

// NewZonesShortDataSource creates the short-name alias proxmox_sdn_zones (data source).
func NewZonesShortDataSource() datasource.DataSource {
	return &zonesDSShort{}
}

func (d *zonesDSShort) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_zones"
}

func (d *zonesDSShort) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.zonesDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}
