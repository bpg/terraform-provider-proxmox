# Created first and destroyed last, so it applies the staged deletions on destroy.
# Nothing is staged when it is created, hence on_create = false.
resource "proxmox_network_applier" "finalizer" {
  node_name = "pve"
  on_create = false
}

# Staged only: nothing takes effect until the applier below reloads the node.
resource "proxmox_network_linux_bond" "bond0" {
  node_name = "pve"
  name      = "bond0"
  slaves    = ["eno1", "eno2"]
  bond_mode = "802.3ad"
  reload    = false

  depends_on = [proxmox_network_applier.finalizer]
}

resource "proxmox_network_linux_bridge" "vmbr0" {
  node_name = "pve"
  name      = "vmbr0"
  ports     = ["bond0"]
  address   = "10.0.1.10/24"
  reload    = false

  depends_on = [proxmox_network_applier.finalizer]
}

# Activates everything staged above in a single reload.
resource "proxmox_network_applier" "apply" {
  node_name = "pve"

  depends_on = [
    proxmox_network_linux_bond.bond0,
    proxmox_network_linux_bridge.vmbr0,
  ]
}
