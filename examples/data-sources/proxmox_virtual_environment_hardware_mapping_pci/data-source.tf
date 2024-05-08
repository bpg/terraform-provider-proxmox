data "proxmox_virtual_environment_hardware_mapping_pci" "example" {
  name = "example"
}

output "data_proxmox_virtual_environment_hardware_mapping_pci" {
  value = data.proxmox_virtual_environment_hardware_mapping_pci.example
}
