---
layout: page
title: proxmox_virtual_environment_firewall_alias
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_firewall_alias

Aliases are used to see what devices or group of devices are affected by a rule.
We can create aliases to identify an IP address or a network. Aliases can be
created on the cluster level, on VM / Container level.

## Example Usage

```hcl
resource "proxmox_virtual_environment_firewall_alias" "local_network" {
  depends_on = [proxmox_virtual_environment_vm.example]

  node_name = proxmox_virtual_environment_vm.example.node_name
  vm_id     = proxmox_virtual_environment_vm.example.vm_id

  name    = "local_network"
  cidr    = "192.168.0.0/23"
  comment = "Managed by Terraform"
}

resource "proxmox_virtual_environment_firewall_alias" "ubuntu_vm" {
  name    = "ubuntu"
  cidr    = "192.168.0.1"
  comment = "Managed by Terraform"
}
```

## Argument Reference

- `node_name` - (Optional) Node name. Leave empty for cluster level aliases. Changing it forces replacement unless `vm_id` or `container_id` is set, so it can change in place after a migration.
- `vm_id` - (Optional) VM ID. Leave empty for cluster level aliases. Changing this value forces replacement.
- `container_id` - (Optional) Container ID. Leave empty for cluster level aliases. Changing this value forces replacement.
- `name` - (Required) Alias name.
- `cidr` - (Required) Network/IP specification in CIDR format.
- `comment` - (Optional) Alias comment.

## Attribute Reference

There are no attribute references available for this resource.

## Import

### Cluster Alias

Use the import ID: `cluster/<alias_name>`
Example uses alias name `ubuntu`.

**Example:**

```bash
terraform import proxmox_virtual_environment_firewall_alias.ubuntu_vm cluster/ubuntu
```

### VM Alias

Use the import ID format: `vm/<node_name>/<vm_id>/<alias_name>`
Example uses node name `pve` and VM ID `100` and alias name `local_network`.

**Example:**

```bash
terraform import proxmox_virtual_environment_firewall_alias.local_network vm/pve/100/local_network
```

### Container Alias

Use the import ID format: `container/<node_name>/<container_id>/<alias_name>`
Example uses node name `pve` and container ID `100` and alias name `local_network`.

**Example:**

```bash
terraform import proxmox_virtual_environment_firewall_alias.container_alias container/pve/100/local_network
```
