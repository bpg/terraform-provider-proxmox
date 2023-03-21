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
data "proxmox_virtual_environment_vms" "ubuntu_vms" {
  id        = "ubuntu_vms"
  tags      = ["ubuntu"]
}
```

## Argument Reference

- `id` - (Required) The data source identifier, could be any string. This is
  used to identify the data source among other data sources of the same type in
  the terraform state.
- `node_name` - (Optional) The node name.
- `tags` - (Optional) A list of tags to filter the VMs. The VM must have all
  the tags to be included in the result.

## Attribute Reference

- `vms` - The VMs list.
    - `name` - The virtual machine name.
    - `node_name` - The node name.
    - `tags` - A list of tags of the VM.
    - `vm_id` - The VM identifier.
