---
layout: page
title: Pool
permalink: /ressources/virtual-environment/pool
nav_order: 7
parent: Virtual Environment Resources
grand_parent: Resources
---

# Resource: Pool

Manages a resource pool.

## Example Usage

```
resource "proxmox_virtual_environment_pool" "operations_pool" {
  comment = "Managed by Terraform"
  pool_id = "operations-pool"
}
```

## Arguments Reference

* `comment` - (Optional) The pool comment.
* `pool_id` - (Required) The pool identifier.

## Attributes Reference

* `members` - The pool members.
    * `datastore_id` - The datastore identifier.
    * `id` - The member identifier.
    * `node_name` - The node name.
    * `type` - The member type.
    * `vm_id` - The virtual machine identifier.
