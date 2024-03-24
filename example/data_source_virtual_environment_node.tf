data "proxmox_virtual_environment_node" "example" {
  node_name = data.proxmox_virtual_environment_nodes.example.names[0]
}

output "data_proxmox_virtual_environment_node_example_cpu_count" {
  value = data.proxmox_virtual_environment_node.example.cpu_count
}

output "data_proxmox_virtual_environment_node_example_cpu_sockets" {
  value = data.proxmox_virtual_environment_node.example.cpu_sockets
}

output "data_proxmox_virtual_environment_node_example_cpu_model" {
  value = data.proxmox_virtual_environment_node.example.cpu_model
}

output "data_proxmox_virtual_environment_node_example_memory_available" {
  value = data.proxmox_virtual_environment_node.example.memory_available
}

output "data_proxmox_virtual_environment_node_example_memory_used" {
  value = data.proxmox_virtual_environment_node.example.memory_used
}

output "data_proxmox_virtual_environment_node_example_memory_total" {
  value = data.proxmox_virtual_environment_node.example.memory_total
}

output "data_proxmox_virtual_environment_node_example_node_name" {
  value = data.proxmox_virtual_environment_node.example.node_name
}

output "data_proxmox_virtual_environment_node_example_uptime" {
  value = data.proxmox_virtual_environment_node.example.uptime
}
