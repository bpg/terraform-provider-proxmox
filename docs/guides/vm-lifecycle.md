---
layout: page
page_title: "VM Lifecycle Management"
subcategory: Guides
description: |-
    This guide covers creating, updating, and destroying VMs, including hotplug behavior and destroy semantics.
---

# VM Lifecycle Management

This guide covers the full VM lifecycle — from creation through day-2 updates to destruction — focusing on Proxmox-specific behaviors that affect how Terraform manages your VMs.

## Creating a VM

A typical VM starts with a cloud image download and a `proxmox_virtual_environment_vm` resource:

```terraform
resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = var.virtual_environment_node_name
  url          = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
}

resource "proxmox_virtual_environment_vm" "example" {
  name      = "example-vm"
  node_name = var.virtual_environment_node_name

  description = "Managed by Terraform"
  machine     = "q35"
  bios        = "ovmf"
  started     = true

  # Always set stop_on_destroy when started = true,
  # otherwise Terraform will attempt a graceful ACPI shutdown
  # that may hang if the guest agent is not installed.
  stop_on_destroy = true

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
  }

  network_device {
    bridge = "vmbr0"
  }
}
```

Key points:

- **`started` defaults to `true`** — the VM will boot after creation. Set `started = false` for templates.
- **Templates ignore `started`** — when `template = true`, the VM is never started regardless of the `started` value. Set `started = false` explicitly to avoid confusion.
- **`stop_on_destroy = true`** is recommended for started VMs — without it, Terraform sends an ACPI shutdown and waits for the guest to power off. If the guest agent is not installed or the guest hangs, the destroy will time out.

## Updating a VM

Proxmox supports hotplugging some attributes while a VM is running. Other changes require a reboot, and some force full recreation of the VM.

```terraform
resource "proxmox_virtual_environment_vm" "hotplug_example" {
  name      = "hotplug-example"
  node_name = var.virtual_environment_node_name

  started         = true
  stop_on_destroy = true

  # reboot_after_update defaults to true.
  # When a non-hotpluggable attribute changes (e.g. cpu.cores),
  # Terraform will automatically reboot the VM to apply it.
  # Set to false if you prefer to reboot manually.
  reboot_after_update = true

  cpu {
    # Changing cores or sockets requires a reboot.
    cores   = 4
    sockets = 1
  }

  memory {
    # Memory is hotpluggable (increase only) when the VM's hotplug
    # setting includes "memory".
    dedicated = 4096
  }

  disk {
    datastore_id = var.datastore_id
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    # Disks can only grow. Shrinking produces an error.
    size = 40
  }

  network_device {
    bridge = "vmbr0"
  }

  # Adding a second NIC is hotpluggable.
  network_device {
    bridge  = "vmbr0"
    vlan_id = 100
  }
}
```

### Hotplug vs reboot vs recreate

When `reboot_after_update = true` (default), Terraform automatically reboots the VM when a non-hotpluggable attribute changes. Set it to `false` if you prefer to schedule reboots manually — Terraform will log a warning instead.

| Change | Behavior |
| ------ | -------- |
| `memory.dedicated` (increase) | Hotplug (when VM `hotplug` includes `memory`) |
| Adding a network device | Hotplug |
| `cpu.cores`, `cpu.sockets` | Requires reboot |
| `cpu.type` | Requires reboot |
| `machine`, `bios` | Requires reboot |
| `disk.size` (increase) | Applied online, no reboot |
| `disk.size` (decrease) | **Error** — disks can only grow |
| `disk.interface` | Deletes old disk, creates new (data-destructive, not VM recreation) |
| `template` | Forces recreation |

-> **Tip:** Run `terraform plan` after changing VM attributes. The plan output will indicate whether a change triggers an in-place update or forces replacement.

## Destroying a VM

Three attributes control what happens when Terraform destroys a VM:

```terraform
resource "proxmox_virtual_environment_vm" "destroy_example" {
  name      = "destroy-example"
  node_name = var.virtual_environment_node_name

  started = true

  # When false (default), Terraform sends an ACPI shutdown and waits
  # for the guest to power off gracefully. Setting to true force-stops
  # the VM immediately, which is safer for started VMs without a
  # guest agent.
  stop_on_destroy = true

  # purge_on_destroy = true (default): Removes backup jobs, replication
  # entries, and HA configuration for this VM on destroy.
  purge_on_destroy = true

  # delete_unreferenced_disks_on_destroy = true (default for vm resource):
  # Deletes any disks not tracked by Terraform on destroy.
  # The cloned_vm resource defaults to false for safety.
  delete_unreferenced_disks_on_destroy = true

  cpu {
    cores = 2
  }

  memory {
    dedicated = 2048
  }

  disk {
    datastore_id = var.datastore_id
    interface    = "virtio0"
    size         = 20
  }

  network_device {
    bridge = "vmbr0"
  }
}
```

| Attribute | Default (`vm`) | Default (`cloned_vm`) | Effect |
| --------- | :------------: | :-------------------: | ------ |
| `stop_on_destroy` | `false` | `false` | `false` = ACPI shutdown (graceful), `true` = force stop |
| `purge_on_destroy` | `true` | `true` | Remove backup jobs, replication, and HA config |
| `delete_unreferenced_disks_on_destroy` | **`true`** | **`false`** | Delete disks not tracked by Terraform |

~> **Warning:** The `vm` resource defaults `delete_unreferenced_disks_on_destroy` to `true`, which deletes any disk not explicitly declared in your Terraform config. If you attach disks outside Terraform, set this to `false` to prevent data loss.

For the `cloned_vm` resource, see the [Clone VM guide](clone-vm) and the [Choosing Between Clone Resources guide](migration-vm-clone) for details on the `delete` block and inherited device management.

## Quick reference

| Operation | What to set |
| --------- | ----------- |
| Create a running VM | `started = true`, `stop_on_destroy = true` |
| Create a template | `template = true`, `started = false` |
| Resize disk up | Increase `disk.size` — applied online |
| Add memory (hot) | Increase `memory.dedicated` (requires `hotplug` to include `memory`) |
| Change CPU cores | Change `cpu.cores` — auto-reboots by default |
| Prevent auto-reboot | `reboot_after_update = false` |
| Force-stop on destroy | `stop_on_destroy = true` |
| Keep unmanaged disks | `delete_unreferenced_disks_on_destroy = false` |

Full example is available in the [examples/guides/vm-lifecycle](https://github.com/bpg/terraform-provider-proxmox/tree/main/examples/guides/vm-lifecycle) directory.
