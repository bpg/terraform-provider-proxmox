resource "proxmox_virtual_environment_storage_pbs" "example" {
  id        = "example-pbs"
  nodes     = ["pve"]
  server    = "pbs.example.local"
  datastore = "backup"

  username    = "pbs-user"
  password    = "pbs-password"
  fingerprint = "AA:BB:CC:DD:EE:FF"

  content = ["backup"]

  generate_encryption_key = true
}

