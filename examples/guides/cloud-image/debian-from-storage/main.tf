resource "proxmox_virtual_environment_vm" "debian_vm" {
  name      = "test-debian"
  node_name = "pve"

  initialization {
    user_account {
      # do not use this in production, configure your own ssh key instead!
      username = "user"
      password = "password"
    }
  }

  disk {
    datastore_id = "local-lvm"
    file_id      = "local:iso/debian-12-genericcloud-amd64.img"
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    size         = 20
  }
}

