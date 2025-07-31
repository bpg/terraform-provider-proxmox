data "proxmox_virtual_environment_sdn_zone_vlan" "example" {
  id = "vlan1"
}

output "data_proxmox_virtual_environment_sdn_zone_vlan" {
  value = {
    id          = data.proxmox_virtual_environment_sdn_zone_vlan.example.id
    nodes       = data.proxmox_virtual_environment_sdn_zone_vlan.example.nodes
    bridge      = data.proxmox_virtual_environment_sdn_zone_vlan.example.bridge
    mtu         = data.proxmox_virtual_environment_sdn_zone_vlan.example.mtu
    dns         = data.proxmox_virtual_environment_sdn_zone_vlan.example.dns
    dns_zone    = data.proxmox_virtual_environment_sdn_zone_vlan.example.dns_zone
    ipam        = data.proxmox_virtual_environment_sdn_zone_vlan.example.ipam
    reverse_dns = data.proxmox_virtual_environment_sdn_zone_vlan.example.reverse_dns
  }
}
