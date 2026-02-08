data "proxmox_virtual_environment_nodes" "example" {}

output "data_proxmox_virtual_environment_nodes" {
  value = {
    names     = data.proxmox_virtual_environment_nodes.example.names
    cpu_count = data.proxmox_virtual_environment_nodes.example.cpu_count
    online    = data.proxmox_virtual_environment_nodes.example.online
  }
}
