---
layout: page
page_title: "Clone a VM with cloned_vm Resource"
subcategory: Guides
description: |-
    This guide explains how to use the proxmox_virtual_environment_cloned_vm resource to clone VMs from templates with explicit opt-in device management.
---

# Clone a VM with cloned_vm Resource

The `proxmox_virtual_environment_cloned_vm` resource provides explicit opt-in management: only devices and configuration explicitly listed in Terraform are managed. Inherited settings from templates are preserved unless explicitly overridden or deleted. This prevents accidental deletion of inherited devices and provides predictable behavior.

## Create a VM template

First, create a VM template that will be used as the source for cloning:

```terraform
resource "proxmox_virtual_environment_vm" "ubuntu_template" {
  name      = "ubuntu-template"
  node_name = var.virtual_environment_node_name

  template = true
  started  = false

  machine     = "q35"
  bios        = "ovmf"
  description = "Managed by Terraform"

  cpu {
    cores = 2
  }

  memory {
    dedicated = 2048
  }

  efi_disk {
    datastore_id = var.datastore_id
    type         = "4m"
  }

  disk {
    datastore_id = var.datastore_id
    file_id      = proxmox_virtual_environment_download_file.ubuntu_cloud_image.id
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    size         = 20
  }

  initialization {
    ip_config {
      ipv4 {
        address = "dhcp"
      }
    }

    user_data_file_id = proxmox_virtual_environment_file.user_data_cloud_config.id
  }

  network_device {
    bridge = "vmbr0"
  }
}

resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = var.virtual_environment_node_name

  url = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
}
```

## Clone using cloned_vm resource

Once you have a template, clone it using the `cloned_vm` resource:

```terraform
resource "proxmox_virtual_environment_cloned_vm" "ubuntu_clone" {
  node_name = var.virtual_environment_node_name
  name      = "ubuntu-clone"

  clone = {
    source_vm_id = proxmox_virtual_environment_vm.ubuntu_template.vm_id
    full         = true
  }

  # Only explicitly listed devices are managed
  # Network device inherited from template is preserved but not managed
  # To manage it, explicitly list it here:
  network = {
    net0 = {
      bridge = "vmbr0"
      model  = "virtio"
    }
  }

  # Memory configuration using new terminology
  memory = {
    size    = 2048  # Total memory available to VM
    balloon = 512   # Minimum guaranteed memory via balloon device
  }

  cpu = {
    cores = 2
  }
}

output "vm_id" {
  value = proxmox_virtual_environment_cloned_vm.ubuntu_clone.id
}
```

## Key Features

### Map-Based Device Addressing

Devices are addressed using slot names for stable references:

- **Network devices**: `net0`, `net1`, `net2`, etc.
- **Disk devices**: `scsi0`, `virtio0`, `sata0`, `ide0`, etc.

This provides stable, addressable references to specific devices that don't change when other devices are added or removed.

### Explicit Deletion

To delete inherited devices, use the `delete` block:

```terraform
resource "proxmox_virtual_environment_cloned_vm" "example" {
  # ... other configuration ...

  delete = {
    network = ["net1", "net2"]  # Remove inherited network devices
    disk    = ["scsi1"]          # Remove inherited disk devices
  }
}
```

### Opt-In Management

Removing a device from Terraform configuration **does not delete it** from the VM. It simply stops managing it. This prevents accidental data loss.

### Inherited Settings

Initialization and agent settings are inherited from the template. The cloned VM resource currently does not manage the following VM-level settings, so they must be defined on the template (or managed via `proxmox_virtual_environment_vm` with a `clone` block):

- BIOS / machine / boot order
- EFI disk / secure boot settings
- TPM state
- Cloud-init / initialization
- QEMU guest agent configuration
- PCI/USB passthrough, serial/audio devices, watchdog, VirtioFS

If you need to customize any of these after cloning, use the legacy `proxmox_virtual_environment_vm` resource with a `clone` block instead. See the [Clone a VM](clone-vm.md) guide for the legacy approach.

## When to Use

Use `proxmox_virtual_environment_cloned_vm` when:

- You want explicit control over which devices are managed
- You need to prevent accidental deletion of inherited devices
- You prefer map-based device addressing for stability
- You don't need to customize initialization after cloning

Use `proxmox_virtual_environment_vm` with `clone` block when:

- You need to customize initialization or agent settings after cloning
- You want Terraform to manage all VM configuration, including inherited settings
- You're migrating from existing configurations using the `vm` resource

Full example is available in the [examples/guides/cloned-vm](https://github.com/bpg/terraform-provider-proxmox/tree/main/examples/guides/cloned-vm) directory.

## See also

- [VM Lifecycle Management](vm-lifecycle) — destroy semantics and default differences between `vm` and `cloned_vm`
- [Choosing Between Clone Resources](migration-vm-clone) — detailed comparison and migration guide
