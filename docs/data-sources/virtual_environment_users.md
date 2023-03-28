---
layout: page
title: proxmox_virtual_environment_users
permalink: /data-sources/virtual_environment_users
nav_order: 19
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_users

Retrieves information about all the available users.

## Example Usage

```terraform
data "proxmox_virtual_environment_users" "available_users" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

- `comments` - The user comments.
- `emails` - The users' email addresses.
- `enabled` - Whether a user account is enabled.
- `expiration_dates` - The user accounts' expiration dates (RFC 3339).
- `first_names` - The users' first names.
- `groups` - The users' groups.
- `keys` - The users' keys.
- `last_names` - The users' last names.
- `user_ids` - The user identifiers.
