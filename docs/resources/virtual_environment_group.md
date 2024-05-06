---
layout: page
title: proxmox_virtual_environment_group
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_group

Manages a user group.

## Example Usage

```hcl
resource "proxmox_virtual_environment_group" "operations_team" {
  comment  = "Managed by Terraform"
  group_id = "operations-team"
}
```

## Argument Reference

- `acl` - (Optional) The access control list (multiple blocks supported).
    - `path` - The path.
    - `propagate` - Whether to propagate to child paths.
    - `role_id` - The role identifier.
- `comment` - (Optional) The group comment.
- `group_id` - (Required) The group identifier.

## Attribute Reference

- `members` - The group members as a list of `username@realm` entries

## Import

Instances can be imported using the `group_id`, e.g.,

```bash
terraform import proxmox_virtual_environment_group.operations_team operations-team
```
