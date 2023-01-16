---
layout: page
title: proxmox_virtual_environment_dns
permalink: /resources/virtual_environment_dns
nav_order: 5
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_dns

Manages the DNS configuration for a specific node.

## Example Usage

```terraform
resource "proxmox_virtual_environment_dns" "first_node_dns_configuration" {
  domain    = data.proxmox_virtual_environment_dns.first_node_dns_configuration.domain
  node_name = data.proxmox_virtual_environment_dns.first_node_dns_configuration.node_name

  servers = [
    "1.1.1.1",
    "1.0.0.1",
  ]
}

data "proxmox_virtual_environment_dns" "first_node_dns_configuration" {
  node_name = "first-node"
}
```

## Argument Reference

- `domain` - (Required) The DNS search domain.
- `node_name` - (Required) A node name.
- `servers` - (Optional) The DNS servers.

## Attribute Reference

There are no additional attributes available for this resource.

## Important Notes

Be careful not to use this resource multiple times for the same node.
