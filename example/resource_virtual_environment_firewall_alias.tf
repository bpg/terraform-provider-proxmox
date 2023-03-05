resource "proxmox_virtual_environment_firewall_alias" "example" {
  name    = "example"
  cidr    = "192.168.0.0/23"
  comment = "Managed by Terraform"
}

output "resource_proxmox_virtual_environment_firewall_alias_example_name" {
  value = proxmox_virtual_environment_firewall_alias.example.name
}

output "resource_proxmox_virtual_environment_firewall_alias_example_cidr" {
  value = proxmox_virtual_environment_firewall_alias.example.cidr
}
