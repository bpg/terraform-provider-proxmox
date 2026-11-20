resource "proxmox_download_file" "ubuntu_cloud_image" {
  content_type = "import"
  datastore_id = "local"
  node_name    = "pve"
  url          = "https://cloud-images.ubuntu.com/noble/current/noble-server-cloudimg-amd64.img"
  # Ubuntu .img files are qcow2.
  file_name = "noble-server-cloudimg-amd64.qcow2"
}

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
    import_from  = proxmox_download_file.ubuntu_cloud_image.id
    interface    = "virtio0"
    size         = 20
  }

  # Empty drive the seed ISO is attached to. Optional: on a free slot the
  # bootstrap leaves an empty CD-ROM behind.
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
    # Bootstrap a replaced VM again. Use the VM's id: the whole VM would also
    # trigger on in-place VM updates.
    # cloud-init applies the seed only on a VM's first boot, so changing user_data
    # later fails the apply: replace the VM to apply it, with
    # -replace=proxmox_virtual_environment_vm.example.
    replace_triggered_by = [proxmox_virtual_environment_vm.example.id]
  }
}

# The guest's IPv4 addresses by the interface names the guest reports.
output "example_ipv4_addresses" {
  value = proxmox_cloudinit_bootstrap.example.ipv4_addresses
}

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
    import_from  = proxmox_download_file.ubuntu_cloud_image.id
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
