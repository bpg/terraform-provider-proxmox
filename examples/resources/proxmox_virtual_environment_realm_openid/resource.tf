resource "proxmox_virtual_environment_realm_openid" "example" {
  realm      = "example-oidc"
  issuer_url = "https://auth.example.com"
  client_id  = "your-client-id"
  client_key = var.oidc_client_secret

  # Username mapping
  username_claim = "email"

  # User provisioning
  autocreate = true

  # Group mapping (optional)
  groups_claim      = "groups"
  groups_autocreate = true
  groups_overwrite  = false

  # Scopes and prompt
  scopes         = "openid email profile"
  query_userinfo = true

  comment = "Example OpenID Connect realm managed by Terraform"
}
