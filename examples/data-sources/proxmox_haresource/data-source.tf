// This will fetch the set of all HA resource identifiers...
data "proxmox_haresources" "all" {}

// ...which we will go through in order to fetch the whole record for each resource.
data "proxmox_haresource" "example" {
  for_each    = data.proxmox_haresources.all.resource_ids
  resource_id = each.value
}

output "proxmox_haresources_full" {
  value = data.proxmox_haresource.example
}
