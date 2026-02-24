---
layout: page
page_title: "Upgrade Guide"
subcategory: Guides
description: |-
    Breaking changes and migration steps across provider versions.
---

# Upgrade Guide

This guide documents breaking changes across provider versions and the recommended upgrade process.

## General upgrade process

1. **Pin the current version** in your `required_providers` block:

   ```terraform
   terraform {
     required_providers {
       proxmox = {
         source  = "bpg/proxmox"
         version = "~> 0.96.0"
       }
     }
   }
   ```

2. **Review the [CHANGELOG](https://github.com/bpg/terraform-provider-proxmox/blob/main/CHANGELOG.md)** for all versions between your current and target versions. Pay special attention to entries marked **BREAKING CHANGES**.

3. **Update the version constraint** to the target version.

4. Run `terraform init -upgrade` to download the new provider.

5. Run `terraform plan` and review the output. Look for unexpected resource recreation or attribute changes.

6. Apply when satisfied.

-> **Tip:** Upgrade one minor version at a time when crossing multiple breaking changes. This makes it easier to isolate issues.

## v0.95.0

### Node data source `cpu_count` fix

The `proxmox_virtual_environment_node` and `proxmox_virtual_environment_nodes` data sources returned inconsistent values for `cpu_count`. This has been corrected so both return the same value.

**Action required:** If you depend on the previous (incorrect) value in expressions or outputs, update downstream references.

## v0.92.0

### `download_file` Content-Length checking restored

The `proxmox_virtual_environment_download_file` resource now checks the upstream URL's `Content-Length` during every refresh when `overwrite = true` (the default). This restores original behavior from v0.33.0 that was accidentally removed in v0.78.2.

**Before (v0.78.2–v0.91.x):** No upstream check — Terraform would not detect if the remote file had changed.

**After (v0.92.0+):** Terraform detects upstream changes and triggers re-download.

**Action required:** If you intentionally want to skip upstream checking (for example, with URLs that don't return `Content-Length`), set `overwrite = false`:

```terraform
resource "proxmox_virtual_environment_download_file" "image" {
  # ...
  overwrite = false
}
```

### `firewall_options` validation tightened

The `proxmox_virtual_environment_firewall_options` resource now requires exactly one of `vm_id` or `container_id` at validation time.

**Before:** Configurations with only `node_name` would pass validation but fail at runtime.

**After:** Terraform rejects the configuration during `plan` with a clear error.

**Action required:** Ensure your `firewall_options` resources specify either `vm_id` or `container_id`.

## v0.89.0

### `cpu.units` default reverted

The `cpu.units` attribute default was reverted to use the PVE server default instead of a provider-set value.

**Before:** The provider set `cpu.units = 1024` explicitly.

**After:** The provider omits `cpu.units` from API calls unless you set it, letting the PVE server use its own default.

**Action required:** If you rely on `cpu.units = 1024`, set it explicitly in your configuration:

```terraform
resource "proxmox_virtual_environment_vm" "example" {
  # ...
  cpu {
    cores = 2
    units = 1024
  }
}
```

## v0.75.0

### Removed deprecated `initialization` attributes

Two deprecated attributes were removed from the `proxmox_virtual_environment_vm` resource:

- `initialization.dns.server` (singular) — use `initialization.dns.servers` (plural)
- `initialization.upgrade` — removed entirely

**Before:**

```terraform
resource "proxmox_virtual_environment_vm" "example" {
  # ...
  initialization {
    dns {
      server = "1.1.1.1"
    }
    upgrade = true
  }
}
```

**After:**

```terraform
resource "proxmox_virtual_environment_vm" "example" {
  # ...
  initialization {
    dns {
      servers = ["1.1.1.1", "8.8.8.8"]
    }
  }
}
```

### `datastores` data source restructured

The `proxmox_virtual_environment_datastores` data source uses a new structured format with filters.

**Action required:** Review the [datastores data source documentation](https://registry.terraform.io/providers/bpg/proxmox/latest/docs/data-sources/virtual_environment_datastores) and update your configuration to the new format.

## v0.74.0

### Disk behavior in clone operations

Disk management during clone operations was both clarified in documentation and fixed in code — disks moved during a VM clone are now handled correctly during subsequent updates.

**Action required:** If you clone VMs and then modify disk configuration, verify that updates behave as expected. Review the [Clone VM guide](clone-vm) for current disk handling semantics.

## v0.69.0

### `cpu.architecture` handling improved

The `cpu.architecture` attribute handling was improved. Previously, some architecture values were not handled correctly.

**Action required:** If you set `cpu.architecture` explicitly, verify that your configuration still produces the expected result after upgrading.

## v0.65.0

### `vga.enabled` removed

The deprecated `vga.enabled` attribute was removed from the `proxmox_virtual_environment_vm` resource.

**Action required:** Remove `enabled` from any `vga` blocks. To disable VGA, omit the `vga` block entirely or configure it with `type = "none"`.

## Earlier versions

| Version | Breaking change | Action |
| ------- | --------------- | ------ |
| v0.48.0 | File snippets now upload via SSH input stream instead of SCP | No config change; ensure SSH connectivity to nodes |
| v0.40.0 | LXC `features` block is now updatable; mount type support added | Review LXC `features` blocks — previously ignored updates now apply |
| v0.4.0 | `disk.interface` is now required | Add `interface` to all `disk` blocks |
| v0.2.0 | `cloud_init` renamed to `initialization`; `os_type` renamed to `operating_system.type` | Rename attributes in configuration |
