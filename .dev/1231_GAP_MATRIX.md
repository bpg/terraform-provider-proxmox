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

### Top-level scalars

| Capability | Target PR | Status | Notes |
|---|---|---|---|
| description | — | done | In `proxmox_vm` |
| name | — | done | In `proxmox_vm` (DNS validator) |
| node_name | — | done | In `proxmox_vm` (Required) |
| tags | — | done | In `proxmox_vm` (stringset) |
| template | — | done | In `proxmox_vm` (RequiresReplace) |
| id (SDK `vm_id`) | — | done | In `proxmox_vm` (renamed) |
| pool_id | #18 | todo | — |
| protection | #18 | todo | — |
| migrate | #19 | todo | — |
| acpi | #14 | todo | — |
| bios | #8 | todo | — |
| boot_order | #8 | todo | — |
| hook_script_file_id | #18 | todo | — |
| hotplug | #14 | todo | — |
| keyboard_layout | #14 | todo | — |
| kvm_arguments | #14 | todo | — |
| machine | #8 | todo | — |
| scsi_hardware | #9 | todo | — |
| tablet_device | #14 | todo | — |
| stop_on_destroy | — | done | In `proxmox_vm` |
| purge_on_destroy | — | done | In `proxmox_vm` |
| delete_unreferenced_disks_on_destroy | — | done | In `proxmox_vm` |
| power_state (new) | #6 | todo | Per Q5 |
| status (Computed) | #6 | todo | Runtime mirror |
| vcpus (rehomed `cpu.hotplugged`) | #14 | todo | Per design D7/P3 |
| on_boot | #6 | todo | PVE "Start at boot"; `Optional` only per ADR-004 |

### Existing sub-blocks (PR #3 contract port)

| Capability | Target PR | Status | Notes |
|---|---|---|---|
| cpu | #3 | todo | Port to ADR-008; drop long-enum validator (F27); drop `numa`/`hotplugged` (P3); fix Limit-branch bug (P1/F24); drop sentinels (P2/F20–F22) |
| vga | #3 | todo | Port to ADR-008; drop long-enum `type` validator (F33) |
| rng | #3 | todo | Port to ADR-008; address int-zero trap (F37) |
| cdrom | #3 | todo | Already map-keyed; tighten slot regex (F46); verify `file_id` default (F47) |
| memory | #3 | todo | **Blocker fixes:** F39 (provider Defaults), F40–F42 (sentinels), F43 (no FillCreateBody), F44 (no-delete update). Wired into `proxmox_vm` in #6. |

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
| memory + power_state + on_boot (MVP setup) | #6 | todo | First PVE-driven boot end-to-end |
| `started` (legacy) | — | dropped | Replaced by `power_state` (Q5) |
| `reboot` (legacy user-facing) | — | dropped | Provider decides from pending-changes (Q5) |
| `reboot_after_update` (legacy) | — | dropped | Same as above |

### Timeouts (folded into `timeouts` block)

| Capability | Target PR | Status | Notes |
|---|---|---|---|
| `timeouts.create` | — | done | Already in `proxmox_vm` |
| `timeouts.read`/`update`/`delete` | — | done | Already in `proxmox_vm` |
| `timeout_migrate` (SDK) | #19 | todo | Folded into `timeouts.update` |
| `timeout_move_disk` (SDK) | #19 | todo | Folded into migrate timeout |
| `timeout_clone` (SDK) | — | dropped | Belongs to clonedvm |
| `timeout_reboot` (SDK) | — | dropped | Provider-internal reboot uses `timeouts.update` |
| `timeout_shutdown_vm` / `timeout_start_vm` / `timeout_stop_vm` (SDK) | #6 | todo | Provider-internal (not user-facing) |

### Dropped / out of scope

| Capability | Target PR | Status | Notes |
|---|---|---|---|
| `clone` (SDK block) | — | dropped | Belongs to `proxmox_cloned_vm` (D4) |
| `interface` (legacy disk slot field) | — | dropped | Replaced by map key per ADR-008 |

## Per-attribute classification (Optional+Computed cleanup, PR #3)

