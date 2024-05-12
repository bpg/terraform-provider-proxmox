resource "proxmox_virtual_environment_apt_repository" "example" {
  enabled   = true
  file_path = "/etc/apt/sources.list"
  index     = 0
  node      = "pve"
}
