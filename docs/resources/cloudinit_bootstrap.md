---
layout: page
title: proxmox_cloudinit_bootstrap
parent: Resources
subcategory: Virtual Environment
description: |-
  Bootstraps a VM with cloud-init through a temporary NoCloud https://docs.cloud-init.io/en/latest/reference/datasources/nocloud.html seed ISO, without SSH.
  Create the VM with started = false, lifecycle { ignore_changes = [started] } (or the next apply stops it again) and an empty CD-ROM (file_id = "none"). On create this resource builds the seed ISO, uploads it to the datastore on the node the VM is on, attaches it, starts the VM, waits for cloud-init to finish, then ejects and deletes the ISO. Nothing is left on the node afterwards, so the VM can be migrated freely.
  Completion is read from /run/cloud-init/result.json through the QEMU guest agent, which must be installed in the guest (installing it from the same cloud-init config works). The wait survives guest reboots such as package_reboot_if_required. If cloud-init reports errors or the wait times out, creation fails, the ISO stays attached for debugging and the resource is tainted. cloud-init only runs on a VM's first boot, so replace the VM to bootstrap it again. Applying again while the VM is running, e.g. after a failed or interrupted apply, only removes the ISO and fails. Don't stop the VM and apply again: booting it a second time does not run cloud-init's setup again, so the bootstrap could report success. Destroying the resource also removes the ISO.
  cloud-init only reports the errors of its last boot. An error before a reboot it triggers itself (e.g. package_reboot_if_required after a failed package install) is not reported; check /var/log/cloud-init.log in the guest if a bootstrapped VM looks incomplete.
  After cloud-init finishes, creation waits up to 2 minutes for every network device to get an IPv4 address, then warns if one still has none.
  Use lifecycle.replace_triggered_by on the VM so a replaced VM is bootstrapped again. If the VM is deleted outside Terraform, the bootstrap is removed from state on refresh, so it is created again with the VM. Changing any attribute but vm_id later only updates state and does not affect the running VM.
  Seed contents are stored in Terraform state, inside the guest under /var/lib/cloud, and in provider debug logs (TF_LOG=DEBUG logs the uploaded ISO). While the ISO is on the datastore (during the bootstrap, or after a failed one), anyone who can read ISOs there can read it.
---

# Resource: proxmox_cloudinit_bootstrap

