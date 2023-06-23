resource "proxmox_virtual_environment_network_linux_vlan" "vlan99" {
  node_name = "pve"
  name      = "ens18.99"

  comment = "VLAN 99"
}

resource "proxmox_virtual_environment_network_linux_bridge" "vmbr99" {
  depends_on = [
    proxmox_virtual_environment_network_linux_vlan.vlan99
  ]

  node_name = "pve"
  name      = "vmbr99"

  address = "99.99.99.99/16"

  comment = "vmbr99 comment"

  ports = [
    "ens18.99"
  ]
}
