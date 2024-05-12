resource "proxmox_virtual_environment_apt_standard_repository" "example" {
  handle = "no-subscription"
  node   = "pve"
}

resource "proxmox_virtual_environment_apt_repository" "example" {
  enabled   = true
  file_path = proxmox_virtual_environment_apt_standard_repository.example.file_path
  index     = proxmox_virtual_environment_apt_standard_repository.example.index
  node      = proxmox_virtual_environment_apt_standard_repository.example.node
}
