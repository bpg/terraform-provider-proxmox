resource "proxmox_disk_zfs" "example" {
  node_name = "pve"
  name      = "tank"
  devices   = ["/dev/sdb", "/dev/sdc"]
  raidlevel = "mirror"

  ashift      = 12
  compression = "lz4"

  add_storage    = true
  cleanup_config = true
  cleanup_disks  = false
}
