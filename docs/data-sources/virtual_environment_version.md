---
layout: page
title: proxmox_virtual_environment_version
permalink: /data-sources/virtual_environment_version
nav_order: 16
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_version

Retrieves the version information from the API endpoint.

## Example Usage

```terraform
data "proxmox_virtual_environment_version" "current_version" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

- `keyboard_layout` - The keyboard layout.
- `release` - The release number.
- `repository_id` - The repository identifier.
- `version` - The version string.
