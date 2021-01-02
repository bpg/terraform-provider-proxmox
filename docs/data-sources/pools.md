---
layout: page
title: Pools
permalink: /data-sources/pools
nav_order: 10
parent: Data Sources
---

# Data Source: Pools

Retrieves the identifiers for all the available resource pools.

## Example Usage

```
data "proxmox_virtual_environment_pools" "available_pools" {}
```

## Arguments Reference

There are no arguments available for this data source.

## Attributes Reference

* `pool_ids` - The pool identifiers.
