// This will fetch all ACME accounts...
data "proxmox_acme_accounts" "all" {}

// ...which we will go through in order to fetch the whole data on each account.
data "proxmox_acme_account" "example" {
  for_each = data.proxmox_acme_accounts.all.accounts
  name     = each.value
}

output "data_proxmox_acme_account" {
  value = data.proxmox_acme_account.example
}
