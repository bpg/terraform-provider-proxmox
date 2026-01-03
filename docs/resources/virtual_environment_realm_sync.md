---
layout: page
title: proxmox_virtual_environment_realm_sync
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_realm_sync

Triggers synchronization of an existing authentication realm in Proxmox VE.

This resource wraps the `/access/domains/{realm}/sync` API and is intended to be
used alongside realm configuration resources such as
`proxmox_virtual_environment_realm_ldap`.

## Example Usage

```hcl
resource "proxmox_virtual_environment_realm_ldap" "example" {
  realm     = "example.com"
  server1   = "ldap.example.com"
  base_dn   = "ou=users,dc=example,dc=com"
  user_attr = "uid"
}

resource "proxmox_virtual_environment_realm_sync" "example_sync" {
  realm = proxmox_virtual_environment_realm_ldap.example.realm
  scope = "users"
}
```

## Argument Reference

- `realm` - (Required, String) Name of the realm to synchronize.
- `scope` - (Optional, String) Sync scope. Valid values: `"users"`, `"groups"`, `"both"`.
- `remove_vanished` - (Optional, String) How to handle vanished entries. Typically a
  semicolon-separated combination of `acl`, `properties`, `entry`, or `none` (e.g. `acl;properties;entry`).
- `enable_new` - (Optional, Boolean) Enable newly synced users.
- `full` - (Optional, Boolean) Perform a full sync.
- `purge` - (Optional, Boolean) Purge removed entries.
- `dry_run` - (Optional, Boolean) Only simulate the sync without applying changes.

## Behavior Notes

- The sync operation is **one-shot**: applying the resource runs the sync
  with the specified options. Proxmox does not expose a persistent sync
  object, so this resource only records the last requested sync
  configuration in Terraform state.
- Destroying the resource does **not** undo any previously performed sync;
  it simply removes the resource from Terraform state.

## Import

You can import a sync resource by realm name:

```bash
terraform import proxmox_virtual_environment_realm_sync.example example.com
```

Importing only populates the `realm` and `id` attributes; other fields must
be set in configuration.
