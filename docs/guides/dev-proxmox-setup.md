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

### Option1 -  Using virt-manager
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

### Option 2 - Using virsh and an unattended installation

   The Proxmox VE ISO can be installed unattended. The installation is fully automatic, which lets you spin up a Proxmox VM effortlessly and saves you from clicking through the installer each time.

   The installer is driven by an `answer.toml` file. See the [Proxmox automated installation documentation](https://pve.proxmox.com/wiki/Automated_Installation#Assistant_Tool) for details.


   **Steps:**

   1. Install `virsh` (the CLI to manage VMs), virtinst and `proxmox-auto-install-assistant`:

      ```shell
      sudo apt update
      sudo apt install libvirt-clients libvirt-daemon-system virtinst
      # proxmox-auto-install-assistant is not in the standard Debian/Ubuntu repositories. Check the proxmox documentation for you distribution

      sudo apt install proxmox-auto-install-assistant
      ```

   2. Download the latest Proxmox VE ISO from https://www.proxmox.com/en/downloads.
   3. Create a different network than the default one. As you will need a NATed network without DHCP.
      ```text
      # proxmox-test-network.xml
      <network>
         <name>proxmox-test-network</name>
         <bridge name="prmx-test" stp="on" delay="0"/>
         <forward mode='nat'>
            <nat/>
         </forward>
         <ip address="192.168.126.1" netmask="255.255.255.0">
         </ip>
      </network>
      ```
      Create and activate the network
      ```shell
      virsh net-define proxmox-test-network.xml
      virsh start proxmox-test-network
      ```
   3. Create the ISO with the `answer.toml` file:

      ```shell
      proxmox-auto-install-assistant prepare-iso /path/to/proxmox-ve_x.y-z.iso --fetch-from iso --answer-file answer-pve.toml --output /path/where/to/store/the/generated/proxmox-with-answer-file.iso
      ```


      Example of an `answer-pve.toml` file:

      ```toml
      [global]
      keyboard = "en-us"
      country = "fr"
      fqdn = "pve.example.local"
      mailto = "admin@example.local"
      timezone = "Europe/Paris"
      root_password = "APassYouRemember"
      root_ssh_keys = ["ssh-ed25519 AAAAC3NzaC1lZDI1........ user@pve"]

      [network]
      source = "from-answer"
      cidr = "192.168.126.43/24"
      dns = "192.168.126.1"
      gateway = "192.168.126.1"
      filter.ID_NET_NAME_PATH = "enp1s0"

      [disk-setup]
      filesystem = "ext4"
      disk_list = ["vda"]
      ```

   4. Run the `virt-install` command to create the Proxmox VM:

      ```shell
      virt-install \
      --name pve \
      --memory 8192 \
      --vcpus 4 \
      --cpu host-passthrough \
      --os-variant debian13 \
      --disk path=/path/where/to/store/the/vm/disk/pve.cow2,size=30,bus=virtio,format=qcow2 \
      --cdrom /path/where/stored/the/generated/proxmox-with-answer-file.iso \
      # You may choose another network if you have one defined
      --network network=proxmox-test-network,model=virtio \ 
      --graphics vnc,listen=127.0.0.1 \
      --noautoconsole

      # Check the status of the VM with
      watch -n 3 virsh list
      ```
       **The VM shuts down once the installation is complete, so you need to start it again with `virsh`.**
   

   5. Start the VM:

   ```shell
   virsh start pve
   ```

   6. (Optional) Verify connectivity from your host:

   ```shell
   ssh root@192.168.126.43
   ```

   You can also access the web console at `https://192.168.126.43:8006`.

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

### Optional: spare disk for ZFS pool management tests

Some acceptance tests for `proxmox_node_disk_zfs` require a spare block device that Proxmox
can fully wipe and use as a ZFS pool. On a VM-based dev setup you can add a small virtual
disk in virt-manager:

1. In virt-manager, shut down the Proxmox VM, add a new virtual disk (e.g., 8 GB, VirtIO),
   then start it again.

2. Identify the new device name inside Proxmox (it will appear as `/dev/vdb` or similar):

   ```sh
   lsblk --output NAME,SIZE,TYPE | grep disk
   ```

3. Set the environment variable in your `testacc.env`:

   ```env
   PROXMOX_VE_ACC_ZFS_DISK="/dev/vdb"
   ```

Tests that require a spare disk will be skipped if this variable is not set.

> [!WARNING]
> The entire device will be wiped during testing. Do not point this at a disk that contains data.

## SSH access

The default provider configuration uses API token authentication. Since there is no password to inherit for SSH, you need one of the following:

- **ssh-agent (recommended):** Load a key that is authorized on the Proxmox host into your agent (`ssh-add`). The provider's `ssh { agent = true }` setting will use it automatically.
- **Explicit password:** Add `password = var.virtual_environment_root_password` inside the `ssh { }` block in `example/main.tf`.
- **Private key file:** Set `private_key` inside the `ssh { }` block.

For more details, see the [SSH Agent](https://bpg.sh/docs#ssh-agent) section in the provider documentation.

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

You need to enable snippets on the datastore. See the [file resource documentation](https://bpg.sh/docs/resources/virtual_environment_file#snippets) for instructions.

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
