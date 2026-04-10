data "proxmox_acme_accounts" "example" {}

output "data_proxmox_acme_accounts" {
  value = data.proxmox_acme_accounts.example.accounts
}
