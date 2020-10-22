---
layout: page
title: Group
permalink: /resources/group
nav_order: 5
parent: Resources
---

# Resource: Group

Manages a user group.

## Example Usage

```
resource "proxmox_virtual_environment_group" "operations_team" {
  comment  = "Managed by Terraform"
  group_id = "operations-team"
}
```

## Arguments Reference

* `acl` - (Optional) The access control list (multiple blocks supported).
    * `path` - The path.
    * `propagate` - Whether to propagate to child paths.
    * `role_id` - The role identifier.
* `comment` - (Optional) The group comment.
* `group_id` - (Required) The group identifier.

## Attributes Reference

* `members` - The group members as a list of `username@realm` entries
