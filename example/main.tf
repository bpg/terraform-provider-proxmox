provider "proxmox" {
  endpoint = var.virtual_environment_endpoint
  insecure = true
  username = var.virtual_environment_username
  password = var.virtual_environment_password
}
