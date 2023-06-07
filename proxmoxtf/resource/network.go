/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validator"
)

const (
	mkNodeNetworkNodeName = "node_name"

	mkNodeNetworkInterfaceName     = "name"
	mkNodeNetworkInterfaceType     = "type"
	mkNodeNetworkInterfaceAddress  = "address"
	mkNodeNetworkInterfaceAddress6 = "address6"

	mkNodeNetworkInterfaceBondPrimary        = "bond_primary"
	mkNodeNetworkInterfaceBondMode           = "bond_mode"
	mkNodeNetworkInterfaceBondXmitHashPolicy = "bond_xmit_hash_policy"
	mkNodeNetworkInterfaceBridgePorts        = "bridge_ports"
	mkNodeNetworkInterfaceBridgeVLANAware    = "bridge_vlan_aware"
)

/*
  node = "pve01"
  type = "vlan"
  name = "enp0s3.100" # Maps to "iface" in API
  vlan = "100" # Maps to vlan-id in API
  interface = "enp0s3" # Maps to vlan-raw-device in API
  comment = "VLAN 100 Interface" # Maps to comments in API
*/

// LinuxBridge - proxmox_virtual_environment_network_linux_bridge
// LinuxBond   - proxmox_virtual_environment_network_linux_bond
// LinuxVLAN   - proxmox_virtual_environment_network_linux_vlan

func NetworkLinuxBridge() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkNodeNetworkNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
				ForceNew:    true,
			},
			mkNodeNetworkInterfaceName: {
				Type:        schema.TypeString,
				Description: "The interface name",
				Required:    true,
				ForceNew:    true,
			},
			mkNodeNetworkInterfaceAddress: {
				Type:             schema.TypeString,
				Description:      "The interface IPv4 address",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
			},
			mkNodeNetworkInterfaceAddress6: {
				Type:             schema.TypeString,
				Description:      "The interface IPv6 address",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv6Address),
			},
			mkNodeNetworkInterfaceBridgePorts: {
				Type:        schema.TypeList,
				Description: "The interface bridge ports",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func NetworkLinuxBond() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkNodeNetworkNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
				ForceNew:    true,
			},
			mkNodeNetworkInterfaceName: {
				Type:        schema.TypeString,
				Description: "The interface name",
				Required:    true,
				ForceNew:    true,
			},
			mkNodeNetworkInterfaceAddress: {
				Type:             schema.TypeString,
				Description:      "The interface IPv4 address",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
			},
			mkNodeNetworkInterfaceAddress6: {
				Type:             schema.TypeString,
				Description:      "The interface IPv6 address",
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv6Address),
			},
			mkNodeNetworkInterfaceBondPrimary: {
				Type:        schema.TypeString,
				Description: "The interface bond primary for active-backup bond",
				Optional:    true,
			},
			mkNodeNetworkInterfaceBondMode: {
				Type:             schema.TypeString,
				Description:      "The interface bonding mode",
				Optional:         true,
				ValidateDiagFunc: validator.NodeNetworkInterfaceBondingModes(),
			},
			mkNodeNetworkInterfaceBondXmitHashPolicy: {
				Type:             schema.TypeString,
				Description:      "Selects the transmit hash policy to use for slave selection in balance-xor and 802.3ad modes",
				Optional:         true,
				ValidateDiagFunc: validator.NodeNetworkInterfaceBondingTransmitHashPolicies(),
			},
		},
	}
}
