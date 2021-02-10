---
layout: page
title: proxmox_virtual_environment_nodes
permalink: /data-sources/virtual_environment_nodes
nav_order: 8
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_nodes

Retrieves information about all available nodes.

## Example Usage

```
data "proxmox_virtual_environment_nodes" "available_nodes" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

* `cpu_count` - The CPU count for each node.
* `cpu_utilization` - The CPU utilization on each node.
* `memory_available` - The memory available on each node.
* `memory_used` - The memory used on each node.
* `names` - The node names.
* `online` - Whether a node is online.
* `ssl_fingerprints` - The SSL fingerprint for each node.
* `support_levels` - The support level for each node.
* `uptime` - The uptime in seconds for each node.
