<!-- markdownlint-disable MD060 -->

# Audit: `proxmox_vm` against ADRs 001–007 (#1231)

**Issue:** [bpg/terraform-provider-proxmox#1231](https://github.com/bpg/terraform-provider-proxmox/issues/1231)
**Design:** [1231_DESIGN.md](1231_DESIGN.md)
**Tracker:** [1231_TRACKER.md](1231_TRACKER.md)
**Scope:** PR #1 (Phase 1)
**Status:** In progress — frozen after PR #1 merges

## Methodology

| Source | What it produces |
|---|---|
| File:line scan of `fwprovider/nodes/vm/` | ADR compliance findings (Section 1) |
| Walk of `proxmoxtf/resource/vm/vm.go` | Capabilities inventory (Section 2) |
| Walk of `fwprovider/test/resource_vm_*.go` + `proxmoxtf/resource/vm/**/*_test.go` | Legacy test inventory (Section 3) |
| Per-attribute mitmproxy trace + reasoning | Per-attribute classification (Section 4) |
| Validator-by-validator review against ADR-004 amendment | Validator inventory (Section 5) |
| Open questions resolution from design | Q5 power_state resolution (Section 6) |
| Grep of `proxmox/nodes/vms` consumers | Shared-types catalog (Section 7) |

Severity tags used in Section 1:

| Tag | Meaning |
|---|---|
| `blocker` | Must fix before Phase 2 (or ship it as PR #3 alongside the contract port) |
| `should-fix` | Fix during Phase 1 (PRs #3–#5) |
| `nit` | Defer; capture in gap matrix as deferred-cleanup |

---

## Section 1 — ADR compliance findings

### Scope of scan

| Path | Files | Lines (approx) |
|---|---|---|
| `fwprovider/nodes/vm/` (top-level) | `resource.go`, `resource_schema.go`, `resource_short.go`, `model.go`, `datasource.go`, `datasource_schema.go`, `datasource_short.go`, `concurrency_test.go`, `datasource_test.go`, `resource_test.go` | ~1500 |
| `fwprovider/nodes/vm/cpu/` | `resource.go`, `resource_schema.go`, `model.go`, `datasource_schema.go`, `resource_test.go` | ~830 |
| `fwprovider/nodes/vm/cdrom/` | `resource.go`, `resource_schema.go`, `model.go`, `datasource_schema.go`, `resource_test.go` | ~340 |
| `fwprovider/nodes/vm/vga/` | `resource.go`, `resource_schema.go`, `model.go`, `datasource_schema.go`, `resource_test.go` | ~400 |
| `fwprovider/nodes/vm/rng/` | `resource.go`, `resource_schema.go`, `model.go`, `datasource_schema.go`, `resource_test.go` | ~400 |
| `fwprovider/nodes/vm/memory/` | `resource.go`, `resource_schema.go`, `model.go` | ~270 |
| `fwprovider/nodes/vm/network/` | (empty placeholder) | 0 |

### Findings

> Each finding: `area`, `severity`, `ADR`, `file:line`, `description`, `target PR`. Pre-resolved findings from the design grilling (P1–P7) are listed at the bottom of this document; the table below holds *new* findings.

#### Top-level package (`fwprovider/nodes/vm/`)

| # | Area | Severity | ADR | File:line | Description | Target PR |
|---|---|---|---|---|---|---|
| F1 | error msg | should-fix | 005 | `resource.go:116` | `"Failed to generate VM ID"` — should be `"Unable to Generate VM ID"` | #5 |
| F2 | error msg | should-fix | 005 | `resource.go:136` | `"VM does not exist after creation", ""` — empty detail; should be `"Unable to Create VM N"` with detail | #5 |
| F3 | error msg | should-fix | 005 | `resource.go:168, 177, 363, 365, 378` | Generic context strings (`"VM create"`, `"VM template conversion"`, `"VM stop/shutdown"`, `"VM delete"`) for `AddDiags`/`AddDiagsAsWarnings`. Format inconsistent with ADR-005. | #5 |
| F4 | error msg | should-fix | 005 | `resource.go:237` | `"VM does not exist after update", ""` — empty detail | #5 |
| F5 | error msg | should-fix | 005 | `resource.go:295` | `"Failed to update VM"` — should be `"Unable to Update VM N"` | #5 |
| F6 | error msg | should-fix | 005 | `resource.go:310, 353` | `"Failed to get VM status"` — should be `"Unable to Get VM %d Status"` (two callsites) | #5 |
| F7 | error msg | should-fix | 005 | `resource.go:326` | `"Cannot convert template back to VM"` — should be `"Unable to Convert Template Back to VM"` | #5 |
| F8 | error msg | should-fix | 005 | `resource.go:374` | `"Unable to Delete VM"` — correct prefix but missing VM ID per CLAUDE.md ADR-005 note | #5 |
| F9 | error msg | should-fix | 005 | `resource.go:425` | `fmt.Sprintf("VM %d does not exist on node %s", ...)` — should be `"Unable to Import VM N"` with detail | #5 |
| F10 | sub-block contract | should-fix | 008 | `resource.go:256–282` | Hand-rolled `del()` closure for top-level scalars (Description, Name, Tags). Should use `attribute.CheckDelete(plan, state, body, "FieldName")` per ADR-008. | #3 |
| F11 | code clarity | nit | — | `resource.go:432–435` | Comment `"not clear why this is needed, but ImportStateVerify fails without it"` for setting StopOnDestroy/PurgeOnDestroy/DeleteUnreferencedDisksOnDestroy on import. Either replace with proper explanation or fix root cause. | #4 (during rename pass) |
| F12 | resource type name | blocker | 007 | `resource_short.go:17, 33–46, 49–57` | `resourceShort` wrapper + `MoveState` exists. Per design D2, this collapses into single `NewResource()` returning `proxmox_vm` in PR #4. | #4 |
| F13 | datasource type name | blocker | 007 | `datasource_short.go:15, 29–46` | `datasourceShort` wrapper. Collapses in PR #4 alongside resourceShort. | #4 |
| F14 | datasource type name | should-fix | 007 | `datasource.go:41` | `req.ProviderTypeName + "_vm2"` — long datasource name still in use. Parallels F12/F13. | #4 |
| F15 | error msg | should-fix | 005 | `datasource.go:89–92` | `"VM Not Found"` summary + custom detail; should be `"Unable to Read VM N"` per ADR-005 | #5 |
| F16 | error msg | should-fix | 005 | `model.go:74, 82` | `"Unable to Read VM"`, `"Unable to Read VM Status"` — correct prefix, missing VM ID | #5 |
| F17 | error msg | should-fix | 005 | `model.go:87, 144` | `"VM ID is missing in status API response"` — wrong format; should be `"Unable to Read VM N"` with detail (two callsites) | #5 |
| F18 | error msg | should-fix | 005 | `model.go:131` | `"Failed to get VM"` — should be `"Unable to Read VM N"` | #5 |
| F19 | error msg | should-fix | 005 | `model.go:139` | `"Failed to get VM status"` — duplicate of F6 in resource read path | #5 |

#### `cpu/` sub-package

| # | Area | Severity | ADR | File:line | Description | Target PR |
|---|---|---|---|---|---|---|
| F20 | sentinel | should-fix | 004 | `cpu/resource.go:38–42` | (Confirms P2) Nil-substitution sentinel `Cores == nil → 1`. Drop per ADR-004 PVE-defaults rule. | #3 |
| F21 | sentinel | should-fix | 004 | `cpu/resource.go:44–48` | **Additional sentinel** beyond P2: `Sockets == nil → 1`. Drop per ADR-004. | #3 |
| F22 | sentinel | should-fix | 004 | `cpu/resource.go:50–60` | **Additional sentinel** beyond P2: `CPUEmulation == nil → Type "kvm64", Flags null`. Drop per ADR-004. | #3 |
| F23 | sub-block contract | should-fix | 008 | `cpu/resource.go:159–225` | (Confirms P4) Hand-rolled `del()` + `ShouldBeRemoved` + `IsDefined` cascades for 8 fields. Replace with `attribute.CheckDelete`. | #3 |
| F24 | bug | should-fix | 008 | `cpu/resource.go:190` | (Confirms P1) `IsDefined(plan.Sockets)` copy-paste bug in Limit branch — should be `IsDefined(plan.Limit)` | #3 |
| F25 | error msg | nit | 005 | `cpu/resource.go:250` | `"Cannot have CPU flags without explicit definition of CPU type", ""` — empty detail, summary phrasing inconsistent with ADR-005 format | #5 |
| F26 | sub-block contract | should-fix | 008 | `cpu/resource.go:227–255` | Special-case for `CPUEmulation` compound update (delType/delFlags switch). Doesn't fit standard CheckDelete shape — ADR-008 should call out compound types as a recognized pattern (or refactor) | #3 |
| F27 | validator | should-fix | 004 | `cpu/resource_schema.go:124–204` | (Confirms P5) Long enum validator (~75 CPU types) for `cpu.type`. Drop per ADR-004 enum rule. | #3 |
| F28 | classification | should-fix | 004 | `cpu/resource_schema.go:31–122` | All 10 CPU attributes are `Optional+Computed`. Per-attribute classification needed against PVE Read behavior (Section 4) | #3 |
| F29 | rehome | should-fix | — | `cpu/resource_schema.go:83` | (Confirms P3) `cpu.hotplugged` exists. Drop in #3, rehome as top-level `vcpus` in #14. | #3 (drop) / #14 (rehome) |
| F30 | rehome | should-fix | — | `cpu/resource_schema.go:101` | (Confirms P3) `cpu.numa` exists. Drop in #3, rehome as `numa.enabled` in #13. | #3 (drop) / #13 (rehome) |

#### `vga/` sub-package

| # | Area | Severity | ADR | File:line | Description | Target PR |
|---|---|---|---|---|---|---|
| F31 | sub-block contract | should-fix | 008 | `vga/resource.go:107–129` | Hand-rolled cascades for Clipboard/Type/Memory. Replace with CheckDelete pattern. | #3 |
| F32 | sub-block contract | should-fix | 008 | `vga/resource.go:101–105` | `vgaDevice` initialized from state (not zero) and mutated. Different mutation pattern from cpu — ADR-008 should normalize. | #3 |
| F33 | validator | should-fix | 004 | `vga/resource_schema.go:55–72` | Long enum validator (14 VGA types, version-evolving). Drop per ADR-004 enum rule. | #3 |
| F34 | classification | should-fix | 004 | `vga/resource_schema.go:33–82` | All 3 VGA attributes are `Optional+Computed`. Section 4 classification needed. | #3 |

#### `rng/` sub-package

| # | Area | Severity | ADR | File:line | Description | Target PR |
|---|---|---|---|---|---|---|
| F35 | sub-block contract | should-fix | 008 | `rng/resource.go:117–141` | Hand-rolled cascades for Source/MaxBytes/Period. Replace with CheckDelete pattern. | #3 |
| F36 | sub-block contract | should-fix | 008 | `rng/resource.go:115` | `rngDevice = createRNGDevice(state, true)` — same state-initialized mutation pattern as vga (F32) | #3 |
| F37 | int-zero trap | nit | 004 | `rng/resource.go:54, 59` | `MaxBytes.ValueInt64() != 0` and `Period.ValueInt64() != 0` use 0 as "not set" sentinel. Documented as `"Use 0 to disable limiting"` in schema, but the FillCreateBody never sends 0 to PVE — meaning user-set 0 is silently dropped. ADR-004 integer-0 trap. | #3 |
| F38 | classification | should-fix | 004 | `rng/resource_schema.go:31–67` | All 3 RNG attributes `Optional+Computed`. Section 4 classification needed. | #3 |

#### `memory/` sub-package

| # | Area | Severity | ADR | File:line | Description | Target PR |
|---|---|---|---|---|---|---|
| F39 | provider default | **blocker** | 004 | `memory/resource_schema.go:52, 66, 79` | `Default(...)` for `size=512`, `balloon=0`, `shares=1000` — these are PVE's own defaults. Direct violation of ADR-004 amendment PVE-defaults rule. | #3 |
| F40 | sentinel | should-fix | 004 | `memory/resource.go:37–40` | Nil-substitution sentinel: `DedicatedMemory == nil → 512`. Drop per ADR-004. | #3 |
| F41 | sentinel | should-fix | 004 | `memory/resource.go:43–48` | Sentinel: `FloatingMemory == nil → 0`. Drop per ADR-004. | #3 |
| F42 | sentinel | should-fix | 004 | `memory/resource.go:51–56` | Sentinel: `FloatingMemoryShares == nil → 1000`. Drop per ADR-004. | #3 |
| F43 | sub-block contract | **blocker** | 008 | `memory/resource.go` (no `FillCreateBody`) | `memory/` package has **no** `FillCreateBody`. ADR-008 contract (single-nested family) requires both Fill methods. Add or document why omitted. | #3 |
| F44 | sub-block contract | should-fix | 008 | `memory/resource.go:75–122` | `FillUpdateBody` doesn't accept `stateValue` and never deletes fields — only sets if present. Cannot remove `hugepages` or `keep_hugepages` once set. Diverges from ADR-008. | #3 |
| F45 | classification | should-fix | 004 | `memory/resource_schema.go:43–105` | After F39 `Default` removal, all 5 attributes need ADR-004 classification (Section 4). | #3 |

#### `cdrom/` sub-package (reference for map-keyed pattern)

| # | Area | Severity | ADR | File:line | Description | Target PR |
|---|---|---|---|---|---|---|
| F46 | regex bound | nit | 008 | `cdrom/resource_schema.go:33` | Slot regex `^(ide[0-3]\|sata[0-5]\|scsi([0-9]\|1[0-3]))$` — `scsi` only goes to 13, but PVE bound is `MAX_SCSI_DISKS=31`. Restrict-then-relax behavior is OK (additive), but design's slot-regex table specifies `scsi([0-9]\|[12][0-9]\|30)`. Tighten or relax to match. | #3 |
| F47 | provider default | nit | 004 | `cdrom/resource_schema.go:46` | `Default: stringdefault.StaticString("cdrom")` for `file_id`. Verify against PVE Read behavior — does PVE auto-default file_id to `"cdrom"` for an attached storage device, or is this a provider sentinel? | #3 |

### Summary by severity

| Severity | Count |
|---|---|
| blocker | 4 (F12, F13, F39, F43) |
| should-fix | 35 |
| nit | 6 |
| Total new findings | 45 |

Plus 7 pre-resolved findings (P1–P7) from grilling. Combined: 52 findings.

### Summary by target PR

| PR | New findings |
|---|---|
| #3 (port sub-packages) | 25 |
| #4 (rename) | 4 (F11, F12, F13, F14) |
| #5 (error sweep) | 16 |

---

## Section 2 — Capabilities inventory

Every attribute and block of the legacy SDK `proxmox_virtual_environment_vm`
classified as one of:

| Status | Meaning |
|---|---|
| `done` | Already implemented in `proxmox_vm` |
| `planned` | Target Phase 2 PR identified |
| `deliberately dropped` | Out of scope; document why |
| `open question` | Needs maintainer decision before PR can land |

### Top-level attributes

| Attribute | SDK source (file:line) | Status | Target PR | Notes |
|---|---|---|---|---|
| TBD | — | — | — | — |

### Top-level blocks

| Block | SDK source (file:line) | Status | Target PR | Notes |
|---|---|---|---|---|
| TBD | — | — | — | — |

### Disk-family attributes

| Attribute | SDK source (file:line) | Status | Target PR | Notes |
|---|---|---|---|---|
| TBD | — | — | — | — |

### Network-family attributes

| Attribute | SDK source (file:line) | Status | Target PR | Notes |
|---|---|---|---|---|
| TBD | — | — | — | — |

### Cloud-init attributes

| Attribute | SDK source (file:line) | Status | Target PR | Notes |
|---|---|---|---|---|
| TBD | — | — | — | — |

### Runtime / lifecycle attributes (`started`, `reboot`, `on_boot`, …)

| Attribute | SDK source (file:line) | Status | Target PR | Notes |
|---|---|---|---|---|
| TBD | — | — | — | — |

### Misc / advanced

| Attribute | SDK source (file:line) | Status | Target PR | Notes |
|---|---|---|---|---|
| TBD | — | — | — | — |

---

## Section 3 — Legacy test inventory

Every test in `proxmoxtf/resource/vm/**/*_test.go` and
`fwprovider/test/resource_vm_*.go` mapped to the user-visible behavior it
exercises and the Phase 2 PR that will own its port.

| Test name | Source (file:line) | Behavior under test | Tier | Target PR | Port status |
|---|---|---|---|---|---|
| TBD | — | — | — | — | — |

Port status values: `done`, `planned`, `dropped` (with reason), `open question`.

---

## Section 4 — Per-attribute classification (ADR-004 amendment)

Every existing `Optional+Computed` attribute in the current `proxmox_vm` (and
its sub-packages) reclassified per the new ADR-004 PVE-defaults rule.

> **Rule recap (ADR-004 amendment, drafted in PR #2):**
>
> | PVE Read behavior | Schema target |
> |---|---|
> | Auto-populates default value | `Optional + Computed` (no Default) |
> | Returns null/absent when unset | `Optional` only |
> | Provider-only attribute (no PVE counterpart) | `Optional + Default` |

### Verified via mitmproxy trace

| Attribute | Current schema | PVE Read behavior | Target schema | Action | Target PR |
|---|---|---|---|---|---|
| TBD | — | — | — | — | — |

### Pending mitmproxy verification

| Attribute | Current schema | Hypothesized PVE Read behavior | Notes |
|---|---|---|---|
| TBD | — | — | — |

---

## Section 5 — Validator inventory (ADR-004 enum rule)

Every validator currently in use, classified per the ADR-004 amendment enum
rule:

> **Rule recap:** Use `OneOf` for short, stable PVE enums (≤5 values, unlikely
> to extend). For long or version-evolving enums (CPU types, VGA types,
> machine types, BIOS modes, scsi_hardware), defer to PVE; use a regex
> validator only if format-only validation is meaningful.

### Top-level (`resource_schema.go`)

| Attribute | Validator | Type | Decision | Reason | Target PR |
|---|---|---|---|---|---|
| `name` | `stringvalidator.RegexMatches` DNS regex (`resource_schema.go:64`) | format | keep | DNS format check is format-only validation, meaningful and stable | — |

### `cpu/`

| Attribute | Validator | Type | Decision | Reason | Target PR |
|---|---|---|---|---|---|
| `affinity` | `stringvalidator.RegexMatches ^\d+[\d-,]*$` (`cpu/resource_schema.go:38`) | format | keep | Format-only validation is meaningful (catches typos before PVE call) | — |
| `architecture` | `stringvalidator.OneOf("aarch64", "x86_64")` (`cpu/resource_schema.go:50`) | short enum (2) | keep | Short, stable, no growth pressure | — |
| `cores` | `int64validator.Between(1, 1024)` (`cpu/resource_schema.go:59`) | range | keep | PVE-source bound | — |
| `flags` | `setvalidator.AlsoRequires(parent.type)` (`cpu/resource_schema.go:73`) | cross-attribute | keep | Frequently hit; compound `CPUEmulation` requires `Type` when `Flags` set (see F26) | — |
| `flags` (elements) | `stringvalidator.RegexMatches (.\|\s)*\S(.\|\s)*` + `LengthAtLeast(1)` (`cpu/resource_schema.go:75–79`) | format | keep | Non-empty/non-whitespace check | — |
| `hotplugged` | `int64validator.Between(1, 1024)` (`cpu/resource_schema.go:89`) | range | **drop** | Attribute itself being rehomed to top-level `vcpus` (P3); will get its own validator there | #3 (drop with attr) / #14 (replace) |
| `limit` | `float64validator.Between(0, 128)` (`cpu/resource_schema.go:98`) | range | keep | PVE-source bound | — |
| `sockets` | `int64validator.Between(1, 16)` (`cpu/resource_schema.go:113`) | range | keep | PVE-source bound | — |
| `type` | `stringvalidator.OneOf(...75 CPU types...)` (`cpu/resource_schema.go:125–204`) | **long enum** | **drop** | (Confirms F27/P5) Long, version-evolving — defer to PVE per ADR-004 enum rule | #3 |
| `units` | `int64validator.Between(1, 262144)` (`cpu/resource_schema.go:214`) | range | keep | PVE-source bound | — |

### `vga/`

| Attribute | Validator | Type | Decision | Reason | Target PR |
|---|---|---|---|---|---|
| `clipboard` | `stringvalidator.OneOf("vnc")` (`vga/resource_schema.go:47`) | enum (1) | keep | Single accepted value today; revisit if PVE adds | — |
| `type` | `stringvalidator.OneOf(...14 VGA types...)` (`vga/resource_schema.go:56–72`) | **long enum** | **drop** | (Confirms F33) Version-evolving (PVE adds qxl variants over time) | #3 |
| `memory` | `int64validator.Between(4, 512)` (`vga/resource_schema.go:80`) | range | keep | PVE-source bound | — |

### `rng/`

| Attribute | Validator | Type | Decision | Reason | Target PR |
|---|---|---|---|---|---|
| `source` | `stringvalidator.LengthAtLeast(1)` (`rng/resource_schema.go:45`) | format | keep | Trivial non-empty check | — |
| `max_bytes` | `int64validator.AtLeast(0)` (`rng/resource_schema.go:55`) | range | keep | Trivial bound; intersects F37 (int-zero trap is in `FillCreateBody`, not validator) | — |
| `period` | `int64validator.AtLeast(0)` (`rng/resource_schema.go:65`) | range | keep | Trivial bound; same int-zero trap context | — |

### `memory/`

| Attribute | Validator | Type | Decision | Reason | Target PR |
|---|---|---|---|---|---|
| `size` | `int64validator.Between(64, 268435456)` (`memory/resource_schema.go:54`) | range | keep | PVE-source bound | — |
| `balloon` | `int64validator.Between(0, 268435456)` (`memory/resource_schema.go:68`) | range | keep | PVE-source bound | — |
| `shares` | `int64validator.Between(0, 50000)` (`memory/resource_schema.go:81`) | range | keep | PVE-source bound | — |
| `hugepages` | `stringvalidator.OneOf("2", "1024", "any")` (`memory/resource_schema.go:95`) | short enum (3) | keep | Short, stable | — |

### `cdrom/`

| Attribute | Validator | Type | Decision | Reason | Target PR |
|---|---|---|---|---|---|
| (map keys) | `mapvalidator.KeysAre + stringvalidator.RegexMatches ^(ide[0-3]\|sata[0-5]\|scsi([0-9]\|1[0-3]))$` (`cdrom/resource_schema.go:31–37`) | slot regex | **tighten** | (Confirms F46) `scsi` only goes to 13; PVE bound is `MAX_SCSI_DISKS=31`. Update to design slot-regex table value `scsi([0-9]\|[12][0-9]\|30)`. | #3 |
| `file_id` | `stringvalidator.Any(OneOf("cdrom", "none"), validators.FileID())` (`cdrom/resource_schema.go:48–49`) | short enum + format | keep | Sentinel values + file ID format check | — |

### Summary

| Decision | Count | Rationale |
|---|---|---|
| keep | 19 | Short stable enums, cross-attribute, range bounds, format checks |
| drop | 2 | Long version-evolving enums (`cpu.type`, `vga.type`) |
| tighten | 1 | `cdrom` slot regex (`scsi` upper bound) |
| drop-with-attr | 1 | `cpu.hotplugged` validator dies with the attribute (rehome) |

All target PR #3 (sub-package port).

---

## Section 6 — Q5 resolution: `power_state` redesign

| Aspect | Resolution |
|---|---|
| Drop `started` (boolean) | Yes — replaced by `power_state` |
| Add `power_state` (string) | Values: `"running"`, `"stopped"`. Default `"running"`. |
| Add Computed `status` | For runtime drift visibility |
| Drop user-facing `reboot` | Provider decides reboot-vs-restart from pending changes |
| Keep `on_boot` (boolean) | Yes — corresponds to PVE "Start at boot" |
| Implementation PR | PR #6 |
| Audit deliverable | This section + an entry in Section 2 (capabilities) |

Implementation notes (deferred to PR #6 design refinement):

- TBD — exact API call sequence (status check → start/stop → poll)
- TBD — interaction with `stop_on_destroy`
- TBD — reboot-detection heuristic from pending changes diff

---

## Section 7 — Shared-types catalog (R3 mitigation)

Every type in `proxmox/nodes/vms` consumed by the SDK resource. Reviewers
of any client-touching PR check this table to flag cascade risk.

> **Allowed-changes rules (from R3):** add freely; internal cleanup freely;
> rename public fields with same-PR SDK callsite updates; remove public
> fields forbidden until post-Phase-2.

### Public types (consumed by both SDK and Framework — high cascade risk)

> Renaming a field on these types requires updating both `proxmoxtf/resource/vm/` and `fwprovider/nodes/{vm,clonedvm}/` callsites in the same PR.

| Type | SDK callsites | FW VM callsites | FW ClonedVM callsites | Notes |
|---|---|---|---|---|
| `vms.GetResponseData` | `vm.go:1936, 1937, 1950`, `disk/disk.go:34`, `network/network.go:109, 123, 134, 190` | `model.go:120` (read), every sub-package `NewValue` | `clonedvm/resource.go:508, 677, 791, 846, 902, 971` | API Read shape — touched everywhere |
| `vms.UpdateRequestBody` | `disk/disk.go:71, 94, 599`, `vm.go:1962` | `resource.go:252`, every sub-package `FillUpdateBody` | `clonedvm/resource.go:509, 587, 678` | API Update shape — touched everywhere |
| `vms.Client` | `disk/disk.go:63`, `network/network.go:188`, `vm.go:1980, 1998, 2019, 2043, 2076, 2122, 2190, 2246` | `resource.go:166, 250, 348, 443, 471` | `clonedvm/resource.go:508, 791, 1013, 1041` | RPC surface (`vms.Client`) — method renames cascade |
| `vms.CustomStorageDevice` / `vms.CustomStorageDevices` | `disk/disk.go:34, 61, 62, 159, 168, 171, 281, 302, 325, 374, 597–607`, `disk/disk_test.go` (>20 callsites), `vm.go:1937, 1950` | `cdrom/model.go:32–39`, `cdrom/resource.go:25` | `clonedvm/resource.go:688, 695, 902` | Storage shape — disk PR (#7) will touch heavily |
| `vms.CustomNetworkDevice` / `vms.CustomNetworkDevices` | `network/network.go:24, 26, 43, 109, 110, 123, 134` | (none yet — added in #10) | `clonedvm/resource.go:606, 611, 846, 971` | Network shape — network PR (#10) and ClonedVM share |
| `vms.ShutdownRequestBody` | `vm.go:2007` | `resource.go:452` | `clonedvm/resource.go:1022` | Shutdown shape — small surface |
| `vms.CloneRequestBody` | (clone is in clonedvm only on FW side) | `concurrency_test.go:88` (test) | `clonedvm/resource.go:463` | Clone shape |
| `vms.CreateRequestBody` | (none — SDK uses different patterns) | `resource.go:148`, sub-package `FillCreateBody` | (none — clones don't use Create) | Framework-only — safe to rename |

### Framework-only types (low cascade risk)

| Type | FW callsites | SDK callsites | Notes |
|---|---|---|---|
| `vms.CustomCPUEmulation` | `cpu/resource.go:123, 229, 253` | (none) | FW-only — safe to rename in #3 |
| `vms.CustomVGADevice` | `vga/resource.go:57, 72, 101, 131` | (none) | FW-only — safe to rename in #3 |
| `vms.CustomRNGDevice` | `rng/resource.go:48, 86, 143` | (none) | FW-only — safe to rename in #3 |

### Public type aliases / constants

| Identifier | File:line | Consumers | Notes |
|---|---|---|---|
| `vms.StorageInterfaces` | `proxmox/nodes/vms/...` (slice) | `proxmoxtf/resource/vm/disk/disk.go:281, 284` | SDK-only consumer; FW uses regex-based slot validation per ADR-008 |
| `vms.MoveDiskRequestBody` | `proxmox/nodes/vms/...` | `proxmoxtf/resource/vm/disk/disk.go:113` | SDK-only today; will become FW-shared when #19 (migrate) lands |
| `vms.ResizeDiskRequestBody` | `proxmox/nodes/vms/...` | `proxmoxtf/resource/vm/disk/disk.go:128, 362`, `clonedvm/resource.go:680, 685, 737` | Shared between SDK and ClonedVM today; #7 (disk MVP) will add FW VM consumer |
| `vms.WaitForIPConfig` | `proxmox/nodes/vms/...` | `proxmoxtf/resource/vm/network/network.go:192` | SDK-only today; #10 (network) likely adds FW consumer |

### Already-confirmed sensitive shapes (carried forward from prior FWK audits)

| Type | Field | Constraint |
|---|---|---|
| `*vms.GetResponseData.EFIDisk` | `*CustomEFIDisk`, json `efidisk0,omitempty` | Single-instance pointer — preserves D5 architectural-single decision |
| `*vms.GetResponseData.TPMState` | `*CustomTPMState`, json `tpmstate0,omitempty` | Single-instance pointer — same |
| `*vms.UpdateRequestBody.ToDelete(string) error` | method | Used by `attribute.CheckDelete`; do not rename without ADR-008 update |

---

## Gaps surfaced for the gap matrix

> Anything found here that needs Phase-2 tracking is mirrored into
> `1231_GAP_MATRIX.md`. The audit (this file) is frozen after PR #1; the gap
> matrix lives through Phase 2 and becomes PR #20's parity report.

---

## Pre-resolved findings (carried from grilling pass)

These are the findings that the design pass already locked in — the audit
confirms but does not re-litigate:

| # | Finding | Severity | Target PR | Source |
|---|---|---|---|---|
| P1 | `IsDefined(plan.Sockets)` copy-paste bug in cpu Limit branch (`fwprovider/nodes/vm/cpu/resource.go:190`) | should-fix | #3 | Direct read of file during grilling |
| P2 | Nil-substitution sentinel `cpu.Cores == nil → 1` violates ADR-004 PVE-defaults rule | should-fix | #3 | Read of `cpu/resource.go:38–42` |
| P3 | `cpu.numa` and `cpu.hotplugged` are SDK-inherited misnomers; should rehome | should-fix | #3 (drop) + #13 (rehome `numa.enabled`) + #14 (rehome `vcpus`) | Design D7 + ADR-008 single-vs-map rule |
| P4 | Hand-rolled `ShouldBeRemoved` + `IsDefined` cascades in cpu/, vga/, rng/, memory/ instead of `attribute.CheckDelete` | should-fix | #3 | Design ADR-008 contract |
| P5 | Long-enum validators (`cpu.type`, `vga.type`) duplicate PVE validation | should-fix | #3 | Design ADR-004 enum rule (Q4) |
| P6 | `proxmox_virtual_environment_vm2` long resource type name + `MoveState` exist | blocker | #4 | Design D2 |
| P7 | Error message format inconsistency vs ADR-005 `"Unable to [Action] VM %d"` | should-fix | #5 | Design Phase 1 PR #5 |
