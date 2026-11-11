# SDN Controller (BGP) - Example configuration for a BGP controller on a single node
resource "proxmox_sdn_controller_bgp" "example_controller_bgp" {
  id        = "bgp1"
  node_name = "pve"
  asn       = 65000
  peers     = ["10.0.0.1", "10.0.0.2"]

  ebgp                        = true
  ebgp_multihop               = 3
  bgp_multipath_as_path_relax = true

  depends_on = [
    proxmox_sdn_applier.finalizer
  ]
}

resource "proxmox_sdn_applier" "controller_applier" {
  depends_on = [
    proxmox_sdn_controller_bgp.example_controller_bgp
  ]
}

resource "proxmox_sdn_applier" "finalizer" {
}
