---
layout: page
title: User
permalink: /data-sources/user
nav_order: 14
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: User

Retrieves information about a specific user.

## Example Usage

```
data "proxmox_virtual_environment_user" "operations_user" {
  user_id = "operation@pam"
}
```

## Argument Reference

* `user_id` - (Required) The user identifier.

## Attribute Reference

* `acl` - The access control list.
    * `path` - The path.
    * `propagate` - Whether to propagate to child paths.
    * `role_id` - The role identifier.
* `comment` - The user comment.
* `email` - The user's email address.
* `enabled` - Whether the user account is enabled.
* `expiration_date` - The user account's expiration date (RFC 3339).
* `first_name` - The user's first name.
* `groups` - The user's groups.
* `keys` - The user's keys.
* `last_name` - The user's last name.
