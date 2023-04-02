---
layout: page
title: proxmox_virtual_environment_role
permalink: /data-sources/virtual_environment_role
nav_order: 15
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_role

Retrieves information about a specific role.

## Example Usage

```terraform
data "proxmox_virtual_environment_role" "operations_role" {
  role_id = "operations"
}
```

## Argument Reference

- `role_id` - (Required) The role identifier.

## Attribute Reference

- `privileges` - The role privileges
