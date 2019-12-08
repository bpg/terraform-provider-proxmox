data "proxmox_virtual_environment_access_group" "example" {
  group_id = "${proxmox_virtual_environment_access_group.example.id}"
}

output "data_proxmox_virtual_environment_access_group_example_comments" {
  value = "${data.proxmox_virtual_environment_access_group.example.*.comment}"
}

output "data_proxmox_virtual_environment_access_group_example_members" {
  value = "${data.proxmox_virtual_environment_access_group.example.*.members}"
}
