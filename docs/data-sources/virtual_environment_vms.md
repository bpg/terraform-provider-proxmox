---
layout: page
title: proxmox_virtual_environment_vms
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_vms

Retrieves information about all VMs in the Proxmox cluster.

## Example Usage

```terraform
data "proxmox_virtual_environment_vms" "ubuntu_vms" {
  tags      = ["ubuntu"]
}
```

## Argument Reference

- `node_name` - (Optional) The node name.
- `tags` - (Optional) A list of tags to filter the VMs. The VM must have all
  the tags to be included in the result.

## Attribute Reference

- `vms` - The VMs list.
  - `name` - The virtual machine name.
  - `node_name` - The node name.
  - `tags` - A list of tags of the VM.
  - `vm_id` - The VM identifier.
