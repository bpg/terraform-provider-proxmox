// This will fetch the set of all HA resource identifiers...
data "proxmox_virtual_environment_haresources" "all" {}

// ...which we will go through in order to fetch the whole record for each resource.
data "proxmox_virtual_environment_haresource" "example" {
  for_each    = data.proxmox_virtual_environment_haresources.all.resource_ids
  resource_id = each.value
}

output "proxmox_virtual_environment_haresources_full" {
  value = data.proxmox_virtual_environment_haresource.example
}
