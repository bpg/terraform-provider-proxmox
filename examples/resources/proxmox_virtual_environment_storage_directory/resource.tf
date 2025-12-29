resource "proxmox_virtual_environment_storage_directory" "example" {
  id    = "example-dir"
  path  = "/var/lib/vz"
  nodes = ["pve"]

  content = ["images"]
  shared  = true
  disable = false

  backups {
    max_protected_backups = 5
    keep_daily            = 7
  }
}

