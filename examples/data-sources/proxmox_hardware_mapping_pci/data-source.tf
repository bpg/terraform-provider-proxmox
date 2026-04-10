data "proxmox_hardware_mapping_pci" "example" {
  name = "example"
}

output "data_proxmox_hardware_mapping_pci" {
  value = data.proxmox_hardware_mapping_pci.example
}
