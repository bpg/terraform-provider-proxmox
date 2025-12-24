---
layout: page
title: proxmox_virtual_environment_cloned_vm
parent: Resources
subcategory: Virtual Environment
description: |-
  Clone a VM from a source template/VM and manage only explicitly-defined configuration. This resource uses explicit opt-in management: only configuration blocks and devices explicitly listed in your Terraform code are managed. Inherited settings from the template are preserved unless explicitly overridden or deleted. Removing a configuration from Terraform stops managing it but does not delete it from the VM.
---

# Resource: proxmox_virtual_environment_cloned_vm

~> **EXPERIMENTAL**

Clone a VM from a source template/VM and manage only explicitly-defined configuration. This resource uses explicit opt-in management: only configuration blocks and devices explicitly listed in your Terraform code are managed. Inherited settings from the template are preserved unless explicitly overridden or deleted. Removing a configuration from Terraform stops managing it but does not delete it from the VM.

## Limitations

This resource intentionally manages only a subset of VM configuration. The following are currently not managed and must be inherited from the source template (or managed via `proxmox_virtual_environment_vm` with a `clone` block):

- BIOS / machine / boot order
- EFI disk / secure boot settings
- TPM state
- Cloud-init / initialization
- QEMU guest agent configuration
- PCI/USB passthrough, serial/audio devices, watchdog, VirtioFS

## Example Usage

