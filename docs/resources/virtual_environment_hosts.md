---
layout: page
title: proxmox_virtual_environment_hosts
permalink: /resources/virtual_environment_hosts
nav_order: 8
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_hosts

Manages the host entries on a specific node.

## Example Usage

```terraform
resource "proxmox_virtual_environment_hosts" "first_node_host_entries" {
  node_name = "first-node"

  entry {
    address = "127.0.0.1"

    hostnames = [
      "localhost",
      "localhost.localdomain",
    ]
  }
}
```

## Argument Reference

- `node_name` - (Required) A node name.
- `entry` - (Required) A host entry (multiple blocks supported).
    - `address` - (Required) The IP address.
    - `hostnames` - (Required) The hostnames.

## Attribute Reference

- `addresses` - The IP addresses.
- `digest` - The SHA1 digest.
- `entries` - The host entries (conversion of `addresses` and `hostnames` into
  objects).
- `hostnames` - The hostnames associated with each of the IP addresses.
