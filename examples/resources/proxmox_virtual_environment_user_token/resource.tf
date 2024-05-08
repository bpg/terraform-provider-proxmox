# if creating a user token, the user must be created first
resource "proxmox_virtual_environment_user" "user" {
  comment         = "Managed by Terraform"
  email           = "user@pve"
  enabled         = true
  expiration_date = "2034-01-01T22:00:00Z"
  user_id         = "user@pve"
}

resource "proxmox_virtual_environment_user_token" "user_token" {
  comment         = "Managed by Terraform"
  expiration_date = "2033-01-01T22:00:00Z"
  token_name      = "tk1"
  user_id         = proxmox_virtual_environment_user.user.user_id
}
