data "proxmox_virtual_environment_acme_plugins" "example" {}

output "data_proxmox_virtual_environment_acme_plugins" {
  value = data.proxmox_virtual_environment_acme_plugins.example.plugins
}
