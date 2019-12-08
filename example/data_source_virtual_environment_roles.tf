data "proxmox_virtual_environment_roles" "example" {
  depends_on = ["proxmox_virtual_environment_role.example"]
}

output "data_proxmox_virtual_environment_roles_example_privileges" {
  value = "${data.proxmox_virtual_environment_roles.example.privileges}"
}

output "data_proxmox_virtual_environment_roles_example_role_ids" {
  value = "${data.proxmox_virtual_environment_roles.example.role_ids}"
}

output "data_proxmox_virtual_environment_roles_example_special" {
  value = "${data.proxmox_virtual_environment_roles.example.special}"
}
