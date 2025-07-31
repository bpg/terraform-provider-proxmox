# List all SDN zones
data "proxmox_virtual_environment_sdn_zones" "all" {}

# List only EVPN zones
data "proxmox_virtual_environment_sdn_zones" "evpn_only" {
  type = "evpn"
}

# List only Simple zones  
data "proxmox_virtual_environment_sdn_zones" "simple_only" {
  type = "simple"
}

output "data_proxmox_virtual_environment_sdn_zones_all" {
  value = {
    zones = data.proxmox_virtual_environment_sdn_zones.all.zones
  }
}

output "data_proxmox_virtual_environment_sdn_zones_filtered" {
  value = {
    evpn_zones   = data.proxmox_virtual_environment_sdn_zones.evpn_only.zones
    simple_zones = data.proxmox_virtual_environment_sdn_zones.simple_only.zones
  }
}
