---
layout: page
title: proxmox_virtual_environment_cluster_firewall_security_group
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_cluster_firewall_security_group

A security group is a collection of rules, defined at cluster level, which can
be used in all VMs' rules. For example, you can define a group named “webserver”
with rules to open the http and https ports.

## Example Usage

```hcl
resource "proxmox_virtual_environment_cluster_firewall_security_group" "webserver" {
  name    = "webserver"
  comment = "Managed by Terraform"

  rule {
    type    = "in"
    action  = "ACCEPT"
    comment = "Allow HTTP"
    dest    = "192.168.1.5"
    dport   = "80"
    proto   = "tcp"
    log     = "info"
  }

  rule {
    type    = "in"
    action  = "ACCEPT"
    comment = "Allow HTTPS"
    dest    = "192.168.1.5"
    dport   = "443"
    proto   = "tcp"
    log     = "info"
  }
}
```

## Argument Reference

- `name` - (Required) Security group name.
- `comment` - (Optional) Security group comment.
- `rule` - (Optional) Firewall rule block (multiple blocks supported).
    - `action` - (Required) Rule action (`ACCEPT`, `DROP`, `REJECT`).
    - `type` - (Required) Rule type (`in`, `out`).
    - `comment` - (Optional) Rule comment.
    - `dest` - (Optional) Restrict packet destination address. This can refer to
          a single IP address, an IP set ('+ipsetname') or an IP alias
          definition. You can also specify an address range like
          `20.34.101.207-201.3.9.99`, or a list of IP addresses and networks
          (entries are separated by comma). Please do not mix IPv4 and IPv6
          addresses inside such lists.
    - `dport` - (Optional) Restrict TCP/UDP destination port. You can use
        service names or simple numbers (0-65535), as defined in '/etc/
        services'. Port ranges can be specified with '\d+:\d+', for example
        `80:85`, and you can use comma separated list to match several ports or
        ranges.
    - `enable` - (Optional) Enable this rule. Defaults to `true`.
    - `iface` - (Optional) Network interface name. You have to use network
        configuration key names for VMs and containers ('net\d+'). Host related
        rules can use arbitrary strings.
    - `log` - (Optional) Log level for this rule (`emerg`, `alert`, `crit`,
        `err`, `warning`, `notice`, `info`, `debug`, `nolog`).
    - `macro`- (Optional) Macro name. Use predefined standard macro
        from <https://pve.proxmox.com/pve-docs/pve-admin-guide.html#_firewall_macro_definitions>
    - `proto` - (Optional) Restrict packet protocol. You can use protocol names
        as defined in '/etc/protocols'.
    - `source` - (Optional) Restrict packet source address. This can refer
        to a single IP address, an IP set ('+ipsetname') or an IP alias
        definition. You can also specify an address range like
        `20.34.101.207-201.3.9.99`, or a list of IP addresses and networks (
        entries are separated by comma). Please do not mix IPv4 and IPv6
        addresses inside such lists.
    - `sport` - (Optional) Restrict TCP/UDP source port. You can use
        service names or simple numbers (0-65535), as defined in '/etc/
        services'. Port ranges can be specified with '\d+:\d+', for example
        `80:85`, and you can use comma separated list to match several ports or
        ranges.

## Attribute Reference

- `rule`
    - `pos` - Position of the rule in the list.

There are no attribute references available for this resource.

## Import

Instances can be imported using the `name`, e.g.,

```bash
terraform import proxmox_virtual_environment_cluster_firewall_security_group.webserver webserver
```
