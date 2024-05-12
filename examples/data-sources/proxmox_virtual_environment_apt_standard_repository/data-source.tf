data "proxmox_virtual_environment_apt_standard_repository" "example" {
  handle = "no-subscription"
  node   = "pve"
}

output "proxmox_virtual_environment_apt_standard_repository" {
  value = data.proxmox_virtual_environment_apt_standard_repository.example
}
