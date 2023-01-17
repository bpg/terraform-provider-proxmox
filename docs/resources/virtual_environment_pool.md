---
layout: page
title: proxmox_virtual_environment_pool
permalink: /resources/virtual_environment_pool
nav_order: 9
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_pool

Manages a resource pool.

## Example Usage

```terraform
resource "proxmox_virtual_environment_pool" "operations_pool" {
  comment = "Managed by Terraform"
  pool_id = "operations-pool"
}
```

## Argument Reference

- `comment` - (Optional) The pool comment.
- `pool_id` - (Required) The pool identifier.

## Attribute Reference

- `members` - The pool members.
    - `datastore_id` - The datastore identifier.
    - `id` - The member identifier.
    - `node_name` - The node name.
    - `type` - The member type.
    - `vm_id` - The virtual machine identifier.
