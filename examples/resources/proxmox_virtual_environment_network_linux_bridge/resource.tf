resource "proxmox_virtual_environment_network_linux_bridge" "vmbr99" {
  depends_on = [
    proxmox_virtual_environment_network_linux_vlan.vlan99
  ]

  node_name = "pve"
  name      = "vmbr99"

  address = "99.99.99.99/16"

  comment = "vmbr99 comment"

  ports = [
    # Network (or VLAN) interfaces to attach to the bridge, specified by their interface name
    # (e.g. "ens18.99" for VLAN 99 on interface ens18).
    # For VLAN interfaces with custom names, use the interface name without the VLAN tag, e.g. "vlan_lab"
    "ens18.99"
  ]

  # VLAN-aware bridge: permit any tagged VLAN through. Segmentation typically
  # lives upstream (switch / firewall), not on the hypervisor bridge.
  vlan_aware = true
  vids       = "2-4094"
}

resource "proxmox_virtual_environment_network_linux_vlan" "vlan99" {
  node_name = "pve"
  name      = "ens18.99"

  ## or alternatively, use custom name:
  # name      = "vlan_lab"
  # interface = "eno0"
  # vlan      = 98
}
