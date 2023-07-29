resource "proxmox_virtual_environment_network_linux_vlan" "vlan99" {
  node_name = "pve"
  name      = "eno0.99"

  comment = "VLAN 99"
}
