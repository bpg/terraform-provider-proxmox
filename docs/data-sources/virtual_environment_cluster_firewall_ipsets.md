---
layout: page
title: proxmox_virtual_environment_cluster_firewall_ipsets
permalink: /data-sources/virtual_environment_cluster_firewall_ipsets
nav_order: 4
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_cluster_firewall_ipsets

Retrieves the names for all the available IPSets.

## Example Usage

```terraform
data "proxmox_virtual_environment_cluster_firewall_ipsets" "available_ipsets" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

- `ipset_names` - The IPSet names.
