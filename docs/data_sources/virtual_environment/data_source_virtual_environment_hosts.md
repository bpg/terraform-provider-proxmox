---
layout: page
title: Hosts
permalink: /data-sources/virtual-environment/hosts
nav_order: 5
parent: Virtual Environment Data Sources
grand_parent: Data Sources
---

# Data Source: Hosts

Retrieves all the host entries from a specific node.

## Example Usage

```
data "proxmox_virtual_environment_hosts" "first_node_host_entries" {
  node_name = "first-node"
}
```

## Arguments Reference

* `node_name` - (Required) A node name.

## Attributes Reference

* `addresses` - The IP addresses.
* `digest` - The SHA1 digest.
* `entries` - The host entries (conversion of `addresses` and `hostnames` into objects).
* `hostnames` - The hostnames associated with each of the IP addresses.
