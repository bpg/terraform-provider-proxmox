---
layout: page
title: Users
permalink: /data-sources/users
nav_order: 15
parent: Data Sources
---

# Data Source: Users

Retrieves information about all the available users.

## Example Usage

```
data "proxmox_virtual_environment_users" "available_users" {}
```

## Arguments Reference

There are no arguments available for this data source.

## Attributes Reference

* `comments` - The user comments.
* `emails` - The users' email addresses.
* `enabled` - Whether a user account is enabled.
* `expiration_dates` - The user accounts' expiration dates (RFC 3339).
* `first_names` - The users' first names.
* `groups` - The users' groups.
* `keys` - The users' keys.
* `last_names` - The users' last names.
* `user_ids` - The user identifiers.
