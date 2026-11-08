resource "proxmox_network_linux_bond" "bond0" {
  node_name = "pve"
  name      = "bond0"

  # Physical interfaces that make up the bond.
  slaves = ["ens18", "ens19"]

  # LACP bond. The switch ports must be configured as an 802.3ad link aggregation group.
  bond_mode             = "802.3ad"
  bond_xmit_hash_policy = "layer3+4"

  comment = "bond0 comment"
}

resource "proxmox_network_linux_bridge" "vmbr10" {
  depends_on = [
    proxmox_network_linux_bond.bond0
  ]

  node_name = "pve"
  name      = "vmbr10"

  address = "10.10.10.10/24"

  ports = [
    "bond0"
  ]
}
