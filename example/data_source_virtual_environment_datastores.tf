data "proxmox_virtual_environment_datastores" "example" {
  node_name = data.proxmox_virtual_environment_nodes.example.names[0]
}

output "data_proxmox_virtual_environment_datastores_example_active" {
  value = data.proxmox_virtual_environment_datastores.example.active
}

output "data_proxmox_virtual_environment_datastores_example_content_types" {
  value = data.proxmox_virtual_environment_datastores.example.content_types
}

output "data_proxmox_virtual_environment_datastores_example_datastore_ids" {
  value = data.proxmox_virtual_environment_datastores.example.datastore_ids
}

output "data_proxmox_virtual_environment_datastores_example_enabled" {
  value = data.proxmox_virtual_environment_datastores.example.enabled
}

output "data_proxmox_virtual_environment_datastores_example_node_name" {
  value = data.proxmox_virtual_environment_datastores.example.node_name
}

output "data_proxmox_virtual_environment_datastores_example_shared" {
  value = data.proxmox_virtual_environment_datastores.example.shared
}

output "data_proxmox_virtual_environment_datastores_example_space_available" {
  value = data.proxmox_virtual_environment_datastores.example.space_available
}

output "data_proxmox_virtual_environment_datastores_example_space_total" {
  value = data.proxmox_virtual_environment_datastores.example.space_total
}

output "data_proxmox_virtual_environment_datastores_example_space_used" {
  value = data.proxmox_virtual_environment_datastores.example.space_used
}

output "data_proxmox_virtual_environment_datastores_example_types" {
  value = data.proxmox_virtual_environment_datastores.example.types
}
