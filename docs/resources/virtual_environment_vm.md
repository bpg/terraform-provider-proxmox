---
layout: page
title: proxmox_virtual_environment_vm
permalink: /resources/virtual_environment_vm
nav_order: 17
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_vm

Manages a virtual machine.

## Example Usage

```terraform
resource "proxmox_virtual_environment_vm" "ubuntu_vm" {
  name        = "terraform-provider-proxmox-ubuntu-vm"
  description = "Managed by Terraform"
  tags        = ["terraform", "ubuntu"]

  node_name = "first-node"
  vm_id     = 4321

  agent {
    enabled = true
  }

  disk {
    datastore_id = "local-lvm"
    file_id      = proxmox_virtual_environment_file.ubuntu_cloud_image.id
    interface    = "scsi0"
  }

  initialization {
    ip_config {
      ipv4 {
        address = "dhcp"
      }
    }

    user_account {
      keys     = [trimspace(tls_private_key.ubuntu_vm_key.public_key_openssh)]
      password = random_password.ubuntu_vm_password.result
      username = "ubuntu"
    }

    user_data_file_id = proxmox_virtual_environment_file.cloud_config.id
  }

  network_device {
    bridge = "vmbr0"
  }

  operating_system {
    type = "l26"
  }

  serial_device {}
}

resource "proxmox_virtual_environment_file" "ubuntu_cloud_image" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = "first-node"

  source_file {
    path = "http://cloud-images.ubuntu.com/focal/current/focal-server-cloudimg-amd64.img"
  }
}

resource "random_password" "ubuntu_vm_password" {
  length           = 16
  override_special = "_%@"
  special          = true
}

resource "tls_private_key" "ubuntu_vm_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

output "ubuntu_vm_password" {
  value     = random_password.ubuntu_vm_password.result
  sensitive = true
}

output "ubuntu_vm_private_key" {
  value     = tls_private_key.ubuntu_vm_key.private_key_pem
  sensitive = true
}

output "ubuntu_vm_public_key" {
  value = tls_private_key.ubuntu_vm_key.public_key_openssh
}
```

## Argument Reference

- `acpi` - (Optional) Whether to enable ACPI (defaults to `true`).
- `agent` - (Optional) The QEMU agent configuration.
    - `enabled` - (Optional) Whether to enable the QEMU agent (defaults
      to `false`).
    - `timeout` - (Optional) The maximum amount of time to wait for data from
      the QEMU agent to become available ( defaults to `15m`).
    - `trim` - (Optional) Whether to enable the FSTRIM feature in the QEMU agent
      (defaults to `false`).
    - `type` - (Optional) The QEMU agent interface type (defaults to `virtio`).
        - `isa` - ISA Serial Port.
        - `virtio` - VirtIO (paravirtualized).
- `audio_device` - (Optional) An audio device.
    - `device` - (Optional) The device (defaults to `intel-hda`).
        - `AC97` - Intel 82801AA AC97 Audio.
        - `ich9-intel-hda` - Intel HD Audio Controller (ich9).
        - `intel-hda` - Intel HD Audio.
    - `driver` - (Optional) The driver (defaults to `spice`).
        - `spice` - Spice.
    - `enabled` - (Optional) Whether to enable the audio device (defaults
      to `true`).
- `bios` - (Optional) The BIOS implementation (defaults to `seabios`).
    - `ovmf` - OVMF (UEFI).
    - `seabios` - SeaBIOS.
- `boot_order` - (Optional) Specify a list of devices to boot from in the order
  they appear in the list (defaults to `[]`).
- `cdrom` - (Optional) The CDROM configuration.
    - `enabled` - (Optional) Whether to enable the CDROM drive (defaults
      to `false`).
    - `file_id` - (Optional) A file ID for an ISO file (defaults to `cdrom` as
      in the physical drive).
- `clone` - (Optional) The cloning configuration.
    - `datastore_id` - (Optional) The identifier for the target datastore.
    - `node_name` - (Optional) The name of the source node (leave blank, if
      equal to the `node_name` argument).
    - `retries` - (Optional) Number of retries in Proxmox for clone vm.
      Sometimes Proxmox errors with timeout when creating multiple clones at
      once.
    - `vm_id` - (Required) The identifier for the source VM.
