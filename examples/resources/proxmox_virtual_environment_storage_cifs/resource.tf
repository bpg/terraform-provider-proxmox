resource "proxmox_virtual_environment_storage_cifs" "example" {
  id     = "example-cifs"
  nodes  = ["pve"]
  server = "10.0.0.20"
  share  = "proxmox"

  username = "cifs-user"
  password = "cifs-password"

  content                  = ["images"]
  domain                   = "WORKGROUP"
  subdirectory             = "terraform"
  preallocation            = "metadata"
  snapshot_as_volume_chain = true

  backups {
    max_protected_backups = 5
    keep_daily            = 7
  }
}


