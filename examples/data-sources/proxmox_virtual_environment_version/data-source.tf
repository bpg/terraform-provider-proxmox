data "proxmox_virtual_environment_version" "example" {}

output "data_proxmox_virtual_environment_version" {
  value = {
    release       = data.proxmox_virtual_environment_version.example.release
    repository_id = data.proxmox_virtual_environment_version.example.repository_id
    version       = data.proxmox_virtual_environment_version.example.version
  }
}