- `cpu` - (Optional) The CPU configuration.
    - `architecture` - (Optional) The CPU architecture (defaults to `x86_64`).
        - `aarch64` - ARM (64 bit).
        - `x86_64` - x86 (64-bit).
    - `cores` - (Optional) The number of CPU cores (defaults to `1`).
    - `flags` - (Optional) The CPU flags.
        - `+aes`/`-aes` - Activate AES instruction set for HW acceleration.
        - `+amd-no-ssb`/`-amd-no-ssb` - Notifies guest OS that host is not
          vulnerable for Spectre on AMD CPUs.
        - `+amd-ssbd`/`-amd-ssbd` - Improves Spectre mitigation performance with
          AMD CPUs, best used with "virt-ssbd".
        - `+hv-evmcs`/`-hv-evmcs` - Improve performance for nested
          virtualization (only supported on Intel CPUs).
        - `+hv-tlbflush`/`-hv-tlbflush` - Improve performance in overcommitted  
          Windows guests (may lead to guest BSOD on old CPUs).
        - `+ibpb`/`-ibpb` - Allows improved Spectre mitigation on AMD CPUs.
        - `+md-clear`/`-md-clear` - Required to let the guest OS know if MDS is
          mitigated correctly.
        - `+pcid`/`-pcid` - Meltdown fix cost reduction on Westmere, Sandy- and
          Ivy Bridge Intel CPUs.
        - `+pdpe1gb`/`-pdpe1gb` - Allows guest OS to use 1 GB size pages, if
          host HW supports it.
        - `+spec-ctrl`/`-spec-ctrl` - Allows improved Spectre mitigation with
          Intel CPUs.
        - `+ssbd`/`-ssbd` - Protection for "Speculative Store Bypass" for Intel
          models.
        - `+virt-ssbd`/`-virt-ssbd` - Basis for "Speculative Store Bypass"
          protection for AMD models.
    - `hotplugged` - (Optional) The number of hotplugged vCPUs (defaults
      to `0`).
    - `numa` - (Boolean) Enable/disable NUMA. (default to `false`)
    - `sockets` - (Optional) The number of CPU sockets (defaults to `1`).
    - `type` - (Optional) The emulated CPU type (defaults to `qemu64`).
        - `486` - Intel 486.
        - `Broadwell`/`Broadwell-IBRS`/`Broadwell-noTSX`/`Broadwell-noTSX-IBRS`
        - Intel Core Processor (Broadwell, 2014).
        - `Cascadelake-Server` - Intel Xeon 32xx/42xx/52xx/62xx/82xx/92xx (
          2019).
        - `Conroe` - Intel Celeron_4x0 (Conroe/Merom Class Core 2, 2006).
        - `EPYC`/`EPYC-IBPB` - AMD EPYC Processor (2017).
        - `Haswell`/`Haswell-IBRS`/`Haswell-noTSX`/`Haswell-noTSX-IBRS` - Intel
          Core Processor (Haswell, 2013).
        - `IvyBridge`/`IvyBridge-IBRS` - Intel Xeon E3-12xx v2 (Ivy Bridge,
          2012).
        - `KnightsMill` - Intel Xeon Phi 72xx (2017).
        - `Nehalem`/`Nehalem-IBRS` - Intel Core i7 9xx (Nehalem Class Core i7,
          2008).
        - `Opteron_G1` - AMD Opteron 240 (Gen 1 Class Opteron, 2004).
        - `Opteron_G2` - AMD Opteron 22xx (Gen 2 Class Opteron, 2006).
        - `Opteron_G3` - AMD Opteron 23xx (Gen 3 Class Opteron, 2009).
        - `Opteron_G4` - AMD Opteron 62xx class CPU (2011).
        - `Opteron_G5` - AMD Opteron 63xx class CPU (2012).
        - `Penryn` - Intel Core 2 Duo P9xxx (Penryn Class Core 2, 2007).
        - `SandyBridge`/`SandyBridge-IBRS` - Intel Xeon E312xx (Sandy Bridge,
          2011).
        - `Skylake-Client`/`Skylake-Client-IBRS` - Intel Core Processor (
          Skylake, 2015).
        - `Skylake-Server`/`Skylake-Server-IBRS` - Intel Xeon Processor (
          Skylake, 2016).
        - `Westmere`/`Westmere-IBRS` - Intel Westmere E56xx/L56xx/X56xx (
          Nehalem-C, 2010).
        - `athlon` - AMD Athlon.
        - `core2duo` - Intel Core 2 Duo.
        - `coreduo` - Intel Core Duo.
        - `host` - Host pass-through.
        - `kvm32`/`kvm64` - Common KVM processor (32 & 64 bit variants).
        - `max` - Maximum amount of features from host CPU.
        - `pentium` - Intel Pentium (1993).
        - `pentium2` - Intel Pentium 2 (1997-1999).
        - `pentium3` - Intel Pentium 3 (1999-2001).
        - `phenom` - AMD Phenom (2010).
        - `qemu32`/`qemu64` - QEMU Virtual CPU version 2.5+ (32 & 64 bit
          variants).
        - `custom-<model>` - Custom CPU model. All `custom-<model>` values
          should be defined in `/etc/pve/virtual-guest/cpu-models.conf` file.
    - `units` - (Optional) The CPU units (defaults to `1024`).
