data "proxmox_virtual_environment_metrics_server" "example" {}

output "data_proxmox_virtual_environment_metrics_server" {
  value = {
    server = data.proxmox_virtual_environment_metrics_server.example.server
    port   = data.proxmox_virtual_environment_metrics_server.example.port
  }
}
