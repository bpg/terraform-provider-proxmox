resource "proxmox_virtual_environment_realm_ldap" "example" {
  realm = "example-ldap"

  # LDAP server configuration
  server1 = "ldap.example.com"
  port    = 389

  # Base DN and user attribute
  base_dn   = "ou=people,dc=example,dc=com"
  user_attr = "uid"

  # Bind credentials (optional but recommended)
  bind_dn       = "cn=admin,dc=example,dc=com"
  bind_password = "secure-password"

  # SSL/TLS configuration
  mode   = "ldap+starttls"
  verify = true

  # Group synchronization (optional)
  group_dn     = "ou=groups,dc=example,dc=com"
  group_filter = "(objectClass=groupOfNames)"

  comment = "Example LDAP realm managed by Terraform"
}
