data "proxmox_virtual_environment_sdn_zone_vxlan" "example" {
  id = "vxlan1"
}

output "data_proxmox_virtual_environment_sdn_zone_vxlan" {
  value = {
    id          = data.proxmox_virtual_environment_sdn_zone_vxlan.example.id
    nodes       = data.proxmox_virtual_environment_sdn_zone_vxlan.example.nodes
    peers       = data.proxmox_virtual_environment_sdn_zone_vxlan.example.peers
    mtu         = data.proxmox_virtual_environment_sdn_zone_vxlan.example.mtu
    dns         = data.proxmox_virtual_environment_sdn_zone_vxlan.example.dns
    dns_zone    = data.proxmox_virtual_environment_sdn_zone_vxlan.example.dns_zone
    ipam        = data.proxmox_virtual_environment_sdn_zone_vxlan.example.ipam
    reverse_dns = data.proxmox_virtual_environment_sdn_zone_vxlan.example.reverse_dns
  }
}
