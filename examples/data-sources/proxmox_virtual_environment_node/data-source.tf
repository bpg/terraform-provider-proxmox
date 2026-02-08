data "proxmox_virtual_environment_node" "example" {
  node_name = "pve"
}

output "data_proxmox_virtual_environment_node" {
  value = {
    cpu_cores    = data.proxmox_virtual_environment_node.example.cpu_cores
    cpu_count    = data.proxmox_virtual_environment_node.example.cpu_count
    cpu_sockets  = data.proxmox_virtual_environment_node.example.cpu_sockets
    cpu_model    = data.proxmox_virtual_environment_node.example.cpu_model
    memory_total = data.proxmox_virtual_environment_node.example.memory_total
    uptime       = data.proxmox_virtual_environment_node.example.uptime
  }
}
