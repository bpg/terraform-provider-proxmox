/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package controllers

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

const (
	TypeBGP    = "bgp"
	TypeEVPN   = "evpn"
	TypeFaucet = "faucet"
	TypeISIS   = "isis"
)

type Controller struct {
	ID    string  `json:"controller"      url:"controller"`
	Type  *string `json:"type,omitempty"  url:"type,omitempty"`
	State *string `json:"state,omitempty" url:"state,omitempty"`

	// EVPN.
	Fabric *string `json:"fabric,omitempty" url:"fabric,omitempty"`

	// BGP and EVPN.
	ASNumber *int64                          `json:"asn,omitempty"   url:"asn,omitempty"`
	Peers    *types.CustomCommaSeparatedList `json:"peers,omitempty" url:"peers,omitempty"`

	// BGP.
	BgpMultiPathAsRelax *types.CustomBool `json:"bgp-multipath-as-relax,omitempty" url:"bgp-multipath-as-relax,omitempty,int"`
	EBGP                *types.CustomBool `json:"ebgp,omitempty"                   url:"ebgp,omitempty,int"`
	EBPGMultiHop        *int64            `json:"ebgp-multihop,omitempty"          url:"ebgp-multihop,omitempty,int"`
	Loopback            *string           `json:"loopback,omitempty"               url:"loopback,omitempty"`

	// BGP and ISIS.
	Nodes *types.CustomCommaSeparatedList `json:"node,omitempty" url:"node,omitempty"`

	// ISIS.
	ISISDomain *string                         `json:"isis-domain,omitempty" url:"isis-domain,omitempty"`
	ISISIfaces *types.CustomCommaSeparatedList `json:"isis-ifaces,omitempty" url:"isis-ifaces,omitempty"`
	ISISNet    *string                         `json:"isis-net,omitempty"    url:"isis-net,omitempty"`
}

// ControllerData represents a controller with optional pending attribute.
type ControllerData struct {
	Controller

	Digest  *string     `json:"digest,omitempty"  url:"digest,omitempty"`
	Pending *Controller `json:"pending,omitempty" url:"pending,omitempty"`
}

// ControllerUpdate wraps a ControllerData struct with optional delete instructions.
type ControllerUpdate struct {
	Controller

	Delete []string `json:"delete,omitempty" url:"delete,omitempty"`
}
