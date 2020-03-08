---
layout: page
title: Group
permalink: /data-sources/virtual-environment/group
nav_order: 3
parent: Virtual Environment Data Sources
grand_parent: Data Sources
---

# Data Source: Group

Retrieves information about a specific user group.

## Example Usage

```
data "proxmox_virtual_environment_group" "operations_team" {
  group_id = "operations-team"
}
```

## Arguments Reference

* `group_id` - (Required) The group identifier.

## Attributes Reference

* `acl` - The access control list.
    * `path` - The path.
    * `propagate` - Whether to propagate to child paths.
    * `role_id` - The role identifier.
* `comment` - The group comment.
* `members` - The group members as a list with `username@realm` entries.
