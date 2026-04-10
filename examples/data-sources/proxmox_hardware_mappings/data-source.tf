data "proxmox_hardware_mappings" "example-dir" {
  check_node = "pve"
  type       = "dir"
}

data "proxmox_hardware_mappings" "example-pci" {
  check_node = "pve"
  type       = "pci"
}

data "proxmox_hardware_mappings" "example-usb" {
  check_node = "pve"
  type       = "usb"
}

output "data_proxmox_hardware_mappings_pci" {
  value = data.proxmox_hardware_mappings.example-pci
}

output "data_proxmox_hardware_mappings_usb" {
  value = data.proxmox_hardware_mappings.example-usb
}
