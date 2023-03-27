---
layout: page
title: proxmox_virtual_environment_cluster_firewall_alias
permalink: /data-sources/virtual_environment_cluster_firewall_alias
nav_order: 1
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_cluster_firewall_alias

Retrieves information about a specific alias.

## Example Usage

```terraform
data "proxmox_virtual_environment_cluster_firewall_alias" "local_network" {
  name = "local_network"
}
```

## Argument Reference

* `name` - (Required) Alias name.

## Attribute Reference

* `cidr` - (Required) Network/IP specification in CIDR format.
* `comment` - (Optional) Alias comment.
