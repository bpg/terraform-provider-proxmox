/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
)

// --- EVPN ---

var (
	_ resource.Resource                = &evpnZoneShort{}
	_ resource.ResourceWithConfigure   = &evpnZoneShort{}
	_ resource.ResourceWithImportState = &evpnZoneShort{}
	_ resource.ResourceWithMoveState   = &evpnZoneShort{}
)

type evpnZoneShort struct{ *EVPNResource }

// NewEVPNShortResource creates the short-name alias proxmox_sdn_zone_evpn.
func NewEVPNShortResource() resource.Resource {
	inner := NewEVPNResource().(*EVPNResource)
	return &evpnZoneShort{EVPNResource: inner}
}

func (r *evpnZoneShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_zone_evpn"
}

func (r *evpnZoneShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.EVPNResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *evpnZoneShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.EVPNResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_sdn_zone_evpn", &schemaResp.Schema),
	}
}

// --- QinQ ---

var (
	_ resource.Resource                = &qinqZoneShort{}
	_ resource.ResourceWithConfigure   = &qinqZoneShort{}
	_ resource.ResourceWithImportState = &qinqZoneShort{}
	_ resource.ResourceWithMoveState   = &qinqZoneShort{}
)

type qinqZoneShort struct{ *QinQResource }

// NewQinQShortResource creates the short-name alias proxmox_sdn_zone_qinq.
func NewQinQShortResource() resource.Resource {
	inner := NewQinQResource().(*QinQResource)
	return &qinqZoneShort{QinQResource: inner}
}

func (r *qinqZoneShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_zone_qinq"
}

func (r *qinqZoneShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.QinQResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *qinqZoneShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.QinQResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_sdn_zone_qinq", &schemaResp.Schema),
	}
}

// --- Simple ---

var (
	_ resource.Resource                = &simpleZoneShort{}
	_ resource.ResourceWithConfigure   = &simpleZoneShort{}
	_ resource.ResourceWithImportState = &simpleZoneShort{}
	_ resource.ResourceWithMoveState   = &simpleZoneShort{}
)

type simpleZoneShort struct{ *SimpleResource }

// NewSimpleShortResource creates the short-name alias proxmox_sdn_zone_simple.
func NewSimpleShortResource() resource.Resource {
	inner := NewSimpleResource().(*SimpleResource)
	return &simpleZoneShort{SimpleResource: inner}
}

func (r *simpleZoneShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_zone_simple"
}

func (r *simpleZoneShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.SimpleResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *simpleZoneShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.SimpleResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_sdn_zone_simple", &schemaResp.Schema),
	}
}

// --- VLAN ---

var (
	_ resource.Resource                = &vlanZoneShort{}
	_ resource.ResourceWithConfigure   = &vlanZoneShort{}
	_ resource.ResourceWithImportState = &vlanZoneShort{}
	_ resource.ResourceWithMoveState   = &vlanZoneShort{}
)

type vlanZoneShort struct{ *VLANResource }

// NewVLANShortResource creates the short-name alias proxmox_sdn_zone_vlan.
func NewVLANShortResource() resource.Resource {
	inner := NewVLANResource().(*VLANResource)
	return &vlanZoneShort{VLANResource: inner}
}

func (r *vlanZoneShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_zone_vlan"
}

func (r *vlanZoneShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.VLANResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *vlanZoneShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.VLANResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_sdn_zone_vlan", &schemaResp.Schema),
	}
}

// --- VXLAN ---

var (
	_ resource.Resource                = &vxlanZoneShort{}
	_ resource.ResourceWithConfigure   = &vxlanZoneShort{}
	_ resource.ResourceWithImportState = &vxlanZoneShort{}
	_ resource.ResourceWithMoveState   = &vxlanZoneShort{}
)

type vxlanZoneShort struct{ *VXLANResource }

// NewVXLANShortResource creates the short-name alias proxmox_sdn_zone_vxlan.
func NewVXLANShortResource() resource.Resource {
	inner := NewVXLANResource().(*VXLANResource)
	return &vxlanZoneShort{VXLANResource: inner}
}

func (r *vxlanZoneShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_sdn_zone_vxlan"
}

func (r *vxlanZoneShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.VXLANResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *vxlanZoneShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse
	r.VXLANResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_sdn_zone_vxlan", &schemaResp.Schema),
	}
}
