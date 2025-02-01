data "proxmox_virtual_environment_metrics_server" "example" {
  name = proxmox_virtual_environment_metrics_server.influxdb_server.name
}
