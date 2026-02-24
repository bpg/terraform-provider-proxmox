---
layout: page
page_title: "Clone VM: Choosing Between Resources"
subcategory: Guides
description: |-
    Guide to help you choose between proxmox_virtual_environment_vm and proxmox_virtual_environment_cloned_vm for VM cloning
---

# Clone VM: Choosing Between Resources

## Overview

Terraform provider for Proxmox offers two approaches for cloning VMs:

1. **`proxmox_virtual_environment_vm`** (Legacy VM resource) - Simple clone support for straightforward scenarios
2. **`proxmox_virtual_environment_cloned_vm`** (Framework resource) - Advanced clone management with explicit device control

This guide helps you choose the right resource and migrate if needed.

## When to Use Each Resource

### Use `proxmox_virtual_environment_vm` (Legacy) When

- You have simple clone requirements
- Templates have minimal device configuration
- You're okay with Terraform managing all inherited configuration
- You don't need fine-grained control over specific devices
- You want to avoid VM recreation during migration

### Use `proxmox_virtual_environment_cloned_vm` When

- You need explicit control over which devices are managed
- Templates have complex multi-device configurations
- You want to preserve some inherited config without Terraform tracking it
- You need map-based device addressing (`net0`, `scsi0`)
- You want explicit deletion semantics for inherited devices
- You're encountering drift issues with inherited configuration

## Understanding Clone Challenges

When cloning a VM in Proxmox, the cloned VM inherits configuration from the template (network devices, disks, CPU, memory, etc.). This creates challenges in Terraform:

### Issues with Traditional Clone Approach

1. **Unknown inherited state**: After cloning, Terraform doesn't know which settings came from the template vs. which were explicitly configured
2. **List-based device addressing**: Network devices and disks modeled as lists make it difficult to track which specific device to update
3. **Drift detection issues**: Changes to inherited settings detected as drift even though they weren't explicitly managed
4. **Update confusion**: Unclear whether omitting configuration means "don't manage it" or "delete it"

### How `cloned_vm` Solves These Issues

The `proxmox_virtual_environment_cloned_vm` resource addresses these challenges with:

- **Explicit opt-in management**: Only configuration you explicitly declare is managed by Terraform
- **Map-based devices**: Network and disk devices use slot-based keys (`net0`, `scsi0`) instead of lists
- **Clear deletion semantics**: Use the `delete` block to explicitly remove inherited devices
- **No drift for unmanaged config**: Inherited settings that aren't declared in Terraform don't cause drift

## Migration Examples

### Example 1: Simple Clone (Legacy VM resource → cloned_vm)

**Before (legacy VM resource):**

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

1. Resource type: `proxmox_virtual_environment_vm` → `proxmox_virtual_environment_cloned_vm`
2. Clone syntax: `clone { vm_id = X }` → `clone = { source_vm_id = X }` (attribute syntax, not block)
3. Network devices: `network_device { ... }` → `network = { net0 = { ... } }` (map-based, not list)
4. Disk devices: `disk { interface = "scsi0" }` → `disk = { scsi0 = { ... } }` (map-based, not list)

### Example 2: Continue Using Legacy VM Resource

The legacy VM resource (`proxmox_virtual_environment_vm`) still supports cloning and will continue to work. This is a valid option if:

- You cannot recreate your VMs
- You're willing to accept the clone semantic limitations
- You don't need the new explicit management features

**Using legacy VM resource:**

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

## Detailed Migration Examples

### Migrating Network Devices

**Before (legacy VM resource):**

```terraform
resource "proxmox_virtual_environment_vm" "old" {
  node_name = "pve"

  clone {
    vm_id = 100
  }

  # List-based - order matters
  network_device {
    bridge  = "vmbr0"
    model   = "virtio"
  }

  network_device {
    bridge   = "vmbr1"
    model    = "virtio"
    vlan_id  = 100
  }
}
```

**After (cloned_vm):**

```terraform
resource "proxmox_virtual_environment_cloned_vm" "new" {
  node_name = "pve"

  clone = {
    source_vm_id = 100
  }

  # Map-based - explicit slots
  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
    }

    net1 = {
      bridge = "vmbr1"
      model  = "virtio"
      tag    = 100  # vlan_id renamed to tag
    }
  }
}
```

### Migrating Disk Devices

**Before (legacy VM resource):**

