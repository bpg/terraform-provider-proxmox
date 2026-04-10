resource "proxmox_storage_zfspool" "example" {
  id    = "example-zfs"
  nodes = ["pve"]

  zfs_pool       = "rpool/data"
  content        = ["images"]
  thin_provision = true
  blocksize      = "64k"
}

