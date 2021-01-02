data "proxmox_virtual_environment_cluster_alias" "example" {
  name = "example"
}

output "data_proxmox_virtual_environment_cluster_alias_example_cidr" {
    value = proxmox_virtual_environment_cluster_alias.example.cidr
}
