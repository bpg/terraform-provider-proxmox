---
layout: page
title: proxmox_virtual_environment_cluster_firewall
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_cluster_firewall

Manages firewall options on the cluster level.

## Example Usage

```hcl
resource "proxmox_virtual_environment_cluster_firewall" "example" {
  enabled = false

  ebtables       = false
  input_policy   = "DROP"
  output_policy  = "ACCEPT"
  forward_policy = "ACCEPT"
  log_ratelimit {
    enabled = false
    burst   = 10
    rate    = "5/second"
  }
}
```

## Argument Reference

- `enabled` - (Optional) Enable or disable the firewall cluster wide.
- `ebtables` - (Optional) Enable ebtables rules cluster wide.
- `input_policy` - (Optional) The default input policy (`ACCEPT`, `DROP`, `REJECT`). Defaults to `DROP`.
- `output_policy` - (Optional) The default output policy (`ACCEPT`, `DROP`, `REJECT`). Defaults to `ACCEPT`.
- `forward_policy` - (Optional) The default forward policy (`ACCEPT`, `DROP`). Defaults to `ACCEPT`.
- `log_ratelimit` - (Optional) The log rate limit.
    - `enabled` - (Optional) Enable or disable the log rate limit.
    - `burst` - (Optional) Initial burst of packages which will always get
        logged before the rate is applied (defaults to `5`).
    - `rate` - (Optional) Frequency with which the burst bucket gets refilled
        (defaults to `1/second`).

## Attribute Reference

There are no additional attributes available for this resource.

## Important Notes

This resource manages cluster-wide firewall options, so it should be used only
once per cluster. Declaring it multiple times results in conflicting updates.

## Import

Instances can be imported without an ID, but you still need to pass one, e.g.,

```bash
terraform import proxmox_virtual_environment_cluster_firewall.example example
```
