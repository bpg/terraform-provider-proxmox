---
layout: page
title: proxmox_virtual_environment_containers
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_containers

Retrieves information about all containers in the Proxmox cluster.

## Example Usage

```hcl
data "proxmox_virtual_environment_containers" "ubuntu_containers" {
  tags      = ["ubuntu"]
}

data "proxmox_virtual_environment_containers" "ubuntu_templates" {
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

- `node_name` - (Optional) The node name. All cluster nodes will be queried in case this is omitted
- `tags` - (Optional) A list of tags to filter the containers. The container must have all
  the tags to be included in the result.
- `filter` - (Optional) Filter blocks. The container must satisfy all filter blocks to be included in the result.
    - `name` - Name of the container attribute to filter on. One of [`name`, `template`, `status`, `node_name`]
    - `values` - List of values to pass the filter. Container's attribute should match at least one value in the list.

## Attribute Reference

- `containers` - The containers list.
    - `name` - The container name.
    - `node_name` - The node name.
    - `tags` - A list of tags of the container.
    - `vm_id` - The container identifier.
    - `status` - Status of the container
    - `template` - Is container a template (true) or a regular container (false)