---
layout: page
title: proxmox_virtual_environment_cluster_firewall_ipset
permalink: /data-sources/virtual_environment_cluster_firewall_ipset
nav_order: 1
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_cluster_firewall_ipset

Retrieves information about a specific IPSet.

## Example Usage

```terraform
data "proxmox_virtual_environment_cluster_firewall_ipset" "local_network" {
  name = "local_network"
}
```

## Argument Reference

- `name` - (Required) IPSet name.

## Attribute Reference

- `cidr` - (Optional) IP/CIDR list.
    - `name` - Network/IP specification in CIDR format.
    - `comment` - (Optional) Arbitrary string annotation.
    - `nomatch` - (Optional) Entries marked as `nomatch` are skipped as if those
      were not added to the set.
- `comment` - (Optional) IPSet comment.
