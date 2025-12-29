resource "proxmox_virtual_environment_storage_nfs" "example" {
  id     = "example-nfs"
  nodes  = ["pve"]
  server = "10.0.0.10"
  export = "/exports/proxmox"

  content = ["images", "iso", "backup"]
  shared  = true

  options                  = "vers=4.2"
  preallocation            = "metadata"
  snapshot_as_volume_chain = true

  backups {
    max_protected_backups = 5
    keep_daily            = 7
  }
}

