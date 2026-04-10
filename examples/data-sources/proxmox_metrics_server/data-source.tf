data "proxmox_metrics_server" "example" {
  name = "example_influxdb"
}

output "data_proxmox_metrics_server" {
  value = {
    server = data.proxmox_metrics_server.example.server
    port   = data.proxmox_metrics_server.example.port
  }
}
