<!-- markdownlint-disable MD060 -->

# Gap Matrix: `proxmox_vm` parity tracker (#1231)

**Issue:** [bpg/terraform-provider-proxmox#1231](https://github.com/bpg/terraform-provider-proxmox/issues/1231)
**Audit:** [1231_AUDIT.md](1231_AUDIT.md)
**Design:** [1231_DESIGN.md](1231_DESIGN.md)
**Tracker:** [1231_TRACKER.md](1231_TRACKER.md)
**Status:** Living through Phase 2 — becomes PR #20 parity report

## Purpose

Cross-reference of:

- Capabilities inventory (audit Section 2)
- Per-attribute classification (audit Section 4)
- Legacy test inventory (audit Section 3)

…with **status updated as Phase 2 PRs land**. PR #20 (parity report) is this
file, finalized.

## Status legend

| Code | Meaning |
|---|---|
| `todo` | Not yet implemented |
| `wip` | Branch open for the PR that will deliver it |
| `done` | Implemented and merged |
| `dropped` | Deliberately not implementing — see Notes |
| `waived` | Test or attribute carries forward but doesn't need a port |
| `open` | Blocked on a maintainer decision |

## Capabilities

> One row per legacy SDK attribute / block. Columns track ownership through
> Phase 2 and tick `done` as PRs merge.

### Top-level

| Capability | Target PR | Status | Notes |
|---|---|---|---|
| TBD | — | — | — |

### Disk

| Capability | Target PR | Status | Notes |
|---|---|---|---|
| disk (map-keyed block) | #7 | todo | First map-keyed application in `proxmox_vm` |

### Network

| Capability | Target PR | Status | Notes |
|---|---|---|---|
| network_device (map-keyed block) | #10 | todo | — |

### Cloud-init

| Capability | Target PR | Status | Notes |
|---|---|---|---|
| initialization | #11 | todo | Depends on disk + network |

### Boot / firmware

| Capability | Target PR | Status | Notes |
|---|---|---|---|
| bios | #8 | todo | — |
| machine | #8 | todo | — |
| boot_order | #8 | todo | — |
| efi_disk (single-nested) | #9 | todo | Architectural-single per ADR-008 |
| tpm_state (single-nested) | #9 | todo | Architectural-single per ADR-008 |
| scsi_hardware | #9 | todo | — |

### OS / metadata

| Capability | Target PR | Status | Notes |
|---|---|---|---|
| operating_system | #12 | todo | — |
| smbios | #12 | todo | — |

### Advanced hardware

| Capability | Target PR | Status | Notes |
|---|---|---|---|
| agent | #13 | todo | — |
| numa (map-keyed block) | #13 | todo | Includes `numa.enabled` (rehomed `cpu.numa`) |
| watchdog | #13 | todo | — |
| acpi | #14 | todo | — |
| tablet_device | #14 | todo | — |
| keyboard_layout | #14 | todo | — |
| kvm_arguments | #14 | todo | — |
| vcpus | #14 | todo | Rehomed `cpu.hotplugged` |
| hotplug | #14 | todo | — |
| parallel (map-keyed block) | #14 | todo | — |
| usb (map-keyed block) | #15 | todo | — |
| hostpci (map-keyed block) | #16 | todo | — |
| serial_device (map-keyed block) | #17 | todo | — |
| audio_device (map-keyed block, one-key today) | #17 | todo | Map-keyed for forward-compat per ADR-008 |
| virtiofs (map-keyed block) | #17 | todo | — |

### Cluster / placement

| Capability | Target PR | Status | Notes |
|---|---|---|---|
| startup | #18 | todo | — |
| pool_id | #18 | todo | — |
| protection | #18 | todo | — |
| hook_script_file_id | #18 | todo | — |
| amd_sev | #18 | todo | — |
| migrate | #19 | todo | Depends on full device set |

### Lifecycle / runtime

| Capability | Target PR | Status | Notes |
|---|---|---|---|
| memory + power_state + on_boot | #6 | todo | MVP setup — sets up first PVE-driven boot |
| `started` (legacy) | — | dropped | Replaced by `power_state` (Q5 resolution) |
| `reboot` (legacy user-facing) | — | dropped | Provider decides from pending-changes (Q5) |

## Per-attribute classification (Optional+Computed cleanup, PR #3)

> Mirrors audit Section 4. Updated when PR #3 lands.

| Attribute | Current schema | Target schema | Status | Target PR |
|---|---|---|---|---|
| TBD | — | — | — | — |

## Legacy test ports

> Mirrors audit Section 3. Updated as Phase 2 PRs port the relevant tests.

| Test name | Behavior | Target PR | Status |
|---|---|---|---|
| TBD | — | — | — |

## Open questions remaining

| # | Question | Owner | Status |
|---|---|---|---|
| — | All design Q1–Q5 resolved during grilling | — | — |

## PR ownership map

> Quick-view of which PR owns which capability cluster. Sourced from the
> design's Phase 2 sections.

| PR | Phase | Ownership |
|---|---|---|
| 6 | 2A | memory, power_state, on_boot |
| 7 | 2A | disk (map-keyed) |
| 8 | 2B | bios, machine, boot_order |
| 9 | 2B | efi_disk, tpm_state, scsi_hardware |
| 10 | 2C | network_device (map-keyed) |
| 11 | 2C | initialization (cloud-init) |
| 12 | 2C | operating_system, smbios |
| 13 | 2D | agent, numa, watchdog |
| 14 | 2D | acpi, tablet, keyboard, kvm_args, vcpus, hotplug, parallel |
| 15 | 2D | usb (map-keyed) |
| 16 | 2D | hostpci (map-keyed) |
| 17 | 2D | serial, audio, virtiofs (all map-keyed) |
| 18 | 2E | startup, pool_id, protection, hook_script, amd_sev |
| 19 | 2E | migrate |
| 20 | 2E | parity report + SDK migration guide (this file finalized) |

## Datasource parity

> Per ADR-006 + design Q1: the datasource gets the same map-keyed blocks
> with `Computed: true`. Provider-only behavior attributes are excluded.

### Excluded from datasource

| Attribute | Reason |
|---|---|
| `purge_on_destroy` | Provider behavior on Delete only |
| `stop_on_destroy` | Provider behavior on Delete only |
| `delete_unreferenced_disks_on_destroy` | Provider behavior on Delete only |
| `template` | Read-from-PVE, but PVE exposes via `is_template`; check spelling at PR-time |
| `power_state` | Maps to runtime `status`; datasource uses `status` directly |
| `timeouts` | Provider config, not VM state |

### Included (every Resource attribute not excluded)

| Block / attribute | Datasource shape | Status | Target PR |
|---|---|---|---|
| TBD | — | — | — |
