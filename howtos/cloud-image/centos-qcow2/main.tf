resource "proxmox_virtual_environment_vm" "centos_vm" {
  name      = "test-centos"
  node_name = "pve"

  initialization {
    user_account {
      # do not use this in production, cofigure your own ssh key instead!
      username = "user"
      password = "password"
    }
  }

  disk {
    datastore_id = "local-lvm"
    file_id      = proxmox_virtual_environment_file.centos_cloud_image.id
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    size         = 20
  }
}

resource "proxmox_virtual_environment_file" "centos_cloud_image" {
  content_type = "iso"
  datastore_id = "local"
  node_name    = "pve"

  source_file {
    # you may download this image locally on your workstation and then use the local path instead of the remote URL
    path      = "https://cloud.centos.org/centos/8-stream/x86_64/images/CentOS-Stream-GenericCloud-8-20231113.0.x86_64.qcow2"
    file_name = "centos8.img"
  }
}
