data "proxmox_virtual_environment_sdn_zone_qinq" "example" {
  id = "qinq1"
}

output "data_proxmox_virtual_environment_sdn_zone_qinq" {
  value = {
    id                    = data.proxmox_virtual_environment_sdn_zone_qinq.example.id
    nodes                 = data.proxmox_virtual_environment_sdn_zone_qinq.example.nodes
    bridge                = data.proxmox_virtual_environment_sdn_zone_qinq.example.bridge
    service_vlan          = data.proxmox_virtual_environment_sdn_zone_qinq.example.service_vlan
    service_vlan_protocol = data.proxmox_virtual_environment_sdn_zone_qinq.example.service_vlan_protocol
    mtu                   = data.proxmox_virtual_environment_sdn_zone_qinq.example.mtu
    dns                   = data.proxmox_virtual_environment_sdn_zone_qinq.example.dns
    dns_zone              = data.proxmox_virtual_environment_sdn_zone_qinq.example.dns_zone
    ipam                  = data.proxmox_virtual_environment_sdn_zone_qinq.example.ipam
    reverse_dns           = data.proxmox_virtual_environment_sdn_zone_qinq.example.reverse_dns
  }
}
