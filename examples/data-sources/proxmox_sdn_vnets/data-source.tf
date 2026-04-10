# List all SDN VNets
data "proxmox_sdn_vnets" "all" {}

output "data_proxmox_sdn_vnets_all" {
  value = {
    vnets = data.proxmox_sdn_vnets.all.vnets
  }
}
