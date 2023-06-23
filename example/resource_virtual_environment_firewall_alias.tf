resource "proxmox_virtual_environment_firewall_alias" "cluster_alias" {
  name    = "cluster-alias"
  cidr    = "192.168.0.0/23"
  comment = "Managed by Terraform"
}

resource "proxmox_virtual_environment_firewall_alias" "vm_alias" {
  depends_on = [proxmox_virtual_environment_vm.example]

  node_name = proxmox_virtual_environment_vm.example.node_name
  vm_id     = proxmox_virtual_environment_vm.example.vm_id

  name    = "vm-alias"
  cidr    = "192.168.1.0/23"
  comment = "Managed by Terraform"
}

resource "proxmox_virtual_environment_firewall_alias" "container_alias" {
  depends_on = [proxmox_virtual_environment_container.example]

  node_name    = proxmox_virtual_environment_container.example.node_name
  container_id = proxmox_virtual_environment_container.example.vm_id

  name    = "container-alias"
  cidr    = "192.168.2.0/23"
  comment = "Managed by Terraform"
}

output "resource_proxmox_virtual_environment_firewall_alias_cluster" {
  value = proxmox_virtual_environment_firewall_alias.cluster_alias
}

output "resource_proxmox_virtual_environment_firewall_alias_vm" {
  value = proxmox_virtual_environment_firewall_alias.vm_alias
}

output "resource_proxmox_virtual_environment_firewall_alias_container" {
  value = proxmox_virtual_environment_firewall_alias.container_alias
}
