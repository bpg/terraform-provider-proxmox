/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zone

import (
	"context"
	"maps"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
)

func baseAttributesWith(extraAttributes ...map[string]schema.Attribute) map[string]schema.Attribute {
	if len(extraAttributes) > 1 {
		panic("baseAttributesWith expects at most one extraAttributes map")
	}

	if len(extraAttributes) == 0 {
		extraAttributes = append(extraAttributes, make(map[string]schema.Attribute))
	}

	maps.Copy(extraAttributes[0], map[string]schema.Attribute{
		"dns": schema.StringAttribute{
			Optional:    true,
			Description: "DNS API server address.",
		},
		"dns_zone": schema.StringAttribute{
			Optional:    true,
			Description: "DNS domain name. The DNS zone must already exist on the DNS server.",
			MarkdownDescription: "DNS domain name. Used to register hostnames, such as `<hostname>.<domain>`. " +
				"The DNS zone must already exist on the DNS server.",
		},
		"id": schema.StringAttribute{
			Description: "The unique identifier of the SDN zone.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				// https://github.com/proxmox/pve-network/blob/faaf96a8378a3e41065018562c09c3de0aa434f5/src/PVE/Network/SDN/Zones/Plugin.pm#L34
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*[A-Za-z0-9]$`),
					"must be a valid zone identifier",
				),
				stringvalidator.LengthAtMost(8),
			},
		},
		"ipam": schema.StringAttribute{
			Optional:    true,
			Description: "IP Address Management system.",
		},
		"mtu": schema.Int64Attribute{
			Optional:    true,
			Description: "MTU value for the zone.",
		},
		"nodes": stringset.ResourceAttribute("Proxmox node names.", ""),
		"reverse_dns": schema.StringAttribute{
			Optional:    true,
			Description: "Reverse DNS API server address.",
		},
	})

	return extraAttributes[0]
}

func (r *SimpleResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Simple Zone in Proxmox SDN.",
		MarkdownDescription: "Simple Zone in Proxmox SDN. It will create an isolated VNet bridge. " +
			"This bridge is not linked to a physical interface, and VM traffic is only local on each the node. " +
			"It can be used in NAT or routed setups.",
		Attributes: baseAttributesWith(),
	}
}

func (r *VLANResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "VLAN Zone in Proxmox SDN.",
		MarkdownDescription: "VLAN Zone in Proxmox SDN. It uses an existing local Linux or OVS bridge to connect to the " +
			"node's physical interface. It uses VLAN tagging defined in the VNet to isolate the network segments. " +
			"This allows connectivity of VMs between different nodes.",
		Attributes: baseAttributesWith(map[string]schema.Attribute{
			"bridge": schema.StringAttribute{
				Description: "Bridge interface for VLAN.",
				MarkdownDescription: "The local bridge or OVS switch, already configured on _each_ node that allows " +
					"node-to-node connection.",
				Optional: true,
			},
		}),
	}
}

func (r *QinQResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "QinQ Zone in Proxmox SDN.",
		MarkdownDescription: "QinQ Zone in Proxmox SDN. QinQ also known as VLAN stacking, that uses multiple layers of " +
			"VLAN tags for isolation. The QinQ zone defines the outer VLAN tag (the Service VLAN) whereas the inner " +
			"VLAN tag is defined by the VNet. Your physical network switches must support stacked VLANs for this " +
			"configuration. Due to the double stacking of tags, you need 4 more bytes for QinQ VLANs. " +
			"For example, you must reduce the MTU to 1496 if you physical interface MTU is 1500.",
		Attributes: baseAttributesWith(map[string]schema.Attribute{
			"bridge": schema.StringAttribute{
				Description: "A local, VLAN-aware bridge that is already configured on each local node",
				Optional:    true,
			},
			"service_vlan": schema.Int64Attribute{
				Optional:    true,
				Description: "Service VLAN tag for QinQ.",
				Validators: []validator.Int64{
					int64validator.Between(int64(1), int64(4094)),
				},
			},
			"service_vlan_protocol": schema.StringAttribute{
				Optional:    true,
				Description: "Service VLAN protocol for QinQ.",
				Validators: []validator.String{
					stringvalidator.OneOf("802.1ad", "802.1q"),
				},
			},
		}),
	}
}

func (r *VXLANResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "VXLAN Zone in Proxmox SDN.",
		MarkdownDescription: "VXLAN Zone in Proxmox SDN. It establishes a tunnel (overlay) on top of an existing network " +
			"(underlay). This encapsulates layer 2 Ethernet frames within layer 4 UDP datagrams using the default " +
			"destination port 4789. You have to configure the underlay network yourself to enable UDP connectivity " +
			"between all peers. Because VXLAN encapsulation uses 50 bytes, the MTU needs to be 50 bytes lower than the " +
			"outgoing physical interface.",
		Attributes: baseAttributesWith(map[string]schema.Attribute{
			"peers": stringset.ResourceAttribute(
				"A list of IP addresses of each node in the VXLAN zone.",
				"A list of IP addresses of each node in the VXLAN zone. "+
					"This can be external nodes reachable at this IP address. All nodes in the cluster need to be "+
					"mentioned here",
			),
		}),
	}
}

func (r *EVPNResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "EVPN Zone in Proxmox SDN.",
		MarkdownDescription: "EVPN Zone in Proxmox SDN. The EVPN zone creates a routable Layer 3 network, capable of " +
			"spanning across multiple clusters.",
		Attributes: baseAttributesWith(map[string]schema.Attribute{
			"advertise_subnets": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable subnet advertisement for EVPN.",
			},
			"controller": schema.StringAttribute{
				Optional:    true,
				Description: "EVPN controller address.",
			},
			"disable_arp_nd_suppression": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable ARP/ND suppression for EVPN.",
			},
			"exit_nodes": stringset.ResourceAttribute("List of exit nodes for EVPN.", ""),
			"exit_nodes_local_routing": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable local routing for EVPN exit nodes.",
			},
			"primary_exit_node": schema.StringAttribute{
				Optional:    true,
				Description: "Primary exit node for EVPN.",
			},
			"rt_import": schema.StringAttribute{
				Optional:    true,
				Description: "Route target import for EVPN.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(\d+):(\d+)$`),
						"must be in the format '<ASN>:<number>' (e.g., '65000:65000')",
					),
				},
			},
			"vrf_vxlan": schema.Int64Attribute{
				Optional: true,
				Description: "VRF VXLAN-ID used for dedicated routing interconnect between VNets. It must be different " +
					"than the VXLAN-ID of the VNets.",
			},
		}),
	}
}
