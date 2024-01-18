---
layout: page
title: proxmox_virtual_environment_hosts
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_hosts

Retrieves all the host entries from a specific node.

## Example Usage

```terraform
data "proxmox_virtual_environment_hosts" "first_node_host_entries" {
  node_name = "first-node"
}
```

## Argument Reference

- `node_name` - (Required) A node name.

## Attribute Reference

- `addresses` - The IP addresses.
- `digest` - The SHA1 digest.
- `entries` - The host entries (conversion of `addresses` and `hostnames` into
  objects).
- `hostnames` - The hostnames associated with each of the IP addresses.
