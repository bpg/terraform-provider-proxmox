resource "proxmox_ceph_mds" "example" {
  node_name = "pve"
  name      = "pve-mds"
}
