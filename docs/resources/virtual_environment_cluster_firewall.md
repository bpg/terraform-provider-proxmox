---
layout: page
title: proxmox_virtual_environment_cluster_firewall
permalink: /resources/virtual_environment_cluster_firewall
nav_order: 2
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_cluster_firewall

Manages firewall options on the cluster level.

## Example Usage

```terraform
resource "proxmox_virtual_environment_cluster_firewall" "example" {
  enabled = false

  ebtables      = false
  input_policy  = "DROP"
  output_policy = "ACCEPT"
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
- `input_policy` - (Optional) The default input policy (`ACCEPT`, `DROP`, `REJECT`).
- `output_policy` - (Optional) The default output policy (`ACCEPT`, `DROP`, `REJECT`).
- `log_ratelimit` - (Optional) The log rate limit.
    - `enabled` - (Optional) Enable or disable the log rate limit.
    - `burst` - (Optional) Initial burst of packages which will always get
      logged before the rate is applied (defaults to `5`).
    - `rate` - (Optional) Frequency with which the burst bucket gets refilled (defaults to `1/second`).

## Attribute Reference

There are no additional attributes available for this resource.
