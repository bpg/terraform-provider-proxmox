---
layout: page
title: Pool
permalink: /data-sources/pool
nav_order: 9
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: Pool

Retrieves information about a specific resource pool.

## Example Usage

```
data "proxmox_virtual_environment_pool" "operations_pool" {
  pool_id = "operations"
}
```

## Argument Reference

* `pool_id` - (Required) The pool identifier.

## Attribute Reference

* `comment` - The pool comment.
* `members` - The pool members.
    * `datastore_id` - The datastore identifier.
    * `id` - The member identifier.
    * `node_name` - The node name.
    * `type` - The member type.
    * `vm_id` - The virtual machine identifier.
