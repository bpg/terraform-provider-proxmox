---
layout: page
title: Role
permalink: /resources/role
nav_order: 10
parent: Resources
subcategory: Virtual Environment
---

# Resource: Role

Manages a role.

## Example Usage

```
resource "proxmox_virtual_environment_role" "operations_monitoring" {
  role_id = "operations-monitoring"

  privileges = [
    "VM.Monitor",
  ]
}
```

## Argument Reference

* `privileges` - (Required) The role privileges.
* `role_id` - (Required) The role identifier.

## Attribute Reference

There are no additional attributes available for this resource.