```terraform
# Example 1: Basic clone with minimal management
resource "proxmox_virtual_environment_cloned_vm" "basic_clone" {
  node_name = "pve"
  name      = "basic-clone"

  clone = {
    source_vm_id = 100  # Template VM ID
    full         = true # Perform full clone (not linked)
  }

  # Only manage CPU, inherit everything else from template
  cpu = {
    cores = 4
  }
}

# Example 2: Clone with explicit network management
resource "proxmox_virtual_environment_cloned_vm" "network_managed" {
  node_name = "pve"
  name      = "network-clone"

  clone = {
    source_vm_id = 100
  }

  # Map-based network devices - manage specific interfaces
  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
      tag    = 100 # VLAN tag
    }

    net1 = {
      bridge      = "vmbr1"
      model       = "virtio"
      firewall    = true
      mac_address = "BC:24:11:2E:C5:00"
    }
  }

  cpu = {
    cores = 2
  }
}

# Example 3: Clone with disk management
resource "proxmox_virtual_environment_cloned_vm" "disk_managed" {
  node_name = "pve"
  name      = "disk-clone"

  clone = {
    source_vm_id     = 100
    target_datastore = "local-lvm"
  }

  # Map-based disk management
  disk = {
    scsi0 = {
      # Resize the cloned boot disk
      datastore_id = "local-lvm"
      size_gb      = 50
      discard      = "on"
      ssd          = true
    }

    scsi1 = {
      # Add a new data disk
      datastore_id = "local-lvm"
      size_gb      = 100
      backup       = false
    }
  }
}

# Example 4: Clone with explicit device deletion
resource "proxmox_virtual_environment_cloned_vm" "selective_delete" {
  node_name = "pve"
  name      = "minimal-clone"

  clone = {
    source_vm_id = 100
  }

  # Only manage net0
  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
    }
  }

  # Explicitly delete inherited net1 and net2
  delete = {
    network = ["net1", "net2"]
  }
}

# Example 5: Full-featured clone with multiple settings
resource "proxmox_virtual_environment_cloned_vm" "full_featured" {
  node_name   = "pve"
  name        = "production-vm"
  description = "Production VM cloned from template"
  tags        = ["production", "web"]

  clone = {
    source_vm_id     = 100
    source_node_name = "pve" # Source node (defaults to target node if omitted)
    full             = true
    target_datastore = "local-lvm"
    retries          = 3
  }

  cpu = {
    cores        = 8
    sockets      = 1
    architecture = "x86_64"
    type         = "host"
  }

  memory = {
    size    = 8192
    balloon = 2048
    shares  = 2000
  }

  network = {
    net0 = {
      bridge     = "vmbr0"
      model      = "virtio"
      tag        = 100
      firewall   = true
      rate_limit = 100.0 # MB/s
    }
  }

  disk = {
    scsi0 = {
      datastore_id = "local-lvm"
      size_gb      = 100
      discard      = "on"
      iothread     = true
      ssd          = true
      cache        = "writethrough"
    }
  }

  vga = {
    type   = "std"
    memory = 16
  }

  # Option 1: Remove inherited CD-ROM using delete block
  delete = {
    disk = ["ide2"] # Remove inherited CD-ROM device
  }

  # Option 2: Alternatively, manage the CD-ROM and ensure it's empty
  # cdrom = {
  #   ide2 = {
  #     file_id = "none" # Ensure CD-ROM is empty
  #   }
  # }

  # Lifecycle options
  stop_on_destroy                      = false # Shutdown gracefully instead of force stop
  purge_on_destroy                     = true
  delete_unreferenced_disks_on_destroy = false # Safety: don't delete unmanaged disks

  timeouts = {
    create = "30m"
    update = "30m"
    delete = "10m"
  }
}

# Example 6: Linked clone for testing
resource "proxmox_virtual_environment_cloned_vm" "test_clone" {
  node_name = "pve"
  name      = "test-vm"

  clone = {
    source_vm_id = 100
    full         = false # Linked clone - faster, uses less space
  }

  cpu = {
    cores = 2
  }

  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
    }
  }
}

# Example 7: Clone with pool assignment
resource "proxmox_virtual_environment_cloned_vm" "pooled_clone" {
  node_name = "pve"
  name      = "pooled-vm"

  clone = {
    source_vm_id = 100
    pool_id      = "production" # Assign to pool during clone
  }

  cpu = {
    cores = 4
  }
}

# Example 8: Import existing cloned VM
resource "proxmox_virtual_environment_cloned_vm" "imported" {
  # Import with: terraform import proxmox_virtual_environment_cloned_vm.imported pve/123

  id        = 123 # VM ID to manage
  node_name = "pve"

  # After import, define managed configuration
  clone = {
    source_vm_id = 100 # Must match original clone source
  }

  cpu = {
    cores = 4
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `clone` (Attributes) Clone settings. Changes require recreation. (see [below for nested schema](#nestedatt--clone))
- `node_name` (String) Target node for the cloned VM.

### Optional

- `cdrom` (Attributes Map) The CD-ROM configuration. The key is the interface of the CD-ROM, could be one of `ideN`, `sataN`, `scsiN`, where N is the index of the interface. Note that `q35` machine type only supports `ide0` and `ide2` of IDE interfaces. (see [below for nested schema](#nestedatt--cdrom))
- `cpu` (Attributes) The CPU configuration. (see [below for nested schema](#nestedatt--cpu))
- `delete` (Attributes) Explicit deletions to perform after cloning/updating. Entries persist across applies. (see [below for nested schema](#nestedatt--delete))
- `delete_unreferenced_disks_on_destroy` (Boolean) Delete unreferenced disks on destroy. WARNING: When set to true, any disks not explicitly managed by Terraform will be deleted on destroy, potentially causing data loss. Defaults to false for safety.
- `description` (String) Optional VM description applied after cloning.
- `disk` (Attributes Map) Disks keyed by slot (scsi0, virtio0, sata0, ide0, ...). Only listed keys are managed. (see [below for nested schema](#nestedatt--disk))
- `id` (Number) The VM identifier in the Proxmox cluster.
- `memory` (Attributes) Memory configuration for the VM. Uses Proxmox memory ballooning to allow dynamic memory allocation. The `size` sets the total available RAM, while `balloon` sets the guaranteed floor. The host can reclaim memory between these values when needed. (see [below for nested schema](#nestedatt--memory))
- `name` (String) Optional VM name override applied after cloning.
- `network` (Attributes Map) Network devices keyed by slot (net0, net1, ...). Only listed keys are managed. (see [below for nested schema](#nestedatt--network))
- `purge_on_destroy` (Boolean) Purge backup configuration on destroy.
- `rng` (Attributes) Configure the RNG (Random Number Generator) device. The RNG device provides entropy to guests to ensure good quality random numbers for guest applications that require them. Can only be set by `root@pam.` See the [Proxmox documentation](https://pve.proxmox.com/pve-docs/pve-admin-guide.html#qm_virtual_machines_settings) for more information. (see [below for nested schema](#nestedatt--rng))
- `stop_on_destroy` (Boolean) Stop the VM on destroy (instead of shutdown).
- `tags` (Set of String) Tags applied after cloning.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `vga` (Attributes) Configure the VGA Hardware. If you want to use high resolution modes (>= 1280x1024x16) you may need to increase the vga memory option. Since QEMU 2.9 the default VGA display type is `std` for all OS types besides some Windows versions (XP and older) which use `cirrus`. The `qxl` option enables the SPICE display server. For win* OS you can select how many independent displays you want, Linux guests can add displays themself. You can also run without any graphic card, using a serial device as terminal. See the [Proxmox documentation](https://pve.proxmox.com/pve-docs/pve-admin-guide.html#qm_virtual_machines_settings) section 10.2.8 for more information and available configuration parameters. (see [below for nested schema](#nestedatt--vga))

<a id="nestedatt--clone"></a>
### Nested Schema for `clone`

Required:

- `source_vm_id` (Number) Source VM/template ID to clone from.

Optional:

- `bandwidth_limit` (Number) Clone bandwidth limit in MB/s.
- `full` (Boolean) Perform a full clone (true) or linked clone (false).
- `pool_id` (String) Pool to assign the cloned VM to.
- `retries` (Number) Number of retries for clone operations.
- `snapshot_name` (String) Snapshot name to clone from.
- `source_node_name` (String) Source node of the VM/template. Defaults to target node if unset.
- `target_datastore` (String) Target datastore for cloned disks.
- `target_format` (String) Target disk format for clone (e.g., raw, qcow2).


<a id="nestedatt--cdrom"></a>
### Nested Schema for `cdrom`

Optional:

- `file_id` (String) The file ID of the CD-ROM, or `cdrom|none`. Defaults to `none` to leave the CD-ROM empty. Use `cdrom` to connect to the physical drive.


<a id="nestedatt--cpu"></a>
### Nested Schema for `cpu`

Optional:

- `affinity` (String) The CPU cores that are used to run the VM’s vCPU. The value is a list of CPU IDs, separated by commas. The CPU IDs are zero-based.  For example, `0,1,2,3` (which also can be shortened to `0-3`) means that the VM’s vCPUs are run on the first four CPU cores. Setting `affinity` is only allowed for `root@pam` authenticated user.
- `architecture` (String) The CPU architecture `<aarch64 | x86_64>` (defaults to the host). Setting `architecture` is only allowed for `root@pam` authenticated user.
- `cores` (Number) The number of CPU cores per socket (defaults to `1`).
- `flags` (Set of String) Set of additional CPU flags. Use `+FLAG` to enable, `-FLAG` to disable a flag. Custom CPU models can specify any flag supported by QEMU/KVM, VM-specific flags must be from the following set for security reasons: `pcid`, `spec-ctrl`, `ibpb`, `ssbd`, `virt-ssbd`, `amd-ssbd`, `amd-no-ssb`, `pdpe1gb`, `md-clear`, `hv-tlbflush`, `hv-evmcs`, `aes`.
- `hotplugged` (Number) The number of hotplugged vCPUs (defaults to `0`).
- `limit` (Number) Limit of CPU usage (defaults to `0` which means no limit).
- `numa` (Boolean) Enable NUMA (defaults to `false`).
- `sockets` (Number) The number of CPU sockets (defaults to `1`).
- `type` (String) Emulated CPU type, it's recommended to use `x86-64-v2-AES` or higher (defaults to `kvm64`). See https://pve.proxmox.com/pve-docs/pve-admin-guide.html#qm_virtual_machines_settings for more information.
- `units` (Number) CPU weight for a VM. Argument is used in the kernel fair scheduler. The larger the number is, the more CPU time this VM gets. Number is relative to weights of all the other running VMs.


<a id="nestedatt--delete"></a>
### Nested Schema for `delete`

Optional:

- `disk` (List of String) Disk slots to delete (e.g., scsi2).
- `network` (List of String) Network slots to delete (e.g., net1).


<a id="nestedatt--disk"></a>
### Nested Schema for `disk`

Optional:

- `aio` (String) AIO mode (io_uring, native, threads).
- `backup` (Boolean) Include disk in backups.
- `cache` (String) Cache mode.
- `datastore_id` (String) Target datastore for new disks when file is not provided.
- `discard` (String) Discard/trim behavior.
- `file` (String) Existing volume reference (e.g., local-lvm:vm-100-disk-0).
- `format` (String) Disk format (raw, qcow2, vmdk).
- `import_from` (String) Import source volume/file id.
- `iothread` (Boolean) Use IO thread.
- `media` (String) Disk media (e.g., disk, cdrom).
- `replicate` (Boolean) Consider disk for replication.
- `serial` (String) Disk serial number.
- `size_gb` (Number) Disk size (GiB) when creating new disks.
- `ssd` (Boolean) Mark disk as SSD.


<a id="nestedatt--memory"></a>
### Nested Schema for `memory`

Optional:

- `balloon` (Number) Minimum guaranteed memory in MiB via balloon device. This is the floor amount of RAM that is always guaranteed to the VM. Setting to `0` disables the balloon driver entirely (defaults to `0`). 

**How it works:** The host can reclaim memory between `balloon` and `size` when under memory pressure. The VM is guaranteed to always have at least `balloon` MiB available.
- `hugepages` (String) Enable hugepages for VM memory allocation. Hugepages can improve performance for memory-intensive workloads by reducing TLB misses. 

**Options:**
- `2` - Use 2 MiB hugepages
- `1024` - Use 1 GiB hugepages
- `any` - Use any available hugepage size
- `keep_hugepages` (Boolean) Don't release hugepages when the VM shuts down. By default, hugepages are released back to the host when the VM stops. Setting this to `true` keeps them allocated for faster VM startup (defaults to `false`).
- `shares` (Number) CPU scheduler priority for memory ballooning. This is used by the kernel fair scheduler. Higher values mean this VM gets more CPU time during memory ballooning operations. The value is relative to other running VMs (defaults to `1000`).
- `size` (Number) Total memory available to the VM in MiB. This is the total RAM the VM can use. When ballooning is enabled (balloon > 0), memory between `balloon` and `size` can be reclaimed by the host. When ballooning is disabled (balloon = 0), this is the fixed amount of RAM allocated to the VM (defaults to `512` MiB).


<a id="nestedatt--network"></a>
### Nested Schema for `network`

Optional:

- `bridge` (String) Bridge name.
- `firewall` (Boolean) Enable firewall on this interface.
- `link_down` (Boolean) Keep link down.
- `mac_address` (String) MAC address (computed if omitted).
- `model` (String) NIC model (e.g., virtio, e1000).
- `mtu` (Number) Interface MTU.
- `queues` (Number) Number of multiqueue NIC queues.
- `rate_limit` (Number) Rate limit (MB/s).
- `tag` (Number) VLAN tag.
- `trunks` (Set of Number) Trunk VLAN IDs.


<a id="nestedatt--rng"></a>
### Nested Schema for `rng`

Optional:

- `max_bytes` (Number) Maximum bytes of entropy allowed to get injected into the guest every period. Use 0 to disable limiting (potentially dangerous).
- `period` (Number) Period in milliseconds to limit entropy injection to the guest. Use 0 to disable limiting (potentially dangerous).
- `source` (String) The file on the host to gather entropy from. In most cases, `/dev/urandom` should be preferred over `/dev/random` to avoid entropy-starvation issues on the host.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--vga"></a>
### Nested Schema for `vga`

Optional:

- `clipboard` (String) Enable a specific clipboard. If not set, depending on the display type the SPICE one will be added. Currently only `vnc` is available. Migration with VNC clipboard is not supported by Proxmox.
- `memory` (Number) The VGA memory in megabytes (4-512 MB). Has no effect with serial display.
- `type` (String) The VGA type (defaults to `std`).
