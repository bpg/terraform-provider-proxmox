---
layout: page
title: proxmox_virtual_environment_container
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_container

Retrieves information about a specific Container.

## Example Usage

```hcl
data "proxmox_virtual_environment_container" "test_container" {
    node_name = "test"
    vm_id = 100
}
```

## Argument Reference

- `node_name` - (Required) The node name.
- `vm_id` - (Required) The container identifier.

## Attribute Reference

- `name` - The container name.
- `tags` - A list of tags of the container.
- `status` - The status of the container.
- `template` - Whether the container is a template.
