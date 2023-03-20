---
layout: page
title: proxmox_virtual_environment_vms
permalink: /data-sources/virtual_environment_vms
nav_order: 18
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_vms

Retrieves information about all VMs on a specific node.

## Example Usage

```terraform
data "proxmox_virtual_environment_vms" "test_vms" {
    node_name = "test"
}
```

## Argument Reference

- `node_name` - (Required) The node name.

## Attribute Reference

- `vms` - The VMs list.
  - `name` - The virtual machine name.
  - `tags` - A list of tags of the VM.
  - `vm_id` - The VM identifier.
