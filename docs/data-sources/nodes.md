---
layout: page
title: Nodes
permalink: /data-sources/nodes
nav_order: 6
parent: Data Sources
---

# Data Source: Nodes

Retrieves information about all available nodes.

## Example Usage

```
data "proxmox_virtual_environment_nodes" "available_nodes" {}
```

## Arguments Reference

There are no arguments available for this data source.

## Attributes Reference

* `cpu_count` - The CPU count for each node.
* `cpu_utilization` - The CPU utilization on each node.
* `memory_available` - The memory available on each node.
* `memory_used` - The memory used on each node.
* `names` - The node names.
* `online` - Whether a node is online.
* `ssl_fingerprints` - The SSL fingerprint for each node.
* `support_levels` - The support level for each node.
* `uptime` - The uptime in seconds for each node.
