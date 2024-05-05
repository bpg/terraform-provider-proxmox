---
layout: page
title: proxmox_virtual_environment_pools
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_pools

Retrieves the identifiers for all the available resource pools.

## Example Usage

```hcl
data "proxmox_virtual_environment_pools" "available_pools" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

- `pool_ids` - The pool identifiers.
