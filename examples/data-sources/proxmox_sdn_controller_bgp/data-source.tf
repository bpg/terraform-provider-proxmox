data "proxmox_sdn_controller_bgp" "example" {
  id = "bgp1"
}

output "data_proxmox_sdn_controller_bgp" {
  value = {
    id        = data.proxmox_sdn_controller_bgp.example.id
    node_name = data.proxmox_sdn_controller_bgp.example.node_name
    asn       = data.proxmox_sdn_controller_bgp.example.asn
    peers     = data.proxmox_sdn_controller_bgp.example.peers
    ebgp      = data.proxmox_sdn_controller_bgp.example.ebgp
  }
}
