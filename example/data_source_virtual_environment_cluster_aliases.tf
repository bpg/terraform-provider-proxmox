data "proxmox_virtual_environment_cluster_aliases" "example" {
    depends_on = ["proxmox_virtual_environment_cluster_alias.example"]
}

output "proxmox_virtual_environment_cluster_aliases" {
  value = "${map(
    "alias_ids", data.proxmox_virtual_environment_cluster_aliases.example.alias_ids,
  )}"
}
