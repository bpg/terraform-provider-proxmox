---
layout: page
title: Groups
permalink: /data-sources/groups
nav_order: 6
parent: Data Sources
---

# Data Source: Groups

Retrieves basic information about all available user groups.

## Example Usage

```
data "proxmox_virtual_environment_groups" "available_groups" {}
```

## Arguments Reference

There are no arguments available for this data source.

## Attributes Reference

* `comments` - The group comments.
* `group_ids` - The group identifiers.
