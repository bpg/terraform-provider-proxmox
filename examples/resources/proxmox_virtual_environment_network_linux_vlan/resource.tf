# using VLAN tag
resource "proxmox_virtual_environment_network_linux_vlan" "vlan99" {
  node_name = "pve"
  name      = "eno0.99"

  comment = "VLAN 99"
}

# using custom network interface name
resource "proxmox_virtual_environment_network_linux_vlan" "vlan98" {
  node_name = "pve"
  name      = "vlan_lab"

  interface = "eno0"
  vlan      = 98
  comment   = "VLAN 98"
}
