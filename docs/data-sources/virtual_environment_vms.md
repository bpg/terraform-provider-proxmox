---
layout: page
title: proxmox_virtual_environment_vms
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_vms

Retrieves information about all VMs in the Proxmox cluster.

## Example Usage

```hcl
data "proxmox_virtual_environment_vms" "ubuntu_vms" {
  tags      = ["ubuntu"]
}

data "proxmox_virtual_environment_vms" "ubuntu_templates" {
  tags      = ["template", "latest"]
  filter {
    name   = "template"
    values = [true]
  }

  filter {
    name   = "status"
    values = ["stopped"]
  }

  filter {
    name   = "name"
    values = ["^ubuntu-20.*$", "^ubuntu-22.04$"]
  }
}
```

## Argument Reference

- `node_name` - (Optional) The node name.
- `tags` - (Optional) A list of tags to filter the VMs. The VM must have all
  the tags to be included in the result.
- `filter` - (Optional) Filter blocks. All blocks should match to pass the filter (AND logic)
    - `name` - Name of the VM attribute to filter on. One of [`name`, `template`, `status`]
    - `value` - List of values to pass the filter (OR logic)

## Attribute Reference

- `vms` - The VMs list.
    - `name` - The virtual machine name.
    - `node_name` - The node name.
    - `tags` - A list of tags of the VM.
    - `vm_id` - The VM identifier.
    - `status` - Status of the VM
    - `template` - Is VM a template (true) or a regular VM (false)