- `description` - (Optional) The description.
- `disk` - (Optional) A disk (multiple blocks supported).
    - `datastore_id` - (Optional) The identifier for the datastore to create
      the disk in (defaults to `local-lvm`).
    - `discard` - (Optional) Whether to pass discard/trim requests to the
      underlying storage. Supported values are `on`/`ignore` (defaults
      to `ignore`).
    - `file_format` - (Optional) The file format (defaults to `qcow2`).
        - `qcow2` - QEMU Disk Image v2.
        - `raw` - Raw Disk Image.
        - `vmdk` - VMware Disk Image.
    - `file_id` - (Optional) The file ID for a disk image (experimental -
      might cause high CPU utilization during import, especially with large
      disk images).
    - `interface` - (Required) The disk interface for Proxmox, currently `scsi`,
      `sata` and `virtio` interfaces are supported. Append the disk index at
      the end, for example, `virtio0` for the first virtio disk, `virtio1` for
      the second, etc.
    - `iothread` - (Optional) Whether to use iothreads for this disk (defaults
      to `false`).
    - `size` - (Optional) The disk size in gigabytes (defaults to `8`).
    - `speed` - (Optional) The speed limits.
        - `read` - (Optional) The maximum read speed in megabytes per second.
        - `read_burstable` - (Optional) The maximum burstable read speed in
          megabytes per second.
        - `write` - (Optional) The maximum write speed in megabytes per second.
        - `write_burstable` - (Optional) The maximum burstable write speed in
          megabytes per second.
    - `ssd` - (Optional) Whether to use an SSD emulation option for this disk (
      defaults to `false`). Note that SSD emulation is not supported on VirtIO
      Block drives.
- `efi_disk` - (Optional) The efi disk device (required if `bios` is set to `ovmf`)
    - `datastore_id` (String) The datastore id
    - `file_format` (String) The file format
    - `size` (String) The disk size in megabytes
- `hostpci` - (Optional) A host PCI device mapping (multiple blocks supported).
    - `device` - (Required) The PCI device name for Proxmox, in form
      of `hostpciX` where `X` is a sequential number from 0 to 3.
    - `id` - (Required) The PCI device ID.
    - `mdev` - (Optional) The mediated device ID to use.
    - `pcie` - (Optional) Tells Proxmox to use a PCIe or PCI port. Some
      guests/device combination require PCIe rather than PCI. PCIe is only
      available for q35 machine types.
    - `rombar` - (Optional) Makes the firmware ROM visible for the VM (defaults
      to `true`).
    - `rom_file` - (Optional) A path to a ROM file for the device to use. This
      is a relative path under `/usr/share/kvm/`.
    - `xvga` - (Optional) Marks the PCI(e) device as the primary GPU of the VM.
      With this enabled the `vga` configuration argument will be ignored.
