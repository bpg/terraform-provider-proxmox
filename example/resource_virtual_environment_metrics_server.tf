resource "proxmox_virtual_environment_metrics_server" "influxdb_server" {
  name   = "example_influxdb_server"
  server = "192.168.3.2"
  port   = 18089
  type   = "influxdb"

}

resource "proxmox_virtual_environment_metrics_server" "graphite_server" {
  name   = "example_graphite_server"
  server = "192.168.4.2"
  port   = 20033
  type   = "graphite"
}

resource "proxmox_virtual_environment_metrics_server" "graphite_server2" {
  name           = "example_graphite_server2"
  server         = "192.168.4.3"
  port           = 20033
  type           = "graphite"
  mtu            = 60000
  timeout        = 5
  graphite_proto = "udp"
}

resource "proxmox_virtual_environment_metrics_server" "opentelemetry_server" {
  name                = "example_opentelemetry_server"
  server              = "192.168.5.2"
  port                = 4318
  type                = "opentelemetry"
  opentelemetry_proto = "http"
  opentelemetry_path  = "/v1/metrics"
}

resource "proxmox_virtual_environment_metrics_server" "opentelemetry_server_https" {
  name                = "example_opentelemetry_server_https"
  server              = "192.168.5.3"
  port                = 4319
  type                = "opentelemetry"
  opentelemetry_proto = "https"
  opentelemetry_path  = "/v1/metrics"
}
