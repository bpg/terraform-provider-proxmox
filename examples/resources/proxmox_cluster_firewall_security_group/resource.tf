resource "proxmox_cluster_firewall_security_group" "webserver" {
  name    = "webserver"
  comment = "Managed by Terraform"

  rule {
    type    = "in"
    action  = "ACCEPT"
    comment = "Allow HTTP"
    dest    = "192.168.1.5"
    dport   = "80"
    proto   = "tcp"
    log     = "info"
  }

  rule {
    type    = "in"
    action  = "ACCEPT"
    comment = "Allow HTTPS"
    dest    = "192.168.1.5"
    dport   = "443"
    proto   = "tcp"
    log     = "info"
  }
}
