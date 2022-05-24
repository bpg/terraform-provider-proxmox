---
layout: page
title: proxmox_virtual_environment_dns
permalink: /data-sources/virtual_environment_dns
nav_order: 4
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_dns

Retrieves the DNS configuration for a specific node.

## Example Usage

```terraform
data "proxmox_virtual_environment_dns" "first_node" {
  node_name = "first-node"
}
```

## Argument Reference

* `node_name` - (Required) A node name.

## Attribute Reference

* `domain` - The DNS search domain.
* `servers` - The DNS servers.
