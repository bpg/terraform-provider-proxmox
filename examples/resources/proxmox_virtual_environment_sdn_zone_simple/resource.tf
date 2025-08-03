resource "proxmox_virtual_environment_sdn_zone_simple" "example" {
  id    = "simple1"
  nodes = ["pve"]
  mtu   = 1500

  # Optional attributes
  dns         = "1.1.1.1"
  dns_zone    = "example.com"
  ipam        = "pve"
  reverse_dns = "1.1.1.1"
}
