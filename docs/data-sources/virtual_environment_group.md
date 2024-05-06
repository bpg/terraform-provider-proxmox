---
layout: page
title: proxmox_virtual_environment_group
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_group

Retrieves information about a specific user group.

## Example Usage

```hcl
data "proxmox_virtual_environment_group" "operations_team" {
  group_id = "operations-team"
}
```

## Argument Reference

- `group_id` - (Required) The group identifier.

## Attribute Reference

- `acl` - The access control list.
    - `path` - The path.
    - `propagate` - Whether to propagate to child paths.
    - `role_id` - The role identifier.
- `comment` - The group comment.
- `members` - The group members as a list with `username@realm` entries.
