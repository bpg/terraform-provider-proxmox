data "proxmox_node_config" "example" {
  node_name = "pve"
}

output "proxmox_node_config_description" {
  value = data.proxmox_node_config.example.description
}
