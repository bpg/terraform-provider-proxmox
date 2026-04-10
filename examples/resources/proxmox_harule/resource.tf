# Node Affinity Rule: assign VMs to preferred nodes with priorities.
# Non-strict rules allow failover to other nodes; strict rules do not.
resource "proxmox_harule" "prefer_node1" {
  rule      = "prefer-node1"
  type      = "node-affinity"
  comment   = "Prefer node1 for these VMs"
  resources = ["vm:100", "vm:101"]

  nodes = {
    node1 = 2 # Higher priority
    node2 = 1
    node3 = 1
  }

  strict = false
}

# Resource Affinity Rule (Positive): keep resources together on the same node.
resource "proxmox_harule" "keep_together" {
  rule      = "db-cluster-together"
  type      = "resource-affinity"
  comment   = "Keep database replicas on the same node"
  resources = ["vm:200", "vm:201"]
  affinity  = "positive"
}

# Resource Affinity Rule (Negative / Anti-Affinity): keep resources on
# separate nodes for high availability.
resource "proxmox_harule" "keep_apart" {
  rule      = "db-cluster-apart"
  type      = "resource-affinity"
  comment   = "Spread database replicas across nodes"
  resources = ["vm:200", "vm:201", "vm:202"]
  affinity  = "negative"
}
