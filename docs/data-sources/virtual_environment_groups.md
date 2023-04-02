---
layout: page
title: proxmox_virtual_environment_groups
permalink: /data-sources/virtual_environment_groups
nav_order: 10
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_groups

Retrieves basic information about all available user groups.

## Example Usage

```terraform
data "proxmox_virtual_environment_groups" "available_groups" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

- `comments` - The group comments.
- `group_ids` - The group identifiers.
