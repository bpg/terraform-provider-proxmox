---
layout: page
page_title: "Migration Guide: VM2 Clone Deprecation"
subcategory: Guides
description: |-
    Migration guide for users affected by the removal of clone support from proxmox_virtual_environment_vm2
---

# Migration Guide: VM2 Clone Deprecation

## Overview

The `proxmox_virtual_environment_vm2` resource no longer supports VM cloning via the `clone` block. This functionality has been moved to a dedicated resource: `proxmox_virtual_environment_cloned_vm`.

This migration guide explains why this change was made and how to migrate your existing configurations.

## Why This Change?

The clone functionality was removed from `vm2` to address fundamental semantic incompatibilities between Terraform's declarative state model and Proxmox's clone operation:

### The Problem with Clone in `vm2`

When cloning a VM in Proxmox, the cloned VM inherits configuration from the template (network devices, disks, CPU, memory, etc.). This creates ambiguity in Terraform:

1. **Unknown inherited state**: After cloning, Terraform doesn't know which settings came from the template vs. which were explicitly configured
2. **List-based device addressing**: Network devices and disks were modeled as lists, making it difficult to track which specific device to update
3. **Drift detection issues**: Changes to inherited settings would be detected as drift even though they weren't explicitly managed
4. **Update confusion**: It was unclear whether omitting a configuration meant "don't manage it" or "delete it"

### The Solution: Dedicated `cloned_vm` Resource

The new `proxmox_virtual_environment_cloned_vm` resource was designed from the ground up to handle clone semantics correctly:

- **Explicit opt-in management**: Only configuration you explicitly declare is managed by Terraform
- **Map-based devices**: Network and disk devices use slot-based keys (`net0`, `scsi0`) instead of lists
- **Clear deletion semantics**: Use the `delete` block to explicitly remove inherited devices
- **No drift for unmanaged config**: Inherited settings that aren't declared in Terraform don't cause drift

## Migration Paths

### Option 1: Replace with `cloned_vm` (Recommended for New Deployments)

If you're starting fresh or can recreate your VMs, migrate to the new resource:

**Before (vm2 with clone - NO LONGER SUPPORTED):**

```terraform
resource "proxmox_virtual_environment_vm2" "cloned_vm" {
  node_name = "pve"
  name      = "my-cloned-vm"

  clone {
    vm_id = 100
  }

  cpu {
    cores = 4
  }

  network {
    bridge = "vmbr0"
    model  = "virtio"
  }
}
```

**After (using cloned_vm):**

```terraform
resource "proxmox_virtual_environment_cloned_vm" "cloned_vm" {
  node_name = "pve"
  name      = "my-cloned-vm"

  clone = {
    source_vm_id = 100
  }

  cpu = {
    cores = 4
  }

  # Map-based network devices
  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
    }
  }
}
```

**Key Changes:**

1. Resource type: `proxmox_virtual_environment_vm2` → `proxmox_virtual_environment_cloned_vm`
2. Clone syntax: `clone { vm_id = X }` → `clone = { source_vm_id = X }`
3. Network devices: List (`network { ... }`) → Map (`network = { net0 = { ... } }`)
4. Disk devices: List (`disk { interface = "scsi0" }`) → Map (`disk = { scsi0 = { ... } }`)

### Option 2: Continue Using Legacy Provider

The legacy SDK provider (`proxmox_virtual_environment_vm`) still supports cloning and will continue to work. This is a valid option if:

- You cannot recreate your VMs
- You're willing to accept the clone semantic limitations
- You don't need the new explicit management features

**Using legacy provider:**

```terraform
resource "proxmox_virtual_environment_vm" "cloned_vm" {
  node_name = "pve"
  name      = "my-cloned-vm"

  clone {
    vm_id = 100
  }

  cpu {
    cores = 4
  }

  network_device {
    bridge = "vmbr0"
    model  = "virtio"
  }
}
```

Note: This uses the legacy resource without the `2` suffix.

### Option 3: Create VMs Without Cloning

If you were using clone primarily for convenience, consider creating VMs from scratch:

```terraform
resource "proxmox_virtual_environment_vm2" "from_scratch" {
  node_name = "pve"
  name      = "my-vm"

  disk {
    datastore_id = "local-lvm"
    file_id      = proxmox_virtual_environment_download_file.cloud_image.id
    interface    = "virtio0"
    size         = 20
  }

  cpu {
    cores = 4
  }

  memory {
    dedicated = 4096
  }

  network {
    bridge = "vmbr0"
    model  = "virtio"
  }
}
```

## Detailed Migration Examples

### Migrating Network Devices

**Old (vm2 - NO LONGER WORKS):**

```terraform
resource "proxmox_virtual_environment_vm2" "old" {
  clone { vm_id = 100 }

  # List-based - order matters
  network {
    bridge = "vmbr0"
    model  = "virtio"
  }

  network {
    bridge = "vmbr1"
    model  = "virtio"
    vlan   = 100
  }
}
```

**New (cloned_vm):**

```terraform
resource "proxmox_virtual_environment_cloned_vm" "new" {
  clone = { source_vm_id = 100 }

  # Map-based - explicit slots
  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
    }

    net1 = {
      bridge = "vmbr1"
      model  = "virtio"
      tag    = 100  # vlan renamed to tag
    }
  }
}
```

