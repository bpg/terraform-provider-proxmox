provider "proxmox" {
  virtual_environment {
    endpoint = var.virtual_environment_endpoint
    username = var.virtual_environment_username
    password = var.virtual_environment_password
    insecure = true
  }
}
