/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"github.com/bpg/terraform-provider-proxmox/internal/types"
)

// NetworkInterfaceListResponseBody contains the body from a node network interface list response.
type NetworkInterfaceListResponseBody struct {
	Data []*NetworkInterfaceListResponseData `json:"data,omitempty"`
}

// NetworkInterfaceListResponseData contains the data from a node network interface list response.
type NetworkInterfaceListResponseData struct {
	Active          *types.CustomBool `json:"active,omitempty"`
	Address         *string           `json:"address,omitempty"`
	Address6        *string           `json:"address6,omitempty"`
	Autostart       *types.CustomBool `json:"autostart,omitempty"`
	BridgeFD        *string           `json:"bridge_fd,omitempty"`
	BridgePorts     *string           `json:"bridge_ports,omitempty"`
	BridgeSTP       *string           `json:"bridge_stp,omitempty"`
	BridgeVIDs      *string           `json:"bridge_vids,omitempty"`
	BridgeVLANAware *types.CustomBool `json:"bridge_vlan_aware,omitempty"`
	CIDR            *string           `json:"cidr,omitempty"`
	CIDR6           *string           `json:"cidr6,omitempty"`
	Comments        *string           `json:"comments,omitempty"`
	Exists          *types.CustomBool `json:"exists,omitempty"`
	Families        *[]string         `json:"families,omitempty"`
	Gateway         *string           `json:"gateway,omitempty"`
	Gateway6        *string           `json:"gateway6,omitempty"`
	Iface           string            `json:"iface"`
	MethodIPv4      *string           `json:"method,omitempty"`
	MethodIPv6      *string           `json:"method6,omitempty"`
	Netmask         *string           `json:"netmask,omitempty"`
	Priority        int               `json:"priority"`
	Type            string            `json:"type"`
}

// NetworkInterfaceCreateUpdateRequestBody contains the body for a node network interface create / update request.
type NetworkInterfaceCreateUpdateRequestBody struct {
	Iface string `json:"iface" url:"iface"`
	Type  string `json:"type"  url:"type"`

	Address            *string           `json:"address,omitempty"               url:"address,omitempty"`
	Address6           *string           `json:"address6,omitempty"              url:"address6,omitempty"`
	Autostart          *types.CustomBool `json:"autostart,omitempty"             url:"autostart,omitempty,int"`
	BondPrimary        *string           `json:"bond-primary,omitempty"          url:"bond-primary,omitempty"`
	BondMode           *string           `json:"bond_mode,omitempty"             url:"bond_mode,omitempty"`
	BondXmitHashPolicy *string           `json:"bond_xmit_hash_policy,omitempty" url:"bond_xmit_hash_policy,omitempty"`
	BridgePorts        *string           `json:"bridge_ports,omitempty"          url:"bridge_ports,omitempty"`
	BridgeVLANAware    *types.CustomBool `json:"bridge_vlan_aware,omitempty"     url:"bridge_vlan_aware,omitempty,int"`
	CIDR               *string           `json:"cidr,omitempty"                  url:"cidr,omitempty"`
	CIDR6              *string           `json:"cidr6,omitempty"                 url:"cidr6,omitempty"`
	Comments           *string           `json:"comments,omitempty"              url:"comments,omitempty"`
	Comments6          *string           `json:"comments6,omitempty"             url:"comments6,omitempty"`
	Gateway            *string           `json:"gateway,omitempty"               url:"gateway,omitempty"`
	Gateway6           *string           `json:"gateway6,omitempty"              url:"gateway6,omitempty"`
	MTU                *int              `json:"mtu,omitempty"                   url:"mtu,omitempty"`
	Netmask            *string           `json:"netmask,omitempty"               url:"netmask,omitempty"`
	Netmask6           *string           `json:"netmask6,omitempty"              url:"netmask6,omitempty"`
	OVSBonds           *string           `json:"ovs_bonds,omitempty"             url:"ovs_bonds,omitempty"`
	OVSBridge          *string           `json:"ovs_bridge,omitempty"            url:"ovs_bridge,omitempty"`
	OVSOptions         *string           `json:"ovs_options,omitempty"           url:"ovs_options,omitempty"`
	OVSPorts           *string           `json:"ovs_ports,omitempty"             url:"ovs_ports,omitempty"`
	OVSTag             *string           `json:"ovs_tag,omitempty"               url:"ovs_tag,omitempty"`
	Slaves             *string           `json:"slaves,omitempty"                url:"slaves,omitempty"`
	VLANID             *int              `json:"vlan_id,omitempty"               url:"vlan_id,omitempty"`
	VLANRawDevice      *string           `json:"vlan_raw_device,omitempty"       url:"vlan_raw_device,omitempty"`
}
