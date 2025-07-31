resource "proxmox_virtual_environment_sdn_zone_qinq" "example" {
  id                    = "qinq1"
  nodes                 = ["pve"]
  bridge                = "vmbr0"
  service_vlan          = 100
  service_vlan_protocol = "802.1ad"
  mtu                   = 1496

  # Optional attributes
  dns         = "1.1.1.1"
  dns_zone    = "example.com"
  ipam        = "pve"
  reverse_dns = "1.1.1.1"
}
