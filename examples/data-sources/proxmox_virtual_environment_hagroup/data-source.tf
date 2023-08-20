// This will fetch the set of HA group identifiers...
data "proxmox_virtual_environment_hagroups" "all" {}

// ...which we will go through in order to fetch the whole data on each group.
data "proxmox_virtual_environment_hagroup" "example" {
  for_each = data.proxmox_virtual_environment_hagroups.all.group_ids
  group    = each.value
}

output "proxmox_virtual_environment_hagroups_full" {
  value = data.proxmox_virtual_environment_hagroup.example
}
