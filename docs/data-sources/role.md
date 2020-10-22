---
layout: page
title: Role
permalink: /data-sources/role
nav_order: 9
parent: Data Sources
---

# Data Source: Role

Retrieves information about a specific role.

## Example Usage

```
data "proxmox_virtual_environment_role" "operations_role" {
  role_id = "operations"
}
```

## Arguments Reference

* `role_id` - (Required) The role identifier.

## Attributes Reference

* `privileges` - The role privileges
