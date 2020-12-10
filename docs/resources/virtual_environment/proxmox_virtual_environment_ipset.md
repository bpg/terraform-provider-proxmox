---
layout: page
title: IPSet
permalink: /ressources/virtual-environment/ipset
nav_order: 12
parent: Virtual Environment Resources
grand_parent: Resources
---

# Resource: IPSet

An IPSet allows us to group multiple IP addresses, IP subnets and aliases.

## Example Usage

```
resource "proxmox_virtual_environment_cluster_ipset" "ipset" {
	name    = "local_network"
	comment = "Managed by Terraform"
    
    ipset {
        cidr = "192.168.0.0/23"
        comment = "Local network 1"
    }
    
    ipset {
        cidr = "192.168.0.1"
        comment = "Server 1"
        nomatch = true
    }
    
    ipset {
        cidr = "192.168.2.1"
        comment = "Server 1"
    }
}
```

## Arguments Reference

* `name` - (Required) Alias name.
* `comment` - (Optional) Alias comment.
* `ipset` - (Optional) IP/CIDR block (multiple blocks supported).
    * `cidr` - Network/IP specification in CIDR format.
    * `comment` - (Optional) Arbitrary string annotation.
    * `nomatch` -  (Optional) Entries marked as `nomatch` are skipped as if those were not added to the set.

## Attributes Reference

There are no attribute references available for this resource.
