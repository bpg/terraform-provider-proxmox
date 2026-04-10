data "proxmox_apt_repository" "example" {
  file_path = "/etc/apt/sources.list"
  index     = 0
  node      = "pve"
}

output "proxmox_apt_repository" {
  value = data.proxmox_apt_repository.example
}
