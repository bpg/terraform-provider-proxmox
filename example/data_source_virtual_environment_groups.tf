data "proxmox_virtual_environment_groups" "example" {}

output "data_proxmox_virtual_environment_groups_example" {
  value = "${map(
    "comments", data.proxmox_virtual_environment_groups.example.comments,
    "group_ids", data.proxmox_virtual_environment_groups.example.group_ids,
  )}"
}
