data "proxmox_version" "example" {}

output "data_proxmox_version" {
  value = {
    release       = data.proxmox_version.example.release
    repository_id = data.proxmox_version.example.repository_id
    version       = data.proxmox_version.example.version
  }
}
