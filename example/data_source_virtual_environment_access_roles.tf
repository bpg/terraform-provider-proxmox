data "proxmox_virtual_environment_access_roles" "example" {}

output "data_proxmox_virtual_environment_access_roles_example_privileges" {
  value = "${data.proxmox_virtual_environment_access_roles.example.privileges}"
}

output "data_proxmox_virtual_environment_access_roles_example_role_ids" {
  value = "${data.proxmox_virtual_environment_access_roles.example.role_ids}"
}

output "data_proxmox_virtual_environment_access_roles_example_special" {
  value = "${data.proxmox_virtual_environment_access_roles.example.special}"
}
