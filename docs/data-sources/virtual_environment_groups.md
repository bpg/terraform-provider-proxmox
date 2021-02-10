---
layout: page
title: Groups
permalink: /data-sources/groups
nav_order: 6
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: Groups

Retrieves basic information about all available user groups.

## Example Usage

```
data "proxmox_virtual_environment_groups" "available_groups" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

* `comments` - The group comments.
* `group_ids` - The group identifiers.
