data "proxmox_virtual_environment_sdn_zone_simple" "example" {
  id = "simple1"
}

output "data_proxmox_virtual_environment_sdn_zone_simple" {
  value = {
    id          = data.proxmox_virtual_environment_sdn_zone_simple.example.id
    nodes       = data.proxmox_virtual_environment_sdn_zone_simple.example.nodes
    mtu         = data.proxmox_virtual_environment_sdn_zone_simple.example.mtu
    dns         = data.proxmox_virtual_environment_sdn_zone_simple.example.dns
    dns_zone    = data.proxmox_virtual_environment_sdn_zone_simple.example.dns_zone
    ipam        = data.proxmox_virtual_environment_sdn_zone_simple.example.ipam
    reverse_dns = data.proxmox_virtual_environment_sdn_zone_simple.example.reverse_dns
  }
}
