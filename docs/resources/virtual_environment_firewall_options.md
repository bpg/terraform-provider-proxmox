---
layout: page
title: proxmox_virtual_environment_firewall_options
permalink: /resources/virtual_environment_firewall_options
nav_order: 9
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_firewall_options

Manages firewall options on VM / Container level.

## Example Usage

```terraform
resource "proxmox_virtual_environment_firewall_options" "example" {
  depends_on = [proxmox_virtual_environment_vm.example]

  node_name = proxmox_virtual_environment_vm.example.node_name
  vm_id     = proxmox_virtual_environment_vm.example.vm_id

  dhcp          = true
  enabled       = false
  ipfilter      = true
  log_level_in  = "info"
  log_level_out = "info"
  macfilter     = false
  ndp           = true
  input_policy  = "ACCEPT"
  output_policy = "ACCEPT"
  radv          = true
}
```

## Argument Reference

- `node_name` - (Required) Node name.
- `vm_id` - (Optional) VM ID. Leave empty for cluster level aliases.
- `container_id` - (Optional) Container ID. Leave empty for cluster level aliases.
- `dhcp` - (Optional)Enable DHCP.
- `enabled` - (Optional) Enable or disable the firewall.
- `ipfilter` - (Optional) Enable default IP filters. This is equivalent to
  adding an empty ipfilter-net<id> ipset for every interface. Such ipsets
  implicitly contain sane default restrictions such as restricting IPv6 link
  local addresses to the one derived from the interface's MAC address. For
  containers the configured IP addresses will be implicitly added.
- `log_level_in` - (Optional) Log level for incoming
  packets (`emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`, `nolog`).
- `log_level_out` - (Optional) Log level for outgoing
  packets (`emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`, `nolog`).
- `macfilter` - (Optional) Enable/disable MAC address filter.
- `ndp` - (Optional) Enable NDP (Neighbor Discovery Protocol).
- `input_policy` - (Optional) The default input
  policy (`ACCEPT`, `DROP`, `REJECT`).
- `output_policy` - (Optional) The default output
  policy (`ACCEPT`, `DROP`, `REJECT`).
- `radv` - (Optional) Enable Router Advertisement.

## Attribute Reference

There are no additional attributes available for this resource.
