data "proxmox_apt_standard_repository" "example" {
  handle = "no-subscription"
  node   = "pve"
}

output "proxmox_apt_standard_repository" {
  value = data.proxmox_apt_standard_repository.example
}
