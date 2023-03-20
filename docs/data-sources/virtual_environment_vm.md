---
layout: page
title: proxmox_virtual_environment_vm
permalink: /data-sources/virtual_environment_vm
nav_order: 17
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_vms

Retrieves information about a specific VM.

## Example Usage

```terraform
data "proxmox_virtual_environment_vm" "test_vm" {
    node_name = "test"
    vm_id = 100
}
```

## Argument Reference

- `node_name` - (Required) The node name.
- `vm_id` - (Required) The VM identifier.

## Attribute Reference

- `name` - The virtual machine name.
- `tags` - A list of tags of the VM.
