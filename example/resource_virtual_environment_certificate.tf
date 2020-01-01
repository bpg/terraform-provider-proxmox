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

output "resource_proxmox_virtual_environment_certificate_example_expiration_date" {
  value = "${proxmox_virtual_environment_certificate.example.expiration_date}"
}

output "resource_proxmox_virtual_environment_certificate_example_file_name" {
  value = "${proxmox_virtual_environment_certificate.example.file_name}"
}

output "resource_proxmox_virtual_environment_certificate_example_issuer" {
  value = "${proxmox_virtual_environment_certificate.example.issuer}"
}

output "resource_proxmox_virtual_environment_certificate_example_public_key_size" {
  value = "${proxmox_virtual_environment_certificate.example.public_key_size}"
}

output "resource_proxmox_virtual_environment_certificate_example_public_key_type" {
  value = "${proxmox_virtual_environment_certificate.example.public_key_type}"
}

output "resource_proxmox_virtual_environment_certificate_example_ssl_fingerprint" {
  value = "${proxmox_virtual_environment_certificate.example.ssl_fingerprint}"
}

output "resource_proxmox_virtual_environment_certificate_example_start_date" {
  value = "${proxmox_virtual_environment_certificate.example.start_date}"
}

output "resource_proxmox_virtual_environment_certificate_example_subject" {
  value = "${proxmox_virtual_environment_certificate.example.subject}"
}

output "resource_proxmox_virtual_environment_certificate_example_subject_alternative_names" {
  value = "${proxmox_virtual_environment_certificate.example.subject_alternative_names}"
}
