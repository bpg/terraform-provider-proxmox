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

4. Complete the Proxmox installation. Choose your preferred timezone, country, password, domain, and email. **Important:** keep the default hostname `pve` — the example configurations expect a node named `pve`. If you use a different hostname, set `virtual_environment_node_name` in your `terraform.tfvars`.

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

## Proxmox configuration

Before running examples, ensure the following on your Proxmox node:

1. Enable "Snippets" and "Import" content types in the `local` storage (Datacenter -> Storage -> local -> Edit)
2. Make the default Linux Bridge "vmbr0" VLAN aware (Datacenter -> pve -> Network -> vmbr0 -> Edit)
3. Create the bind mount directory: `mkdir -p /mnt/bindmounts/shared`
4. Create an API token (Datacenter -> Permissions -> API Tokens)

### Optional: ZFS storage for disk tests

Some acceptance tests (e.g., `TestAccResourceVMDisks/clone_with_moving_disk`) require a ZFS-backed storage pool. To set one up:

1. SSH into your Proxmox node and create a ZFS pool. If you don't have a spare disk, you can use a file-backed pool for testing:

   ```sh
   dd if=/dev/zero of=/tank.img bs=1M count=4096
   zpool create tank /tank.img
   ```

2. Add the ZFS pool as storage in Proxmox (Datacenter -> Storage -> Add -> ZFS, set ID to `tank`, select the `tank` pool).

3. Set the environment variable in your `testacc.env`:

   ```env
   PROXMOX_VE_ACC_ZFS_DATASTORE_ID="tank"
   ```

Tests that require ZFS storage will be skipped if this variable is not set.

## SSH access

The default provider configuration uses API token authentication. Since there is no password to inherit for SSH, you need one of the following:

- **ssh-agent (recommended):** Load a key that is authorized on the Proxmox host into your agent (`ssh-add`). The provider's `ssh { agent = true }` setting will use it automatically.
- **Explicit password:** Add `password = var.virtual_environment_root_password` inside the `ssh { }` block in `example/main.tf`.
- **Private key file:** Set `private_key` inside the `ssh { }` block.

For more details, see the [SSH Agent](https://registry.terraform.io/providers/bpg/proxmox/latest/docs#ssh-agent) section in the provider documentation.

## Configuring the provider

Create a `terraform.tfvars` file in the `example/` directory (this file is git-ignored):

```hcl
virtual_environment_endpoint      = "https://192.168.122.43:8006/"
virtual_environment_ssh_username  = "root"
virtual_environment_api_token     = "terraform@pve!provider=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
virtual_environment_root_password = "your-root-password"
```

Replace the IP address, token, and password with your actual values.

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

### Timezone not set on fresh Proxmox installs

On some fresh Proxmox installations (commonly with European timezones), the timezone may not be properly configured, causing the `proxmox_virtual_environment_time` resource to fail with:

```text
Error: failed to update node time: received an HTTP 400 response -
  Reason: Parameter verification failed. (timezone: No such timezone)
```

**Workaround:** Set the timezone manually in the Proxmox web UI (Node -> Time -> Edit) or via the command line:

```sh
timedatectl set-timezone Europe/Berlin  # replace with your timezone
```

Then re-run `make example`. See [#2460](https://github.com/bpg/terraform-provider-proxmox/issues/2460) for details.

### Trunks example requires DHCP

The trunks example VM uses `address = "dhcp"` and waits for the QEMU guest agent to report an IP address. If your network does not have a DHCP v4 server, the VM will start but the provider will time out waiting for an IP. Either set up DHCP on the network connected to `vmbr0`, or move `resource_virtual_environment_trunks.tf` to the `example/exclude/` directory to skip it.

### Network connectivity

If the VM isn't reachable, check that:

- The `virbr0` bridge exists on your host (`ip a`).
- The VM's network is set to use the default NAT network in virt-manager.
- No firewall rules are blocking traffic on the `192.168.122.0/24` subnet.
