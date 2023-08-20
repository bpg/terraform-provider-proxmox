data "proxmox_virtual_environment_hagroups" "example" {}

output "data_proxmox_virtual_environment_hagroups" {
  value = data.proxmox_virtual_environment_hagroups.example.group_ids
}
