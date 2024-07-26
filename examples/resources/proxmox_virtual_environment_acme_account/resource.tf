resource "proxmox_virtual_environment_acme_account" "example" {
  name      = "example"
  contact   = "example@email.com"
  directory = "https://acme-staging-v02.api.letsencrypt.org/directory"
  tos       = "https://letsencrypt.org/documents/LE-SA-v1.3-September-21-2022.pdf"
}