- `initialization` - (Optional) The cloud-init configuration.
    - `datastore_id` - (Optional) The identifier for the datastore to create the
      cloud-init disk in (defaults to `local-lvm`).
    - `dns` - (Optional) The DNS configuration.
        - `domain` - (Optional) The DNS search domain.
        - `server` - (Optional) The DNS server.
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
    - `user_account` - (Optional) The user account configuration (conflicts
      with `user_data_file_id`).
        - `keys` - (Optional) The SSH keys.
        - `password` - (Optional) The SSH password.
        - `username` - (Optional) The SSH username.
    - `network_data_file_id` - (Optional) The identifier for a file containing
      network configuration data passed to the VM via cloud-init (conflicts
      with `ip_config`).
    - `user_data_file_id` - (Optional) The identifier for a file containing
      custom user data (conflicts with `user_account`).
    - `vendor_data_file_id` - (Optional) The identifier for a file containing
      all vendor data passed to the VM via cloud-init. 
    - `meta_data_file_id` - (Optional) The identifier for a file containing
      all meta data passed to the VM via cloud-init.
- `keyboard_layout` - (Optional) The keyboard layout (defaults to `en-us`).
    - `da` - Danish.
    - `de` - German.
    - `de-ch` - Swiss German.
    - `en-gb` - British English.
    - `en-us` - American English.
    - `es` - Spanish.
    - `fi` - Finnish.
    - `fr` - French.
    - `fr-be` - Belgian French.
    - `fr-ca` - French Canadian.
    - `fr-ch` - Swish French.
    - `hu` - Hungarian.
    - `is` - Icelandic.
    - `it` - Italian.
    - `ja` - Japanese.
    - `lt` - Lithuanian.
    - `mk` - Macedonian.
    - `nl` - Dutch.
    - `no` - Norwegian.
    - `pl` - Polish.
    - `pt` - Portuguese.
    - `pt-br` - Brazilian Portuguese.
    - `sl` - Slovenian.
    - `sv` - Swedish.
    - `tr` - Turkish.
- `kvm_arguments` - (Optional) Arbitrary arguments passed to kvm.
- `machine` - (Optional) The VM machine type (defaults to `i440fx`).
    - `i440fx` - Standard PC (i440FX + PIIX, 1996).
    - `q35` - Standard PC (Q35 + ICH9, 2009).
- `memory` - (Optional) The memory configuration.
    - `dedicated` - (Optional) The dedicated memory in megabytes (defaults
      to `512`).
    - `floating` - (Optional) The floating memory in megabytes (defaults
      to `0`).
    - `shared` - (Optional) The shared memory in megabytes (defaults to `0`).
- `name` - (Optional) The virtual machine name.
- `network_device` - (Optional) A network device (multiple blocks supported).
    - `bridge` - (Optional) The name of the network bridge (defaults
      to `vmbr0`).
    - `enabled` - (Optional) Whether to enable the network device (defaults
      to `true`).
    - `firewall` - (Optional) Whether this interface's firewall rules should be
      used (defaults to `false`).
    - `mac_address` - (Optional) The MAC address.
    - `model` - (Optional) The network device model (defaults to `virtio`).
        - `e1000` - Intel E1000.
        - `rtl8139` - Realtek RTL8139.
        - `virtio` - VirtIO (paravirtualized).
        - `vmxnet3` - VMware vmxnet3.
    - `mtu` - (Optional) Force MTU, for VirtIO only. Set to 1 to use the bridge
      MTU. Cannot be larger than the bridge MTU.
    - `rate_limit` - (Optional) The rate limit in megabytes per second.
    - `vlan_id` - (Optional) The VLAN identifier.
- `node_name` - (Required) The name of the node to assign the virtual machine
  to.
- `on_boot` - (Optional) Specifies whether a VM will be started during system
  boot. (defaults to `true`)
- `operating_system` - (Optional) The Operating System configuration.
    - `type` - (Optional) The type (defaults to `other`).
        - `l24` - Linux Kernel 2.4.
        - `l26` - Linux Kernel 2.6 - 5.X.
        - `other` - Unspecified OS.
        - `solaris` - OpenIndiania, OpenSolaris og Solaris Kernel.
        - `w2k` - Windows 2000.
        - `w2k3` - Windows 2003.
        - `w2k8` - Windows 2008.
        - `win7` - Windows 7.
        - `win8` - Windows 8, 2012 or 2012 R2.
        - `win10` - Windows 10 or 2016.
        - `wvista` - Windows Vista.
        - `wxp` - Windows XP.
- `pool_id` - (Optional) The identifier for a pool to assign the virtual machine
  to.
- `reboot` - (Optional) Reboot the VM after initial creation. (defaults
  to `false`)
