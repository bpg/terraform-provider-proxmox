data "proxmox_metrics_server" "example" {
  name = proxmox_metrics_server.influxdb_server.name
}
