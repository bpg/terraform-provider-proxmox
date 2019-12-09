resource "proxmox_virtual_environment_user" "example" {
  comment  = "Managed by Terraform"
  password = "Test1234!"
  user_id  = "terraform-provider-proxmox-example@pve"
}

output "resource_proxmox_virtual_environment_user_example_comment" {
  value = "${proxmox_virtual_environment_user.example.comment}"
}

output "resource_proxmox_virtual_environment_user_example_email" {
  value = "${proxmox_virtual_environment_user.example.email}"
}

output "resource_proxmox_virtual_environment_user_example_enabled" {
  value = "${proxmox_virtual_environment_user.example.enabled}"
}

output "resource_proxmox_virtual_environment_user_example_expiration_date" {
  value = "${proxmox_virtual_environment_user.example.expiration_date}"
}

output "resource_proxmox_virtual_environment_user_example_first_name" {
  value = "${proxmox_virtual_environment_user.example.first_name}"
}

output "resource_proxmox_virtual_environment_user_example_groups" {
  value = "${proxmox_virtual_environment_user.example.groups}"
}

output "resource_proxmox_virtual_environment_user_example_keys" {
  value = "${proxmox_virtual_environment_user.example.keys}"
}

output "resource_proxmox_virtual_environment_user_example_last_name" {
  value = "${proxmox_virtual_environment_user.example.last_name}"
}

output "resource_proxmox_virtual_environment_user_example_password" {
  value     = "${proxmox_virtual_environment_user.example.password}"
  sensitive = true
}

output "resource_proxmox_virtual_environment_user_example_user_id" {
  value = "${proxmox_virtual_environment_user.example.id}"
}
