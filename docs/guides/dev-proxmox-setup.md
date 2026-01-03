---
layout: page
page_title: "Development Proxmox Setup"
subcategory: Guides
description: |-
  How to set up Proxmox VE in a local VM for running examples and acceptance tests.
---

# Setting up Proxmox in a VM for development

To test provider changes, you need access to a real Proxmox cluster. This guide walks you through setting up Proxmox VE inside a local VM using [virt-manager](https://virt-manager.org/).

## Prerequisites

- Go and Terraform installed on your system.
- Linux with KVM/QEMU support (tested on Debian 12, but other distros should work).
- At least 4 GB RAM and 30 GB disk space available for the VM.

## Installation steps

1. Install virt-manager:

   ```sh
   sudo apt-get install virt-manager
   ```

2. Download the latest Proxmox VE ISO from <https://www.proxmox.com/en/downloads>.

3. Open virt-manager and create a new virtual machine:
   - Use the downloaded ISO.
   - Select Debian as the operating system.
   - Allocate at least 4 GB RAM and 30 GB disk.
   - Keep default network settings.

4. Complete the Proxmox installation. Choose your preferred timezone, country, password, domain, and email. Keep other settings at their defaults.

5. After installation, log in as `root` with the password you set. Run `ip a` to find the assigned IP address:

   ```text
   root@proxmox:~# ip a
   ...
   3: vmbr0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 ...
       inet 192.168.122.43/24 scope global vmbr0
   ...
   ```

6. (Optional) Verify connectivity from your host:

   ```sh
   ssh root@192.168.122.43
   ```

   You can also access the web console at `https://192.168.122.43:8006`.

## Configuring the provider

Create a `terraform.tfvars` file in the `example/` directory (this file is git-ignored):

```hcl
virtual_environment_ssh_username = "your-ssh-username"
virtual_environment_endpoint     = "https://192.168.122.43:8006/"
virtual_environment_password     = "your-password"
virtual_environment_api_token    = "root@pam!<token>=<value>"
```

Replace the IP address and password with your actual values.

## Running examples

Run the example configurations with:

```sh
make example
```

This applies all resources defined in the `example/` directory against your Proxmox instance.

## Common issues

### Snippets content type not enabled

If you see an error like:

```text
the datastore "local" does not support content type "snippets"
```

You need to enable snippets on the datastore. See the [file resource documentation](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/resources/virtual_environment_file#snippets) for instructions.

### Network connectivity

If the VM isn't reachable, check that:

- The `virbr0` bridge exists on your host (`ip a`).
- The VM's network is set to use the default NAT network in virt-manager.
- No firewall rules are blocking traffic on the `192.168.122.0/24` subnet.

