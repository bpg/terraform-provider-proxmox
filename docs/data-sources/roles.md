---
layout: page
title: Roles
permalink: /data-sources/roles
nav_order: 12
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: Roles

Retrieves information about all the available roles.

## Example Usage

```
data "proxmox_virtual_environment_roles" "available_roles" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

* `privileges` - The role privileges.
* `role_ids` - The role identifiers.
* `special` - Whether the role is special (built-in).
