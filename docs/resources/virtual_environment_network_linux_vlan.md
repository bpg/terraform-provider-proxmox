---
layout: page
title: proxmox_virtual_environment_network_linux_vlan
permalink: /resources/virtual_environment_network_linux_vlan
nav_order: 13
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_network_linux_vlan

Manages a Linux VLAN network interface in a Proxmox VE node.

## Example Usage

```terraform
resource "proxmox_virtual_environment_network_linux_vlan" "vlan21" {
  node_name = "pve"
  iface     = "ens18.21"
  comment   = "created by terraform"
}
```

## Argument Reference

- `node_name` - (Required) The name of the node to manage the interface on.
- `name` - (Required) The interface name. Add the VLAN tag number to an
  existing interface name, e.g. "ens18.21".

- `address` - (Optional) The interface IPv4/CIDR address.
- `address6` - (Optional) The interface IPv6/CIDR address.
- `autostart` - (Optional) Automatically start interface on boot (defaults
  to `true`).
- `comment` - (Optional) Comment for the interface.
- `gateway` - (Optional) Default gateway address.
- `gateway6` - (Optional) Default IPv6 gateway address.
- `mtu` - (Optional) The interface MTU.

### Read-Only

- `id` (String) A unique identifier with format '<node name>:<iface>'
- `interface` (String) The VLAN raw device.
- `vlan` (Number) The VLAN tag
