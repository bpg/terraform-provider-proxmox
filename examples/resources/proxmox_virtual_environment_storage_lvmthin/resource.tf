resource "proxmox_virtual_environment_storage_lvmthin" "example" {
  id    = "example-lvmthin"
  nodes = ["pve"]

  volume_group = "vg0"
  thin_pool    = "data"

  content = ["images"]
}

