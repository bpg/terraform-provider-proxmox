# Issue #1465 — Creation of Standalone VM Disks

## Status: Open, P1, 45+ upvotes

## Quick Context

Issue: <https://github.com/bpg/terraform-provider-proxmox/issues/1465>
Author: @legege
Date: 2024-08-04
Priority: P1
Labels: enhancement, lifecycle:acknowledged

## What's Requested

A new resource (`proxmox_virtual_environment_disk` or similar) for managing VM disks
as standalone resources, decoupled from the VM lifecycle.

**Proxmox API endpoint**: `POST /api2/json/nodes/{node}/storage/{storage}/content`
(documented at <https://pve.proxmox.com/pve-docs/api-viewer/index.html#/nodes/{node}/storage/{storage}/content>)

## Use Cases

### 1. Kubernetes CSI Integration (original author @legege)
Using the Proxmox CSI Plugin to provision PersistentVolumes. Currently disks are named
`vm-9999-pvc-<UUID>` making them hard to identify. Wants Terraform-managed named disks
that the CSI plugin can use.

Already working around this with the `Mastercard/restapi` provider to call the Proxmox
API directly.

### 2. Persistent Data Volumes (@DataLabTechTV, Sep 2025)
Workflow:
1. Create external disk (e.g., for `/data`)
2. Create VM/CT with rootfs + mount the data disk
3. Destroy and recreate VM/CT without losing the data disk

Example: MinIO deployment where the data volume survives VM recreation.
Applies to both VMs and LXC containers.

### 3. Ephemeral Boot Disk Rotation (from Discussion #2403)
@Vaneixus's use case: replace boot disk image without destroying the VM. A standalone
disk resource could enable: create new disk from image → attach to VM → detach old disk.

## Connection to Discussion #2403

Discussion #2403 asks for in-place disk re-import when `import_from` changes. The
decision there was to keep `import_from` as create-only (aligning with industry
standards). A standalone disk resource would be the **clean long-term solution** for
the underlying need — independent disk lifecycle management.

See: [.dev/2403_DISK_REIMPORT_REVIEW.md](2403_DISK_REIMPORT_REVIEW.md) for the full
analysis of the re-import discussion.

## Technical Notes from Codebase Research

### Proxmox API for Standalone Disks

**Create disk**: `POST /nodes/{node}/storage/{storage}/content`
- Parameters: `vmid`, `filename`, `size`, `format`
- Creates a disk volume in the specified storage
- Returns the volume ID (e.g., `local-lvm:vm-100-disk-0`)

**Delete disk**: `DELETE /nodes/{node}/storage/{storage}/content/{volume}`

**List disks**: `GET /nodes/{node}/storage/{storage}/content`

### Current Provider Architecture

Disks are currently **inline** within VM/CT resources — there's no independent disk
resource. The relevant existing code:

| Component | Location | Notes |
|-----------|----------|-------|
| SDK disk schema | `proxmoxtf/resource/vm/disk/schema.go` | Inline in VM resource |
| SDK disk CRUD | `proxmoxtf/resource/vm/disk/disk.go` | Tied to VM lifecycle |
| Framework cloned_vm disk | `fwprovider/nodes/clonedvm/` | Inline, map-based |
| API disk types | `proxmox/nodes/vms/custom_storage_device.go` | Shared across providers |
| Storage API client | `proxmox/nodes/storage/` | Has content listing, may need create/delete |

### What Would Need to Be Built

1. **New resource**: `proxmox_virtual_environment_disk`
   - Create: `POST /nodes/{node}/storage/{storage}/content`
   - Read: `GET /nodes/{node}/storage/{storage}/content/{volume}`
   - Delete: `DELETE /nodes/{node}/storage/{storage}/content/{volume}`
   - Update: resize only (same as current inline resize)

2. **New resource (optional)**: `proxmox_virtual_environment_vm_disk_attachment`
   - Attach: `PUT /nodes/{node}/qemu/{vmid}/config` with `scsi0=<volume-id>`
   - Detach: `PUT /nodes/{node}/qemu/{vmid}/config` with `delete=scsi0`
   - This decouples attachment from the disk itself

3. **API client extensions**: `proxmox/nodes/storage/` may need:
   - `CreateContent()` method
   - `DeleteContent()` method
   - Content read/list may already exist

4. **Framework only** — new resources go in `fwprovider/`

### Design Considerations

- **Naming**: Proxmox names disks as `vm-{vmid}-disk-{N}`. A standalone disk still needs
  a `vmid` in the API call — but this is just a namespace, not actual VM ownership.
- **Import support**: Allow importing existing Proxmox disks into Terraform state.
- **LXC support**: @DataLabTechTV's use case applies to containers too. The mount point
  attachment would be a separate resource from VM disk attachment.
- **Relationship to inline disks**: Need to decide if inline disk blocks and standalone
  disk resources can coexist on the same VM, or if they're mutually exclusive.

## Next Steps

- [ ] Decide on resource design (disk + attachment vs single resource)
- [ ] Check if storage API client already has create/delete content methods
- [ ] Determine scope: VMs only, or also LXC containers?
- [ ] Consider relationship to `proxmox_virtual_environment_file` (which already handles
      some storage content operations)
- [ ] Create implementation plan if proceeding

## Related

- Discussion #2403 — Disk re-import (decided: keep create-only)
- Discussion #2324 — Clone block vs explicit device config (architectural)
- Issue #1231 — Migrate VM resource to Plugin Framework (umbrella)
- Issue #1501 — Support seamless datastore_id changes for disks
- Issue #2285 — Disks of imported VM not properly matched to interfaces
