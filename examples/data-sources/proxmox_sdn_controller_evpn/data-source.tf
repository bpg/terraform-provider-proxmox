data "proxmox_sdn_controller_evpn" "example" {
  id = "evpn1"
}

output "data_proxmox_sdn_controller_evpn" {
  value = {
    id     = data.proxmox_sdn_controller_evpn.example.id
    asn    = data.proxmox_sdn_controller_evpn.example.asn
    fabric = data.proxmox_sdn_controller_evpn.example.fabric
    peers  = data.proxmox_sdn_controller_evpn.example.peers
  }
}
