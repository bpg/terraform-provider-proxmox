---
layout: page
title: proxmox_virtual_environment_container
permalink: /resources/virtual_environment_container
nav_order: 4
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_container

Manages a container.

## Example Usage

```terraform
resource "proxmox_virtual_environment_container" "ubuntu_container" {
  description = "Managed by Terraform"

  node_name = "first-node"
  vm_id     = 1234

  initialization {
    hostname = "terraform-provider-proxmox-ubuntu-container"

    ip_config {
      ipv4 {
        address = "dhcp"
      }
    }

    user_account {
      keys = [
        trimspace(tls_private_key.ubuntu_container_key.public_key_openssh)
      ]
      password = random_password.ubuntu_container_password.result
    }
  }

  network_interface {
    name = "veth0"
  }

  operating_system {
    template_file_id = proxmox_virtual_environment_file.ubuntu_container_template.id
    type             = "ubuntu"
  }
}

resource "proxmox_virtual_environment_file" "ubuntu_container_template" {
  content_type = "vztmpl"
  datastore_id = "local"
  node_name    = "first-node"

  source_file {
    path = "http://download.proxmox.com/images/system/ubuntu-20.04-standard_20.04-1_amd64.tar.gz"
  }
}

resource "random_password" "ubuntu_container_password" {
  length           = 16
  override_special = "_%@"
  special          = true
}

resource "tls_private_key" "ubuntu_container_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

output "ubuntu_container_password" {
  value     = random_password.ubuntu_container_password.result
  sensitive = true
}

output "ubuntu_container_private_key" {
  value     = tls_private_key.ubuntu_container_key.private_key_pem
  sensitive = true
}

output "ubuntu_container_public_key" {
  value = tls_private_key.ubuntu_container_key.public_key_openssh
}
```

## Argument Reference

- `clone` - (Optional) The cloning configuration.
    - `datastore_id` - (Optional) The identifier for the target datastore.
    - `node_name` - (Optional) The name of the source node (leave blank, if
      equal to the `node_name` argument).
    - `vm_id` - (Required) The identifier for the source container.
- `console` - (Optional) The console configuration.
    - `enabled` - (Optional) Whether to enable the console device (defaults
      to `true`).
    - `mode` - (Optional) The console mode (defaults to `tty`).
        - `console` - Console.
        - `shell` - Shell.
        - `tty` - TTY.
    - `tty_count` - (Optional) The number of available TTY (defaults to `2`).
- `cpu` - (Optional) The CPU configuration.
    - `architecture` - (Optional) The CPU architecture (defaults to `amd64`).
        - `amd64` - x86 (64 bit).
        - `arm64` - ARM (64-bit).
        - `armhf` - ARM (32 bit).
        - `i386` - x86 (32 bit).
    - `cores` - (Optional) The number of CPU cores (defaults to `1`).
    - `units` - (Optional) The CPU units (defaults to `1024`).
- `description` - (Optional) The description.
- `disk` - (Optional) The disk configuration.
    - `datastore_id` - (Optional) The identifier for the datastore to create the
      disk in (defaults to `local`).
      -`size` - (Optional) The size of the root filesystem in gigabytes (
      defaults to `4`). Requires `datastore_id` to be set.
- `initialization` - (Optional) The initialization configuration.
    - `dns` - (Optional) The DNS configuration.
        - `domain` - (Optional) The DNS search domain.
        - `server` - (Optional) The DNS server.
    - `hostname` - (Optional) The hostname.
    - `ip_config` - (Optional) The IP configuration (one block per network
      device).
        - `ipv4` - (Optional) The IPv4 configuration.
            - `address` - (Optional) The IPv4 address (use `dhcp` for
              autodiscovery).
            - `gateway` - (Optional) The IPv4 gateway (must be omitted
              when `dhcp` is used as the address).
        - `ipv6` - (Optional) The IPv4 configuration.
            - `address` - (Optional) The IPv6 address (use `dhcp` for
              autodiscovery).
            - `gateway` - (Optional) The IPv6 gateway (must be omitted
              when `dhcp` is used as the address).
    - `user_account` - (Optional) The user account configuration.
        - `keys` - (Optional) The SSH keys for the root account.
        - `password` - (Optional) The password for the root account.
- `memory` - (Optional) The memory configuration.
    - `dedicated` - (Optional) The dedicated memory in megabytes (defaults
      to `512`).
    - `swap` - (Optional) The swap size in megabytes (defaults to `0`).
- `network_interface` - (Optional) A network interface (multiple blocks
  supported).
    - `bridge` - (Optional) The name of the network bridge (defaults
      to `vmbr0`).
    - `enabled` - (Optional) Whether to enable the network device (defaults
      to `true`).
    - `mac_address` - (Optional) The MAC address.
    - `mtu` - (Optional) Maximum transfer unit of the interface. Cannot be
      larger than the bridge's MTU.
    - `name` - (Required) The network interface name.
    - `rate_limit` - (Optional) The rate limit in megabytes per second.
    - `vlan_id` - (Optional) The VLAN identifier.
- `node_name` - (Required) The name of the node to assign the container to.
- `operating_system` - (Required) The Operating System configuration.
    - `template_file_id` - (Required) The identifier for an OS template file.
    - `type` - (Optional) The type (defaults to `unmanaged`).
        - `alpine` - Alpine.
        - `archlinux` - Arch Linux.
        - `centos` - CentOS.
        - `debian` - Debian.
        - `fedora` - Fedora.
        - `gentoo` - Gentoo.
        - `opensuse` - openSUSE.
        - `ubuntu` - Ubuntu.
        - `unmanaged` - Unmanaged.
- `pool_id` - (Optional) The identifier for a pool to assign the container to.
- `started` - (Optional) Whether to start the container (defaults to `true`).
- `tags` - (Optional) A list of tags of the container. This is only meta
  information (defaults to `[]`). Note: Proxmox always sorts the container tags.
  If the list in template is not sorted, then Proxmox will always report a
  difference on the resource. You may use the `ignore_changes` lifecycle
  meta-argument to ignore changes to this attribute.
- `template` - (Optional) Whether to create a template (defaults to `false`).
- `unprivileged` - (Optional) Whether the container runs as unprivileged on
the host (defaults to `false`).
- `vm_id` - (Optional) The virtual machine identifier

## Attribute Reference

There are no additional attributes available for this resource.
