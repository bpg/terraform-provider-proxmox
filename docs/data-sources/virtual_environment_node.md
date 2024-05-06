---
layout: page
title: proxmox_virtual_environment_node
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_node

Retrieves information about node.

## Example Usage

```hcl
data "proxmox_virtual_environment_node" "node" {}
```

## Argument Reference

- `node_name` - (Required) The node name.

## Attribute Reference

- `cpu_count` - The CPU count on the node.
- `cpu_sockets` - The CPU utilization on the node.
- `cpu_model` - The CPU model on the node.
- `memory_available` - The memory available on the node.
- `memory_used` - The memory used on the node.
- `memory_total` - The total memory on the node.
- `uptime` - The uptime in seconds on the node.
