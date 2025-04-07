data "proxmox_virtual_environment_datastores" "example" {
  node_name = data.proxmox_virtual_environment_nodes.example.names[0]
}

output "data_proxmox_virtual_environment_datastores_example" {
  value = data.proxmox_virtual_environment_datastores.example
}
