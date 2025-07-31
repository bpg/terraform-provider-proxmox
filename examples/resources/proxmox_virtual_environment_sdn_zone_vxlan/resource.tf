resource "proxmox_virtual_environment_sdn_zone_vxlan" "example" {
  id    = "vxlan1"
  nodes = ["pve"]
  peers = ["10.0.0.1", "10.0.0.2", "10.0.0.3"]
  mtu   = 1450

  # Optional attributes
  dns         = "1.1.1.1"
  dns_zone    = "example.com"
  ipam        = "pve"
  reverse_dns = "1.1.1.1"
}
