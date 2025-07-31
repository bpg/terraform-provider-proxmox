resource "proxmox_virtual_environment_sdn_zone_vlan" "example" {
  id     = "vlan1"
  nodes  = ["pve"]
  bridge = "vmbr0"
  mtu    = 1500

  # Optional attributes
  dns         = "1.1.1.1"
  dns_zone    = "example.com"
  ipam        = "pve"
  reverse_dns = "1.1.1.1"
}
