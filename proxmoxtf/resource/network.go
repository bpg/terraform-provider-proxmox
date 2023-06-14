/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/internal/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validator"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

const (
	mkNodeNetworkNodeName = "node_name"

	mkNodeNetworkInterfaceName      = "name"
	mkNodeNetworkInterfaceType      = "type"
	mkNodeNetworkInterfaceAddress   = "address"
	mkNodeNetworkInterfaceGateway   = "gateway"
	mkNodeNetworkInterfaceAddress6  = "address6"
	mkNodeNetworkInterfaceGateway6  = "gateway6"
	mkNodeNetworkInterfaceAutostart = "autostart"
	mkNodeNetworkInterfaceComment   = "comment"
	mkNodeNetworkInterfaceMTU       = "mtu"

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

// common:
//  - name
//  - ipv4/CIDR
//  - gateway4
//  - ipv6/CIDR
//  - gateway6
//  - autostart
//  - comments
//  - MTU

// LinuxBridge - proxmox_virtual_environment_network_linux_bridge
// 	- VLAN aware
//  - bridge_ports

// LinuxBond   - proxmox_virtual_environment_network_linux_bond
//  - slaves
//  - bond_mode
//  - bond_primary
//  - bond_xmit_hash_policy

// LinuxVLAN   - proxmox_virtual_environment_network_linux_vlan
//  - vlan_raw_device
//  - vlan_tag

func baseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
			Description:      "The interface IPv4/CIDR address",
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
		},
		mkNodeNetworkInterfaceGateway: {
			Type:             schema.TypeString,
			Description:      "Default gateway address",
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
		},
		mkNodeNetworkInterfaceAddress6: {
			Type:             schema.TypeString,
			Description:      "The interface IPv6/CIDR address",
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
		},
		mkNodeNetworkInterfaceGateway6: {
			Type:             schema.TypeString,
			Description:      "Default IPv6 gateway address",
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv6Address),
		},
		mkNodeNetworkInterfaceAutostart: {
			Type:        schema.TypeBool,
			Description: "Automatically start interface on boot.",
			Optional:    true,
			Default:     true,
		},
		mkNodeNetworkInterfaceComment: {
			Type:        schema.TypeString,
			Description: "Comment for the interface.",
			Optional:    true,
		},
	}
}

func NetworkLinuxBridge() *schema.Resource {
	s := map[string]*schema.Schema{
		mkNodeNetworkInterfaceBridgePorts: {
			Type:        schema.TypeList,
			Description: "The interface bridge ports",
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		mkNodeNetworkInterfaceBridgeVLANAware: {
			Type:        schema.TypeBool,
			Description: "The interface bridge VLAN aware",
			Optional:    true,
			Default:     false,
		},
	}

	structure.MergeSchema(s, baseSchema())

	return &schema.Resource{
		Schema:        s,
		CreateContext: CreateNetworkInterfaceBridge,
		ReadContext:   ReadNetworkInterfaceBridge,
		UpdateContext: UpdateNetworkInterfaceBridge,
		DeleteContext: DeleteNetworkInterfaceBridge,
	}
}

func CreateNetworkInterfaceBridge(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkNodeNetworkNodeName).(string)

	body := getNetworkInterfaceCreateUpdateRequestBody(d)

	bridgePorts := getBridgePorts(d)
	bridgeVLANAware := types.CustomBool(d.Get(mkNodeNetworkInterfaceBridgeVLANAware).(bool))

	body.Type = "bridge"
	body.BridgePorts = &bridgePorts
	body.BridgeVLANAware = &bridgeVLANAware

	err = api.Node(nodeName).CreateNetworkInterface(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s-%s", nodeName, body.Iface))

	return ReadNetworkInterfaceBridge(ctx, d, m)
}

func getNetworkInterfaceCreateUpdateRequestBody(d *schema.ResourceData) *nodes.NetworkInterfaceCreateUpdateRequestBody {
	iface := d.Get(mkNodeNetworkInterfaceName).(string)
	autostart := types.CustomBool(d.Get(mkNodeNetworkInterfaceAutostart).(bool))

	addr := d.Get(mkNodeNetworkInterfaceAddress).(string)
	gw := d.Get(mkNodeNetworkInterfaceGateway).(string)
	addr6 := d.Get(mkNodeNetworkInterfaceAddress6).(string)
	gw6 := d.Get(mkNodeNetworkInterfaceGateway6).(string)
	comments := d.Get(mkNodeNetworkInterfaceComment).(string)

	body := &nodes.NetworkInterfaceCreateUpdateRequestBody{
		Iface:     iface,
		Autostart: &autostart,
	}

	if addr != "" {
		body.CIDR = &addr
	}

	if gw != "" {
		body.Gateway = &gw
	}

	if addr6 != "" {
		body.CIDR6 = &addr6
	}

	if gw6 != "" {
		body.Gateway6 = &gw6
	}

	if comments != "" {
		body.Comments = &comments
	}

	return body
}

func ReadNetworkInterfaceBridge(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkNodeNetworkNodeName).(string)

	ifaces, err := api.Node(nodeName).ListNetworkInterfaces(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	iface := d.Get(mkNodeNetworkInterfaceName).(string)
	for _, ii := range ifaces {
		if ii.Iface == iface {
			d.Set(mkNodeNetworkInterfaceAddress, ii.CIDR)
			d.Set(mkNodeNetworkInterfaceGateway, ii.Gateway)
			d.Set(mkNodeNetworkInterfaceAddress6, ii.CIDR)
			d.Set(mkNodeNetworkInterfaceGateway6, ii.Gateway6)
		}
	}
	return nil
}

func UpdateNetworkInterfaceBridge(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func DeleteNetworkInterfaceBridge(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
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

func getBridgePorts(d *schema.ResourceData) string {
	ports := d.Get(mkNodeNetworkInterfaceBridgePorts).([]interface{})
	var sanitizedPorts []string
	for i := 0; i < len(ports); i++ {
		tag := strings.TrimSpace(ports[i].(string))
		if len(tag) > 0 {
			sanitizedPorts = append(sanitizedPorts, tag)
		}
	}
	sort.Strings(sanitizedPorts)
	return strings.Join(sanitizedPorts, " ")
}