### Migrating Disk Devices

**Old (vm2 - NO LONGER WORKS):**

```terraform
resource "proxmox_virtual_environment_vm2" "old" {
  clone { vm_id = 100 }

  disk {
    interface    = "scsi0"
    datastore_id = "local-lvm"
    size         = 50
  }

  disk {
    interface    = "scsi1"
    datastore_id = "local-lvm"
    size         = 100
  }
}
```

**New (cloned_vm):**

```terraform
resource "proxmox_virtual_environment_cloned_vm" "new" {
  clone = { source_vm_id = 100 }

  disk = {
    scsi0 = {
      datastore_id = "local-lvm"
      size_gb      = 50
    }

    scsi1 = {
      datastore_id = "local-lvm"
      size_gb      = 100
    }
  }
}
```

### Handling Inherited Devices

One of the key benefits of `cloned_vm` is explicit control over inherited configuration.

**Scenario: Template has 3 NICs, you only want to manage 1**

```terraform
# Template VM has net0, net1, net2
resource "proxmox_virtual_environment_vm" "template" {
  template = true

  network_device { bridge = "vmbr0" }
  network_device { bridge = "vmbr1" }
  network_device { bridge = "vmbr2" }
}

# Clone and manage only net0, leave net1 and net2 as-is
resource "proxmox_virtual_environment_cloned_vm" "partial_management" {
  clone = { source_vm_id = proxmox_virtual_environment_vm.template.id }

  network = {
    net0 = {
      bridge = "vmbr0"
      tag    = 100  # Only net0 is managed
    }
    # net1 and net2 are inherited but not managed
  }
}
```

**Scenario: Remove specific inherited devices**

```terraform
# Clone but delete net1 from the template
resource "proxmox_virtual_environment_cloned_vm" "selective_delete" {
  clone = { source_vm_id = proxmox_virtual_environment_vm.template.id }

  network = {
    net0 = {
      bridge = "vmbr0"
    }
  }

  # Explicitly delete inherited devices
  delete = {
    network = ["net1", "net2"]
  }
}
```

## Migration Checklist

When migrating from `vm2` with clone to `cloned_vm`:

- [ ] Update resource type from `proxmox_virtual_environment_vm2` to `proxmox_virtual_environment_cloned_vm`
- [ ] Change `clone { vm_id = X }` to `clone = { source_vm_id = X }`
- [ ] Convert network devices from list to map format
  - [ ] Determine slot names (`net0`, `net1`, etc.)
  - [ ] Convert `vlan` attribute to `tag`
- [ ] Convert disk devices from list to map format
  - [ ] Determine slot names (`scsi0`, `virtio0`, etc.)
  - [ ] Use `size_gb` instead of `size`
- [ ] Decide which inherited devices to manage vs. preserve
- [ ] Add `delete` block if you need to remove inherited devices
- [ ] Test in a non-production environment first
- [ ] Plan for VM recreation (cloned_vm will create a new VM)

## State Migration

**Important**: Changing from `vm2` to `cloned_vm` requires recreating the VM. Terraform will:

1. Create the new cloned VM
2. Destroy the old vm2 resource

To minimize downtime:

1. Export any important data from the old VM
2. Apply the new configuration (Terraform will create new VM)
3. Verify the new VM works correctly
4. Let Terraform destroy the old VM (or manually destroy if needed)

If you need to preserve the VM ID, you can:

1. Manually import the existing VM into the new resource
2. Use the `id` attribute to specify the desired VM ID

```terraform
resource "proxmox_virtual_environment_cloned_vm" "imported" {
  id        = 123  # Preserve existing VM ID
  node_name = "pve"

  clone = { source_vm_id = 100 }
  # ... rest of config
}
```

Then import the existing VM:

```bash
terraform import proxmox_virtual_environment_cloned_vm.imported pve/123
```

## FAQs

**Q: Can I still use `vm2` for non-cloned VMs?**

A: Yes! The `proxmox_virtual_environment_vm2` resource still works for creating VMs from scratch. Only the clone functionality was removed.

**Q: Will the legacy provider's clone support be removed?**

A: No current plans to remove clone from `proxmox_virtual_environment_vm` (without the `2`). However, new features will focus on the Framework-based `cloned_vm` resource.

**Q: What if I have complex clone configurations?**

A: The `cloned_vm` resource supports all clone options:
- Full vs. linked clones (`full = true/false`)
- Target datastore (`target_datastore`)
- Target format (`target_format`)
- Snapshot cloning (`snapshot_name`)
- Pool assignment (`pool_id`)
- Bandwidth limits (`bandwidth_limit`)

**Q: Can I mix managed and unmanaged devices?**

A: Yes! This is one of the key features. You can manage `net0` while leaving `net1` inherited from the template without Terraform tracking it.

**Q: What happens if I remove a device from my config?**

A: Terraform stops managing it but does NOT delete it from the VM. Use the `delete` block to explicitly remove devices.

## Getting Help

- [cloned_vm resource documentation](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/resources/virtual_environment_cloned_vm)
- [Clone VM guide](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/guides/clone-vm)
- [GitHub Issues](https://github.com/bpg/terraform-provider-proxmox/issues)