> Mirrors audit Section 4. Depends on mitmproxy verification per attribute.
> Skeletal list below — classification column fills in after Section 4 work.

### `cpu` attributes (`fwprovider/nodes/vm/cpu/resource_schema.go`)

| Attribute | Current schema | Target schema (post mitmproxy) | Status | Target PR |
|---|---|---|---|---|
| `cpu.affinity` | Optional+Computed | TBD | open | #3 |
| `cpu.architecture` | Optional+Computed | TBD | open | #3 |
| `cpu.cores` | Optional+Computed | TBD | open | #3 |
| `cpu.flags` | Optional+Computed | TBD | open | #3 |
| `cpu.hotplugged` | Optional+Computed | dropped (rehomed to `vcpus`) | — | #3 / #14 |
| `cpu.limit` | Optional+Computed | TBD | open | #3 |
| `cpu.numa` | Optional+Computed | dropped (rehomed to `numa.enabled`) | — | #3 / #13 |
| `cpu.sockets` | Optional+Computed | TBD | open | #3 |
| `cpu.type` | Optional+Computed | TBD | open | #3 |
| `cpu.units` | Optional+Computed | TBD | open | #3 |

### `vga` attributes (`fwprovider/nodes/vm/vga/resource_schema.go`)

| Attribute | Current schema | Target schema (post mitmproxy) | Status | Target PR |
|---|---|---|---|---|
| `vga.clipboard` | Optional+Computed | TBD | open | #3 |
| `vga.type` | Optional+Computed | TBD | open | #3 |
| `vga.memory` | Optional+Computed | TBD | open | #3 |

### `rng` attributes (`fwprovider/nodes/vm/rng/resource_schema.go`)

| Attribute | Current schema | Target schema (post mitmproxy) | Status | Target PR |
|---|---|---|---|---|
| `rng.source` | Optional+Computed | TBD | open | #3 |
| `rng.max_bytes` | Optional+Computed | TBD | open | #3 |
| `rng.period` | Optional+Computed | TBD | open | #3 |

### `memory` attributes (`fwprovider/nodes/vm/memory/resource_schema.go`)

| Attribute | Current schema | Target schema (post mitmproxy) | Status | Target PR |
|---|---|---|---|---|
| `memory.size` | Optional+Computed+Default=512 | TBD (drop Default per F39; classify) | open | #3 |
| `memory.balloon` | Optional+Computed+Default=0 | TBD (drop Default per F39; classify) | open | #3 |
| `memory.shares` | Optional+Computed+Default=1000 | TBD (drop Default per F39; classify) | open | #3 |
| `memory.hugepages` | Optional+Computed | TBD | open | #3 |
| `memory.keep_hugepages` | Optional+Computed | TBD | open | #3 |

### `cdrom` map-level + per-slot (`fwprovider/nodes/vm/cdrom/resource_schema.go`)

| Attribute | Current schema | Target schema (post mitmproxy) | Status | Target PR |
|---|---|---|---|---|
| `cdrom` (map-level) | Optional+Computed | TBD (verify PVE auto-populates default cdrom devices) | open | #3 |
| `cdrom[slot].file_id` | Optional+Computed+Default="cdrom" | TBD (F47: verify if PVE auto-defaults) | open | #3 |

### Top-level (existing)

| Attribute | Current schema | Target schema (post mitmproxy) | Status | Target PR |
|---|---|---|---|---|
| `description` | Optional | keep Optional | confirmed | — |
| `name` | Optional | keep Optional | confirmed | — |
| `tags` | stringset (Optional) | keep Optional | confirmed | — |
| `template` | Optional (with RequiresReplace) | keep | confirmed | — |
| `stop_on_destroy` | Optional+Computed+Default=false | keep (provider-only per ADR-004) | confirmed | — |
| `purge_on_destroy` | Optional+Computed+Default=true | keep (provider-only per ADR-004) | confirmed | — |
| `delete_unreferenced_disks_on_destroy` | Optional+Computed+Default=true | keep (provider-only per ADR-004) | confirmed | — |

