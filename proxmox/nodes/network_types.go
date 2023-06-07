/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// NetworkDeviceListResponseBody contains the body from a node network device list response.
type NetworkDeviceListResponseBody struct {
	Data []*NetworkDeviceListResponseData `json:"data,omitempty"`
}

// NetworkDeviceListResponseData contains the data from a node network device list response.
type NetworkDeviceListResponseData struct {
	Active          *types.CustomBool `json:"active,omitempty"`
	Address         *string           `json:"address,omitempty"`
	Autostart       *types.CustomBool `json:"autostart,omitempty"`
	BridgeFD        *string           `json:"bridge_fd,omitempty"`
	BridgePorts     *string           `json:"bridge_ports,omitempty"`
	BridgeSTP       *string           `json:"bridge_stp,omitempty"`
	BridgeVIDs      *string           `json:"bridge_vids,omitempty"`
	BridgeVLANAWARE *string           `json:"bridge_vlan_aware,omitempty"`
	CIDR            *string           `json:"cidr,omitempty"`
	Comments        *string           `json:"comments,omitempty"`
	Exists          *types.CustomBool `json:"exists,omitempty"`
	Families        *[]string         `json:"families,omitempty"`
	Gateway         *string           `json:"gateway,omitempty"`
	Iface           string            `json:"iface"`
	MethodIPv4      *string           `json:"method,omitempty"`
	MethodIPv6      *string           `json:"method6,omitempty"`
	Netmask         *string           `json:"netmask,omitempty"`
	Priority        int               `json:"priority"`
	Type            string            `json:"type"`
}

// NetworkDeviceCreateUpdateRequestBody contains the body for a node network device create / update request.
type NetworkDeviceCreateUpdateRequestBody struct {
	Iface string `json:"iface"`
	Type  string `json:"type"`

	Address            *string           `json:"address,omitempty"`
	Address6           *string           `json:"address6,omitempty"`
	Autostart          *types.CustomBool `json:"autostart,omitempty"`
	BondPrimary        *string           `json:"bond-primary,omitempty"`
	BondMode           *string           `json:"bond_mode,omitempty"`
	BondXmitHashPolicy *string           `json:"bond_xmit_hash_policy,omitempty"`
	BridgePorts        *string           `json:"bridge_ports,omitempty"`
	BridgeVLANAware    *types.CustomBool `json:"bridge_vlan_aware,omitempty"`
	CIDR               *string           `json:"cidr,omitempty"`
	CIDR6              *string           `json:"cidr6,omitempty"`
	Comments           *string           `json:"comments,omitempty"`
	Comments6          *string           `json:"comments6,omitempty"`
	Gateway            *string           `json:"gateway,omitempty"`
	Gateway6           *string           `json:"gateway6,omitempty"`
	MTU                *int              `json:"mtu,omitempty"`
	Netmask            *string           `json:"netmask,omitempty"`
	Netmask6           *string           `json:"netmask6,omitempty"`
	OVSBonds           *string           `json:"ovs_bonds,omitempty"`
	OVSBridge          *string           `json:"ovs_bridge,omitempty"`
	OVSOptions         *string           `json:"ovs_options,omitempty"`
	OVSPorts           *string           `json:"ovs_ports,omitempty"`
	OVSTag             *string           `json:"ovs_tag,omitempty"`
	Slaves             *string           `json:"slaves,omitempty"`
	VLANID             *int              `json:"vlan_id,omitempty"`
	VLANRawDevice      *string           `json:"vlan_raw_device,omitempty"`
}
