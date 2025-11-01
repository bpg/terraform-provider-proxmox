# Example: Basic ACME certificate with HTTP-01 challenge (standalone)
resource "proxmox_virtual_environment_acme_account" "example" {
  name      = "production"
  contact   = "admin@example.com"
  directory = "https://acme-v02.api.letsencrypt.org/directory"
  tos       = "https://letsencrypt.org/documents/LE-SA-v1.3-September-21-2022.pdf"
}

resource "proxmox_virtual_environment_acme_certificate" "http_example" {
  node_name = "pve-node-01"
  account   = proxmox_virtual_environment_acme_account.example.name

  domains = [
    {
      domain = "pve.example.com"
      # No plugin specified = HTTP-01 challenge
    }
  ]
}

# Example: ACME certificate with DNS-01 challenge using Cloudflare
resource "proxmox_virtual_environment_acme_dns_plugin" "cloudflare" {
  plugin = "cloudflare"
  api    = "cf"

  # Wait 2 minutes for DNS propagation
  validation_delay = 120

  data = {
    CF_Account_ID = "your-cloudflare-account-id"
    CF_Token      = "your-cloudflare-api-token"
    CF_Zone_ID    = "your-cloudflare-zone-id"
  }
}

resource "proxmox_virtual_environment_acme_certificate" "dns_example" {
  node_name = "pve-node-01"
  account   = proxmox_virtual_environment_acme_account.example.name

  domains = [
    {
      domain = "pve.example.com"
      plugin = proxmox_virtual_environment_acme_dns_plugin.cloudflare.plugin
    }
  ]

  depends_on = [
    proxmox_virtual_environment_acme_account.example,
    proxmox_virtual_environment_acme_dns_plugin.cloudflare
  ]
}

# Example: Force certificate renewal
resource "proxmox_virtual_environment_acme_certificate" "force_renew" {
  node_name = "pve-node-01"
  account   = proxmox_virtual_environment_acme_account.example.name
  force     = true

  domains = [
    {
      domain = "pve.example.com"
      plugin = proxmox_virtual_environment_acme_dns_plugin.cloudflare.plugin
    }
  ]

  depends_on = [
    proxmox_virtual_environment_acme_account.example,
    proxmox_virtual_environment_acme_dns_plugin.cloudflare
  ]
}