```terraform
resource "proxmox_virtual_environment_vm" "old" {
  node_name = "pve"

  clone {
    vm_id = 100
  }

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

**After (cloned_vm):**

```terraform
resource "proxmox_virtual_environment_cloned_vm" "new" {
  node_name = "pve"

  clone = {
    source_vm_id = 100
  }

  disk = {
    scsi0 = {
      datastore_id = "local-lvm"
      size_gb      = 50  # size renamed to size_gb
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

#### Scenario: Template has 3 NICs, you only want to manage 1

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
  clone = {
    source_vm_id = proxmox_virtual_environment_vm.template.id
  }

  network = {
    net0 = {
      bridge = "vmbr0"
      tag    = 100  # Only net0 is managed
    }
    # net1 and net2 are inherited but not managed
  }
}
```

#### Scenario: Remove specific inherited devices

```terraform
# Clone but delete net1 from the template
resource "proxmox_virtual_environment_cloned_vm" "selective_delete" {
  clone = {
    source_vm_id = proxmox_virtual_environment_vm.template.id
  }

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

### Migrating Clone Options

All clone options from the legacy VM resource are supported in `cloned_vm`:

**Before (legacy VM resource):**

```terraform
resource "proxmox_virtual_environment_vm" "old" {
  node_name = "pve"

  clone {
    vm_id         = 100
    full          = true
    datastore_id  = "local-lvm"
    node_name     = "pve-source"
    retries       = 3
  }
}
```

**After (cloned_vm):**

```terraform
resource "proxmox_virtual_environment_cloned_vm" "new" {
  node_name = "pve"

  clone = {
    source_vm_id      = 100
    full              = true
    target_datastore  = "local-lvm"      # datastore_id → target_datastore
    source_node_name  = "pve-source"    # node_name → source_node_name
    retries           = 3
  }
}
```

**Clone Option Mapping:**

| Legacy Attribute | cloned_vm Attribute | Notes |
| ---------------- | ------------------- | ----- |
| `vm_id` | `source_vm_id` | Required in both |
| `full` | `full` | Defaults to `true` in both |
| `datastore_id` | `target_datastore` | Renamed for clarity |
| `node_name` | `source_node_name` | Renamed for clarity |
| `retries` | `retries` | Defaults to `3` in both |

**Additional clone options available in `cloned_vm` (not in legacy):**

- `target_format` - Target disk format (e.g., raw, qcow2)
- `snapshot_name` - Clone from a specific snapshot
- `pool_id` - Assign cloned VM to a pool
- `bandwidth_limit` - Clone bandwidth limit in MB/s

## Migration Checklist

When migrating from legacy VM resource (`proxmox_virtual_environment_vm`) with clone to `cloned_vm`:

- [ ] Update resource type from `proxmox_virtual_environment_vm` to `proxmox_virtual_environment_cloned_vm`
- [ ] Change `clone { vm_id = X }` to `clone = { source_vm_id = X }`
- [ ] Update clone block attributes:
  - [ ] `datastore_id` → `target_datastore`
  - [ ] `node_name` (in clone block) → `source_node_name`
- [ ] Convert network devices from `network_device` blocks to `network` map
  - [ ] Determine slot names (`net0`, `net1`, etc.)
  - [ ] Convert `vlan_id` attribute to `tag`
- [ ] Convert disk devices from `disk` blocks to `disk` map
  - [ ] Determine slot names (`scsi0`, `virtio0`, etc.)
  - [ ] Use `size_gb` instead of `size`
- [ ] Decide which inherited devices to manage vs. preserve
- [ ] Add `delete` block if you need to remove inherited devices
- [ ] Test in a non-production environment first
- [ ] Plan for VM recreation (cloned_vm will create a new VM)

## State Migration

~> Changing from legacy VM resource (`proxmox_virtual_environment_vm`) to `cloned_vm` requires recreating the VM. Terraform will:

1. Create the new cloned VM
2. Destroy the old legacy VM resource

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
  id        = 123
  node_name = "pve"

  clone = {
    source_vm_id = 100
  }

  # ... rest of config
}
```

Then import the existing VM:

```bash
terraform import proxmox_virtual_environment_cloned_vm.imported pve/123
```

## FAQs

**Q: Will the legacy VM resource's clone support be removed?**

A: No current plans to remove clone from `proxmox_virtual_environment_vm`. However, new features will focus on the Framework-based `cloned_vm` resource.

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
- [VM Lifecycle Management guide](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/guides/vm-lifecycle) — destroy semantics and default differences between `vm` and `cloned_vm`
- [GitHub Issues](https://github.com/bpg/terraform-provider-proxmox/issues)
