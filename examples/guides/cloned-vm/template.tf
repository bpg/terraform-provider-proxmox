# Create a template VM with multiple network interfaces
resource "proxmox_virtual_environment_vm" "multi_nic_template" {
  name      = "multi-nic-template"
  node_name = var.virtual_environment_node_name

  template = true
  started  = false

  machine     = "q35"
  bios        = "ovmf"
  description = "Template with multiple NICs for cloned_vm example"

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

  # Template has three network interfaces
  network_device {
    bridge = "vmbr0"
    model  = "virtio"
  }

  network_device {
    bridge = "vmbr0"
    model  = "virtio"
  }

  network_device {
    bridge = "vmbr0"
    model  = "virtio"
  }
}

resource "proxmox_virtual_environment_download_file" "ubuntu_cloud_image" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = var.virtual_environment_node_name

  url = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
}
