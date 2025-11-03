data "proxmox_virtual_environment_sdn_zone_simple" "test_zone_1" {
  depends_on = [proxmox_virtual_environment_sdn_zone_simple.test_zone_1]
  id         = proxmox_virtual_environment_sdn_zone_simple.test_zone_1.id
}

data "proxmox_virtual_environment_sdn_zones" "simple" {
  depends_on = [
    proxmox_virtual_environment_sdn_zone_simple.test_zone_1,
    proxmox_virtual_environment_sdn_zone_simple.test_zone_2,
  ]
  type = "simple"
}

data "proxmox_virtual_environment_sdn_zones" "all" {
  depends_on = [
    proxmox_virtual_environment_sdn_zone_simple.test_zone_1,
    proxmox_virtual_environment_sdn_zone_simple.test_zone_2,
  ]
}

data "proxmox_virtual_environment_sdn_vnet" "test_vnet_1" {
  depends_on = [proxmox_virtual_environment_sdn_vnet.test_vnet_1]
  id = proxmox_virtual_environment_sdn_vnet.test_vnet_1.id
}

data "proxmox_virtual_environment_sdn_subnet" "test_subnet_dhcp" {
  depends_on = [proxmox_virtual_environment_sdn_subnet.test_subnet_dhcp]
  cidr = proxmox_virtual_environment_sdn_subnet.test_subnet_dhcp.cidr
  vnet = proxmox_virtual_environment_sdn_vnet.test_vnet_1.id
}

output "proxmox_virtual_environment_sdn_zone_simple" {
  value = data.proxmox_virtual_environment_sdn_zone_simple.test_zone_1
}

output "proxmox_virtual_environment_sdn_zones_simple" {
  value = data.proxmox_virtual_environment_sdn_zones.simple.zones
}

output "proxmox_virtual_environment_sdn_zones_all" {
  value = data.proxmox_virtual_environment_sdn_zones.all.zones
}

output "proxmox_virtual_environment_sdn_vnet" {
  value = data.proxmox_virtual_environment_sdn_vnet.test_vnet_1
}

output "proxmox_virtual_environment_sdn_subnet" {
  value = data.proxmox_virtual_environment_sdn_subnet.test_subnet_dhcp
}
