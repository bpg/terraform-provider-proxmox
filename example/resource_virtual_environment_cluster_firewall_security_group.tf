resource "proxmox_virtual_environment_cluster_firewall_security_group" "example" {
  name    = "example-sg"
  comment = "Managed by Terraform"

  rule {
    type    = "in"
    action  = "ACCEPT"
    comment = "Allow FTP"
    dest    = "192.168.1.5"
    dport   = "21"
    proto   = "tcp"
    log     = "info"
  }

  rule {
    type    = "in"
    action  = "DROP"
    comment = "Drop SSH"
    dest    = "192.168.1.5"
    dport   = "22"
    proto   = "udp"
    log     = "info"
  }
}

output "resource_proxmox_virtual_environment_cluster_firewall_security_group_example" {
  value = proxmox_virtual_environment_cluster_firewall_security_group.example
}
