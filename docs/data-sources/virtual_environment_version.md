---
layout: page
title: proxmox_virtual_environment_version
parent: Data Sources
subcategory: Virtual Environment
description: |-
  Retrieves API version details.
---

# Data Source: proxmox_virtual_environment_version

Retrieves API version details.

## Example Usage

```terraform
data "proxmox_virtual_environment_version" "example" {}

output "data_proxmox_virtual_environment_version" {
  value = {
    release       = data.proxmox_virtual_environment_version.example.release
    repository_id = data.proxmox_virtual_environment_version.example.repository_id
    version       = data.proxmox_virtual_environment_version.example.version
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) Placeholder identifier attribute.
- `release` (String) The current Proxmox VE point release in `x.y` format.
- `repository_id` (String) The short git revision from which this version was build.
- `version` (String) The full pve-manager package version of this node.
