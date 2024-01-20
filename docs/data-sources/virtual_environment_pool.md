---
layout: page
title: proxmox_virtual_environment_pool
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_pool

Retrieves information about a specific resource pool.

## Example Usage

```terraform
data "proxmox_virtual_environment_pool" "operations_pool" {
  pool_id = "operations"
}
```

## Argument Reference

- `pool_id` - (Required) The pool identifier.

## Attribute Reference

- `comment` - The pool comment.
- `members` - The pool members.
  - `datastore_id` - The datastore identifier.
  - `id` - The member identifier.
  - `node_name` - The node name.
  - `type` - The member type.
  - `vm_id` - The virtual machine identifier.
