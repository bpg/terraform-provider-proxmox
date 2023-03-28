---
layout: page
title: proxmox_virtual_environment_cluster_firewall_security_group
permalink: /data-sources/virtual_environment_cluster_firewall_security_group
nav_order: 5
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_cluster_firewall_security_group

Retrieves information about a specific security group.

## Example Usage

```terraform
data "proxmox_virtual_environment_cluster_firewall_security_group" "webserver" {
  name = "webserver"
}
```

## Argument Reference

- `name` - (Required) Security group name.

## Attribute Reference

- `comment` - (Optional) Security group comment.
- `rules` - (Optional) List of firewall rules.
    - `action` - Rule action (`ACCEPT`, `DROP`, `REJECT`).
    - `type` - Rule type (`in`, `out`).
    - `comment` - (Optional) Rule comment.
    - `dest` - (Optional) Packet destination address.
    - `dport` - (Optional) TCP/UDP destination port.
    - `enable` - (Optional) Enable this rule.
    - `iface` - (Optional) Network interface name.
    - `log` - (Optional) Log level for this rule (`emerg`, `alert`, `crit`,
      `err`, `warning`, `notice`, `info`, `debug`, `nolog`).
    - `macro`- (Optional) Macro name. Use predefined standard macro.
    - `proto` - (Optional) Packet protocol.
    - `source` - (Optional) Packet source address.
    - `sport` - (Optional) TCP/UDP source port.
