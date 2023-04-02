---
layout: page
title: proxmox_virtual_environment_cluster_firewall_alias
permalink: /resources/virtual_environment_cluster_firewall_alias
nav_order: 7
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_cluster_firewall_alias

Aliases are used to see what devices or group of devices are affected by a rule.
We can create aliases to identify an IP address or a network. Aliases can be
created on the cluster level, on VM / Container level.

## Example Usage

```terraform
resource "proxmox_virtual_environment_cluster_firewall_alias" "local_network" {
  depends_on = [proxmox_virtual_environment_vm.example]

  node_name = proxmox_virtual_environment_vm.example.node_name
  vm_id     = proxmox_virtual_environment_vm.example.vm_id

  name    = "local_network"
  cidr    = "192.168.0.0/23"
  comment = "Managed by Terraform"
}

resource "proxmox_virtual_environment_cluster_firewall_alias" "ubuntu_vm" {
  name    = "ubuntu"
  cidr    = "192.168.0.1"
  comment = "Managed by Terraform"
}
```

## Argument Reference

- `node_name` - (Optional) Node name. Leave empty for cluster level aliases.
- `vm_id` - (Optional) VM ID. Leave empty for cluster level aliases.
- `container_id` - (Optional) Container ID. Leave empty for cluster level aliases.
- `name` - (Required) Alias name.
- `cidr` - (Required) Network/IP specification in CIDR format.
- `comment` - (Optional) Alias comment.

## Attribute Reference

There are no attribute references available for this resource.
