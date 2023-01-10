---
layout: page
title: proxmox_virtual_environment_cluster_aliases
permalink: /data-sources/virtual_environment_cluster_aliases
nav_order: 2
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_cluster_aliases

Retrieves the identifiers for all the available aliases.

## Example Usage

```terraform
data "proxmox_virtual_environment_cluster_aliases" "available_aliases" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

- `alias_ids` - The pool identifiers.
