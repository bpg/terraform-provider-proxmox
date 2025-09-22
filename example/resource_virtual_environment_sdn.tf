resource "proxmox_virtual_environment_sdn_zone_simple" "test_zone_1" {
  id = "tZone1"
  nodes= data.proxmox_virtual_environment_nodes.example.names
  mtu = 1496

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

resource "proxmox_virtual_environment_sdn_zone_simple" "test_zone_2" {
  id = "tZone2"
  nodes= data.proxmox_virtual_environment_nodes.example.names
  mtu = 1496
  
  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

resource "proxmox_virtual_environment_sdn_vnet" "test_vnet_1" {
  id            = "tstVNet1"
  zone          = proxmox_virtual_environment_sdn_zone_simple.test_zone_1.id
  alias         = "Test Virtual Network 1"
  isolate_ports = true
  vlan_aware    = false

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

resource "proxmox_virtual_environment_sdn_vnet" "test_vnet_2" {
  id         = "tstVNet2"
  zone       = proxmox_virtual_environment_sdn_zone_simple.test_zone_2.id
  alias      = "Test Virtual Network 2"
  vlan_aware = true

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

resource "proxmox_virtual_environment_sdn_subnet" "test_subnet_dhcp" {
  cidr            = "10.100.0.0/24"
  vnet            = proxmox_virtual_environment_sdn_vnet.test_vnet_1.id
  gateway         = "10.100.0.1"
  dhcp_dns_server = "10.100.0.53"
  snat            = true

  dhcp_range = {
    start_address = "10.100.0.100"
    end_address   = "10.100.0.200"
  }

  depends_on = [
    proxmox_virtual_environment_sdn_applier.finalizer
  ]
}

resource "proxmox_virtual_environment_sdn_applier" "applier" {
   lifecycle {
    replace_triggered_by = [
      proxmox_virtual_environment_sdn_zone_simple.test_zone_1,
      proxmox_virtual_environment_sdn_zone_simple.test_zone_2,
      proxmox_virtual_environment_sdn_vnet.test_vnet_1,
      proxmox_virtual_environment_sdn_vnet.test_vnet_2,
      proxmox_virtual_environment_sdn_subnet.test_subnet_dhcp,
      ]
  }

  depends_on = [
    proxmox_virtual_environment_sdn_zone_simple.test_zone_1,
    proxmox_virtual_environment_sdn_zone_simple.test_zone_2,
    proxmox_virtual_environment_sdn_vnet.test_vnet_1,
    proxmox_virtual_environment_sdn_vnet.test_vnet_2,
    proxmox_virtual_environment_sdn_subnet.test_subnet_dhcp,
  ]
}

resource "proxmox_virtual_environment_sdn_applier" "finalizer" {
}
