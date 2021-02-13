data "proxmox_virtual_environment_users" "example" {
  depends_on = [proxmox_virtual_environment_user.example]
}

output "data_proxmox_virtual_environment_users_example_comments" {
  value = data.proxmox_virtual_environment_users.example.comments
}

output "data_proxmox_virtual_environment_users_example_emails" {
  value = data.proxmox_virtual_environment_users.example.emails
}

output "data_proxmox_virtual_environment_users_example_enabled" {
  value = data.proxmox_virtual_environment_users.example.enabled
}

output "data_proxmox_virtual_environment_users_example_expiration_dates" {
  value = data.proxmox_virtual_environment_users.example.expiration_dates
}

output "data_proxmox_virtual_environment_users_example_first_names" {
  value = data.proxmox_virtual_environment_users.example.first_names
}

output "data_proxmox_virtual_environment_users_example_groups" {
  value = data.proxmox_virtual_environment_users.example.groups
}

output "data_proxmox_virtual_environment_users_example_keys" {
  value = data.proxmox_virtual_environment_users.example.keys
}

output "data_proxmox_virtual_environment_users_example_last_names" {
  value = data.proxmox_virtual_environment_users.example.last_names
}

output "data_proxmox_virtual_environment_users_example_user_ids" {
  value = data.proxmox_virtual_environment_users.example.user_ids
}
