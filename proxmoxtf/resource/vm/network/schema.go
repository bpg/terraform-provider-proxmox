package network

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validators"
)

const (
	// MaxNetworkDevices is the maximum number of network devices supported by the resource.
	MaxNetworkDevices = 32

	dvNetworkDeviceBridge    = "vmbr0"
	dvNetworkDeviceEnabled   = true
	dvNetworkDeviceFirewall  = false
	dvNetworkDeviceMTU       = 0
	dvNetworkDeviceModel     = "virtio"
	dvNetworkDeviceQueues    = 0
	dvNetworkDeviceRateLimit = 0
	dvNetworkDeviceVLANID    = 0

	mkIPv4Addresses = "ipv4_addresses"
	mkIPv6Addresses = "ipv6_addresses"
	mkMACAddresses  = "mac_addresses"

	// MkNetworkDevice is the name of the network device.
	MkNetworkDevice             = "network_device"
	mkNetworkDeviceBridge       = "bridge"
	mkNetworkDeviceDisconnected = "disconnected"
	mkNetworkDeviceEnabled      = "enabled"
	mkNetworkDeviceFirewall     = "firewall"
	mkNetworkDeviceMACAddress   = "mac_address"
	mkNetworkDeviceMTU          = "mtu"
	mkNetworkDeviceModel        = "model"
	mkNetworkDeviceQueues       = "queues"
	mkNetworkDeviceRateLimit    = "rate_limit"
	mkNetworkDeviceTrunks       = "trunks"
	mkNetworkDeviceVLANID       = "vlan_id"
	mkNetworkInterfaceNames     = "network_interface_names"
)

// Schema returns the schema for the network resource.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		mkIPv4Addresses: {
			Type:        schema.TypeList,
			Description: "The IPv4 addresses published by the QEMU agent",
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
		},
		mkIPv6Addresses: {
			Type:        schema.TypeList,
			Description: "The IPv6 addresses published by the QEMU agent",
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
		},
		mkMACAddresses: {
			Type:        schema.TypeList,
			Description: "The MAC addresses for the network interfaces",
			Computed:    true,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		MkNetworkDevice: {
			Type:        schema.TypeList,
			Description: "The network devices",
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					mkNetworkDeviceBridge: {
						Type:        schema.TypeString,
						Description: "The bridge",
						Optional:    true,
						Default:     dvNetworkDeviceBridge,
					},
					mkNetworkDeviceDisconnected: {
						Type:        schema.TypeBool,
						Description: "Whether the network device should be disconnected from the network",
						Optional:    true,
					},
					mkNetworkDeviceEnabled: {
						Type:        schema.TypeBool,
						Description: "Whether to enable the network device",
						Optional:    true,
						Default:     dvNetworkDeviceEnabled,
					},
					mkNetworkDeviceFirewall: {
						Type:        schema.TypeBool,
						Description: "Whether this interface's firewall rules should be used",
						Optional:    true,
						Default:     dvNetworkDeviceFirewall,
					},
					mkNetworkDeviceMACAddress: {
						Type:             schema.TypeString,
						Description:      "The MAC address",
						Optional:         true,
						Computed:         true,
						ValidateDiagFunc: validators.MACAddress(),
					},
					mkNetworkDeviceModel: {
						Type:        schema.TypeString,
						Description: "The model",
						Optional:    true,
						Default:     dvNetworkDeviceModel,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
							"e1000",
							"e1000e",
							"rtl8139",
							"virtio",
							"vmxnet3",
						}, false)),
					},
					mkNetworkDeviceQueues: {
						Type:             schema.TypeInt,
						Description:      "Number of packet queues to be used on the device",
						Optional:         true,
						Default:          dvNetworkDeviceQueues,
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 64)),
					},
					mkNetworkDeviceRateLimit: {
						Type:        schema.TypeFloat,
						Description: "The rate limit in megabytes per second",
						Optional:    true,
						Default:     dvNetworkDeviceRateLimit,
					},
					mkNetworkDeviceVLANID: {
						Type:        schema.TypeInt,
						Description: "The VLAN identifier",
						Optional:    true,
						Default:     dvNetworkDeviceVLANID,
					},
					mkNetworkDeviceTrunks: {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "List of VLAN trunks for the network interface",
					},
					mkNetworkDeviceMTU: {
						Type:        schema.TypeInt,
						Description: "Maximum transmission unit (MTU)",
						Optional:    true,
						Default:     dvNetworkDeviceMTU,
					},
				},
			},
			MaxItems: MaxNetworkDevices,
			MinItems: 0,
		},
		mkNetworkInterfaceNames: {
			Type:        schema.TypeList,
			Description: "The network interface names published by the QEMU agent",
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

// CustomizeDiff returns the custom diff functions for the network resource.
func CustomizeDiff() []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		customdiff.ComputedIf(
			mkIPv4Addresses,
			func(_ context.Context, d *schema.ResourceDiff, _ interface{}) bool {
				return d.HasChange("started") ||
					d.HasChange(MkNetworkDevice)
			},
		),
		customdiff.ComputedIf(
			mkIPv6Addresses,
			func(_ context.Context, d *schema.ResourceDiff, _ interface{}) bool {
				return d.HasChange("started") ||
					d.HasChange(MkNetworkDevice)
			},
		),
		customdiff.ComputedIf(
			mkNetworkInterfaceNames,
			func(_ context.Context, d *schema.ResourceDiff, _ interface{}) bool {
				return d.HasChange("started") ||
					d.HasChange(MkNetworkDevice)
			},
		),
	}
}
