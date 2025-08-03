data "proxmox_virtual_environment_sdn_zone_evpn" "example" {
  id = "evpn1"
}

output "data_proxmox_virtual_environment_sdn_zone_evpn" {
  value = {
    id                         = data.proxmox_virtual_environment_sdn_zone_evpn.example.id
    nodes                      = data.proxmox_virtual_environment_sdn_zone_evpn.example.nodes
    controller                 = data.proxmox_virtual_environment_sdn_zone_evpn.example.controller
    vrf_vxlan                  = data.proxmox_virtual_environment_sdn_zone_evpn.example.vrf_vxlan
    advertise_subnets          = data.proxmox_virtual_environment_sdn_zone_evpn.example.advertise_subnets
    disable_arp_nd_suppression = data.proxmox_virtual_environment_sdn_zone_evpn.example.disable_arp_nd_suppression
    exit_nodes                 = data.proxmox_virtual_environment_sdn_zone_evpn.example.exit_nodes
    exit_nodes_local_routing   = data.proxmox_virtual_environment_sdn_zone_evpn.example.exit_nodes_local_routing
    primary_exit_node          = data.proxmox_virtual_environment_sdn_zone_evpn.example.primary_exit_node
    rt_import                  = data.proxmox_virtual_environment_sdn_zone_evpn.example.rt_import
    mtu                        = data.proxmox_virtual_environment_sdn_zone_evpn.example.mtu
    dns                        = data.proxmox_virtual_environment_sdn_zone_evpn.example.dns
    dns_zone                   = data.proxmox_virtual_environment_sdn_zone_evpn.example.dns_zone
    ipam                       = data.proxmox_virtual_environment_sdn_zone_evpn.example.ipam
    reverse_dns                = data.proxmox_virtual_environment_sdn_zone_evpn.example.reverse_dns
  }
}
