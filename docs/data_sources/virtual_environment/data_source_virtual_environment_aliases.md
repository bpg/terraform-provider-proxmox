---
layout: page
title: Aliases
permalink: /data-sources/virtual-environment/aliases
nav_order: 1
parent: Virtual Environment Data Sources
parent: Data Sources
---

# Data Source: Aliases

Retrieves the identifiers for all the available aliases.

## Example Usage

```
data "proxmox_virtual_environment_cluster_aliases" "available_aliases" {}
```

## Arguments Reference

There are no arguments available for this data source.

## Attributes Reference

* `alias_ids` - The pool identifiers.
