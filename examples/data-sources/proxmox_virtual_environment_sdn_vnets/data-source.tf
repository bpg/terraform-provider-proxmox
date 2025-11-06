# List all SDN VNets
data "proxmox_virtual_environment_sdn_vnets" "all" {}

output "data_proxmox_virtual_environment_sdn_vnets_all" {
  value = {
    vnets = data.proxmox_virtual_environment_sdn_vnets.all.vnets
  }
}
