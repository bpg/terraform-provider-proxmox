---
layout: page
title: Time
permalink: /data-sources/time
nav_order: 13
parent: Data Sources
---

# Data Source: Time

Retrieves the current time for a specific node.

## Example Usage

```
data "proxmox_virtual_environment_time" "first_node_time" {
  node_name = "first-node"
}
```

## Arguments Reference

* `node_name` - (Required) A node name.

## Attributes Reference

* `local_time` - The node's local time.
* `time_zone` - The node's time zone.
* `utc_time` - The node's local time formatted as UTC.
