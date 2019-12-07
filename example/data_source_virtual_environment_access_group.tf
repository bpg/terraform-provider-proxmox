data "proxmox_virtual_environment_access_group" "example" {
  count = length(data.proxmox_virtual_environment_access_groups.example.ids)
  id    = element(data.proxmox_virtual_environment_access_groups.example.ids, count.index)
}

output "data_proxmox_virtual_environment_access_group_example_comments" {
  value = data.proxmox_virtual_environment_access_group.example.*.comment
}

output "data_proxmox_virtual_environment_access_group_example_members" {
  value = data.proxmox_virtual_environment_access_group.example.*.members
}
