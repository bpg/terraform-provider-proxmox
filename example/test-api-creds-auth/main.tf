provider "proxmox" {
  insecure = true

  auth_ticket           = var.auth_ticket
  csrf_prevention_token = var.csrf_prevention_token

  api_token = var.api_token

  otp = var.otp

  username = var.username
  password = var.password

}
