# SDN Zone (Simple) - Basic zone for simple vnets
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

# SDN Zone (Simple) - Second zone for demonstration
resource "proxmox_virtual_environment_sdn_zone_simple" "example_zone_2" {
  id    = "zone2"
  nodes = ["pve"]
  mtu   = 1500

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

# Basic VNet (Simple)
resource "proxmox_virtual_environment_sdn_vnet" "basic_vnet" {
  id   = "vnet1"
  zone = proxmox_virtual_environment_sdn_zone_simple.example_zone_1.id

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

# VNet with Alias and Port Isolation
resource "proxmox_virtual_environment_sdn_vnet" "isolated_vnet" {
  id            = "vnet2"
  zone          = proxmox_virtual_environment_sdn_zone_simple.example_zone_2.id
  alias         = "Isolated VNet"
  isolate_ports = true
  vlan_aware    = false

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

# SDN Applier for all resources
resource "proxmox_virtual_environment_sdn_applier" "vnet_applier" {
  depends_on = [
    proxmox_virtual_environment_sdn_zone_simple.example_zone_1,
    proxmox_virtual_environment_sdn_zone_simple.example_zone_2,
    proxmox_virtual_environment_sdn_vnet.basic_vnet,
    proxmox_virtual_environment_sdn_vnet.isolated_vnet
  ]
}

resource "proxmox_virtual_environment_sdn_applier" "finalizer" {
}
