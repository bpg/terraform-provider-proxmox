---
layout: page
title: proxmox_virtual_environment_cluster_firewall_security_groups
permalink: /data-sources/virtual_environment_cluster_firewall_security_groups
nav_order: 6
parent: Data Sources
subcategory: Virtual Environment
---

# Data Source: proxmox_virtual_environment_cluster_firewall_security_groups

Retrieves the names for all the available security groups.

## Example Usage

```terraform
data "proxmox_virtual_environment_cluster_firewall_security_groups" "available_sgs" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

- `security_group_names` - The security group names.
