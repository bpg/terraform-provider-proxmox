data "proxmox_virtual_environment_version" "ve_version" {}

output "proxmox_virtual_environment_version" {
  value = "${map(
    "keyboard", data.proxmox_virtual_environment_version.ve_version.keyboard,
    "release", data.proxmox_virtual_environment_version.ve_version.release,
    "repository_id", data.proxmox_virtual_environment_version.ve_version.repository_id,
    "version", data.proxmox_virtual_environment_version.ve_version.version,
  )}"
}
