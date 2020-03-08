---
layout: page
title: Groups
permalink: /data-sources/virtual-environment/groups
nav_order: 4
parent: Virtual Environment Data Sources
grand_parent: Data Sources
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
