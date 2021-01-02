---
layout: page
title: Alias
permalink: /data-sources/alias
nav_order: 1
parent: Data Sources
---

# Data Source: Alias

Retrieves information about a specific alias.

## Example Usage

```
data "proxmox_virtual_environment_cluster_alias" "local_network" {
  name    = "local_network"
}
```

## Arguments Reference

* `name` - (Required) Alias name.

## Attributes Reference

* `cidr` - (Required) Network/IP specification in CIDR format.
* `comment` - (Optional) Alias comment.
