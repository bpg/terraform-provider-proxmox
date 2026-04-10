# Replication
resource "proxmox_replication" "example_replication_1" {
  id     = "100-0"
  target = "pve-02"
  type   = "local"

  # Optional attributes
  disable  = false
  comment  = "Replication to pve-02 every 30 min"
  schedule = "*/30"
}
