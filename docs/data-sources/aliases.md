---
layout: page
title: Aliases
permalink: /data-sources/aliases
nav_order: 2
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: Aliases

Retrieves the identifiers for all the available aliases.

## Example Usage

```
data "proxmox_virtual_environment_cluster_aliases" "available_aliases" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

* `alias_ids` - The pool identifiers.
