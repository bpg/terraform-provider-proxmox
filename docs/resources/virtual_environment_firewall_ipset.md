---
layout: page
title: proxmox_virtual_environment_firewall_ipset
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_firewall_ipset

An IPSet allows us to group multiple IP addresses, IP subnets and aliases. IPSets can be
created on the cluster level, on VM / Container level.

## Example Usage

```hcl
resource "proxmox_virtual_environment_firewall_ipset" "ipset" {
  depends_on = [proxmox_virtual_environment_vm.example]

  node_name = proxmox_virtual_environment_vm.example.node_name
  vm_id     = proxmox_virtual_environment_vm.example.vm_id

  name    = "local_network"
  comment = "Managed by Terraform"

  cidr {
    name    = "192.168.0.0/23"
    comment = "Local network 1"
  }

  cidr {
    name    = "192.168.0.1"
    comment = "Server 1"
    nomatch = true
  }

  cidr {
    name    = "192.168.2.1"
    comment = "Server 1"
  }
}
```

## Argument Reference

- `node_name` - (Optional) Node name. Leave empty for cluster level ipsets.
- `vm_id` - (Optional) VM ID. Leave empty for cluster level ipsets.
- `container_id` - (Optional) Container ID. Leave empty for cluster level ipsets.
- `name` - (Required) IPSet name.
- `comment` - (Optional) IPSet comment.
- `cidr` - (Optional) IP/CIDR block (multiple blocks supported).
    - `name` - Network/IP specification in CIDR format.
    - `comment` - (Optional) Arbitrary string annotation.
    - `nomatch` - (Optional) Entries marked as `nomatch` are skipped as if those
        were not added to the set.

## Attribute Reference

There are no attribute references available for this resource.


## Import

### Cluster IPSet
Use the import ID: `cluster/<ipset_name>`
Example uses ipset name `local_network`.

**Example:**
```bash
terraform import proxmox_virtual_environment_firewall_ipset.cluster_ipset cluster/local_network
```

### VM IPSet
Use the import ID format: `vm/<node_name>/<vm_id>/<ipset_name>`
Example uses node name `pve` and VM ID `100` and ipset name `local_network`.

**Example:**
```bash
terraform import proxmox_virtual_environment_firewall_ipset.vm_ipset vm/pve/100/local_network
```

### Container IPSet
Use the import ID format: `container/<node_name>/<container_id>/<ipset_name>`
Example uses node name `pve` and container ID `100` and ipset name `local_network`.

**Example:**
```bash
terraform import proxmox_virtual_environment_firewall_ipset.container_ipset container/pve/100/local_network
```
