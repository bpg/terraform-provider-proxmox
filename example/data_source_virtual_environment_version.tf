data "proxmox_virtual_environment_version" "example" {}

output "data_proxmox_virtual_environment_version" {
  value = "${map(
    "keyboard_layout", data.proxmox_virtual_environment_version.example.keyboard_layout,
    "release", data.proxmox_virtual_environment_version.example.release,
    "repository_id", data.proxmox_virtual_environment_version.example.repository_id,
    "version", data.proxmox_virtual_environment_version.example.version,
  )}"
}
