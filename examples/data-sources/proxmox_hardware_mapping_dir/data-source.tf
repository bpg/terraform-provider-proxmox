data "proxmox_hardware_mapping_dir" "example" {
  name = "example"
}

output "data_proxmox_hardware_mapping_dir" {
  value = data.proxmox_hardware_mapping_dir.example
}
