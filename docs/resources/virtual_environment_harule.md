---
layout: page
title: proxmox_virtual_environment_harule
parent: Resources
subcategory: Virtual Environment
description: |-
  Manages a High Availability rule in a Proxmox VE cluster.
---

# Resource: proxmox_virtual_environment_harule

Manages a High Availability rule in a Proxmox VE cluster.

~> **Note:** This resource requires Proxmox VE 9.0 or later. In PVE 9, HA groups
have been replaced by HA rules, which provide node affinity and resource affinity
capabilities. For PVE 8 and earlier, use
[`proxmox_virtual_environment_hagroup`](virtual_environment_hagroup.md) instead.

## Example Usage

### Node Affinity Rule

Assign VMs to preferred nodes with priorities. Non-strict rules allow failover
to other nodes; strict rules do not.

```hcl
resource "proxmox_virtual_environment_harule" "prefer_node1" {
  rule      = "prefer-node1"
  type      = "node-affinity"
  comment   = "Prefer node1 for these VMs"
  resources = ["vm:100", "vm:101"]

  nodes = {
    node1 = 2  # Higher priority
    node2 = 1
    node3 = 1
  }

  strict = false
}
```

### Resource Affinity Rule (Positive)

Keep resources together on the same node.

```hcl
resource "proxmox_virtual_environment_harule" "keep_together" {
  rule      = "db-cluster-together"
  type      = "resource-affinity"
  comment   = "Keep database replicas on the same node"
  resources = ["vm:200", "vm:201"]
  affinity  = "positive"
}
```

### Resource Affinity Rule (Negative / Anti-Affinity)

Keep resources on separate nodes for high availability.

```hcl
resource "proxmox_virtual_environment_harule" "keep_apart" {
  rule      = "db-cluster-apart"
  type      = "resource-affinity"
  comment   = "Spread database replicas across nodes"
  resources = ["vm:200", "vm:201", "vm:202"]
  affinity  = "negative"
}
```

## Argument Reference

- `rule` - (Required) The identifier of the HA rule. Must start with a letter,
  end with a letter or number, and be composed of letters, numbers, `-`, `_` and
  `.`. Must be at least 2 characters long. Changing this forces a new resource.
- `type` - (Required) The HA rule type. Must be `node-affinity` or
  `resource-affinity`. Changing this forces a new resource.
- `resources` - (Required) A set of HA resource IDs that this rule applies to
  (e.g. `vm:100`, `ct:101`). The resources must already be managed by HA
  (i.e. they must exist as
  [`proxmox_virtual_environment_haresource`](virtual_environment_haresource.md)).
- `comment` - (Optional) A comment associated with the rule.
- `disable` - (Optional) Whether the HA rule is disabled. Defaults to `false`.

### Node Affinity Arguments

These arguments are only applicable when `type` is `node-affinity`:

- `nodes` - (Required) A map of cluster node names to their priorities. Higher
  values indicate higher priority. Use `null` for unset priorities.
- `strict` - (Optional) Whether the node affinity rule is strict. When strict,
  resources cannot run on nodes not listed. Defaults to `false`.

### Resource Affinity Arguments

These arguments are only applicable when `type` is `resource-affinity`:

- `affinity` - (Required) The affinity type. `positive` keeps resources on the
  same node. `negative` keeps resources on separate nodes.

## Attribute Reference

- `id` - The unique identifier of the HA rule.

## Import

An existing HA rule can be imported using its identifier:

```bash
terraform import proxmox_virtual_environment_harule.example rule-name
```
