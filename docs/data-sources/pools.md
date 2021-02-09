---
layout: page
title: Pools
permalink: /data-sources/pools
nav_order: 10
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: Pools

Retrieves the identifiers for all the available resource pools.

## Example Usage

```
data "proxmox_virtual_environment_pools" "available_pools" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

* `pool_ids` - The pool identifiers.
