data "proxmox_virtual_environment_pool" "example" {
  pool_id = proxmox_virtual_environment_pool.example.id
}

output "data_proxmox_virtual_environment_pool_example_comment" {
  value = data.proxmox_virtual_environment_pool.example.comment
}

output "data_proxmox_virtual_environment_pool_example_members" {
  value = data.proxmox_virtual_environment_pool.example.members
}

output "data_proxmox_virtual_environment_pool_example_pool_id" {
  value = data.proxmox_virtual_environment_pool.example.id
}
