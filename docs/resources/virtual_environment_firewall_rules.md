---
layout: page
title: proxmox_virtual_environment_firewall_rules
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_firewall_rules

Manages cluster-level and VM/container-level firewall rules.

## Example Usage

```hcl
resource "proxmox_virtual_environment_firewall_rules" "inbound" {
  depends_on = [
    proxmox_virtual_environment_vm.example,
    proxmox_virtual_environment_cluster_firewall_security_group.example,
  ]

  node_name = proxmox_virtual_environment_vm.example.node_name
  vm_id     = proxmox_virtual_environment_vm.example.vm_id

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

  rule {
    security_group = proxmox_virtual_environment_cluster_firewall_security_group.example.name
    comment        = "From security group"
    iface          = "net0"
  }
}
```

## Argument Reference

- `node_name` - (Optional) Node name. Leave empty for cluster level rules.
- `vm_id` - (Optional) VM ID. Leave empty for cluster level rules.
- `container_id` - (Optional) Container ID. Leave empty for cluster level
    rules.
- `rule` - (Optional) Firewall rule block (multiple blocks supported).
    The provider supports two types of the `rule` blocks:
    - A rule definition block, which includes the following arguments:
        - `action` - (Required) Rule action (`ACCEPT`, `DROP`, `REJECT`).
        - `type` - (Required) Rule type (`in`, `out`, `forward`).
        - `comment` - (Optional) Rule comment.
        - `dest` - (Optional) Restrict packet destination address. This can
            refer to a single IP address, an IP set ('+ipsetname') or an IP
            alias definition. You can also specify an address range
            like `20.34.101.207-201.3.9.99`, or a list of IP addresses and
            networks (entries are separated by comma). Please do not mix IPv4
            and IPv6 addresses inside such lists.
        - `dport` - (Optional) Restrict TCP/UDP destination port. You can use
            service names or simple numbers (0-65535), as defined
            in `/etc/services`. Port ranges can be specified with '\d+:\d+', for
            example `80:85`, and you can use comma separated list to match
            several ports or ranges.
        - `enabled` - (Optional) Enable this rule. Defaults to `true`.
        - `iface` - (Optional) Network interface name. You have to use network
            configuration key names for VMs and containers ('net\d+'). Host
            related rules can use arbitrary strings.
        - `log` - (Optional) Log level for this rule (`emerg`, `alert`, `crit`,
            `err`, `warning`, `notice`, `info`, `debug`, `nolog`).
        - `macro`- (Optional) Macro name. Use predefined standard macro
            from <https://pve.proxmox.com/pve-docs/pve-admin-guide.html#_firewall_macro_definitions>
        - `proto` - (Optional) Restrict packet protocol. You can use protocol
            names as defined in '/etc/protocols'.
        - `source` - (Optional) Restrict packet source address. This can refer
            to a single IP address, an IP set ('+ipsetname') or an IP alias
            definition. You can also specify an address range
            like `20.34.101.207-201.3.9.99`, or a list of IP addresses and
            networks (entries are separated by comma). Please do not mix IPv4
            and IPv6 addresses inside such lists.
        - `sport` - (Optional) Restrict TCP/UDP source port. You can use
            service names or simple numbers (0-65535), as defined
            in `/etc/services`. Port ranges can be specified with '\d+:\d+', for
            example `80:85`, and you can use comma separated list to match
            several ports or ranges.
    - a security group insertion block, which includes the following arguments:
        - `comment` - (Optional) Rule comment.
        - `enabled` - (Optional) Enable this rule. Defaults to `true`.
        - `iface` - (Optional) Network interface name. You have to use network
            configuration key names for VMs and containers ('net\d+'). Host
            related rules can use arbitrary strings.
        - `security_group` - (Required) Security group name.

## Attribute Reference

- `rule`
    - `pos` - Position of the rule in the list.

## Import

### Cluster Rules
Use the import ID: `cluster`

**Example:**
```bash
terraform import proxmox_virtual_environment_firewall_rules.cluster_rules cluster
```

### VM Rules
Use the import ID format: `vm/<node_name>/<vm_id>`
Example uses node name `pve` and VM ID `100`.

**Example:**
```bash
terraform import proxmox_virtual_environment_firewall_rules.vm_rules vm/pve/100
```

### Container Rules
Use the import ID format: `container/<node_name>/<container_id>`
Example uses node name `pve` and container ID `100`.

**Example:**
```bash
terraform import proxmox_virtual_environment_firewall_rules.container_rules container/pve/100
```
