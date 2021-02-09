---
layout: page
title: Time
permalink: /resources/time
nav_order: 11
parent: Resources
subcategory: Virtual Environment
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

## Argument Reference

* `node_name` - (Required) A node name.
* `time_zone` - (Required) The node's time zone.

## Attribute Reference

* `local_time` - The node's local time.
* `utc_time` - The node's local time formatted as UTC.
