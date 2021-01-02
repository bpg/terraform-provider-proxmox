---
layout: page
title: Certificate
permalink: /resources/certificate
nav_order: 2
parent: Resources
---

# Resource: Certificate

Manages the custom SSL/TLS certificate for a specific node.

## Example Usage

```
resource "proxmox_virtual_environment_certificate" "example" {
  certificate = "${tls_self_signed_cert.proxmox_virtual_environment_certificate.cert_pem}"
  node_name   = "first-node"
  private_key = "${tls_private_key.proxmox_virtual_environment_certificate.private_key_pem}"
}

resource "tls_private_key" "proxmox_virtual_environment_certificate" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "proxmox_virtual_environment_certificate" {
  key_algorithm   = "${tls_private_key.proxmox_virtual_environment_certificate.algorithm}"
  private_key_pem = "${tls_private_key.proxmox_virtual_environment_certificate.private_key_pem}"

  subject {
    common_name  = "example.com"
    organization = "Terraform Provider for Proxmox"
  }

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}
```

## Arguments Reference

* `certificate` - (Required) The PEM encoded certificate.
* `certificate_chain` - (Optional) The PEM encoded certificate chain.
* `node_name` - (Required) A node name.
* `private_key` - (Required) The PEM encoded private key.

## Attributes Reference

* `expiration_date` - The expiration date (RFC 3339).
* `file_name` - The file name.
* `issuer` - The issuer.
* `public_key_size` - The public key size.
* `public_key_type` - The public key type.
* `ssl_fingerprint` - The SSL fingerprint.
* `start_date` - The start date (RFC 3339).
* `subject` - The subject.
* `subject_alternative_names` - The subject alternative names.
