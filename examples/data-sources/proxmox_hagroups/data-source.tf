data "proxmox_hagroups" "example" {}

output "data_proxmox_hagroups" {
  value = data.proxmox_hagroups.example.group_ids
}
