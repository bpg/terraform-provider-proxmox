provider "proxmox" {
  endpoint  = var.virtual_environment_endpoint
  api_token = var.virtual_environment_api_token
  insecure  = true
  ssh {
    agent    = true
    username = var.virtual_environment_ssh_username
    # When using api_token, there is no provider password to inherit for SSH.
    # Ensure ssh-agent has the appropriate keys loaded, or set password / private_key here.
  }
}

provider "proxmox" {
  alias    = "root"
  endpoint = var.virtual_environment_endpoint
  username = "root@pam"
  password = var.virtual_environment_root_password
  insecure = true
}
