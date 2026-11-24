data "proxmox_cluster_firewall_security_groups" "all" {}

output "security_group_names" {
  value = data.proxmox_cluster_firewall_security_groups.all.groups[*].name
}
