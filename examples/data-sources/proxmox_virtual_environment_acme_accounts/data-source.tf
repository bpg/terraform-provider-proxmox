data "proxmox_virtual_environment_acme_accounts" "example" {}

output "data_proxmox_virtual_environment_acme_accounts" {
  value = data.proxmox_virtual_environment_acme_accounts.example.accounts
}
