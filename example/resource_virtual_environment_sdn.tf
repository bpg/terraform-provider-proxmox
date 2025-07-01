# --- SDN Zones ---

resource "proxmox_virtual_environment_sdn_zone" "zone_simple" {
  name  = "zoneS"
  type  = "simple"
  nodes = var.virtual_environment_node_name
  mtu   = 1496
}

resource "proxmox_virtual_environment_sdn_zone" "zone_vlan" {
  name   = "zoneVLAN"
  type   = "vlan"
  nodes  = var.virtual_environment_node_name
  mtu    = 1500
  bridge = "vmbr0"
}

# --- SDN Vnets ---

resource "proxmox_virtual_environment_sdn_vnet" "vnet_simple" {
  name          = "vnetM"
  zone          = proxmox_virtual_environment_sdn_zone.zone_simple.name
  alias         = "vnet in zoneM"
  isolate_ports = "0"
  vlanaware     = "0"
  zonetype      = proxmox_virtual_environment_sdn_zone.zone_simple.type
  depends_on = [ proxmox_virtual_environment_sdn_zone.zone_simple ]
}

resource "proxmox_virtual_environment_sdn_vnet" "vnet_vlan" {
  name     = "vnetVLAN"
  zone     = proxmox_virtual_environment_sdn_zone.zone_vlan.name
  alias    = "vnet in zoneVLAN"
  tag      = 1000
  zonetype = proxmox_virtual_environment_sdn_zone.zone_vlan.type
  depends_on = [ proxmox_virtual_environment_sdn_zone.zone_vlan ]
}

# --- SDN Subnets ---

resource "proxmox_virtual_environment_sdn_subnet" "subnet_simple" {
  subnet          = "10.10.0.0/24"
  vnet            = proxmox_virtual_environment_sdn_vnet.vnet_simple.name
  dhcp_dns_server = "10.10.0.53"
  dhcp_range = [
    {
      start_address = "10.10.0.10"
      end_address   = "10.10.0.100"
    }
  ]
  gateway    = "10.10.0.1"
  snat       = true
  depends_on = [ proxmox_virtual_environment_sdn_vnet.vnet_simple ]
}

resource "proxmox_virtual_environment_sdn_subnet" "subnet_simple2" {
  subnet          = "10.40.0.0/24"
  vnet            = proxmox_virtual_environment_sdn_vnet.vnet_simple.name
  dhcp_dns_server = "10.40.0.53"
  dhcp_range = [
    {
      start_address = "10.40.0.10"
      end_address   = "10.40.0.100"
    }
  ]
  gateway    = "10.40.0.1"
  snat       = true
  depends_on = [ proxmox_virtual_environment_sdn_vnet.vnet_simple ]
}

resource "proxmox_virtual_environment_sdn_subnet" "subnet_vlan" {
  subnet          = "10.20.0.0/24"
  vnet            = proxmox_virtual_environment_sdn_vnet.vnet_vlan.name
  dhcp_dns_server = "10.20.0.53"
  dhcp_range = [
    {
      start_address = "10.20.0.10"
      end_address   = "10.20.0.100"
    }
  ]
  gateway = "10.20.0.100"
  snat    = false
  depends_on = [ proxmox_virtual_environment_sdn_vnet.vnet_vlan ]
}

# --- Data Sources ---

data "proxmox_virtual_environment_sdn_zone" "zone_ex" {
  name = "zoneS"
  depends_on = [ proxmox_virtual_environment_sdn_zone.zone_simple ]
}

data "proxmox_virtual_environment_sdn_vnet" "vnet_ex" {
  name = "vnetM"
  depends_on = [ proxmox_virtual_environment_sdn_vnet.vnet_simple ]
}

data "proxmox_virtual_environment_sdn_subnet" "subnet_ex" {
  subnet = "zoneS-10.10.0.0-24"
  vnet   = data.proxmox_virtual_environment_sdn_vnet.vnet_ex.id
  depends_on = [ proxmox_virtual_environment_sdn_subnet.subnet_simple ]
}

# --- Outputs ---

output "sdn_zone" {
  value = data.proxmox_virtual_environment_sdn_zone.zone_ex
}

output "sdn_vnet" {
  value = data.proxmox_virtual_environment_sdn_vnet.vnet_ex
}

output "sdn_subnet" {
  value = data.proxmox_virtual_environment_sdn_subnet.subnet_ex
}