Bootstraps a VM with cloud-init through a temporary [NoCloud](https://docs.cloud-init.io/en/latest/reference/datasources/nocloud.html) seed ISO, without SSH.

Create the VM with `started = false`, `lifecycle { ignore_changes = [started] }` (or the next apply stops it again) and an empty CD-ROM (`file_id = "none"`). On create this resource builds the seed ISO, uploads it to the datastore on the node the VM is on, attaches it, starts the VM, waits for cloud-init to finish, then ejects and deletes the ISO. Nothing is left on the node afterwards, so the VM can be migrated freely.

Completion is read from `/run/cloud-init/result.json` through the QEMU guest agent, which must be installed in the guest (installing it from the same cloud-init config works). The wait survives guest reboots such as `package_reboot_if_required`. If cloud-init reports errors or the wait times out, creation fails, the ISO stays attached for debugging and the resource is tainted. cloud-init only runs on a VM's first boot, so replace the VM to bootstrap it again. Applying again while the VM is running, e.g. after a failed or interrupted apply, only removes the ISO and fails. Don't stop the VM and apply again: booting it a second time does not run cloud-init's setup again, so the bootstrap could report success. Destroying the resource also removes the ISO.

cloud-init only reports the errors of its last boot. An error before a reboot it triggers itself (e.g. `package_reboot_if_required` after a failed package install) is not reported; check `/var/log/cloud-init.log` in the guest if a bootstrapped VM looks incomplete.

After cloud-init finishes, creation waits up to 2 minutes for every network device to get an IPv4 address, then warns if one still has none.

Use `lifecycle.replace_triggered_by` on the VM so a replaced VM is bootstrapped again. If the VM is deleted outside Terraform, the bootstrap is removed from state on refresh, so it is created again with the VM. Changing any attribute but `vm_id` later only updates state and does not affect the running VM.

Seed contents are stored in Terraform state, inside the guest under `/var/lib/cloud`, and in provider debug logs (`TF_LOG=DEBUG` logs the uploaded ISO). While the ISO is on the datastore (during the bootstrap, or after a failed one), anyone who can read ISOs there can read it.

## Example Usage

```terraform
resource "proxmox_virtual_environment_vm" "example" {
  node_name = "pve"
  name      = "example"
  machine   = "q35"

  # proxmox_cloudinit_bootstrap starts the VM once the seed ISO is attached.
  started = false

  # The VM ends up running, so force-stop it on destroy; a graceful ACPI shutdown
  # may hang if the guest agent is not installed.
  stop_on_destroy = true

  disk {
    datastore_id = "local-lvm"
    import_from  = "local:import/noble-server-cloudimg-amd64.qcow2"
    interface    = "virtio0"
    size         = 20
  }

  # Empty drive the seed ISO is attached to.
  cdrom {
    interface = "ide0"
    file_id   = "none"
  }

  network_device {
    bridge = "vmbr0"
  }

  agent {
    enabled = true
  }

  lifecycle {
    ignore_changes = [started]
  }
}

resource "proxmox_cloudinit_bootstrap" "example" {
  vm_id        = proxmox_virtual_environment_vm.example.vm_id
  datastore_id = "local"

  user_data = <<-EOT
    #cloud-config
    packages: [qemu-guest-agent]
    runcmd:
      - [systemctl, enable, --now, qemu-guest-agent]
  EOT

  lifecycle {
    # Bootstrap a replaced VM again.
    replace_triggered_by = [proxmox_virtual_environment_vm.example.id]
  }
}

# The guest's IPv4 addresses by interface name, e.g.
# proxmox_cloudinit_bootstrap.example.ipv4_addresses["eth0"][0]

# A VM with two NICs. Without `initialization`, Proxmox generates no network config,
# so `network_config` names the NICs by the MACs Proxmox assigned.
resource "proxmox_virtual_environment_vm" "multi_nic" {
  node_name = "pve"
  name      = "multi-nic"
  machine   = "q35"

  started         = false
  stop_on_destroy = true

  disk {
    datastore_id = "local-lvm"
    import_from  = "local:import/noble-server-cloudimg-amd64.qcow2"
    interface    = "virtio0"
    size         = 20
  }

  cdrom {
    interface = "ide0"
    file_id   = "none"
  }

  network_device {
    bridge = "vmbr0"
  }

  network_device {
    bridge = "vmbr1"
  }

  agent {
    enabled = true
  }

  lifecycle {
    ignore_changes = [started]
  }
}

resource "proxmox_cloudinit_bootstrap" "multi_nic" {
  vm_id        = proxmox_virtual_environment_vm.multi_nic.vm_id
  datastore_id = "local"

  user_data = <<-EOT
    #cloud-config
    packages: [qemu-guest-agent]
    runcmd:
      - [systemctl, enable, --now, qemu-guest-agent]
  EOT

  network_config = yamlencode({
    version = 1
    config = [for i, nic in proxmox_virtual_environment_vm.multi_nic.network_device : {
      type        = "physical"
      name        = "eth${i}"
      mac_address = lower(nic.mac_address)
      subnets     = [{ type = "dhcp" }]
    }]
  })

  lifecycle {
    replace_triggered_by = [proxmox_virtual_environment_vm.multi_nic.id]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `datastore_id` (String) The datastore on the VM's node to upload the temporary ISO to. It must allow `iso` content.
- `user_data` (String, Sensitive) The `user-data` file contents, e.g. a `#cloud-config` document.
- `vm_id` (Number) The VM to bootstrap. Its node is looked up at apply time.

### Optional

- `interface` (String) The CD-ROM interface to attach the ISO to. Defaults to `ide0` (`ide2` is the `initialization` default).
- `meta_data` (String, Sensitive) The `meta-data` file contents. Defaults to `instance-id` and `local-hostname` set to the VM name.
- `network_config` (String, Sensitive) The `network-config` file contents (network config version 1 or 2).
- `timeout` (Number) How long to wait for cloud-init to finish, in seconds. Defaults to 1800.
- `vendor_data` (String, Sensitive) The `vendor-data` file contents.

### Read-Only

- `id` (String) The VM ID.
- `ipv4_addresses` (Map of List of String) IPv4 addresses of the VM's network devices reported by the guest agent, keyed by guest interface name. Refreshed while the VM is running and every device has an address; the last known addresses are kept otherwise.
