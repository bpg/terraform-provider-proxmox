resource "proxmox_ceph_pool" "example" {
  node_name = "pve"
  name      = "tank"

  application       = "rbd"
  size              = 3
  min_size          = 2
  pg_num            = 128
  pg_autoscale_mode = "warn"
}
