resource "proxmox_virtual_environment_sdn_zone_evpn" "example" {
  id         = "evpn1"
  nodes      = ["pve"]
  controller = "evpn-controller1"
  vrf_vxlan  = 4000

  # Optional attributes
  advertise_subnets          = true
  disable_arp_nd_suppression = false
  exit_nodes                 = ["pve-exit1", "pve-exit2"]
  exit_nodes_local_routing   = true
  primary_exit_node          = "pve-exit1"
  rt_import                  = "65000:65000"
  mtu                        = 1450

  # Generic optional attributes
  dns         = "1.1.1.1"
  dns_zone    = "example.com"
  ipam        = "pve"
  reverse_dns = "1.1.1.1"
}
