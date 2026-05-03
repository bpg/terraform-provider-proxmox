# SDN Controller (EVPN) - Example configuration for SDN Controller with fabric
resource "proxmox_sdn_fabric_openfabric" "example_fabric" {
  id        = "main"
  ip_prefix = "10.0.0.0/16"
  depends_on = [
    proxmox_sdn_applier.finalizer
  ]
}

resource "proxmox_sdn_controller_evpn" "example_controller_evpn" {
  id     = "evpn1"
  asn    = 65000
  fabric = proxmox_sdn_fabric_openfabric.example_fabric.id
  depends_on = [
    proxmox_sdn_applier.finalizer
  ]
}

resource "proxmox_sdn_applier" "controller_applier" {
  depends_on = [
    proxmox_sdn_controller_evpn.example_controller_evpn
  ]
}

resource "proxmox_sdn_applier" "finalizer" {
}
