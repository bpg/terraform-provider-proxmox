resource "proxmox_virtual_environment_vm" "centos_vm" {
  name      = "test-ubuntu"
  node_name = "pve"

  initialization {
    user_account {
      keys     = [trimspace(tls_private_key.ubuntu_vm_key.public_key_openssh)]

      # do not use this in production, cofigure your own ssh key instead!
      username = "ubuntu"
    }
  }

  disk {
    datastore_id = "local-lvm"
    file_id      = proxmox_virtual_environment_file.ubuntu_cloud_image.id
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    size         = 20
  }
}

resource "proxmox_virtual_environment_file" "ubuntu_cloud_image" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = "pve"

  source_file {
    path      = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
  }
}
