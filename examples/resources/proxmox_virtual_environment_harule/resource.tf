resource "proxmox_virtual_environment_harule" "example" {
  rule      = "prefer-node1"
  type      = "node-affinity"
  comment   = "Prefer node1 for HA placement"
  resources = ["vm:100", "vm:101"]

  nodes = {
    node1 = 2
    node2 = 1
    node3 = 1
  }

  strict = false
}
