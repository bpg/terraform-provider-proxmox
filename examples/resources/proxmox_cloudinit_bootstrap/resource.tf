resource "proxmox_virtual_environment_vm" "example" {
  node_name = "pve"
  name      = "example"
  machine   = "q35"

  # proxmox_cloudinit_bootstrap starts the VM once the seed ISO is attached.
  started = false

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

output "example_ipv4" {
  value = proxmox_cloudinit_bootstrap.example.ipv4_addresses
}
