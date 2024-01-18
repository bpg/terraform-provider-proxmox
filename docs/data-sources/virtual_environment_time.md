---
layout: page
title: proxmox_virtual_environment_time
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_time

Retrieves the current time for a specific node.

## Example Usage

```terraform
data "proxmox_virtual_environment_time" "first_node_time" {
  node_name = "first-node"
}
```

## Argument Reference

- `node_name` - (Required) A node name.

## Attribute Reference

- `local_time` - The node's local time.
- `time_zone` - The node's time zone.
- `utc_time` - The node's local time formatted as UTC.
