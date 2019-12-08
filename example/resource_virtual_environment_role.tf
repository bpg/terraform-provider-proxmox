resource "proxmox_virtual_environment_role" "example" {
  privileges = [
    "VM.Console",
    "VM.Monitor",
  ]
  role_id = "terraform-provider-proxmox-example"
}

output "resource_proxmox_virtual_environment_role_example_privileges" {
  value = "${proxmox_virtual_environment_role.example.privileges}"
}

output "resource_proxmox_virtual_environment_role_example_role_id" {
  value = "${proxmox_virtual_environment_role.example.role_id}"
}
