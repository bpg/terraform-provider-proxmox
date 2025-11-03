# SDN Zone (Simple) - First zone for applier demonstration
resource "proxmox_virtual_environment_sdn_zone_simple" "example_zone_1" {
  id    = "zone1"
  nodes = ["pve"]
  mtu   = 1500

  # Optional attributes
  dns         = "1.1.1.1"
  dns_zone    = "example.com"
  ipam        = "pve"
  reverse_dns = "1.1.1.1"

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

# SDN Zone (Simple) - Second zone for applier demonstration
resource "proxmox_virtual_environment_sdn_zone_simple" "example_zone_2" {
  id    = "zone2"
  nodes = ["pve"]
  mtu   = 1500

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

# SDN Applier - Applies SDN configuration changes
resource "proxmox_virtual_environment_sdn_applier" "example_applier" {
  lifecycle {
    replace_triggered_by = [
      proxmox_virtual_environment_sdn_zone_simple.example_zone_1,
      proxmox_virtual_environment_sdn_zone_simple.example_zone_2,
    ]
  }

  depends_on = [
    proxmox_virtual_environment_sdn_zone_simple.example_zone_1,
    proxmox_virtual_environment_sdn_zone_simple.example_zone_2,
  ]
}

resource "proxmox_virtual_environment_sdn_applier" "finalizer" {
}
