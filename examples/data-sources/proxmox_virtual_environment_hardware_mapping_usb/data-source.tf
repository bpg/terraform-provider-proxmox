data "proxmox_virtual_environment_hardware_mapping_usb" "example" {
  name = "example"
}

output "data_proxmox_virtual_environment_hardware_mapping_usb" {
  value = data.proxmox_virtual_environment_hardware_mapping_usb.example
}