## Legacy test ports

> Mirrors audit Section 3. Updated as Phase 2 PRs port the relevant tests.

### Acceptance tests (in `fwprovider/test/`)

| Test | Target PR | Status | Notes |
|---|---|---|---|
| `TestAccResourceVM` (omnibus, ~40 sub-cases) | spread across #3, #6, #13, #14, #17 | todo | Split by sub-case cluster |
| `TestAccResourceVMImport` | every Phase 2 PR adds coverage | todo | Import round-trip is mandatory per design |
| `TestAccResourceVMInitialization` | #11 | todo | Cloud-init |
| `TestAccResourceVMNetwork` | #10 | todo | network_device |
| `TestAccResourceVMClone` | — | dropped | clonedvm domain (D4) |
| `TestAccResourceVMVirtioSCSISingleWithAgent` | #7 + #13 | todo | Split |
| `TestAccResourceVMUpdateWhileStopped` | #6 | todo | power_state interaction |
| `TestAccResourceVMDisks` (+ 10 disk variants) | #7 | todo | All disk tests land in #7 |
| `TestAccResourceVMEFIDiskStorageMigration` | #9 + #19 | todo | — |
| `TestAccResourceVMHotplug` | #14 | todo | — |
| `TestAccResourceVMPool*` (6 pool variants) | #18 | todo | All pool tests land in #18 |
| `TestAccResourceVMRebootAfterCreationWithAgent` | #6 | todo | **rewritten** — no user-facing `reboot` |
| `TestAccResourceVMRebootAfterUpdateTPMStatePolicy` | #6 + #9 | todo | **rewritten** |
| `TestAccResourceVMRebootAfterUpdateCloudInitMovePolicy` | #6 + #11 | todo | **rewritten** |
| `TestAccResourceVMRebootAfterUpdateTemplatePolicy` | #6 | todo | **rewritten** |
| `TestAccResourceVMRebootAfterUpdateDiskMovePolicy` | #6 + #7 | todo | **rewritten** |
| `TestAccResourceVMRebootAfterUpdateDiskResizePolicy` | #6 + #7 | todo | **rewritten** |
| `TestAccResourceVMTemplateConversion` | — | done | Already covered in `proxmox_vm` |
| `TestAccResourceVMTpmState` | #9 | todo | — |
| `TestAccResourceVMCDROM` | — | done | Already in `cdrom/` |
| `TestAccDatasourceSDKVMNotFound` | — | waived | SDK datasource; unaffected |

### Unit tests (in `proxmoxtf/resource/vm/`)

| Test | Target PR | Status | Notes |
|---|---|---|---|
| `TestVMInstantiation` | — | dropped | FW schema is source of truth |
| `TestVMSchema` (top-level + disk subpkg) | — | dropped | FW schema is source of truth |
| `TestHotplugContains` | #14 | todo | Port as FW unit test |
| `Test_parseImportIDWIthNodeName` | — | done | Equivalent at `fwprovider/nodes/vm/resource.go:397` |
| `TestCPUType` | — | dropped | Long enum validator dropped (F27) |
| `TestMachineType` | — | dropped | Long enum validator dropped |
| `TestVmHostname` | — | done | `name` validator kept; FW test exists |
| `TestDiskOrderingDeterministic` | — | dropped | Map-keyed eliminates ordering |
| `TestDiskOrderingVariousInterfaces` | — | dropped | Same |
| `TestDiskDevicesEqual` | #7 | open | Port only if FW impl needs equality helper |
| `TestDiskUpdateSkipsUnchangedDisks` | #7 | todo | Covered by mandatory "plan empty diff" scenario |
| `TestImportFromDiskNotReimportedOnSizeChange` | #7 | todo | Behavior-specific |
| `TestDiskDeletionDetectionInGetDiskDeviceObjects` | #7 | todo | — |
| `TestDiskDeletionWithBootDiskProtection` | #7 | todo | — |
| `TestOriginalBugScenario` | #7 | todo | Regression test |
| `TestDiskSpeedSettingsPerDisk` | #7 | todo | — |
| `TestNetworkSchema` | — | dropped | FW schema is source of truth |

