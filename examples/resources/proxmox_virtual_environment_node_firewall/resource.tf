resource "proxmox_virtual_environment_node_firewall" "node-pve1" {
  node_name = "pve1"
  enabled   = false
}

resource "proxmox_virtual_environment_node_firewall" "pve2" {
  node_name           = "pve2"
  enabled             = true
  log_level_in        = "alert"
  log_level_out       = "alert"
  log_level_forward   = "alert"
  ndp                 = true
  nftables            = true
  nosmurfs            = true
  smurf_log_level     = "alert"
  tcp_flags_log_level = "alert"
}
