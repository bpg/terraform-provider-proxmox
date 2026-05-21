// Cluster-wide Ceph status (hits GET /cluster/ceph/status).
data "proxmox_ceph_status" "cluster" {}

// Node-scoped Ceph status (hits GET /nodes/{node}/ceph/status).
data "proxmox_ceph_status" "node" {
  node_name = "pve"
}

output "ceph_health" {
  value = data.proxmox_ceph_status.cluster.health_status
}

output "ceph_fsid" {
  value = data.proxmox_ceph_status.cluster.fsid
}

output "ceph_quorum" {
  value = data.proxmox_ceph_status.cluster.quorum_names
}