## Open questions remaining

> Design Q1–Q5 resolved during grilling. These are *new* open questions
> surfaced by the audit that need maintainer decision before the relevant
> PR can land.

| # | Question | Target PR | Status |
|---|---|---|---|
| OQ1 | Granular `timeout_*` controls (SDK) → adopt framework `timeouts` block only, or preserve per-phase granularity? Recommend: framework `timeouts` block, fold internal phase timeouts into parent op. | #6 / #19 | open |
| OQ2 | Read-only network attributes (`ipv4_addresses`, `ipv6_addresses`, `mac_addresses`, `network_interface_names`) — resource Computed, datasource only, or both? Recommend: datasource only. | #10 | open |
| OQ3 | Cloud-init `ip_config` — keep ordered list or map-keyed by interface name? Recommend: map-keyed, matches `network_device`. | #11 | open |
| OQ4 | Agent `timeout` / `wait_for_ip.{ipv4,ipv6}` — keep as PVE pass-through or fold into provider timeouts? Recommend: PVE pass-through. | #13 | open |

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
| `id` | Required | done | — |
| `node_name` | Required | done | — |
| `description` | Computed string | done | — |
| `name` | Computed string | done | — |
| `tags` | Computed stringset | done | — |
| `status` | Computed string (runtime mirror) | done | — |
| `cpu` | Computed SingleNested (DataSourceSchema) | done | — |
| `vga` | Computed SingleNested | done | — |
| `rng` | Computed SingleNested | done | — |
| `cdrom` | Computed MapNested | done | — |
| `memory` | Computed SingleNested | todo | #6 |
| `on_boot` | Computed bool | todo | #6 |
| `disk` | Computed MapNested | todo | #7 |
| `bios` | Computed string | todo | #8 |
| `machine` | Computed string | todo | #8 |
| `boot_order` | Computed list | todo | #8 |
| `efi_disk` | Computed SingleNested | todo | #9 |
| `tpm_state` | Computed SingleNested | todo | #9 |
| `scsi_hardware` | Computed string | todo | #9 |
| `network_device` | Computed MapNested | todo | #10 |
| `initialization` | Computed SingleNested | todo | #11 |
| `operating_system` | Computed SingleNested | todo | #12 |
| `smbios` | Computed SingleNested | todo | #12 |
| `agent` | Computed SingleNested | todo | #13 |
| `numa` | Computed MapNested | todo | #13 |
| `watchdog` | Computed SingleNested | todo | #13 |
| `acpi` | Computed bool | todo | #14 |
| `tablet_device` | Computed bool | todo | #14 |
| `keyboard_layout` | Computed string | todo | #14 |
| `kvm_arguments` | Computed string | todo | #14 |
| `vcpus` | Computed int | todo | #14 |
| `hotplug` | Computed list | todo | #14 |
| `parallel` | Computed MapNested | todo | #14 |
| `usb` | Computed MapNested | todo | #15 |
| `hostpci` | Computed MapNested | todo | #16 |
| `serial_device` | Computed MapNested | todo | #17 |
| `audio_device` | Computed MapNested | todo | #17 |
| `virtiofs` | Computed MapNested | todo | #17 |
| `startup` | Computed SingleNested | todo | #18 |
| `pool_id` | Computed string | todo | #18 |
| `protection` | Computed bool | todo | #18 |
| `hook_script_file_id` | Computed string | todo | #18 |
| `amd_sev` | Computed SingleNested | todo | #18 |

### Datasource-specific (not in resource)

| Attribute | Shape | Status | Notes |
|---|---|---|---|
| `ipv4_addresses` (per-network-slot) | Computed list of strings | todo | Per OQ2 resolution (resource vs datasource) |
| `ipv6_addresses` (per-network-slot) | Computed list of strings | todo | Per OQ2 |
| `mac_addresses` (per-network-slot) | Computed list of strings | todo | Per OQ2 |
| `network_interface_names` (per-network-slot) | Computed list of strings | todo | Per OQ2 |
