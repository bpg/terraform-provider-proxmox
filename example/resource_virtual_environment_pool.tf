resource "proxmox_virtual_environment_pool" "example" {
  comment  = "Managed by Terraform"
  pool_id = "terraform-provider-proxmox-example"
}

output "resource_proxmox_virtual_environment_pool_example_comment" {
  value = "${proxmox_virtual_environment_pool.example.comment}"
}

output "resource_proxmox_virtual_environment_pool_example_members" {
  value = "${proxmox_virtual_environment_pool.example.members}"
}

output "resource_proxmox_virtual_environment_pool_example_pool_id" {
  value = "${proxmox_virtual_environment_pool.example.id}"
}
