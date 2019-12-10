data "proxmox_virtual_environment_user" "example" {
  user_id = "${proxmox_virtual_environment_user.example.id}"
}

output "data_proxmox_virtual_environment_user_example_acl" {
  value = "${data.proxmox_virtual_environment_user.example.acl}"
}

output "data_proxmox_virtual_environment_user_example_comment" {
  value = "${data.proxmox_virtual_environment_user.example.comment}"
}

output "data_proxmox_virtual_environment_user_example_email" {
  value = "${data.proxmox_virtual_environment_user.example.email}"
}

output "data_proxmox_virtual_environment_user_example_enabled" {
  value = "${data.proxmox_virtual_environment_user.example.enabled}"
}

output "data_proxmox_virtual_environment_user_example_expiration_date" {
  value = "${data.proxmox_virtual_environment_user.example.expiration_date}"
}

output "data_proxmox_virtual_environment_user_example_first_name" {
  value = "${data.proxmox_virtual_environment_user.example.first_name}"
}

output "data_proxmox_virtual_environment_user_example_groups" {
  value = "${data.proxmox_virtual_environment_user.example.groups}"
}

output "data_proxmox_virtual_environment_user_example_keys" {
  value = "${data.proxmox_virtual_environment_user.example.keys}"
}

output "data_proxmox_virtual_environment_user_example_last_name" {
  value = "${data.proxmox_virtual_environment_user.example.last_name}"
}

output "data_proxmox_virtual_environment_user_example_user_id" {
  value = "${data.proxmox_virtual_environment_user.example.id}"
}
