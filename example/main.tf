provider "proxmox" {
  endpoint = var.dev_virtual_environment_endpoint
  api_token = var.dev_virtual_environment_api_token
  insecure = true
  ssh {
    agent = true
    username = var.dev_virtual_environment_ssh_username
  }
}
