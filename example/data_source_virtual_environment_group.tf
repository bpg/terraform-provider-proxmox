data "proxmox_virtual_environment_group" "example" {
  group_id = "${proxmox_virtual_environment_group.example.id}"
}

output "data_proxmox_virtual_environment_group_example_comment" {
  value = "${data.proxmox_virtual_environment_group.example.comment}"
}

output "data_proxmox_virtual_environment_group_example_members" {
  value = "${data.proxmox_virtual_environment_group.example.members}"
}
