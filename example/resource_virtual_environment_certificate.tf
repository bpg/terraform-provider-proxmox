resource "proxmox_virtual_environment_certificate" "example" {
  certificate = "${tls_self_signed_cert.proxmox_virtual_environment_certificate.cert_pem}"
  node_name   = "${data.proxmox_virtual_environment_nodes.example.names[0]}"
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
