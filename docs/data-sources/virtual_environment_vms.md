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
    regex  = true
    values = ["^ubuntu-20.*$"]
  }

  filter {
    name   = "node_name"
    regex  = true
    values = ["node_us_[1-3]", "node_eu_[1-3]"]
  }
}
```

## Argument Reference

- `node_name` - (Optional) The node name. If omitted, all cluster nodes are queried.
- `tags` - (Optional) A list of tags to filter the VMs. The VM must have all
  the tags to be included in the result.
- `filter` - (Optional) Filter blocks. The VM must satisfy all filter blocks to be included in the result.
    - `name` - Name of the VM attribute to filter on. One of [`name`, `template`, `status`, `node_name`]
    - `values` - List of values to pass the filter. VM's attribute should match at least one value in the list.
    - `regex` - (Optional) Whether to treat the `values` as regex patterns (defaults to `false`).

## Attribute Reference

- `vms` - The VMs list.
    - `name` - The virtual machine name.
    - `node_name` - The node name.
    - `tags` - A list of tags of the VM.
    - `vm_id` - The VM identifier.
    - `status` - The status of the VM.
    - `template` - Whether the VM is a template.
