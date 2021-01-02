---
layout: page
title: Roles
permalink: /data-sources/roles
nav_order: 12
parent: Data Sources
---

# Data Source: Roles

Retrieves information about all the available roles.

## Example Usage

```
data "proxmox_virtual_environment_roles" "available_roles" {}
```

## Arguments Reference

There are no arguments available for this data source.

## Attributes Reference

* `privileges` - The role privileges.
* `role_ids` - The role identifiers.
* `special` - Whether the role is special (built-in).
