<!-- markdownlint-disable MD013 MD060 -->

# Tracker: Migrate VM Resource to Plugin Framework (#1231)

**Issue:** [bpg/terraform-provider-proxmox#1231](https://github.com/bpg/terraform-provider-proxmox/issues/1231)
**Design:** [1231_DESIGN.md](1231_DESIGN.md)
**Started:** 2026-04-18

## At-a-glance

| Phase | PRs | Merged | In-flight | Blocked |
|---|---|---|---|---|
| 1 — Audit & Redesign | 5 | 0 | 1 | 0 |
| 2A — MVP setup + MVP | 2 | 0 | 0 | 0 |
| 2B — Boot + UEFI | 2 | 0 | 0 | 0 |
| 2C — Network/cloud-init/OS | 3 | 0 | 0 | 0 |
| 2D — Advanced hardware | 5 | 0 | 0 | 0 |
| 2E — Cluster + parity | 3 | 0 | 0 | 0 |
| Floating client-refactor | 0–1 | 0 | 0 | 0 |
| **Total** | **20 (+1)** | **0** | **1** | **0** |

**Currently active:** PR #1 (audit) on `chore/1231-audit-proxmox-vm`

## Status legend

| Code | Meaning |
|---|---|
| `todo` | Not started |
| `wip` | Branch open, in progress |
| `review` | PR open, awaiting review |
| `merged` | Merged to main |
| `blocked` | See blocker column |

## PRs

### Phase 1 — Audit & Redesign

| # | Title | Status | Branch | PR | Session | Blocker |
|---|---|---|---|---|---|---|
| 1 | `chore(vm2): audit proxmox_vm against ADRs 001–007` | wip | `chore/1231-audit-proxmox-vm` | — | `1231_SESSION_STATE.md` | — |
| 2 | `docs(adr): ADR-008 sub-block contract + ADR-004 amendment` | todo | — | — | — | — |
| 3 | `refactor(vm2): port cpu/vga/rng/cdrom/memory to ADR-008` | todo | — | — | — | — |
| 4 | `refactor(vm2)!: rename to proxmox_vm; delete MoveState` | todo | — | — | — | — |
| 5 | `refactor(vm2): ADR-005 error format sweep` | todo | — | — | — | — |

### Phase 2A — MVP setup + MVP

| # | Title | Status | Branch | PR | Session | Blocker |
|---|---|---|---|---|---|---|
| 6 | `feat(vm2): add memory + power_state + on_boot scalars` | todo | — | — | — | — |
| 7 | `feat(vm2): add disk map-keyed block` | todo | — | — | — | — |

### Phase 2B — Boot config + UEFI

| # | Title | Status | Branch | PR | Session | Blocker |
|---|---|---|---|---|---|---|
| 8 | `feat(vm2): add bios + machine + boot_order scalars` | todo | — | — | — | — |
| 9 | `feat(vm2): add efi_disk + tpm_state + scsi_hardware` | todo | — | — | — | — |

### Phase 2C — Network, cloud-init, OS

| # | Title | Status | Branch | PR | Session | Blocker |
|---|---|---|---|---|---|---|
| 10 | `feat(vm2): add network_device map-keyed block` | todo | — | — | — | — |
| 11 | `feat(vm2): add initialization (cloud-init)` | todo | — | — | — | — |
| 12 | `feat(vm2): add operating_system + smbios` | todo | — | — | — | — |

### Phase 2D — Advanced hardware

| # | Title | Status | Branch | PR | Session | Blocker |
|---|---|---|---|---|---|---|
| 13 | `feat(vm2): add agent + numa (with numa.enabled) + watchdog` | todo | — | — | — | — |
| 14 | `feat(vm2): add acpi + tablet + keyboard + kvm_args + vcpus + hotplug + parallel` | todo | — | — | — | — |
| 15 | `feat(vm2): add usb map-keyed block` | todo | — | — | — | — |
| 16 | `feat(vm2): add hostpci map-keyed block` | todo | — | — | — | — |
| 17 | `feat(vm2): add serial_device + audio_device + virtiofs (all map-keyed)` | todo | — | — | — | — |

### Phase 2E — Cluster + parity

| # | Title | Status | Branch | PR | Session | Blocker |
|---|---|---|---|---|---|---|
| 18 | `feat(vm2): add startup + pool_id + protection + hook_script + amd_sev` | todo | — | — | — | — |
| 19 | `feat(vm2): add migrate` | todo | — | — | — | — |
| 20 | `docs(vm2): feature parity report + SDK migration guide` | todo | — | — | — | — |

### Floating client-refactor slot

| # | Title | Status | Branch | PR | Session | Blocker |
|---|---|---|---|---|---|---|
| F | `refactor(code)!: cleanup proxmox/nodes/vms types` | todo (lands ad hoc in Phase 2) | — | — | — | — |

## Mid-execution decisions

Amendments to the design that surface during implementation. Folded into
`1231_DESIGN.md` at phase boundaries. Each entry: date, what changed, which
PR(s) affected, rationale.

| Date | Decision | Affects | Rationale |
|---|---|---|---|
| 2026-04-19 | Audit Section 4 (per-attribute classification) deferred to dedicated mitmproxy session — code-only sections (1, 5, 7) front-loaded | PR #1 | User chose option 2 (sequence audit by infra dependency) |

## Active blockers

| PR | Blocker | Owner | Opened | Resolved |
|---|---|---|---|---|
| — | — | — | — | — |

## Quick links

- [Design doc](1231_DESIGN.md)
- [ADR-003 (file org)](../docs/adr/003-resource-file-organization.md)
- [ADR-004 (schema conventions)](../docs/adr/004-schema-design-conventions.md)
- [ADR-005 (error handling)](../docs/adr/005-error-handling.md)
- [ADR-006 (testing)](../docs/adr/006-testing-requirements.md)
- [ADR-007 (resource rename migration)](../docs/adr/007-resource-type-name-migration.md)
- ADR-008 (sub-block contract) — to be created in PR #2

## Update protocol

- Bump status when a PR opens, enters review, merges, or blocks.
- Fill `Branch` when branch is created, `PR` when PR is opened, `Session` with the per-PR session-state filename when work begins.
- Recompute the at-a-glance counts on each status change.
- New mid-execution decisions append to the decisions log immediately; design doc absorbs them at phase boundary (end of Phase 1 / end of Phase 2).
- Blockers stay in the active table until resolved; on resolve, fill `Resolved` date and leave the row for history.
