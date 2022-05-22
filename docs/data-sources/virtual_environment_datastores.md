---
layout: page
title: proxmox_virtual_environment_datastores
permalink: /data-sources/virtual_environment_datastores
nav_order: 3
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_datastores

Retrieves information about all the datastores available to a specific node.

## Example Usage

```terraform
data "proxmox_virtual_environment_datastores" "first_node" {
  node_name = "first-node"
}
```

## Argument Reference

* `node_name` - (Required) A node name.

## Attribute Reference

* `active` - Whether the datastore is active.
* `content_types` - The allowed content types.
* `datastore_ids` - The datastore identifiers.
* `enabled` - Whether the datastore is enabled.
* `shared` - Whether the datastore is shared.
* `space_available` - The available space in bytes.
* `space_total` - The total space in bytes.
* `space_used` - The used space in bytes.
* `types` - The storage types.