- `serial_device` - (Optional) A serial device (multiple blocks supported).
    - `device` - (Optional) The device (defaults to `socket`).
        - `/dev/*` - A host serial device.
        - `socket` - A unix socket.
- `scsi_hardware` - (Optional) The SCSI hardware type (defaults
  to `virtio-scsi-pci`).
    - `lsi` - LSI Logic SAS1068E.
    - `lsi53c810` - LSI Logic 53C810.
    - `virtio-scsi-pci` - VirtIO SCSI.
    - `virtio-scsi-single` - VirtIO SCSI (single queue).
    - `megasas` - LSI Logic MegaRAID SAS.
    - `pvscsi` - VMware Paravirtual SCSI.
- `started` - (Optional) Whether to start the virtual machine (defaults
  to `true`).
- `tablet_device` - (Optional) Whether to enable the USB tablet device (defaults
  to `true`).
- `tags` - (Optional) A list of tags of the VM. This is only meta information (
  defaults to `[]`). Note: Proxmox always sorts the VM tags. If the list in
  template is not sorted, then Proxmox will always report a difference on the
  resource. You may use the `ignore_changes` lifecycle meta-argument to ignore
  changes to this attribute.
- `template` - (Optional) Whether to create a template (defaults to `false`).
- `timeout_clone` - (Optional) Timeout for cloning a VM in seconds (defaults to
  1800).
- `timeout_move_disk` - (Optional) Timeout for moving the disk of a VM in
  seconds (defaults to 1800).
- `timeout_reboot` - (Optional) Timeout for rebooting a VM in seconds (defaults
  to 1800).
- `timeout_shutdown_vm` - (Optional) Timeout for shutting down a VM in seconds (
  defaults to 1800).
- `timeout_start_vm` - (Optional) Timeout for starting a VM in seconds (defaults
  to 1800).
- `timeout_stop_vm` - (Optional) Timeout for stopping a VM in seconds (defaults
  to 300).
- `vga` - (Optional) The VGA configuration.
    - `enabled` - (Optional) Whether to enable the VGA device (defaults
      to `true`).
    - `memory` - (Optional) The VGA memory in megabytes (defaults to `16`).
    - `type` - (Optional) The VGA type (defaults to `std`).
        - `cirrus` - Cirrus (deprecated since QEMU 2.2).
        - `qxl` - SPICE.
        - `qxl2` - SPICE Dual Monitor.
        - `qxl3` - SPICE Triple Monitor.
        - `qxl4` - SPICE Quad Monitor.
        - `serial0` - Serial Terminal 0.
        - `serial1` - Serial Terminal 1.
        - `serial2` - Serial Terminal 2.
        - `serial3` - Serial Terminal 3.
        - `std` - Standard VGA.
        - `virtio` - VirtIO-GPU.
        - `vmware` - VMware Compatible.
- `vm_id` - (Optional) The VM identifier.

## Attribute Reference

- `ipv4_addresses` - The IPv4 addresses per network interface published by the
  QEMU agent (empty list when `agent.enabled` is `false`)
- `ipv6_addresses` - The IPv6 addresses per network interface published by the
  QEMU agent (empty list when `agent.enabled` is `false`)
- `mac_addresses` - The MAC addresses published by the QEMU agent with fallback
  to the network device configuration, if the agent is disabled
- `network_interface_names` - The network interface names published by the QEMU
  agent (empty list when `agent.enabled` is `false`)

## Important Notes

When cloning an existing virtual machine, whether it's a template or not, the
resource will only detect changes to the arguments which are not set to their
default values.

Furthermore, when cloning from one node to a different one, the behavior changes
depening on the datastores of the source VM. If at least one non-shared
datastore is used, the VM is first cloned to the source node before being
migrated to the target node. This circumvents a limitation in the Proxmox clone
API.

**Note:** Because the migration step after the clone tries to preserve the used
datastores by their name, it may fail if a datastore used in the source VM is
not available on the target node (e.g. `local-lvm` is used on the source node in
the VM but no `local-lvm` datastore is available on the target node). In this
case, it is recommended to set the `datastore_id` argument in the `clone` block
to force the migration step to migrate all disks to a specific datastore on the
target node. If you need certain disks to be on specific datastores, set
the `datastore_id` argument of the disks in the `disks` block to move the disks
to the correct datastore after the cloning and migrating succeeded.
