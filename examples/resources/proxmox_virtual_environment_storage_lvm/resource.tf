resource "proxmox_virtual_environment_storage_lvm" "example" {
  id    = "example-lvm"
  nodes = ["pve"]

  volume_group = "vg0"
  content      = ["images"]

  wipe_removed_volumes = false
}

