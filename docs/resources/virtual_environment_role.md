---
layout: page
title: proxmox_virtual_environment_role
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_role

Manages a role.

## Example Usage

```hcl
resource "proxmox_virtual_environment_role" "operations_monitoring" {
  role_id = "operations-monitoring"

  privileges = [
    "VM.Monitor",
  ]
}
```

## Argument Reference

- `privileges` - (Required) The role privileges.
- `role_id` - (Required) The role identifier.

## Attribute Reference

There are no additional attributes available for this resource.

## Import

Instances can be imported using the `role_id`, e.g.,

```bash
terraform import proxmox_virtual_environment_role.operations_monitoring operations-monitoring
```
