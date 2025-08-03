/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zones

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

const (
	TypeSimple = "simple"
	TypeVLAN   = "vlan"
	TypeQinQ   = "qinq"
	TypeVXLAN  = "vxlan"
	TypeEVPN   = "evpn"
)

type ZoneData struct {
	ID         string  `json:"zone"                 url:"zone"`
	Type       *string `json:"type,omitempty"       url:"type,omitempty"`
	IPAM       *string `json:"ipam,omitempty"       url:"ipam,omitempty"`
	DNS        *string `json:"dns,omitempty"        url:"dns,omitempty"`
	ReverseDNS *string `json:"reversedns,omitempty" url:"reversedns,omitempty"`
	DNSZone    *string `json:"dnszone,omitempty"    url:"dnszone,omitempty"`
	Nodes      *string `json:"nodes,omitempty"      url:"nodes,omitempty"`
	MTU        *int64  `json:"mtu,omitempty"        url:"mtu,omitempty"`

	// VLAN.
	Bridge *string `json:"bridge,omitempty" url:"bridge,omitempty"`

	// QinQ.
	ServiceVLAN         *int64  `json:"tag,omitempty"           url:"tag,omitempty"`
	ServiceVLANProtocol *string `json:"vlan-protocol,omitempty" url:"vlan-protocol,omitempty"`

	// VXLAN.
	Peers *string `json:"peers,omitempty" url:"peers,omitempty"`

	// EVPN.
	Controller              *string           `json:"controller,omitempty"                 url:"controller,omitempty"`
	VRFVXLANID              *int64            `json:"vrf-vxlan,omitempty"                  url:"vrf-vxlan,omitempty"`
	ExitNodes               *string           `json:"exitnodes,omitempty"                  url:"exitnodes,omitempty"`
	ExitNodesPrimary        *string           `json:"exitnodes-primary,omitempty"          url:"exitnodes-primary,omitempty"`
	ExitNodesLocalRouting   *types.CustomBool `json:"exitnodes-local-routing,omitempty"    url:"exitnodes-local-routing,omitempty,int"`
	AdvertiseSubnets        *types.CustomBool `json:"advertise-subnets,omitempty"          url:"advertise-subnets,omitempty,int"`
	DisableARPNDSuppression *types.CustomBool `json:"disable-arp-nd-suppression,omitempty" url:"disable-arp-nd-suppression,omitempty,int"`
	RouteTargetImport       *string           `json:"rt-import,omitempty"                  url:"rt-import,omitempty"`
}

// ZoneRequestData wraps a ZoneData struct with optional delete instructions.
type ZoneRequestData struct {
	ZoneData

	Delete []string `url:"delete,omitempty"`
}

// ZoneResponseBody represents the response for a single zone.
type ZoneResponseBody struct {
	Data *ZoneData `json:"data"`
}

// ZonesResponseBody represents the response for a list of zones.
type ZonesResponseBody struct {
	Data *[]ZoneData `json:"data"`
}
