data "proxmox_virtual_environment_hardware_mapping_dir" "example" {
  name = "example"
}

output "data_proxmox_virtual_environment_hardware_mapping_dir" {
  value = data.proxmox_virtual_environment_hardware_mapping_dir.example
}
