data "proxmox_virtual_environment_role" "example" {
  role_id = "${proxmox_virtual_environment_role.example.id}"
}

output "data_proxmox_virtual_environment_role_example_privileges" {
  value = "${data.proxmox_virtual_environment_role.example.privileges}"
}

output "data_proxmox_virtual_environment_role_example_role_id" {
  value = "${data.proxmox_virtual_environment_role.example.id}"
}
