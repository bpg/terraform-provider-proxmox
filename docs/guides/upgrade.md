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
         version = "~> 0.104.1" # x-release-please-version
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

## v0.102.0

### `proxmox_metrics_server.disable` now defaults to `false`

The `disable` attribute on `proxmox_metrics_server` (and its alias `proxmox_virtual_environment_metrics_server`) now has a schema default of `false`. The Proxmox API omits `disable` from GET responses when the value is `false`, and the resource now reads state back from the API after Create and Update — so the default is required to keep state consistent.

**Before:** When `disable` was not set in configuration, it was stored as null in state.

**After:** When `disable` is not set in configuration, it is stored as `false` in state. This matches the Proxmox server default (metrics servers are enabled by default).

**Action required:** None. On the first `terraform apply` after upgrading, Terraform may show a one-time diff of `disable: null -> false` for existing metrics server resources that did not set `disable`. The apply is a no-op against the Proxmox API.

### LXC `cpu.units` no longer defaults to 1024

The `cpu.units` attribute on `proxmox_virtual_environment_container` no longer defaults to `1024`. Instead, the Proxmox server's own default is used (`100` on cgroup2, `1024` on cgroup1).

**Before:** When `cpu.units` was not specified, the provider always sent `cpuunits=1024` to the API, overriding the server default.

**After:** When `cpu.units` is not specified, the provider does not send `cpuunits` at all, letting the server use its own default.

**Action required:** If you rely on `cpu.units = 1024` and are on a cgroup2 system, explicitly set it in your configuration:

```terraform
resource "proxmox_virtual_environment_container" "example" {
  # ...
  cpu {
    units = 1024
  }
}
```

## v0.101.0

### VM datasources error on not-found

Both VM datasources now return an error when the specified VM does not exist, instead of silently returning empty or broken state.

- **Framework** (`proxmox_virtual_environment_vm2` / `proxmox_vm`): Was previously non-functional (`Value Conversion Error` on every call). Now works correctly but errors on missing VMs.
- **SDK** (`proxmox_virtual_environment_vm`): Previously returned an unclear raw error. Now returns a proper "not found" diagnostic. This datasource is also now **deprecated** in favor of `proxmox_vm`.

**Action required:** If your configuration depends on a VM datasource not erroring when the VM doesn't exist, switch to the `proxmox_virtual_environment_vms` list datasource — it returns an empty list when no VMs match. Note that this datasource supports filtering by `name`, `template`, `status`, and `node_name`, but not by `vm_id` directly. To find a specific VM by ID, post-filter the results:

```terraform
data "proxmox_virtual_environment_vms" "all" {
  node_name = "pve"
}

locals {
  my_vm = [for vm in data.proxmox_virtual_environment_vms.all.vms : vm if vm.vm_id == 999]
}
```

## v0.100.0

### `template` attribute no longer forces recreation

The `template` attribute on `proxmox_virtual_environment_vm` no longer forces resource recreation. Converting a VM to a template (or back) now happens in-place.

**Before:** Changing `template = true` caused Terraform to destroy and recreate the VM.

**After:** The conversion is applied in-place without recreation.

**Action required:** If you relied on the destroy-and-recreate behavior (for example, to get a fresh VM when toggling `template`), you may need to use `terraform taint` or add a `replace_triggered_by` lifecycle rule.

### Shutdown operations fail when `reboot_after_update = false`

Operations that require a VM shutdown — TPM state migration, cloud-init disk move, template conversion, and disk move — now fail with an error when `reboot_after_update = false`, instead of silently powering off the VM.

**Before:** The provider would silently shut down a running VM even when `reboot_after_update = false`.

**After:** The provider returns an error, respecting the `reboot_after_update = false` setting.

**Action required:** If you perform operations that require a shutdown (TPM, cloud-init move, template conversion, disk move), either:

- Set `reboot_after_update = true` (default), or
- Shut down the VM before applying changes

### Short-name resource aliases (`proxmox_*` prefix)

All existing Framework resources now have shorter aliases using the `proxmox_` prefix instead of `proxmox_virtual_environment_`. Both names work simultaneously — the old names emit a deprecation warning. New resources added from v0.99.0 onward (such as `proxmox_backup_job` and `proxmox_harule`) use the short prefix exclusively. See [ADR-007](https://github.com/bpg/terraform-provider-proxmox/blob/main/docs/adr/007-resource-type-name-migration.md) for full details.

Examples of the new short names:

| Old name                                           | New name                       |
|----------------------------------------------------|--------------------------------|
| `proxmox_virtual_environment_vm2` (resource)       | `proxmox_vm`                   |
| `proxmox_virtual_environment_download_file`        | `proxmox_download_file`        |
| `proxmox_virtual_environment_network_linux_bridge` | `proxmox_network_linux_bridge` |
| `proxmox_virtual_environment_hagroup`              | `proxmox_hagroup`              |
| `proxmox_virtual_environment_acl`                  | `proxmox_acl`                  |

**Action required:** No immediate action needed — old names continue to work. To migrate at your own pace, use Terraform's `moved` block (requires Terraform >= 1.8):

```terraform
moved {
  from = proxmox_virtual_environment_network_linux_bridge.example
  to   = proxmox_network_linux_bridge.example
}
```

Then update the resource block to use the new name. After a successful `terraform apply`, the `moved` block can be removed.

For Terraform < 1.8, use `terraform state mv`:

```shell
terraform state mv proxmox_virtual_environment_vm.example proxmox_vm.example
```

## v0.98.1

### `network_device.enabled` deprecated

The `enabled` attribute on the `network_device` block in `proxmox_virtual_environment_vm` is deprecated and will be removed in a future release.

**Action required:**

- If you use `enabled = true` (or rely on the default): remove the `enabled` line — the device is enabled by default.
- If you use `enabled = false` to disable a device: remove the entire `network_device` block instead.

```terraform
# Before
resource "proxmox_virtual_environment_vm" "example" {
  # ...
  network_device {
    bridge  = "vmbr0"
    enabled = true  # drop this line
  }
}

# After
resource "proxmox_virtual_environment_vm" "example" {
  # ...
  network_device {
    bridge = "vmbr0"
  }
}
```

## v0.97.1

### VM and container name validation

The `name` attribute on `proxmox_virtual_environment_vm` and the `hostname` in `proxmox_virtual_environment_container` `initialization` block now validate that the value is a valid DNS name. Previously, invalid names were accepted by Terraform and rejected at apply time by the Proxmox API.

**Action required:** If you have VMs or containers with names that are not valid DNS names (for example, names starting with `.` or containing special characters), update them to valid DNS names.

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

**Action required:** If you clone VMs and then modify disk configuration, verify that updates behave as expected. Review the [Clone VM guide](clone-vm.md) for current disk handling semantics.

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
