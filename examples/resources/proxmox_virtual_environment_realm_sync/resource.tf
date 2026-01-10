resource "proxmox_virtual_environment_realm_ldap" "example" {
  realm = "example-ldap"

  server1   = "ldap.example.com"
  port      = 389
  base_dn   = "ou=people,dc=example,dc=com"
  user_attr = "uid"

  # Enable group sync
  group_dn     = "ou=groups,dc=example,dc=com"
  group_filter = "(objectClass=groupOfNames)"
}

resource "proxmox_virtual_environment_realm_sync" "example" {
  realm = proxmox_virtual_environment_realm_ldap.example.realm

  # Sync both users and groups
  scope = "both"

  # Remove entries that no longer exist in LDAP
  remove_vanished = "acl;entry;properties"

  # Enable new users/groups by default
  enable_new = true
}
