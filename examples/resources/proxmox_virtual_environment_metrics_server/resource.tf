resource "proxmox_virtual_environment_metrics_server" "influxdb_server" {
  name   = "example_influxdb_server"
  server = "192.168.3.2"
  port   = 8089
  type   = "influxdb"
}

resource "proxmox_virtual_environment_metrics_server" "graphite_server" {
  name   = "example_graphite_server"
  server = "192.168.4.2"
  port   = 2003
  type   = "graphite"
}
