---
layout: page
title: DNS
permalink: /data-sources/dns
nav_order: 4
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: DNS

Retrieves the DNS configuration for a specific node.

## Example Usage

```
data "proxmox_virtual_environment_dns" "first_node" {
  node_name = "first-node"
}
```

## Argument Reference

* `node_name` - (Required) A node name.

## Attribute Reference

* `domain` - The DNS search domain.
* `servers` - The DNS servers.
