data "proxmox_hardware_mapping_usb" "example" {
  name = "example"
}

output "data_proxmox_hardware_mapping_usb" {
  value = data.proxmox_hardware_mapping_usb.example
}
