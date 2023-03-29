data "proxmox_virtual_environment_cluster_firewall_alias" "example" {
  depends_on = [proxmox_virtual_environment_cluster_firewall_alias.example]

  name = proxmox_virtual_environment_cluster_firewall_alias.example.name
}

output "data_proxmox_virtual_environment_cluster_firewall_alias_example_cidr" {
  value = proxmox_virtual_environment_cluster_firewall_alias.example.cidr
}

data "proxmox_virtual_environment_cluster_firewall_aliases" "example" {
  depends_on = [proxmox_virtual_environment_cluster_firewall_alias.example]
}

output "data_proxmox_virtual_environment_cluster_firewall_aliases" {
  value = {
    "alias_names" = data.proxmox_virtual_environment_cluster_firewall_aliases.example.alias_names
  }
}



data "proxmox_virtual_environment_cluster_firewall_ipset" "example" {
  depends_on = [proxmox_virtual_environment_cluster_firewall_ipset.example]

  name = proxmox_virtual_environment_cluster_firewall_ipset.example.name
}

output "data_proxmox_virtual_environment_cluster_firewall_ipset_example_cidr" {
  value = {
    "cidrs" = data.proxmox_virtual_environment_cluster_firewall_ipset.example.cidr
  }
}

data "proxmox_virtual_environment_cluster_firewall_ipsets" "example" {
  depends_on = [proxmox_virtual_environment_cluster_firewall_ipset.example]
}

output "data_proxmox_virtual_environment_cluster_firewall_ipsets" {
  value = {
    "ipset_names" = data.proxmox_virtual_environment_cluster_firewall_ipsets.example.ipset_names
  }
}


data "proxmox_virtual_environment_cluster_firewall_security_group" "example" {
  depends_on = [proxmox_virtual_environment_cluster_firewall_security_group.example]

  name = proxmox_virtual_environment_cluster_firewall_security_group.example.name
}

output "data_proxmox_virtual_environment_cluster_firewall_security_group" {
  value = data.proxmox_virtual_environment_cluster_firewall_security_group.example
}

data "proxmox_virtual_environment_cluster_firewall_security_groups" "example" {
  depends_on = [proxmox_virtual_environment_cluster_firewall_security_group.example]
}

output "data_proxmox_virtual_environment_cluster_firewall_security_groups" {
  value = data.proxmox_virtual_environment_cluster_firewall_security_groups.example.security_group_names
}
