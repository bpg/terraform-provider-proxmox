data "proxmox_virtual_environment_hardware_mappings" "example-pci" {
  check_node = "pve"
  type       = "pci"
}

data "proxmox_virtual_environment_hardware_mappings" "example-usb" {
  check_node = "pve"
  type       = "usb"
}

output "data_proxmox_virtual_environment_hardware_mappings_pci" {
  value = data.proxmox_virtual_environment_hardware_mappings.example-pci
}

output "data_proxmox_virtual_environment_hardware_mappings_usb" {
  value = data.proxmox_virtual_environment_hardware_mappings.example-usb
}
