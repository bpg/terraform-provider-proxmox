---
layout: page
title: DNS
permalink: /data-sources/virtual-environment/dns
nav_order: 2
parent: Virtual Environment Data Sources
grand_parent: Data Sources
---

# Data Source: DNS

Retrieves the DNS configuration for a specific node.

## Example Usage

```
data "proxmox_virtual_environment_dns" "first_node" {
  node_name = "first-node"
}
```

## Arguments Reference

* `node_name` - (Required) A node name.

## Attributes Reference

* `domain` - The DNS search domain.
* `servers` - The DNS servers.
