---
layout: page
title: proxmox_virtual_environment_realm_ldap
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_realm_ldap

Manages an LDAP authentication realm in Proxmox VE.

LDAP realms allow Proxmox to authenticate users against an LDAP directory service such as OpenLDAP, FreeIPA, or Active Directory (using LDAP protocol).

## Privileges Required

| Path | Attribute |
|-----------------|----------------|
| /access/domains | Realm.Allocate |

## Example Usage

### Basic LDAP Configuration

```hcl
resource "proxmox_virtual_environment_realm_ldap" "example" {
  realm     = "example.com"
  server1   = "ldap.example.com"
  base_dn   = "ou=users,dc=example,dc=com"
  user_attr = "uid"
  comment   = "Corporate LDAP directory"
}
```

### LDAP with Authentication

```hcl
resource "proxmox_virtual_environment_realm_ldap" "secure_example" {
  realm         = "secure.example.com"
  server1       = "ldap.example.com"
  server2       = "ldap-backup.example.com"
  base_dn       = "ou=users,dc=example,dc=com"
  bind_dn       = "cn=readonly,dc=example,dc=com"
  bind_password = var.ldap_password
  user_attr     = "uid"
  secure        = true
  verify        = true
  capath        = "/etc/pve/priv/ca.crt"
  comment       = "Secure LDAP with failover"
}
```

### LDAP with User/Group Sync

```hcl
resource "proxmox_virtual_environment_realm_ldap" "sync_example" {
  realm                 = "lab.example.com"
  server1               = "ldap.example.com"
  base_dn               = "cn=users,cn=accounts,dc=lab,dc=example,dc=com"
  bind_dn               = "uid=ldapread,cn=users,cn=accounts,dc=lab,dc=example,dc=com"
  bind_password         = var.ldap_password
  user_attr             = "uid"
  
  # Sync configuration
  sync_attributes       = "email=mail,firstname=givenName,lastname=sn"
  sync_defaults_options = "scope=users,enable-new=1"
  
  # Group configuration
  group_dn              = "cn=groups,cn=accounts,dc=lab,dc=example,dc=com"
  group_filter          = "(objectClass=groupOfNames)"
  group_name_attr       = "cn"
  
  secure                = true
  comment               = "LDAP realm with user/group synchronization"
}
```

### LDAP with StartTLS

```hcl
resource "proxmox_virtual_environment_realm_ldap" "starttls_example" {
  realm      = "starttls.example.com"
  server1    = "ldap.example.com"
  base_dn    = "dc=example,dc=com"
  user_attr  = "uid"
  mode       = "ldap+starttls"
  sslversion = "tlsv1_3"
  verify     = true
  capath     = "/etc/pve/priv/ca.crt"
  comment    = "LDAP with StartTLS"
}
```

## Argument Reference

### Required

- `realm` - (String) Unique identifier for the realm (e.g., "example.com"). Maximum length: 32 characters. Cannot be changed after creation.
- `server1` - (String) Primary LDAP server hostname or IP address.
- `base_dn` - (String) LDAP base DN for user searches (e.g., "ou=users,dc=example,dc=com").

### Optional

#### Connection Settings

- `server2` - (String) Fallback LDAP server hostname or IP address. Used if `server1` is unavailable.
- `port` - (Number) LDAP server port. Valid range: 1-65535. Default: 389 for LDAP, 636 for LDAPS.
- `secure` - (Boolean) Use LDAPS (LDAP over SSL/TLS) instead of plain LDAP. Default: `false`.
- `mode` - (String) LDAP connection mode. Valid values: `ldap`, `ldaps`, `ldap+starttls`. Overrides `secure` when set.
- `sslversion` - (String) SSL/TLS version to use. Valid values: `tlsv1`, `tlsv1_1`, `tlsv1_2`, `tlsv1_3`.
- `verify` - (Boolean) Verify LDAP server SSL certificate. Default: `false`.
- `capath` - (String) Path to CA certificate file for SSL verification (e.g., "/etc/pve/priv/ca.crt").
- `cert` - (String) Path to client certificate for SSL authentication.
- `certkey` - (String) Path to client certificate private key.

#### Authentication Settings

- `bind_dn` - (String) LDAP bind DN for authentication (e.g., "cn=admin,dc=example,dc=com"). If not set, anonymous binding is used.
- `bind_password` - (String, Sensitive) Password for the bind DN. **Note:** This value is stored in Proxmox but is never returned by the API, so Terraform cannot detect external changes to the password.
- `user_attr` - (String) LDAP attribute representing the username. Default: `"uid"`.
- `filter` - (String) LDAP filter for user searches (e.g., "(objectClass=person)").
- `user_classes` - (String) LDAP objectClasses for users, comma-separated (e.g., "inetOrgPerson,posixAccount").

#### Group Settings

- `group_dn` - (String) LDAP base DN for group searches (e.g., "ou=groups,dc=example,dc=com").
- `group_filter` - (String) LDAP filter for group searches (e.g., "(objectClass=groupOfNames)").
- `group_classes` - (String) LDAP objectClasses for groups, comma-separated (e.g., "groupOfNames,posixGroup").
- `group_name_attr` - (String) LDAP attribute representing the group name (e.g., "cn").

#### Synchronization Settings

- `sync_attributes` - (String) Comma-separated list of attributes to sync from LDAP to Proxmox (e.g., "email=mail,firstname=givenName,lastname=sn").
- `sync_defaults_options` - (String) Default synchronization options. Format: `key=value,key=value`. Example: `"scope=users,enable-new=1"`. Valid keys:
  - `scope` - Sync scope: `users`, `groups`, or `both`
  - `enable-new` - Enable newly synced users: `1` or `0`
  - `remove-vanished` - How to handle vanished entries (e.g. `entry;acl`).
  - `full` / `purge` - (Deprecated by Proxmox) Use `remove_vanished` instead.

#### General Settings

- `comment` - (String) Description of the realm.
- `default` - (Boolean) Use this realm as the default for login. Default: `false`.
- `case_sensitive` - (Boolean) Enable case-sensitive username matching. Default: `true`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` - (String) The realm identifier (same as `realm`).

## Import

LDAP realms can be imported using the realm identifier:

```bash
terraform import proxmox_virtual_environment_realm_ldap.example example.com
```

**Note:** When importing, the `bind_password` attribute cannot be imported since it's not returned by the Proxmox API. You'll need to set this attribute in your Terraform configuration after the import to manage it with Terraform.

## Notes

### Password Security

The `bind_password` is sent to Proxmox and stored securely, but it's never returned by the API. This means:
- Terraform cannot detect if the password was changed outside of Terraform
- You must maintain the password in your Terraform configuration or use a variable
- The password will be marked as sensitive in Terraform state

### LDAP vs LDAPS

- **LDAP (port 389)**: Unencrypted connection. Not recommended for production.
- **LDAPS (port 636)**: Encrypted connection using SSL/TLS. Recommended for production.
- **LDAP+StartTLS**: Upgrades plain LDAP connection to TLS. Alternative to LDAPS.

### User Synchronization

To trigger synchronization, use the `proxmox_virtual_environment_realm_sync` resource.

## See Also

- [Proxmox VE User Management](https://pve.proxmox.com/wiki/User_Management)
- [Proxmox VE LDAP Authentication](https://pve.proxmox.com/wiki/User_Management#pveum_ldap)
- [Proxmox API: /access/domains](https://pve.proxmox.com/pve-docs/api-viewer/index.html#/access/domains)
