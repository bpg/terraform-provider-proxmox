data "proxmox_virtual_environment_container" "example" {
  depends_on = [proxmox_virtual_environment_container.example]
  vm_id      = proxmox_virtual_environment_container.example.vm_id
  node_name  = data.proxmox_virtual_environment_nodes.example.names[0]
}

output "proxmox_virtual_environment_container_example" {
  value = data.proxmox_virtual_environment_container.example
}
