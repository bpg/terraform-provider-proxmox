---
layout: page
page_title: "Known Issues"
subcategory: Guides
description: |-
  Known limitations of the provider and of Proxmox VE that affect it, with workarounds.
---

# Known Issues

This page lists limitations that are known, reproducible, and either by design or outside the provider's control.
For bugs, see the [issue tracker](https://github.com/bpg/terraform-provider-proxmox/issues).

## HA-managed VMs and containers

When Proxmox VE HA moves a VM or container to another node, the provider detects the new node on the next read and reports a drift on `node_name`.
Terraform then wants to migrate the guest back (if `migrate = true`) or recreate it on the configured node.
Use `lifecycle { ignore_changes = [node_name] }` to let HA own placement.
See [HA clusters and `node_name` drift](multi-node.md#ha-clusters-and-node_name-drift) in the multi-node guide for the full recipe and its trade-offs.

## Kernel panic when resizing a cloud image boot disk

Debian and Ubuntu cloud images expect a serial console. Resizing their boot disk without one can leave the guest kernel panicking on boot.
Add a serial device to the VM:

```hcl
  serial_device {
    device = "socket"
  }
```

For more context, see [#1639](https://github.com/bpg/terraform-provider-proxmox/issues/1639) and [#1770](https://github.com/bpg/terraform-provider-proxmox/issues/1770).

## Lock errors when creating multiple VMs or containers

Creating many guests at once can fail with `can't lock file '/var/lock/qemu-server/lock-<id>.conf' - got timeout` or similar.
This is an I/O bottleneck in Proxmox VE, not a provider defect.
Create guests sequentially, or run Terraform with `-parallelism=1`.

Additional information and sample error messages are in [#1929](https://github.com/bpg/terraform-provider-proxmox/issues/1929) and [#995](https://github.com/bpg/terraform-provider-proxmox/issues/995).
An OpenTofu feature request for per-provider parallelism is tracked at [opentofu/opentofu#2466](https://github.com/opentofu/opentofu/issues/2466).

## Snippets require a PAM account and SSH

The Proxmox VE API has no endpoint for uploading snippets, so the provider uploads them over SFTP.
This requires SSH access to the node with a PAM account (a regular Linux user), as described in [SSH Connection](../index.md#ssh-connection).
API-only setups cannot manage snippets until Proxmox VE adds an API for it, which is tracked upstream in [Bugzilla #2208](https://bugzilla.proxmox.com/show_bug.cgi?id=2208).

## VMware disk images

Proxmox VE cannot boot a `.vmdk` image directly; it must be imported.
Upload the image to a datastore with the `import` content type enabled, either with `proxmox_virtual_environment_file` (the content type is detected from the `.vmdk` extension) or with `proxmox_virtual_environment_download_file` and `content_type = "import"`, and reference it from the VM's `disk` block through `import_from`.
The import happens through the API and does not require SSH.
