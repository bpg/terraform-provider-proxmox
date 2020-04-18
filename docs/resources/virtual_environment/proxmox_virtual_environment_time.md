---
layout: page
title: Time
permalink: /ressources/virtual-environment/time
nav_order: 9
parent: Virtual Environment Resources
grand_parent: Resources
---

# Resource: Time

Manages the time for a specific node.

## Example Usage

```
resource "proxmox_virtual_environment_time" "first_node_time" {
  node_name = "first-node"
  time_zone = "UTC"
}
```

## Arguments Reference

* `node_name` - (Required) A node name.
* `time_zone` - (Required) The node's time zone.

## Attributes Reference

* `local_time` - The node's local time.
* `utc_time` - The node's local time formatted as UTC.
