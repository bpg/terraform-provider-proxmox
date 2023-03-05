data "proxmox_virtual_environment_firewall_aliases" "example" {
  depends_on = [proxmox_virtual_environment_firewall_alias.example]
}

output "data_proxmox_virtual_environment_firewall_aliases" {
  value = {
    "alias_names" = data.proxmox_virtual_environment_firewall_aliases.example.alias_names
  }
}
