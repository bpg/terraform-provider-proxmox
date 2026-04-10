data "proxmox_acme_plugins" "example" {}

output "data_proxmox_acme_plugins" {
  value = data.proxmox_acme_plugins.example.plugins
}
