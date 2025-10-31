data "proxmox_virtual_environment_sdn_vnet" "example" {
  id = "vnet1"
}

output "data_proxmox_virtual_environment_sdn_vnet" {
  value = {
    id            = data.proxmox_virtual_environment_sdn_vnet.example.id
    zone          = data.proxmox_virtual_environment_sdn_vnet.example.zone
    alias         = data.proxmox_virtual_environment_sdn_vnet.example.alias
    isolate_ports = data.proxmox_virtual_environment_sdn_vnet.example.isolate_ports
    tag           = data.proxmox_virtual_environment_sdn_vnet.example.tag
    vlan_aware    = data.proxmox_virtual_environment_sdn_vnet.example.vlan_aware
  }
}
