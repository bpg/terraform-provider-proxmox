---
layout: page
title: proxmox_virtual_environment_time
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_time

Manages the time for a specific node.

## Example Usage

```hcl
resource "proxmox_virtual_environment_time" "first_node_time" {
  node_name = "first-node"
  time_zone = "UTC"
}
```

## Argument Reference

- `node_name` - (Required) A node name.
- `time_zone` - (Required) The node's time zone.

## Attribute Reference

- `local_time` - The node's local time.
- `utc_time` - The node's local time formatted as UTC.

## Import

Instances can be imported using the `node_name`, e.g.,

```bash
terraform import proxmox_virtual_environment_dns.first_node first-node
```
