---
layout: page
title: proxmox_virtual_environment_time
permalink: /resources/virtual_environment_time
nav_order: 12
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_time

Manages the time for a specific node.

## Example Usage

```terraform
resource "proxmox_virtual_environment_time" "first_node_time" {
  node_name = "first-node"
  time_zone = "UTC"
}
```

## Argument Reference

- `node_name` - (Required) A node name.
- `time_zone` - (Required) The node's time zone.

## Attribute Reference

- `local_time` - The node's local time.
- `utc_time` - The node's local time formatted as UTC.
