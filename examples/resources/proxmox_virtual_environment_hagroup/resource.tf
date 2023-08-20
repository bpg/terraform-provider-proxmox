resource "proxmox_virtual_environment_hagroup" "example" {
  group   = "example"
  comment = "This is a comment."

  # Member nodes, with or without priority.
  nodes = {
    node1 = null
    node2 = 2
    node3 = 1
  }

  restricted  = true
  no_failback = false
}
