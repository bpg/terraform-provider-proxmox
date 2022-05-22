---
layout: page
title: proxmox_virtual_environment_roles
permalink: /data-sources/virtual_environment_roles
nav_order: 12
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_roles

Retrieves information about all the available roles.

## Example Usage

```terraform
data "proxmox_virtual_environment_roles" "available_roles" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

* `privileges` - The role privileges.
* `role_ids` - The role identifiers.
* `special` - Whether the role is special (built-in).
