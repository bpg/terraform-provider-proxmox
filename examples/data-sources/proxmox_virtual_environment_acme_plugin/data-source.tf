data "proxmox_virtual_environment_acme_plugin" "example" {
  plugin = "standalone"
}

output "data_proxmox_virtual_environment_acme_plugin" {
  value = data.proxmox_virtual_environment_acme_plugin.example
}
