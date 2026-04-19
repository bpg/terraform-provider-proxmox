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
| 2026-04-19 | Audit Sections 2 (capabilities), 3 (legacy tests), 6 (Q5 power_state notes) added — only Section 4 (mitmproxy) remains | PR #1 | Continued after first checkpoint commit |
| 2026-04-19 | Section 4 complete — mitmproxy + qemu-server.git source cross-validated. Major finding: PVE returns absent for nearly all unset config fields; ~23 attributes drop Computed (Optional+Computed → Optional only) | PR #1 / #3 | Empirical mitmproxy data is much more aggressive than the design predicted |
| 2026-04-19 | Per-attribute classification in gap matrix updated from `open` to confirmed targets | PR #3 | Section 4 results applied to gap matrix |
| 2026-04-19 | Scrutiny pass — independent reviewer + 2 additional mitmproxy passes (cpu, vga+rng+cdrom). Section 4 corrected: `cpu.cores`/`cpu.sockets` keep Optional+Computed (PVE auto-populates), other 21 still drop Computed | PR #1 / #3 | Original wholesale claim was over-generalized; carve-out documented |
| 2026-04-19 | F39, F43 severity downgraded from `blocker` to `should-fix (PR-#3-blocker)` and `should-fix (PR-#6-blocker)` | PR #1 | Reviewer flagged that "blocker" without context implies production regression; clonedvm runs fine today |
| 2026-04-19 | F40-F42 redescribed as `NewValue` sentinels (not `FillCreateBody` sentinels — there is no FillCreateBody for memory) | PR #1 / #3 | PR #3 fix is two-part: rewrite NewValue + add FillCreateBody |
| 2026-04-19 | F32a, F36a added: `reflect.DeepEqual` zero-struct anti-pattern in vga + rng FillCreateBody/FillUpdateBody | PR #3 | ADR-008 should explicitly reject |
| 2026-04-19 | Section 2 expanded with sub-attribute tables for watchdog, agent, amd_sev, audio_device, numa | PR #1 | Reviewer flagged inconsistency vs disk/network/cloud-init coverage |
| 2026-04-19 | F12 line citation corrected (`:17` const line → `:28-34, 49-57` wrapper struct + MoveState) | PR #1 | Citation accuracy |
| 2026-04-19 | `cpu.units` validator decision changed from `keep` to `open question` — `Between(1, 262144)` rejects `0` which PVE allows on cgroup v2 | PR #3 | Reviewer flagged |
| 2026-04-19 | Second scrutiny pass: Section 4 expanded with "Implementation implications" subsection — `NewValue` functions must return `types.ObjectNull(...)` when underlying API device is nil; otherwise schema change creates permanent drift | PR #1 / PR #3 | First-pass scrutiny missed this implementation gap |
| 2026-04-19 | F44 expanded: also covers "always re-sends fields" issue (in addition to "never deletes"); both fixed by `stateValue` + `plan.Equal(state)` guard | PR #1 / PR #3 | Same root cause |
| 2026-04-19 | F47 explicitly resolved: keep `Default("cdrom")` as provider UX convenience (Section 4 verified file_id always present in PVE response when slot exists; default isn't duplicating PVE auto-populate) | PR #1 / PR #3 | Was previously left as "verify" |
| 2026-04-19 | Datasource schema files verified clean per CLAUDE.md "Datasource Schema Attributes" rule (no `Optional: true` in any datasource_schema.go); recorded in Section 1 Scope of scan | PR #1 | Audit coverage gap closed |

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
