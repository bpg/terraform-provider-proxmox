---
layout: page
title: Version
permalink: /data-sources/version
nav_order: 16
parent: Data Sources
---

# Data Source: Version

Retrieves the version information from the API endpoint.

## Example Usage

```
data "proxmox_virtual_environment_version" "current_version" {}
```

## Arguments Reference

There are no arguments available for this data source.

## Attributes Reference

* `keyboard_layout` - The keyboard layout.
* `release` - The release number.
* `repository_id` - The repository identifier.
* `version` - The version string.
