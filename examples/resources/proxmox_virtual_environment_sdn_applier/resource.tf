resource "proxmox_virtual_environment_sdn_zone_simple" "test_zone_1" {
  id    = "tZone1"
  nodes = [data.proxmox_virtual_environment_nodes.example.names]
  mtu   = 1496

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

resource "proxmox_virtual_environment_sdn_zone_simple" "test_zone_2" {
  id    = "tZone2"
  nodes = [data.proxmox_virtual_environment_nodes.example.names]
  mtu   = 1496

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

resource "proxmox_virtual_environment_sdn_applier" "applier" {
  lifecycle {
    replace_triggered_by = [
      proxmox_virtual_environment_sdn_zone_simple.test_zone_1,
      proxmox_virtual_environment_sdn_zone_simple.test_zone_2,
    ]
  }

  depends_on = [
    proxmox_virtual_environment_sdn_zone_simple.test_zone_1,
    proxmox_virtual_environment_sdn_zone_simple.test_zone_2,
  ]
}

resource "proxmox_virtual_environment_sdn_applier" "finalizer" {
}
