---
layout: page
title: proxmox_virtual_environment_vm
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_vm

Manages a virtual machine.

## Example Usage

```hcl
resource "proxmox_virtual_environment_vm" "ubuntu_vm" {
  name        = "terraform-provider-proxmox-ubuntu-vm"
  description = "Managed by Terraform"
  tags        = ["terraform", "ubuntu"]

  node_name = "first-node"
  vm_id     = 4321

  agent {
    # read 'Qemu guest agent' section, change to true only when ready
    enabled = false
  }
  # if agent is not enabled, the VM may not be able to shutdown properly, and may need to be forced off
  stop_on_destroy = true

  startup {
    order      = "3"
    up_delay   = "60"
    down_delay = "60"
  }

  cpu {
    cores        = 2
    type         = "x86-64-v2-AES"  # recommended for modern CPUs
  }

  memory {
    dedicated = 2048
    floating  = 2048 # set equal to dedicated to enable ballooning
  }

  disk {
    datastore_id = "local-lvm"
    import_from  = proxmox_virtual_environment_download_file.latest_ubuntu_22_jammy_qcow2_img.id
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

  tpm_state {
    version = "v2.0"
  }

  serial_device {}

  virtiofs {
    mapping = "data_share"
    cache = "always"
    direct_io = true
  }
}

resource "proxmox_virtual_environment_download_file" "latest_ubuntu_22_jammy_qcow2_img" {
  content_type = "import"
  datastore_id = "local"
  node_name    = "pve"
  url = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
  # need to rename the file to *.qcow2 to indicate the actual file format for import
  file_name = "jammy-server-cloudimg-amd64.qcow2"
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
- `amd_sev` - (Optional) Secure Encrypted Virtualization (SEV) features by AMD CPUs.
    - `type` - (Optional) Enable standard SEV with `std` or enable experimental SEV-ES with the `es` option or enable experimental SEV-SNP with the `snp` option (defaults to `std`).
    - `allow_smt` - (Optional) Sets policy bit to allow Simultaneous Multi Threading (SMT)
        (Ignored unless for SEV-SNP) (defaults to `true`).
    - `kernel_hashes` - (Optional) Add kernel hashes to guest firmware for measured linux kernel launch (defaults to `false`).
    - `no_debug` - (Optional) Sets policy bit to disallow debugging of guest (defaults
        to `false`).
    - `no_key_sharing` - (Optional) Sets policy bit to disallow key sharing with other guests (Ignored for SEV-SNP) (defaults to `false`).

    The `amd_sev` setting is only allowed for a `root@pam` authenticated user.
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
- `cdrom` - (Optional) The CD-ROM configuration.
    - `enabled` - (Optional) Whether to enable the CD-ROM drive (defaults
        to `false`). *Deprecated*. The attribute will be removed in the next version of the provider.
        Set `file_id` to `none` to leave the CD-ROM drive empty.
    - `file_id` - (Optional) A file ID for an ISO file (defaults to `cdrom` as
        in the physical drive). Use `none` to leave the CD-ROM drive empty.
    - `interface` - (Optional) A hardware interface to connect CD-ROM drive to (defaults to `ide3`).
      "Must be one of `ideN`, `sataN`, `scsiN`, where N is the index of the interface. " +
      "Note that `q35` machine type only supports `ide0` and `ide2` of IDE interfaces.
- `clone` - (Optional) The cloning configuration.
    - `datastore_id` - (Optional) The identifier for the target datastore.
    - `node_name` - (Optional) The name of the source node (leave blank, if
        equal to the `node_name` argument).
    - `retries` - (Optional) Number of retries in Proxmox for clone vm.
        Sometimes Proxmox errors with timeout when creating multiple clones at
        once.
    - `vm_id` - (Required) The identifier for the source VM.
    - `full` - (Optional) Full or linked clone (defaults to `true`).
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
    - `limit` - (Optional) Limit of CPU usage, `0...128`. (defaults to `0` -- no limit).
    - `numa` - (Boolean) Enable/disable NUMA. (default to `false`)
    - `sockets` - (Optional) The number of CPU sockets (defaults to `1`).
    - `type` - (Optional) The emulated CPU type, it's recommended to
        use `x86-64-v2-AES` (defaults to `qemu64`).
        - `486` - Intel 486.
        - `Broadwell`/`Broadwell-IBRS`/`Broadwell-noTSX`/`Broadwell-noTSX-IBRS` - Intel Core Processor (Broadwell, 2014).
        - `Cascadelake-Server`/`Cascadelake-Server-noTSX`/`Cascadelake-Server-v2`/`Cascadelake-Server-v4`/`Cascadelake-Server-v5` - Intel Xeon 32xx/42xx/52xx/62xx/82xx/92xx (2019).
        - `Conroe` - Intel Celeron_4x0 (Conroe/Merom Class Core 2, 2006).
        - `Cooperlake`/`Cooperlake-v2`
        - `EPYC`/`EPYC-Genoa`/`EPYC-IBPB`/`EPYC-Milan`/`EPYC-Rome`/`EPYC-Rome-v2`/`EPYC-v3`/`EPYC-v4`/ -
            AMD EPYC Processor (2017).
        - `Haswell`/`Haswell-IBRS`/`Haswell-noTSX`/`Haswell-noTSX-IBRS` - Intel
            Core Processor (Haswell, 2013).
        - `Icelake-Client`/`Icelake-Client-noTSX`
        - `Icelake-Server`/`Icelake-Server-noTSX`/`Icelake-Server-v3`/
            `Icelake-Server-v4`/`Icelake-Server-v5`/`Icelake-Server-v6`
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
        - `SapphireRapids`
        - `Skylake-Client`/`Skylake-Client-IBRS`/`Skylake-Client-noTSX-IBRS`/`Skylake-Client-v4` -
            Intel Core Processor (Skylake, 2015).
        - `Skylake-Server`/`Skylake-Server-IBRS`/`Skylake-Server-noTSX-IBRS`/`Skylake-Server-v4`/`Skylake-Server-v5` -
            Intel Xeon Processor (Skylake, 2016).
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
        - `x86-64-v2`/`x86-64-v2-AES`/`x86-64-v3`/`x86-64-v4`
            See <https://en.wikipedia.org/wiki/X86-64#Microarchitecture_levels>
        - `custom-<model>` - Custom CPU model. All `custom-<model>` values
            should be defined in `/etc/pve/virtual-guest/cpu-models.conf` file.
    - `units` - (Optional) The CPU units (defaults to `1024`).
    - `affinity` - (Optional) The CPU cores that are used to run the VM’s vCPU. The
        value is a list of CPU IDs, separated by commas. The CPU IDs are zero-based.
        For example, `0,1,2,3` (which also can be shortened to `0-3`) means that the VM’s vCPUs are run on the first four
        CPU cores. Setting `affinity` is only allowed for `root@pam` authenticated user.
- `description` - (Optional) The description.
- `disk` - (Optional) A disk (multiple blocks supported).
    - `aio` - (Optional) The disk AIO mode (defaults to `io_uring`).
        - `io_uring` - Use io_uring.
        - `native` - Use native AIO. Should be used with to unbuffered, O_DIRECT, raw block storage only,
            with the disk `cache` must be set to `none`. Raw block storage types include iSCSI, CEPH/RBD, and NVMe.
        - `threads` - Use thread-based AIO.
    - `backup` - (Optional) Whether the drive should be included when making backups (defaults to `true`).
    - `cache` - (Optional) The cache type (defaults to `none`).
        - `none` - No cache.
        - `directsync` - Write to the host cache and wait for completion.
        - `writethrough` - Write to the host cache, but write through to
            the guest.
        - `writeback` - Write to the host cache, but write back to the
            guest when possible.
        - `unsafe` - Write directly to the disk bypassing the host cache.
    - `datastore_id` - (Optional) The identifier for the datastore to create
        the disk in (defaults to `local-lvm`).
    - `path_in_datastore` - (Optional) The in-datastore path to the disk image.
        ***Experimental.***Use to attach another VM's disks,
        or (as root only) host's filesystem paths (`datastore_id` empty string).
        See "*Example: Attached disks*".
    - `discard` - (Optional) Whether to pass discard/trim requests to the
        underlying storage. Supported values are `on`/`ignore` (defaults
        to `ignore`).
    - `file_format` - (Optional) The file format.
        - `qcow2` - QEMU Disk Image v2.
        - `raw` - Raw Disk Image.
        - `vmdk` - VMware Disk Image.
    - `file_id` - (Optional) The file ID for a disk image when importing a disk into VM. The ID format is
          `<datastore_id>:<content_type>/<file_name>`, for example `local:iso/centos8.img`. Can be also taken from
          `proxmox_virtual_environment_download_file` resource. *Deprecated*, use `import_from` instead.
    - `import_from` - (Optional) The file ID for a disk image to import into VM. The image must be of `import` content type.
       The ID format is `<datastore_id>:import/<file_name>`, for example `local:import/centos8.qcow2`. Can be also taken from
       `proxmox_virtual_environment_download_file` resource.
    - `interface` - (Required) The disk interface for Proxmox, currently `scsi`,
        `sata` and `virtio` interfaces are supported. Append the disk index at
        the end, for example, `virtio0` for the first virtio disk, `virtio1` for
        the second, etc.
    - `iothread` - (Optional) Whether to use iothreads for this disk (defaults
        to `false`).
    - `replicate` - (Optional) Whether the drive should be considered for replication jobs (defaults to `true`).
    - `serial` - (Optional) The serial number of the disk, up to 20 bytes long.
    - `size` - (Optional) The disk size in gigabytes (defaults to `8`).
    - `speed` - (Optional) The speed limits.
        - `iops_read` - (Optional) The maximum read I/O in operations per second.
        - `iops_read_burstable` - (Optional) The maximum unthrottled read I/O pool in operations per second.
        - `iops_write` - (Optional) The maximum write I/O in operations per second.
        - `iops_write_burstable` - (Optional) The maximum unthrottled write I/O pool in operations per second.
        - `read` - (Optional) The maximum read speed in megabytes per second.
        - `read_burstable` - (Optional) The maximum burstable read speed in
            megabytes per second.
        - `write` - (Optional) The maximum write speed in megabytes per second.
        - `write_burstable` - (Optional) The maximum burstable write speed in
            megabytes per second.
    - `ssd` - (Optional) Whether to use an SSD emulation option for this disk (
        defaults to `false`). Note that SSD emulation is not supported on VirtIO
        Block drives.
- `efi_disk` - (Optional) The efi disk device (required if `bios` is set
    to `ovmf`)
    - `datastore_id` (Optional) The identifier for the datastore to create
        the disk in (defaults to `local-lvm`).
    - `file_format` (Optional) The file format (defaults to `raw`).
    - `type` (Optional) Size and type of the OVMF EFI disk. `4m` is newer and
        recommended, and required for Secure Boot. For backwards compatibility
        use `2m`. Ignored for VMs with cpu.architecture=`aarch64` (defaults
        to `2m`).
    - `pre_enrolled_keys` (Optional) Use am EFI vars template with
        distribution-specific and Microsoft Standard keys enrolled, if used with
        EFI type=`4m`. Ignored for VMs with cpu.architecture=`aarch64` (defaults
        to `false`).
- `tpm_state` - (Optional) The TPM state device.
    - `datastore_id` (Optional) The identifier for the datastore to create
        the disk in (defaults to `local-lvm`).
    - `version` (Optional) TPM state device version. Can be `v1.2` or `v2.0`.
        (defaults to `v2.0`).
- `hostpci` - (Optional) A host PCI device mapping (multiple blocks supported).
    - `device` - (Required) The PCI device name for Proxmox, in form
        of `hostpciX` where `X` is a sequential number from 0 to 15.
    - `id` - (Optional) The PCI device ID. This parameter is not compatible
        with `api_token` and requires the root `username` and `password`
        configured in the proxmox provider. Use either this or `mapping`.
    - `mapping` - (Optional) The resource mapping name of the device, for
        example gpu. Use either this or `id`.
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
- `usb` - (Optional) A host USB device mapping (multiple blocks supported).
    - `host` - (Optional) The Host USB device or port or the value `spice`. Use either this or `mapping`.
    - `mapping` - (Optional) The cluster-wide resource mapping name of the device, for example "usbdevice". Use either this or `host`.
    - `usb3` - (Optional) Makes the USB device a USB3 device for the VM
        (defaults to `false`).
- `initialization` - (Optional) The cloud-init configuration.
    - `datastore_id` - (Optional) The identifier for the datastore to create the
        cloud-init disk in (defaults to `local-lvm`).
    - `interface` - (Optional) The hardware interface to connect the cloud-init
        image to. Must be one of `ide0..3`, `sata0..5`, `scsi0..30`. Will be
        detected if the setting is missing but a cloud-init image is present,
        otherwise defaults to `ide2`.
    - `dns` - (Optional) The DNS configuration.
        - `domain` - (Optional) The DNS search domain.
        - `server` - (Optional) The DNS server. The `server` attribute is
            deprecated and will be removed in a future release. Please use the
            `servers` attribute instead.
        - `servers` - (Optional) The list of DNS servers.
    - `ip_config` - (Optional) The IP configuration (one block per network
        device).
        - `ipv4` - (Optional) The IPv4 configuration.
            - `address` - (Optional) The IPv4 address in CIDR notation
                (e.g. 192.168.2.2/24). Alternatively, set this to `dhcp` for
                autodiscovery.
            - `gateway` - (Optional) The IPv4 gateway (must be omitted
                when `dhcp` is used as the address).
        - `ipv6` - (Optional) The IPv6 configuration.
            - `address` - (Optional) The IPv6 address in CIDR notation
                (e.g. fd1c:000:0000::0000:000:7334/64). Alternatively, set this
                to `dhcp` for autodiscovery.
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
- `machine` - (Optional) The VM machine type (defaults to `pc`).
    - `pc` - Standard PC (i440FX + PIIX, 1996).
    - `q35` - Standard PC (Q35 + ICH9, 2009). Optionally, you can enable VIOMMU by adding `viommu=virtio|intel` to the value, for example `q35,viommu=virtio`.
- `memory` - (Optional) The memory configuration.
    - `dedicated` - (Optional) The dedicated memory in megabytes (defaults to `512`).
    - `floating` - (Optional) The floating memory in megabytes. The default is `0`, which disables "ballooning device" for the VM.
        Please note that Proxmox has ballooning enabled by default. To enable it, set `floating` to the same value as `dedicated`.
        See [Proxmox documentation](https://pve.proxmox.com/pve-docs/pve-admin-guide.html#qm_memory) section 10.2.6 for more information.
    - `shared` - (Optional) The shared memory in megabytes (defaults to `0`).
    - `hugepages` - (Optional) Enable/disable hugepages memory (defaults to disable).
        - `2` - 2MB hugepages.
        - `1024` - 1GB hugepages.
        - `any` - Any hugepages.
    - `keep_hugepages` - (Optional) Keep hugepages memory after the VM is stopped (defaults to `false`).

    Settings `hugepages` and `keep_hugepages` are only allowed for `root@pam` authenticated user.
    And required `cpu.numa` to be enabled.
- `numa` - (Optional) The NUMA configuration.
    - `device` - (Required) The NUMA device name for Proxmox, in form
        of `numaX` where `X` is a sequential number from 0 to 7.
    - `cpus` - (Required) The CPU cores to assign to the NUMA node (format is `0-7;16-31`).
    - `memory` - (Required) The memory in megabytes to assign to the NUMA node.
    - `hostnodes` - (Optional) The NUMA host nodes.
    - `policy` - (Optional) The NUMA policy (defaults to `preferred`).
        - `interleave` - Interleave memory across nodes.
        - `preferred` - Prefer the specified node.
        - `bind` - Only use the specified node.

- `migrate` - (Optional) Migrate the VM on node change instead of re-creating
    it (defaults to `false`).
- `name` - (Optional) The virtual machine name.
- `network_device` - (Optional) A network device (multiple blocks supported).
    - `bridge` - (Optional) The name of the network bridge (defaults to `vmbr0`).
    - `disconnected` - (Optional) Whether to disconnect the network device from the network (defaults to `false`).
    - `enabled` - (Optional) Whether to enable the network device (defaults to `true`).
    - `firewall` - (Optional) Whether this interface's firewall rules should be used (defaults to `false`).
    - `mac_address` - (Optional) The MAC address.
    - `model` - (Optional) The network device model (defaults to `virtio`).
        - `e1000` - Intel E1000.
        - `e1000e` - Intel E1000E.
        - `rtl8139` - Realtek RTL8139.
        - `virtio` - VirtIO (paravirtualized).
        - `vmxnet3` - VMware vmxnet3.
    - `mtu` - (Optional) Force MTU, for VirtIO only. Set to 1 to use the bridge MTU. Cannot be larger than the bridge MTU.
    - `queues` - (Optional) The number of queues for VirtIO (1..64).
    - `rate_limit` - (Optional) The rate limit in megabytes per second.
    - `vlan_id` - (Optional) The VLAN identifier.
    - `trunks` - (Optional) String containing a `;` separated list of VLAN trunks
        ("10;20;30"). Note that the VLAN-aware feature need to be enabled on the PVE
        Linux Bridge to use trunks.
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
        - `win11` - Windows 11
        - `wvista` - Windows Vista.
        - `wxp` - Windows XP.
- `pool_id` - (Optional) The identifier for a pool to assign the virtual machine to.
- `protection` - (Optional) Sets the protection flag of the VM. This will disable the remove VM and remove disk operations (defaults to `false`).
- `reboot` - (Optional) Reboot the VM after initial creation (defaults to `false`).
- `reboot_after_update` - (Optional) Reboot the VM after update if needed (defaults to `true`).
- `rng` - (Optional) The random number generator configuration. Can only be set by `root@pam.`
    - `source` - The file on the host to gather entropy from. In most cases, `/dev/urandom` should be preferred over `/dev/random` to avoid entropy-starvation issues on the host.
    - `max_bytes` - (Optional) Maximum bytes of entropy allowed to get injected into the guest every `period` milliseconds (defaults to `1024`). Prefer a lower value when using `/dev/random` as source.
    - `period` - (Optional) Every `period` milliseconds the entropy-injection quota is reset, allowing the guest to retrieve another `max_bytes` of entropy (defaults to `1000`).
- `serial_device` - (Optional) A serial device (multiple blocks supported).
    - `device` - (Optional) The device (defaults to `socket`).
        - `/dev/*` - A host serial device.
        - `socket` - A unix socket.
- `scsi_hardware` - (Optional) The SCSI hardware type (defaults to
    `virtio-scsi-pci`).
    - `lsi` - LSI Logic SAS1068E.
    - `lsi53c810` - LSI Logic 53C810.
    - `virtio-scsi-pci` - VirtIO SCSI.
    - `virtio-scsi-single` - VirtIO SCSI (single queue).
    - `megasas` - LSI Logic MegaRAID SAS.
    - `pvscsi` - VMware Paravirtual SCSI.
- `smbios` - (Optional) The SMBIOS (type1) settings for the VM.
    - `family`- (Optional) The family string.
    - `manufacturer` - (Optional) The manufacturer.
    - `product` - (Optional) The product ID.
    - `serial` - (Optional) The serial number.
    - `sku` - (Optional) The SKU number.
    - `uuid` - (Optional) The UUID (defaults to randomly generated UUID).
    - `version` - (Optional) The version.
- `started` - (Optional) Whether to start the virtual machine (defaults
    to `true`).
- `startup` - (Optional) Defines startup and shutdown behavior of the VM.
    - `order` - (Required) A non-negative number defining the general startup
        order.
    - `up_delay` - (Optional) A non-negative number defining the delay in
        seconds before the next VM is started.
    - `down_delay` - (Optional) A non-negative number defining the delay in
        seconds before the next VM is shut down.
- `tablet_device` - (Optional) Whether to enable the USB tablet device (defaults
    to `true`).
- `tags` - (Optional) A list of tags of the VM. This is only meta information (
    defaults to `[]`). Note: Proxmox always sorts the VM tags. If the list in
    template is not sorted, then Proxmox will always report a difference on the
    resource. You may use the `ignore_changes` lifecycle meta-argument to ignore
    changes to this attribute.
- `template` - (Optional) Whether to create a template (defaults to `false`).
- `stop_on_destroy` - (Optional) Whether to stop rather than shutdown on VM destroy (defaults to `false`)
- `timeout_clone` - (Optional) Timeout for cloning a VM in seconds (defaults to
    1800).
- `timeout_create` - (Optional) Timeout for creating a VM in seconds (defaults to
    1800).
- `timeout_migrate` - (Optional) Timeout for migrating the VM (defaults to
    1800).
- `timeout_reboot` - (Optional) Timeout for rebooting a VM in seconds (defaults
    to 1800).
- `timeout_shutdown_vm` - (Optional) Timeout for shutting down a VM in seconds (
    defaults to 1800).
- `timeout_start_vm` - (Optional) Timeout for starting a VM in seconds (defaults
    to 1800).
- `timeout_stop_vm` - (Optional) Timeout for stopping a VM in seconds (defaults
    to 300).
- `vga` - (Optional) The VGA configuration.
    - `memory` - (Optional) The VGA memory in megabytes (defaults to `16`).
    - `type` - (Optional) The VGA type (defaults to `std`).
        - `cirrus` - Cirrus (deprecated since QEMU 2.2).
        - `none` - No VGA device.
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
        - `virtio-gl` - VirtIO-GPU with 3D acceleration (VirGL). VirGL support needs some extra libraries that aren’t installed by default. See the [Proxmox documentation](https://pve.proxmox.com/pve-docs/pve-admin-guide.html#qm_virtual_machines_settings) section 10.2.8 for more information.
        - `vmware` - VMware Compatible.
    - `clipboard` - (Optional) Enable VNC clipboard by setting to `vnc`. See the [Proxmox documentation](https://pve.proxmox.com/pve-docs/pve-admin-guide.html#qm_virtual_machines_settings) section 10.2.8 for more information.
- `virtiofs` - (Optional) Virtiofs share
    - `mapping` - Identifier of the directory mapping
    - `cache` - (Optional) The caching mode
        - `auto`
        - `always`
        - `metadata`
        - `never`
    - `direct_io` - (Optional) Whether to allow direct io
    - `expose_acl` - (Optional) Enable POSIX ACLs, implies xattr support
    - `expose_xattr` - (Optional) Enable support for extended attributes
- `vm_id` - (Optional) The VM identifier.
- `hook_script_file_id` - (Optional) The identifier for a file containing a hook script (needs to be executable, e.g. by using the `proxmox_virtual_environment_file.file_mode` attribute).
- `watchdog` - (Optional) The watchdog configuration. Once enabled (by a guest action), the watchdog must be periodically polled by an agent inside the guest or else the watchdog will reset the guest (or execute the respective action specified).
    - `enabled` - (Optional) Whether the watchdog is enabled (defaults to `false`).
    - `model` - (Optional) The watchdog type to emulate (defaults to `i6300esb`).
        - `i6300esb` - Intel 6300ESB.
        - `ib700` - iBase IB700.
    - `action` - (Optional) The action to perform if after activation the guest fails to poll the watchdog in time  (defaults to `none`).
        - `debug`
        - `none`
        - `pause`
        - `poweroff`
        - `reset`
        - `shutdown`

## Attribute Reference

- `ipv4_addresses` - The IPv4 addresses per network interface published by the
    QEMU agent (empty list when `agent.enabled` is `false`)
- `ipv6_addresses` - The IPv6 addresses per network interface published by the
    QEMU agent (empty list when `agent.enabled` is `false`)
- `mac_addresses` - The MAC addresses published by the QEMU agent with fallback
    to the network device configuration, if the agent is disabled
- `network_interface_names` - The network interface names published by the QEMU
    agent (empty list when `agent.enabled` is `false`)

## Qemu guest agent

Qemu-guest-agent is an application which can be installed inside guest VM, see
[Proxmox Wiki](https://pve.proxmox.com/wiki/Qemu-guest-agent) and [Proxmox
Documentation](https://pve.proxmox.com/pve-docs/pve-admin-guide.html#qm_qemu_agent)

For VM with `agent.enabled = false`, Proxmox uses ACPI for `Shutdown` and
`Reboot`, and `qemu-guest-agent` is not needed inside the VM. For some VMs,
the shutdown process may not work, causing the VM to be stuck on destroying.
Add `stop_on_destroy = true` to the VM configuration to stop the VM instead of
shutting it down.

Setting `agent.enabled = true` informs Proxmox that the guest agent is expected
to be *running* inside the VM. Proxmox then uses `qemu-guest-agent` instead of
ACPI to control the VM. If the agent is not running, Proxmox operations
`Shutdown` and `Reboot` time out and fail. The failing operation gets a lock on
the VM, and until the operation times out, other operations like `Stop` and
`Reboot` cannot be used.

Do **not** run VM with `agent.enabled = true`, unless the VM is configured to
automatically **start** `qemu-guest-agent` at some point.

"Monitor" tab in Proxmox GUI can be used to send low-level commands to `qemu`.
See the [documentation](https://www.qemu.org/docs/master/system/monitor.html).
Commands `system_powerdown` and `quit` have proven useful in shutting down VMs
with `agent.enabled = true` and no agent running.

Cloud images usually do not have `qemu-guest-agent` installed. It is possible to
install and *start* it using cloud-init, e.g. using custom `user_data_file_id`
file.

This provider requires `agent.enabled = true` to populate `ipv4_addresses`,
`ipv6_addresses` and `network_interface_names` output attributes.

Setting `agent.enabled = true` without running `qemu-guest-agent` in the VM will
also result in long timeouts when using the provider, both when creating VMs,
and when refreshing resources.  The provider has no way to distinguish between
"qemu-guest-agent not installed" and "very long boot due to a disk check", it
trusts the user to set `agent.enabled` correctly and waits for
`qemu-guest-agent` to start.

## AMD SEV

AMD SEV (-ES, -SNP) are security features for AMD processors. SEV-SNP support
is included in Proxmox version **8.4**, see [Proxmox Wiki](
https://pve.proxmox.com/wiki/Qemu/KVM_Virtual_Machines#qm_virtual_machines_settings)
and [Proxmox Documentation](https://pve.proxmox.com/pve-docs/pve-admin-guide.html#qm_memory_encryption)
for more information.

`amd-sev` requires root and therefore `root@pam` auth.

SEV-SNP requires `bios = OVMF` and a supported AMD CPU (`EPYC-v4` for instance), `machine = q35` is also advised. No EFI disk is required since SEV-SNP uses consolidated read-only firmware. A configured EFI will be ignored.

All changes made to `amd_sev` will trigger reboots. Removing or adding the `amd_sev` block will force a replacement of the resource. Modifying the `amd_sev` block will not trigger replacements.

`allow_smt` is by default set to `true` even if `snp` is not the selected type. Proxmox will ignore this value when `snp` is not in use. Likewise `no_key_sharing` is `false` by default but ignored by Proxmox when `snp` is in use.

## High Availability

When managing a virtual machine in a multi-node cluster, the VM's HA settings can
be managed using the `proxmox_virtual_environment_haresource` resource.

```hcl
resource "proxmox_virtual_environment_vm" "ubuntu_vm" {
  name  = "terraform-provider-proxmox-ubuntu-vm"
  vm_id = 4321
  # ...
}

resource "proxmox_virtual_environment_haresource" "ubuntu_vm" {
  resource_id  = "vm:${proxmox_virtual_environment_vm.ubuntu_vm.vm_id}"
  group        = "node1"
  state        = "started"
  comment      = "Managed by Terraform"
}
```

## Important Notes

### `local-lvm` Datastore

The `local-lvm` is the **default datastore** for many configuration blocks, including `initialization` and `tpm_state`, which may not seem to be related to "storage".
If you do not have `local-lvm` configured in your environment, you may need to explicitly set the `datastore_id` in such blocks to a different value.

### Cloning

When cloning an existing virtual machine, whether it's a template or not, the
resource will inherit the disks and other configuration from the source VM.

*If* you modify any attributes of an existing disk in the clone, you also need to  
explicitly provide values for any other attributes that differ from the schema defaults  
in the source (e.g., `size`, `discard`, `cache`, `aio`).  
Otherwise, the schema defaults will take effect and override the source values.

Furthermore, when cloning from one node to a different one, the behavior changes
depening on the datastores of the source VM. If at least one non-shared
datastore is used, the VM is first cloned to the source node before being
migrated to the target node. This circumvents a limitation in the Proxmox clone
API.

Because the migration step after the clone tries to preserve the used
datastores by their name, it may fail if a datastore used in the source VM is
not available on the target node (e.g. `local-lvm` is used on the source node in
the VM but no `local-lvm` datastore is available on the target node). In this
case, it is recommended to set the `datastore_id` argument in the `clone` block
to force the migration step to migrate all disks to a specific datastore on the
target node. If you need certain disks to be on specific datastores, set
the `datastore_id` argument of the disks in the `disks` block to move the disks
to the correct datastore after the cloning and migrating succeeded.

## Example: Attached disks

In this example VM `data_vm` holds two data disks, and is not used as an actual VM,
but only as a container for the disks.
It does not have any OS installation, it is never started.

VM `data_user_vm` attaches those disks as `scsi1` and `scsi2`.
**VM `data_user_vm` can be *re-created/replaced* without losing data stored on disks
owned by `data_vm`.**

~> **Experimental** Please test your configuration first in an environment where you can tolerate the potential data loss.

~> Do *not* simultaneously run more than one VM using same disk. For most filesystems,
attaching one disk to multiple VM will cause errors or even data corruption.

~> Do *not* move or resize `data_vm` disks.
(Resource `data_user_vm` should reject attempts to move or resize non-owned disks.)

```hcl
resource "proxmox_virtual_environment_vm" "data_vm" {
  node_name = "first-node"
  started = false
  on_boot = false

  disk {
    datastore_id = "local-zfs"
    interface    = "scsi0"
    size         = 1
  }

  disk {
    datastore_id = "local-zfs"
    interface    = "scsi1"
    size         = 4
  }
}

resource "proxmox_virtual_environment_vm" "data_user_vm" {
  # boot disk
  disk {
    datastore_id = "local-zfs"
    interface    = "scsi0"
    size         = 8
  }

  # attached disks from data_vm
  dynamic "disk" {
    for_each = { for idx, val in proxmox_virtual_environment_vm.data_vm.disk : idx => val }
    iterator = data_disk
    content {
      datastore_id      = data_disk.value["datastore_id"]
      path_in_datastore = data_disk.value["path_in_datastore"]
      file_format       = data_disk.value["file_format"]
      size              = data_disk.value["size"]
      # assign from scsi1 and up
      interface         = "scsi${data_disk.key + 1}"
    }
  }

  # remainder of VM configuration
  ...
}
````

## Example: Disk pass-through

You can attach another physical disk from the PVE host to a VM.
This is done by setting the `path_in_datastore` to the path of the block device on the host.

~> **Experimental** Please test your configuration first in an environment where you can tolerate the potential data loss.

~> Do *not* attach the same disk to more than one VM as it may cause data corruption.

```hcl
resource "proxmox_virtual_environment_vm" "test_vm" {
  ...

  # boot disk
  disk {
    ...
  }
  
  # attached disk
  disk {
    datastore_id = ""
    path_in_datastore  = "/dev/path/to/block/device"
    file_format = "raw"
  }

  ...
}
```

## Import

Instances can be imported using the `node_name` and the `vm_id`, e.g.,

```bash
terraform import proxmox_virtual_environment_vm.ubuntu_vm first-node/4321
```
