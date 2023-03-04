data "proxmox_virtual_environment_firewall_alias" "example" {
  depends_on = [proxmox_virtual_environment_firewall_alias.example]

  name = proxmox_virtual_environment_firewall_alias.example.name
}

output "data_proxmox_virtual_environment_firewall_alias_example_cidr" {
  value = proxmox_virtual_environment_firewall_alias.example.cidr
}
