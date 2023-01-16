---
layout: page
title: proxmox_virtual_environment_user
permalink: /resources/virtual_environment_user
nav_order: 12
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_user

Manages a user.

## Example Usage

```terraform
resource "proxmox_virtual_environment_user" "operations_automation" {
  acl {
    path      = "/vms/1234"
    propagate = true
    role_id   = proxmox_virtual_environment_role.operations_monitoring.role_id
  }

  comment  = "Managed by Terraform"
  password = "a-strong-password"
  user_id  = "operations-automation@pve"
}

resource "proxmox_virtual_environment_role" "operations_monitoring" {
  role_id = "operations-monitoring"

  privileges = [
    "VM.Monitor",
  ]
}
```

## Argument Reference

- `acl` - (Optional) The access control list (multiple blocks supported).
    - `path` - The path.
    - `propagate` - Whether to propagate to child paths.
    - `role_id` - The role identifier.
- `comment` - (Optional) The user comment.
- `email` - (Optional) The user's email address.
- `enabled` - (Optional) Whether the user account is enabled.
- `expiration_date` - (Optional) The user account's expiration date (RFC 3339).
- `first_name` - (Optional) The user's first name.
- `groups` - (Optional) The user's groups.
- `keys` - (Optional) The user's keys.
- `last_name` - (Optional) The user's last name.
- `password` - (Required) The user's password.
- `user_id` - (Required) The user identifier.

## Attribute Reference

There are no additional attributes available for this resource.
