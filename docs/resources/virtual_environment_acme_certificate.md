---
layout: page
title: proxmox_virtual_environment_acme_certificate
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_acme_certificate

Manages ACME SSL certificates for Proxmox VE nodes. This resource orders and renews certificates from an ACME Certificate Authority for a specific node.

## Example Usage

### Basic ACME Certificate with HTTP-01 Challenge

```terraform
# First, create an ACME account
resource "proxmox_virtual_environment_acme_account" "example" {
  name      = "production"
  contact   = "admin@example.com"
  directory = "https://acme-v02.api.letsencrypt.org/directory"
  tos       = "https://letsencrypt.org/documents/LE-SA-v1.3-September-21-2022.pdf"
}

# Order a certificate for the node
resource "proxmox_virtual_environment_acme_certificate" "example" {
  node_name = "pve"
  account   = proxmox_virtual_environment_acme_account.example.name

  domains = [
    {
      domain = "pve.example.com"
      # No plugin specified implies HTTP-01 challenge
    }
  ]
}
```

### ACME Certificate with DNS-01 Challenge

```terraform
# Create an ACME account
resource "proxmox_virtual_environment_acme_account" "example" {
  name      = "production"
  contact   = "admin@example.com"
  directory = "https://acme-v02.api.letsencrypt.org/directory"
  tos       = "https://letsencrypt.org/documents/LE-SA-v1.3-September-21-2022.pdf"
}

# Configure a DNS plugin (Desec example)
resource "proxmox_virtual_environment_acme_dns_plugin" "desec" {
  plugin = "desec"
  api    = "desec"

  data = {
    DEDYN_TOKEN = var.dedyn_token
  }
}

# Order a certificate using the DNS plugin
resource "proxmox_virtual_environment_acme_certificate" "test" {
  node_name = "pve"
  account   = proxmox_virtual_environment_acme_account.example.name
  force     = false

  domains = [
    {
      domain = "pve.example.dedyn.io"
      plugin = proxmox_virtual_environment_acme_dns_plugin.desec.plugin
    }
  ]

  depends_on = [
    proxmox_virtual_environment_acme_account.example,
    proxmox_virtual_environment_acme_dns_plugin.desec
  ]
}
```

### Force Certificate Renewal

```terraform
resource "proxmox_virtual_environment_acme_certificate" "example_force" {
  node_name = "pve"
  force     = true  # This will trigger renewal on every apply
}
```

## Import

ACME certificates can be imported using the node name:

```shell
#!/usr/bin/env sh
# ACME certificates can be imported using the node name, e.g.:
terraform import proxmox_virtual_environment_acme_certificate.example pve
```

## Related Resources

- [`proxmox_virtual_environment_acme_account`](virtual_environment_acme_account) - Manages ACME accounts
- [`proxmox_virtual_environment_acme_dns_plugin`](virtual_environment_acme_dns_plugin) - Manages ACME DNS plugins for DNS-01 challenges
- [`proxmox_virtual_environment_certificate`](virtual_environment_certificate) - Manages custom SSL/TLS certificates (non-ACME)
