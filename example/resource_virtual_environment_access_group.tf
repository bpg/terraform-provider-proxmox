resource "proxmox_virtual_environment_access_group" "example" {
  comment  = "Created by Terraform"
  group_id = "terraform-provider-proxmox-example"
}

output "resource_proxmox_virtual_environment_access_group_example_comment" {
  value = "${proxmox_virtual_environment_access_group.example.comment}"
}

output "resource_proxmox_virtual_environment_access_group_example_id" {
  value = "${proxmox_virtual_environment_access_group.example.id}"
}

output "resource_proxmox_virtual_environment_access_group_example_members" {
  value = "${proxmox_virtual_environment_access_group.example.members}"
}
