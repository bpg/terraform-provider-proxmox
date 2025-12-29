resource "proxmox_virtual_environment_storage_smb" "example" {
  id     = "example-smb"
  nodes  = ["pve"]
  server = "10.0.0.20"
  share  = "proxmox"

  username = "smb-user"
  password = "smb-password"

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

