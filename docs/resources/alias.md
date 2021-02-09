---
layout: page
title: Alias
permalink: /ressources/alias
nav_order: 1
parent: Resources
subcategory: Virtual Environment
---

# Resource: Alias

Aliases are used to see what devices or group of devices are affected by a rule.
We can create aliases to identify an IP address or a network.

## Example Usage

```
resource "proxmox_virtual_environment_cluster_alias" "local_network" {
	name    = "local_network"
	cidr    = "192.168.0.0/23"
	comment = "Managed by Terraform"
}

resource "proxmox_virtual_environment_cluster_alias" "ubuntu_vm" {
	name    = "ubuntu"
	cidr    = "192.168.0.1"
	comment = "Managed by Terraform"
}
```

## Argument Reference

* `name` - (Required) Alias name.
* `cidr` - (Required) Network/IP specification in CIDR format.
* `comment` - (Optional) Alias comment.

## Attribute Reference

There are no attribute references available for this resource.
