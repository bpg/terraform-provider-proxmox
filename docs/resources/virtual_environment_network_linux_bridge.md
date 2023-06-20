---
layout: page
title: proxmox_virtual_environment_network_linux_bridge
permalink: /resources/virtual_environment_network_linux_bridge
nav_order: 12
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_network_linux_bridge

Manages a Linux Bridge network interface in a Proxmox VE node.

## Example Usage

```terraform
resource "proxmox_virtual_environment_network_linux_bridge" "bridge99" {
  node_name = "pve"
  iface     = "vmbr99"
  address   = "3.3.3.3/24"
  comment   = "created by terraform"
  mtu       = 1499
}
```

## Argument Reference

- `iface` - (Required) The interface name. Must be "vmbrN", where N is a number between 0 and 9999.
- `node_name` - (Required) The name of the node to manage the interface on.

- `address` - (Optional) The interface IPv4/CIDR address.
- `address6` - (Optional) The interface IPv6/CIDR address.
- `autostart` - (Optional) Automatically start interface on boot (defaults to `true`).
- `bridge_ports` - (Optional) Specify the list of the interface bridge ports.
- `bridge_vlan_aware` - (Optional) Whether the interface bridge is VLAN aware (defaults to `true`).
- `comment` - (Optional) Comment for the interface.
- `gateway` - (Optional) Default gateway address.
- `gateway6` - (Optional) Default IPv6 gateway address.
- `mtu` - (Optional) The interface MTU.
