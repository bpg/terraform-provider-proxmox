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

# SDN VNet - Basic vnet
resource "proxmox_virtual_environment_sdn_vnet" "example_vnet_1" {
  id   = "vnet1"
  zone = proxmox_virtual_environment_sdn_zone_simple.example_zone_1.id

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

# SDN VNet - VNet with alias and port isolation
resource "proxmox_virtual_environment_sdn_vnet" "example_vnet_2" {
  id            = "vnet2"
  zone          = proxmox_virtual_environment_sdn_zone_simple.example_zone_2.id
  alias         = "Example VNet 2"
  isolate_ports = true
  vlan_aware    = false

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

# Basic Subnet
resource "proxmox_virtual_environment_sdn_subnet" "basic_subnet" {
  cidr    = "192.168.1.0/24"
  vnet    = proxmox_virtual_environment_sdn_vnet.example_vnet_1.id
  gateway = "192.168.1.1"

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

# Subnet with DHCP Configuration
resource "proxmox_virtual_environment_sdn_subnet" "dhcp_subnet" {
  cidr            = "192.168.2.0/24"
  vnet            = proxmox_virtual_environment_sdn_vnet.example_vnet_2.id
  gateway         = "192.168.2.1"
  dhcp_dns_server = "192.168.2.53"
  dns_zone_prefix = "internal.example.com"
  snat            = true

  dhcp_range = {
    start_address = "192.168.2.10"
    end_address   = "192.168.2.100"
  }

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

# SDN Applier for all resources
resource "proxmox_virtual_environment_sdn_applier" "subnet_applier" {
  depends_on = [
    proxmox_virtual_environment_sdn_zone_simple.example_zone_1,
    proxmox_virtual_environment_sdn_zone_simple.example_zone_2,
    proxmox_virtual_environment_sdn_vnet.example_vnet_1,
    proxmox_virtual_environment_sdn_vnet.example_vnet_2,
    proxmox_virtual_environment_sdn_subnet.basic_subnet,
    proxmox_virtual_environment_sdn_subnet.dhcp_subnet
  ]
}

resource "proxmox_virtual_environment_sdn_applier" "finalizer" {
}


