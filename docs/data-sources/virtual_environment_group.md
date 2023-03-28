---
layout: page
title: proxmox_virtual_environment_group
permalink: /data-sources/virtual_environment_group
nav_order: 9
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_group

Retrieves information about a specific user group.

## Example Usage

```terraform
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
