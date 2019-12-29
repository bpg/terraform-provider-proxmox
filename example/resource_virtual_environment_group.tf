resource "proxmox_virtual_environment_group" "example" {
  acl {
    path    = "/vms/${proxmox_virtual_environment_vm.example.id}"
    role_id = "${proxmox_virtual_environment_role.example.id}"
  }

  comment  = "Managed by Terraform"
  group_id = "terraform-provider-proxmox-example"
}

output "resource_proxmox_virtual_environment_group_example_acl" {
  value = "${proxmox_virtual_environment_group.example.acl}"
}

output "resource_proxmox_virtual_environment_group_example_comment" {
  value = "${proxmox_virtual_environment_group.example.comment}"
}

output "resource_proxmox_virtual_environment_group_example_id" {
  value = "${proxmox_virtual_environment_group.example.id}"
}

output "resource_proxmox_virtual_environment_group_example_members" {
  value = "${proxmox_virtual_environment_group.example.members}"
}
