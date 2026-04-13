data "proxmox_vm" "example" {
  depends_on = [proxmox_virtual_environment_vm.example]
  id         = proxmox_virtual_environment_vm.example.vm_id
  node_name  = data.proxmox_virtual_environment_nodes.example.names[0]
}

output "proxmox_vm_example" {
  value = data.proxmox_vm.example
}
