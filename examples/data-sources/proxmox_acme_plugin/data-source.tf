data "proxmox_acme_plugin" "example" {
  plugin = "standalone"
}

output "data_proxmox_acme_plugin" {
  value = data.proxmox_acme_plugin.example
}
