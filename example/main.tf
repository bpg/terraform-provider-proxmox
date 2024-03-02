provider "proxmox" {
  endpoint = var.virtual_environment_endpoint
  api_token = var.virtual_environment_api_token
  insecure = true
  ssh {
    agent = true
    username = var.virtual_environment_ssh_username
  }
}
