data "proxmox_virtual_environment_access_groups" "example" {}

output "data_proxmox_virtual_environment_access_groups" {
  value = "${map(
    "comments", data.proxmox_virtual_environment_access_groups.example.comments,
    "ids", data.proxmox_virtual_environment_access_groups.example.ids,
  )}"
}
