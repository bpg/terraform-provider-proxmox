resource "proxmox_virtual_environment_access_role" "example" {
  privileges = [
    "VM.Monitor",
  ]
  role_id = "terraform-provider-proxmox-example"
}

output "resource_proxmox_virtual_environment_access_role_example_privileges" {
  value = "${proxmox_virtual_environment_access_role.example.privileges}"
}

output "resource_proxmox_virtual_environment_access_role_example_role_id" {
  value = "${proxmox_virtual_environment_access_role.example.role_id}"
}
