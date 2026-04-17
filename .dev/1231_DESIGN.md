<!-- markdownlint-disable MD060 -->

# Design: Migrate VM Resource to Plugin Framework (#1231)

**Issue:** [bpg/terraform-provider-proxmox#1231](https://github.com/bpg/terraform-provider-proxmox/issues/1231)
**Status:** Draft
**Created:** 2026-04-17

## Summary

Complete the Plugin Framework migration of the VM resource by:

1. **Phase 1 — Audit & Redesign** (5 PRs): Inventory the current
   `proxmox_vm` (formerly `proxmox_virtual_environment_vm2`) resource against
   ADRs 001–007, codify the existing sub-block contract in a new ADR-008,
   amend ADR-004 with the PVE-defaults rule, port the five existing
   sub-packages (`cpu`, `vga`, `rng`, `cdrom`, `memory`) to the ADR-008
   contract, rename to `proxmox_vm`, and sweep error messages to the ADR-005
   format. Fully breaking changes allowed; the resource is explicitly
   experimental.

2. **Phase 2 — Feature Rollout** (15 PRs): Implement every missing attribute
   and block from the legacy SDK `proxmox_virtual_environment_vm` resource,
   starting with `memory` + `power_state` (MVP setup), then `disk` (the
   real MVP — VM with a bootable disk), then UEFI/network/cloud-init/OS,
   then advanced hardware, and finally cluster concerns + parity report.

Plus a **floating "client refactor" PR slot** in Phase 2 to bundle any
breaking cleanups to `proxmox/nodes/vms` types that surface during
implementation.

Total: 20 PRs (+ 1 floating).

## Goals

- Reach feature parity with the legacy SDK `proxmox_virtual_environment_vm`
  on the API surface.
- Codify the **existing Value-centric sub-block contract** (function-based,
  not method-based) in a new ADR-008 so future sub-packages follow the
  pattern uniformly. Include both single-nested and map-keyed signature
  families.
- Amend ADR-004 with the **PVE-defaults rule**: provider does not duplicate
  PVE's own defaults via schema `Default(...)`. Schema choice (`Optional`
  vs `Optional+Computed`) follows PVE's Read behavior, not user-intent
  guesses.
- Switch new device collections to the **map-keyed pattern** (keyed by PVE
  slot name: `scsi0`, `net0`, `virtio0`, `usb1`, …) — already established
  by `cdrom` and `clonedvm`. Eliminates reorder-churn when devices are
  inserted or removed. Slot regex per family is bounded by PVE source.
- Produce a formal audit of the existing resource against all ADRs with
  per-finding file:line citations and target-PR assignments.
- Produce a per-attribute classification table mapping every existing
  `Optional+Computed` attribute to its target schema shape under the new
  ADR-004 rule, plus a capabilities/test inventory and gap matrix that
  forms the parity checklist.
- Keep individual PRs small and independently reviewable.

## Non-Goals

- **`proxmox_cloned_vm` is out of scope** — its audit and feature catch-up
  is a separate follow-up epic. However, the five sub-packages (`cpu`,
  `vga`, `rng`, `cdrom`, `memory`) are **jointly owned** with `clonedvm`.
  ADR-008 changes that touch any sub-package must keep `clonedvm`
  acceptance tests green.
- **No SchemaVersion/UpgradeState helpers** for users of the experimental
  `proxmox_virtual_environment_vm2` name. Users re-import.
- **No automatic SDK→Framework state migration** for users of
  `proxmox_virtual_environment_vm`. Users migrate manually (aligned with
  ADR-007 Phase 3).
- **No graduation from "experimental" to stable** during this work.
  Decision is tracked outside this issue, post-Phase 2.
- **No shared-client refactors** in `proxmox/nodes/vms/` beyond additive
  fields needed by new attributes. Renames touching SDK callsites are
  allowed and ship in the same PR (or in the floating client-refactor
  slot). See R3.
- **No performance benchmarks** or plan/apply-time SLOs.

## Decisions Log

| # | Decision | Alternative rejected |
|---|---|---|
| D1 | Phase 1 = 5 PRs (audit; ADR docs; sub-block port; rename; error sweep). Audit produces 2 files: frozen `1231_AUDIT.md` and living `1231_GAP_MATRIX.md`. | 6+ PRs with 4-file audit split |
| D2 | Fully breaking changes allowed, no upgrade paths. MoveState deleted in PR #4. SDK→Framework migration guide ships once with PR #20's parity report — that's the moment users have an incentive to migrate. | Provide MoveState/UpgradeState helpers |
| D3 | Small-PR approach: many focused PRs, each independently reviewable | Single large branch / mega-PR |
| D4 | Scope = resource + datasource + long-name cleanup. `proxmox_cloned_vm` is out of scope **for new features** but five shared sub-packages (`cpu`, `vga`, `rng`, `cdrom`, `memory`) are jointly owned; their tests gate any port. | Touch `cloned_vm` features now / treat sub-packages as vm-only |
| D5 | New device families use map-keyed pattern (already in production via `cdrom` + `clonedvm`). Per-family slot regex tightened to PVE source bounds (see Map-Keyed Device Pattern §). `SingleNestedAttribute` only for **architecturally** single devices (`efidisk0`, `tpmstate0` — VM hardware model precludes multiples). Conventionally-single-today devices (`audio0`) use map-keyed with a one-key regex so future PVE expansion is additive, not a breaking schema-shape change. | User-chosen key + explicit `slot` attribute / SingleNested for everything currently single |
| D6 | PR #1 audit confirms there is no clone scaffolding in `proxmox_vm` (verified). The real "remove clone scaffolding" cleanup is the per-attribute audit of `Optional+Computed` flags vs the new ADR-004 rule (PR #3 applies the cleanup). | Search-and-destroy for non-existent clone code |
| D7 | Feature order: `memory` + `power_state` setup → `disk` (MVP) → UEFI → `network_device` → cloud-init → OS → advanced hardware → cluster concerns → parity report | network_device first / scalars first |
| D8 | Use `vm2` as commit scope for migration-epic implementation PRs; `docs(adr)` for new/amended ADRs. `vm2` documented in `commit-scopes` memory as transient (retires on #1231 close). | `vm` (collides with SDK changelog entries) |
| D9 | Audit produces 2 files: `1231_AUDIT.md` (ADR findings + capabilities inventory + legacy test inventory; frozen after PR #1) and `1231_GAP_MATRIX.md` (cross-reference, living through Phase 2; becomes PR #20's parity report) | 4 separate files / single monolithic file |

## Phase 1 — Audit & Redesign (5 PRs)

### PR breakdown

| # | Title | Type |
|---|---|---|
| 1 | `chore(vm2): audit proxmox_vm against ADRs 001–007` | Doc only |
| 2 | `docs(adr): ADR-008 sub-block contract + ADR-004 amendment for PVE-default rule` | Doc only |
| 3 | `refactor(vm2): port cpu/vga/rng/cdrom/memory to ADR-008; CheckDelete refactor; sentinel removal; bug fix; drop cpu.numa/cpu.hotplugged` | Code |
| 4 | `refactor(vm2)!: rename proxmox_virtual_environment_vm2 to proxmox_vm; delete MoveState` | Breaking |
| 5 | `refactor(vm2): ADR-005 error format sweep` | Code |

### PR #1: Audit (2 artifacts)

| Artifact | File | Lifecycle | Source |
|---|---|---|---|
| Audit | `.dev/1231_AUDIT.md` | Frozen after PR #1 | File:line scan of `fwprovider/nodes/vm/`; walk of `proxmoxtf/resource/vm/vm.go`; walk of `fwprovider/test/resource_vm_*.go` + `proxmoxtf/resource/vm/**/*_test.go` |
| Gap matrix | `.dev/1231_GAP_MATRIX.md` | Living through Phase 2; finalized as parity report | Cross-reference of audit + test inventory + capabilities inventory |

The audit (single doc, three sections) covers:

1. **ADR compliance findings** — severity-tagged: `blocker` (must fix
   before Phase 2), `should-fix` (fix during Phase 1), `nit` (defer).
2. **Capabilities inventory** — every SDK attribute classified as `done`
   (already in vm2), `planned` (target Phase 2 PR), `deliberately
   dropped`, or `open question`.
3. **Legacy test inventory** — every legacy test mapped to a behavior and
   target PR.

Plus two cross-cutting tables resolved in this audit:

- **Per-attribute classification** (per ADR-004 amendment): every existing
  `Optional+Computed` attribute classified by PVE Read behavior:
  `auto-populates default → keep O+C, no Default`,
  `returns null when unset → drop Computed`, or
  `provider-only attribute → keep Default`. Verified against mitmproxy
  traces.
- **Validator inventory**: every existing validator classified per the
  ADR-004 enum rule (short stable enum → keep; long version-evolving enum
  → drop).
- **Q5 resolution**: confirms `power_state` design and removal of `reboot`
  as a user-facing attribute (provider decides via pending-changes).
- **Shared-types catalog** (per R3): every `proxmox/nodes/vms` type
  consumed by the SDK resource, so reviewers can flag client refactor
  cascades.

### PR #2: ADR-008 + ADR-004 amendment

Document-only. Two ADRs land together because the contract (ADR-008) is
already in production code and ADR-004's PVE-default rule is the schema-
shape decision the contract relies on.

**ADR-008 — Sub-block Contract**

Codifies the *existing* Value-centric pattern from `cpu/`, `vga/`, `rng/`,
`cdrom/`, `memory/`. Two signature families:

```go
// Single-nested family (cpu, vga, rng, memory):
type Value = types.Object

func ResourceSchema() schema.Attribute
func DataSourceSchema() dsschema.Attribute
func NullValue() Value

func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value
func FillCreateBody(ctx context.Context, planValue Value, body *vms.CreateRequestBody, diags *diag.Diagnostics)
func FillUpdateBody(ctx context.Context, planValue, stateValue Value, body *vms.UpdateRequestBody, diags *diag.Diagnostics)

// Map-keyed family (cdrom, future: disk, network_device, usb, hostpci, serial, parallel, virtiofs):
type Value = types.Map  // (other functions identical signatures)
```

Rules codified in the ADR:

- Internal `Model` struct holds the framework Values; not exported.
- Use `attribute.StringPtrFromValue`, `Int64PtrFromValue`,
  `Float64PtrFromValue`, `CustomBoolPtrFromValue` in `FillCreateBody`. No
  manual `ValueXxxPointer()` with `IsUnknown` cascades.
- Use `attribute.CheckDelete(plan, state, body, "ApiFieldName")` for
  scalar update-time deletion. Replaces hand-rolled `ShouldBeRemoved` +
  `IsDefined` cascades. The helper calls `body.ToDelete(...)` directly;
  no `*[]string` parameter passed through call graphs.
- For map-keyed sub-packages: use `utils.MapDiff(plan, state)` to compute
  `(toCreate, toUpdate, toDelete)`. Slot deletes go to `body.Delete = append(...)`
  — distinct from scalar field deletes which go via `body.ToDelete(...)`.
- `NewValue` (FromAPI direction) uses `types.*PointerValue()` so `nil`
  maps to null. Provider does **not** substitute PVE default values for
  nil API fields (this is the ADR-004 PVE-default rule).
- File layout per package per ADR-003: `resource_schema.go`,
  `datasource_schema.go`, `model.go`, plus colocated test files.
- **Single-vs-map rule:**
  - Use `SingleNestedAttribute` only when PVE's underlying hardware
    model **architecturally** precludes multiple instances (e.g.,
    `efidisk0` — one EFI variable store per VM is firmware-defined;
    `tpmstate0` — TPM spec is single-instance per system).
  - Use `MapNestedAttribute` (map-keyed by PVE slot name) for everything
    else, **including devices PVE currently limits to one slot** (e.g.,
    `audio0`). Tighten the slot regex to PVE's current bounds; relax in
    a future additive PR if PVE expands. This keeps schema-shape changes
    out of the breaking-change menu when PVE grows a feature (e.g., a
    pending feature request to support multiple cdroms is already
    additive because cdrom is map-keyed; audio_device should follow the
    same pattern).
- Sub-block names normally correspond to a single PVE concept. Exception:
  `cpu` historically includes `cpu.numa` (PVE: `numa=1`) and
  `cpu.hotplugged` (PVE: `vcpus=N`); these are SDK-inherited and being
  relocated in PR #3 (drop) + PR #13 (`numa.enabled`) + PR #14 (`vcpus`).
  No "virtual sub-block" pattern in ADR-008.

**ADR-004 amendment — PVE-defaults rule + enum/cross-attribute rules**

Three new sections:

1. **Provider Defaults vs PVE Defaults.** Provider does not duplicate
   PVE's defaults via schema `Default(...)`. Classification table:

   | PVE Read behavior | Schema | Examples |
   |---|---|---|
   | Auto-populates default value | `Optional + Computed` (no Default) | `cpu.cores`, `cpu.sockets`, `cpu.type`, `cpu.numa`, `cpu.units`, `cpu.limit` |
   | Returns null/absent when unset | `Optional` only | `cpu.affinity`, `cpu.flags`, `description` |
   | Provider-only attribute (no PVE counterpart) | `Optional + Default` | `purge_on_destroy`, `stop_on_destroy`, `delete_unreferenced_disks_on_destroy` |

   PVE Read behavior must be verified empirically (mitmproxy + audit
   table) per attribute.

2. **Enum validators.** Use `OneOf` for short, stable PVE enums (≤5
   values, unlikely to extend). For long or version-evolving enums (CPU
   types, VGA types, machine types, BIOS modes, scsi_hardware), defer to
   PVE; use a regex validator only if format-only validation is
   meaningful.

3. **Cross-attribute constraints.** Document in `MarkdownDescription` by
   default. Promote to plan-time validators only when (a) the constraint
   is hit frequently in support issues, or (b) PVE's apply-time error is
   unhelpful.

### PR #3: Port cpu/vga/rng/cdrom/memory to ADR-008

Mechanical application of the contract across **five** sub-packages
(memory included; it's currently shared with `clonedvm` and used by
`proxmox_vm` from PR #6 onward). Also:

- Replaces hand-rolled `ShouldBeRemoved` + `IsDefined` cascades with
  `attribute.CheckDelete(...)` (cpu's ~100 lines become ~15).
- Fixes the `IsDefined(plan.Sockets)` copy-paste bug in cpu's Limit branch
  at `fwprovider/nodes/vm/cpu/resource.go:190`.
- Removes nil-substitution sentinels (e.g., `cpu.Cores == nil → 1`) per
  ADR-004 PVE-defaults rule.
- Applies the per-attribute classification table from PR #1 audit:
  downgrade `Optional+Computed` to `Optional` only where PVE returns
  null; drop long-enum validators (`cpu.type`, `vga.type`, etc.).
- **Drops `cpu.numa` and `cpu.hotplugged`** (rehomed in PR #13 and PR
  #14). Acknowledge the temporary regression in PR body and gap matrix.
- Adds unit tests on Model methods per §Testing below.
- **Definition-of-done includes `clonedvm` acceptance tests passing.**

**Size note.** Combined diff expected at 500–700 LOC net. Within the
~500 LOC threshold's ballpark; acceptable because the PR applies a
single pattern to five sub-packages (reviewers read it as one pattern
with five applications). If any sub-package port surfaces a contract
revision, split that sub-package out as a separate PR.

**This PR is the reality-check on ADR-008.** Five sub-packages, one of
them (`memory/`) shared with a second consumer. If the contract doesn't
fit, it gets revised here before Phase 2 starts.

### PR #4: Rename to `proxmox_vm`

Breaking. The `proxmox_virtual_environment_vm2` long name goes away
entirely. `resourceShort` / `NewShortResource` collapse into a single
`NewResource()` returning `proxmox_vm`. **`MoveState` is deleted** —
users on `_vm2` re-import. Example HCL for the old name is also removed.

Also in this PR:

- Move/rename `examples/resources/proxmox_virtual_environment_vm2/` to
  `examples/resources/proxmox_vm/`.
- Update the `go:generate cp` line in `main.go` that wires the example
  file into docs generation (per memory: `docs-generation-pipeline.md`).
- Delete `docs/resources/virtual_environment_vm2.md` if present.
- Verify `proxmox_cloned_vm` still compiles (defensive; the load-bearing
  joint-ownership constraint is enforced by PR #3, not this rename).

### PR #5: ADR-005 error format sweep

Mechanical: `"Unable to [Action] VM %d"` everywhere across `vm/`. ~50
LOC of error-string changes.

**Schema description stays "experimental, MUST NOT use in production"**
(per D2's no-quarter framing).

**No file split.** `fwprovider/nodes/vm/resource.go` (currently 487
lines, projected ~700 by end of Phase 2) stays as one file — well below
the codebase's existing precedent (`fwprovider/nodes/clonedvm/resource.go`
at 1057 lines, kept as one file). Splitting becomes a separate
discussion if/when concrete merge-conflict pain surfaces.

## Phase 2 — Feature Rollout (15 PRs)

### Phase 2A — MVP setup + MVP (2 PRs)

| # | Title |
|---|---|
| 6 | `feat(vm2): add memory + power_state + on_boot scalars` |
| 7 | `feat(vm2): add disk map-keyed block` |

**Milestone after PR #7: MVP** — VM with a bootable disk, end-to-end.

### Phase 2B — Boot config + UEFI (2 PRs)

| # | Title |
|---|---|
| 8 | `feat(vm2): add bios + machine + boot_order scalars` |
| 9 | `feat(vm2): add efi_disk + tpm_state + scsi_hardware` |

**Milestone after PR #9: First credible SDK replacement** — full
UEFI/SeaBIOS VM with bootable disk.

### Phase 2C — Network, cloud-init, OS (3 PRs)

| # | Title |
|---|---|
| 10 | `feat(vm2): add network_device map-keyed block` |
| 11 | `feat(vm2): add initialization (cloud-init)` |
| 12 | `feat(vm2): add operating_system + smbios` |

### Phase 2D — Advanced hardware (5 PRs)

| # | Title |
|---|---|
| 13 | `feat(vm2): add agent + numa (with numa.enabled) + watchdog` |
| 14 | `feat(vm2): add acpi + tablet_device + keyboard_layout + kvm_arguments + vcpus + hotplug + parallel` |
| 15 | `feat(vm2): add usb map-keyed block` |
| 16 | `feat(vm2): add hostpci map-keyed block` |
| 17 | `feat(vm2): add serial_device + audio_device + virtiofs (all map-keyed)` |

`vcpus` in PR #14 is the rehomed `cpu.hotplugged`. `numa.enabled` in
PR #13 is the rehomed `cpu.numa`.

PRs #15–#17 are independent applications of established patterns. Each
still touches top-level `resource.go` to wire the new block, so merge
conflicts are unavoidable with many branches open at once — realistic
parallelism is 2 branches concurrently. Conflicts are trivially
resolvable on non-overlapping lines.

### Phase 2E — Cluster + parity (3 PRs)

| # | Title |
|---|---|
| 18 | `feat(vm2): add startup + pool_id + protection + hook_script_file_id + amd_sev` |
| 19 | `feat(vm2): add migrate` |
| 20 | `docs(vm2): feature parity report + SDK migration guide` |

**Milestone after PR #20: Feature parity reached.** Graduation-from-
experimental decision point (tracked outside this issue). PR #20 also
ships a one-shot SDK→Framework migration guide.

### Floating client-refactor PR

If breaking refactors to `proxmox/nodes/vms` types surface during Phase
2 implementation (poorly named fields, type mismatches, etc.), bundle
them into a single PR rather than scattering across feature PRs. Updates
both Framework and SDK callsites in the same PR. Lands at any point in
Phase 2.

### Escape hatch

Any combined PR may be split during execution if it exceeds ~500 LOC net
change, if a specific sub-block blocks review, or if the audit reveals
deeper issues than expected. Splitting in-flight is cheap;
over-splitting upfront is wasteful.

### PR dependency graph

| PR | Depends on |
|---|---|
| 1 | — |
| 2 | 1 (audit decisions feed ADR-008 + ADR-004 amendment) |
| 3 | 2 |
| 4 | 3 (rename after sub-blocks are on the contract) |
| 5 | 3 (error sweep can land before or after rename) |
| 6 | 5 (memory wired in via new contract; needs Phase 1 done) |
| 7 | 6 (disk MVP needs memory + power_state) |
| 8 | 7 |
| 9 | 8 (UEFI builds on bios/machine) |
| 10 | 7 (network needs map-keyed pattern; cdrom + disk both established it) |
| 11 | 7 + 10 (cloud-init needs disks; ipconfig needs network) |
| 12 | 11 |
| 13 | 7 (independent feature block) |
| 14 | 7 (independent; `vcpus` rehome is independent of numa rehome) |
| 15–17 | 7 (independent map-keyed blocks) |
| 18 | 9 (pool/protection touch fully-configured VMs) |
| 19 | 13–17 (migrate semantics depend on full device set) |
| 20 | 19 |

Anything else may land in any order that CI and review capacity permit.

## Map-Keyed Device Pattern

### Status

Already in production via `cdrom` (`fwprovider/nodes/vm/cdrom/`) and
`clonedvm` (two `MapNestedAttribute` blocks). ADR-008 codifies the
pattern from existing code; PR #7 (`disk`) is the first *new* application
in `proxmox_vm`.

### Schema shape (cdrom as reference)

```hcl
cdrom = {
  ide0 = { file_id = "local:iso/debian.iso" }
  ide2 = { file_id = "cdrom" }
}
```

```go
func ResourceSchema() schema.Attribute {
    return schema.MapNestedAttribute{
        Description: "...",
        Optional:    true,
        Computed:    true,  // map-level Computed iff PVE auto-populates
        Validators: []validator.Map{
            mapvalidator.KeysAre(stringvalidator.RegexMatches(
                regexp.MustCompile(`^(ide[0-3]|sata[0-5]|scsi([0-9]|1[0-3]))$`),
                "...",
            )),
        },
        NestedObject: schema.NestedAttributeObject{
            Attributes: map[string]schema.Attribute{ /* ... */ },
        },
    }
}
```

**Map-level Computed rule:** `Optional+Computed` iff PVE auto-populates a
default set of devices (cdrom-style — PVE reports IDE devices it
auto-attaches). `Optional`-only otherwise (network-device-style — PVE
does not auto-attach network devices).

### Diff semantics in FillUpdateBody

Use `utils.MapDiff(plan, state)` for the three-way diff:

- Slot in plan AND state with different spec → update the slot
- Slot in plan NOT in state → new device, add to body
- Slot in state NOT in plan → removed device, append slot key to
  `body.Delete` (distinct from `body.ToDelete(...)` which is for scalar
  field deletes)

Order-independent. Inserting `net2` does not churn `net0` or `net1`.
Removing `net1` does not renumber the rest.

### Per-family slot regex (PVE-source bounds)

Bounds verified from `qemu-server.git` Perl source:

| Device | PVE constant | Regex |
|---|---|---|
| `network_device` | `MAX_NETS=32` | `^net([0-9]\|[12][0-9]\|3[01])$` |
| `disk` (combined) | `MAX_IDE_DISKS=4`, `MAX_SATA_DISKS=6`, `MAX_SCSI_DISKS=31`, `MAX_VIRTIO_DISKS=16` | `^(ide[0-3]\|sata[0-5]\|scsi([0-9]\|[12][0-9]\|30)\|virtio([0-9]\|1[0-5]))$` |
| `usb` | `MAX_USB_DEVICES=14` | `^usb([0-9]\|1[0-3])$` |
| `hostpci` | `MAX_HOSTPCI_DEVICES=16` | `^hostpci([0-9]\|1[0-5])$` |
| `numa` | `MAX_NUMA=8` | `^numa[0-7]$` |
| `serial_device` | `MAX_SERIAL_PORTS=4` | `^serial[0-3]$` |
| `parallel` | `MAX_PARALLEL_PORTS=3` | `^parallel[0-2]$` |
| `virtiofs` | `max_virtiofs()` (verify in audit) | `^virtiofs[0-N]$` |
| `audio_device` | `audio0` only today (Perl source `$id //= 0`) | `^audio0$` (relax if PVE adds slots) |

### Single-instance devices (SingleNestedAttribute, not map-keyed)

Per ADR-008 single-vs-map rule — only **architecturally** single devices
(VM hardware model precludes multiples):

- `efi_disk` → PVE `efidisk0` (one EFI variable store per VM,
  firmware-defined)
- `tpm_state` → PVE `tpmstate0` (TPM spec is single-instance per system)

`audio_device` is *conventionally* single today but uses the map-keyed
pattern with a one-key regex (see table above) so future PVE expansion
to `audio1+` is additive, not a breaking schema-shape change. Same
forward-looking choice we already get for free with `cdrom` (a pending
PVE feature request adds support for multiple cdroms; cdrom is already
map-keyed, so accommodating it requires no schema change).

## Testing Strategy

### Required per PR

| Layer | Requirement | Location |
|---|---|---|
| Unit tests | Sub-package functions with non-trivial logic (`NewValue` sentinels, `FillUpdateBody` diff, `utils.MapDiff` integration for map-keyed) | `fwprovider/nodes/vm/<sub>/model_test.go` |
| Acceptance tests | Full CRUD + import round-trip per sub-block | `fwprovider/nodes/vm/<sub>/resource_test.go` |
| Datasource coverage | Datasource has equivalent functional coverage (per ADR-006) for every attribute/block added to the resource | `fwprovider/nodes/vm/datasource_test.go` |
| Functional coverage | Per ADR-006 — every distinct user-visible behavior has a scenario | Same |
| API verification | `/bpg:debug-api` (mitmproxy) at least once per sub-block PR, referenced in PR body | Manual |
| Docs regeneration | `make docs` run after any schema change; regenerated `docs/resources/proxmox_vm.md` committed in the same PR | `docs/`, `templates/` |

### Tier classification

- `//testacc:tier=medium` for most sub-block tests
- `//testacc:tier=heavy` for tests requiring cloud images or shared templates
- `//testacc:resource=vm` on all new tests

### Mandatory map-keyed device scenarios

Every map-keyed device sub-block **must** include:

| Scenario | What it proves |
|---|---|
| Create with 2+ slots | Multi-device support |
| Update adds a new slot | Insert without renumbering |
| **Update removes a middle slot** | **Reorder immunity — the payoff of maps** |
| Update renames a slot | Delete-then-add (slot name is identity) |
| Update modifies a slot's field | In-place attribute update |
| Import round-trip; **plan after import shows empty diff** | State matches config after import; Read populates correctly |
| **Apply, re-plan with same config — assert empty diff** | Plan stability; ADR-004 classification correct |
| **Plan with out-of-range slot key fails with validator error** | Tight slot regex actually fires |
| **Mixed-interface map** (disk only) | Multi-prefix regex handled correctly |
| **Rename across interface families** (disk only: `scsi0` → `virtio0`) | Semantic delete+create, not in-place |

### Unit tests on sub-block functions

Mandatory coverage:

- `NewValue` — every nil-API-pointer path (most should leave field null
  per ADR-004; the ones that don't are sentinels documented per field)
- `FillUpdateBody` — every `CheckDelete` path
- Map-keyed `FillUpdateBody` — full three-way diff matrix
  (add/update/delete × in-plan/in-state)

These are pure Go, no Terraform harness, milliseconds instead of minutes.

### Legacy test port-as-you-go

Each Phase 2 PR's definition-of-done includes porting the relevant rows
from the audit's legacy test inventory section. Reviewers block the PR
if scenarios targeted for it are missing.

Ports are spec-level, not byte-level:

- Schema differences (map-keyed devices, renamed attributes) mean HCL is
  rewritten; scenarios carry over.
- Validation-layer unit tests from `proxmoxtf/resource/vm/*_test.go` are
  ported to Framework validators where the validator survives the new
  ADR-004 enum rule; otherwise dropped (long enums are no longer
  validated client-side).
- Domain-layer unit tests on `proxmox/nodes/vms/` parsing may apply
  as-is with import path updates.

### Phase 1 test expansion

PR #3 (cpu/vga/rng/cdrom/memory port) closes the test gaps surfaced in
PR #1's audit for those five sub-packages. cpu specifically needs tests
for:

- `Cores = nil → null` (post-sentinel-removal verification)
- `CPUEmulation.Type` without `Flags`
- The `Type` + `Flags` coupling on update

### Non-goals

- No mocked-API unit tests on the top-level resource.
- No automated SDK-vs-Framework behavioral comparison tests; PR #20 is
  hand-written.
- No performance benchmarks.

## File Layout for Phase 1 Artifacts

Active work (through end of Phase 2):

```text
.dev/
├── 1231_DESIGN.md                      # this document
├── 1231_AUDIT.md                       # ADR findings + capabilities + tests (frozen after PR #1)
├── 1231_GAP_MATRIX.md                  # cross-reference, living doc, becomes parity report
└── 1231_{PR_NUMBER}_SESSION_STATE.md   # per-PR state files
```

After issue close: all moved to `.dev/archive/1231_*.md`, matching the
existing `FWK_AUDIT_*.md` pattern already in the archive.

## Risks and Mitigations

| # | Risk | Mitigation |
|---|---|---|
| R1 | ADR-008 contract ossifies wrong — every subsequent PR compounds the mistake | PR #3 ports 5 existing sub-blocks at once (one of them — `memory/` — shared with `clonedvm`). If the contract doesn't fit all five, revise it before Phase 2. |
| R2 | Map-keyed device semantics have Framework edge cases | Already known: `cdrom` and `clonedvm` (2 instances) have map-keyed in production. ADR-008 codifies from this existing code. PR #7 (`disk`) is the first new application, not a pattern-establishing PR. |
| R3 | Shared API client changes cascade — structs shared with SDK resource | **Allowed-changes table:** add freely; internal cleanup freely; rename public fields with same-PR SDK callsite updates; remove public fields forbidden until post-Phase-2. **Enforcement:** PR #1 audit catalogs all `proxmox/nodes/vms` types consumed by SDK. **Floating client-refactor PR slot** in Phase 2 bundles genuine breaking cleanups. |
| R4 | Long-running Phase 2 branches conflict on `resource.go` | Conflicts are trivially resolvable: each sub-block adds an import + a Schema entry + a `FillCreateBody`/`FillUpdateBody` call, on non-overlapping lines. Realistic parallelism is 2 concurrent branches. |
| R5 | Acceptance-test cost grows with PR count | Unit tests on sub-package functions carry most logic. Acceptance focuses on PVE-boundary behavior. CI gates on full suite at merge only. |
| R6 | Audit PR #1 reveals more rot than expected | Report is doc-only; severity-tagged findings let us defer nits. Breaking-changes-allowed means aggressive fixes OK. |
| R7 | `memory/` package status | Resolved: `memory/` is shared with `clonedvm`, **not orphaned**. Wired into `proxmox_vm` in Phase 2 PR #6 (memory + power_state + on_boot) using the new ADR-008 contract from PR #3. |
| R8 | Joint sub-package ownership constraint with `proxmox_cloned_vm` | Five sub-packages (cpu/vga/rng/cdrom/memory) are jointly owned with `clonedvm`. ADR-008 changes touching any sub-package must keep `clonedvm` acceptance tests green. PR #3's definition-of-done lists `clonedvm` tests passing. |

## Open Questions

Resolved in PR #1's audit before Phase 2 starts. Each needs a maintainer
decision, not more research.

| # | Question | Resolution |
|---|---|---|
| Q1 | Does the `proxmox_vm` datasource expose devices as map-keyed attributes too? | Yes — symmetric schemas, with carve-out: provider-only behavior attributes (`purge_on_destroy`, `stop_on_destroy`, `delete_unreferenced_disks_on_destroy`, `template`, `power_state`) are omitted from the datasource. Lookup keys (`id`, `node_name`) are Required. |
| Q2 | Does ADR-008 formally document the "virtual sub-block" pattern? | **No.** Pattern dropped. cpu's smell (containing `cpu.numa` and `cpu.hotplugged`) is one-off SDK debt, being relocated to `numa.enabled` (PR #13) and `vcpus` (PR #14). ADR-008 instead states: sub-block names normally correspond to a single PVE concept. |
| Q3 | Expose PVE runtime status fields (vcpu count, uptime) via a separate read-only block? | No — out of scope; separate datasource if demand arises |
| Q4 | Validator strictness | **Strict where duplication of PVE's validation is cheap and stable** (e.g., short stable enums, slot regexes per PVE-source bounds). **Drop long version-evolving enums** (cpu.type, vga.type, machine, bios, scsi_hardware, audio_device.driver) — defer to PVE. **Cross-attribute constraints** are doc-only by default; promoted to validators only when frequently hit or PVE error is unhelpful. Codified in ADR-004 amendment. |
| Q5 | Keep current `started`-field semantics or redesign? | **Redesign.** Replace `started` with **`power_state`** (`"running"` / `"stopped"`, default `"running"`). Add Computed `status` for runtime drift. Drop `reboot` as user-facing — provider decides from pending changes. `on_boot` stays (PVE "Start at boot"). Resolved in PR #1 audit; implemented in PR #6. |

## Success Criteria

### Phase 1 complete

- PR #1 audit + gap matrix merged
- PR #2 ADR-008 + ADR-004 amendment merged
- PRs #3–#5 merged; all five existing sub-packages on the ADR-008 contract
- `clonedvm` acceptance tests green throughout
- `make lint` and `make test` green
- `./testacc --resource vm` passes
- Resource type name `proxmox_virtual_environment_vm2` no longer exists
- All error messages follow ADR-005 `"Unable to [Action] VM %d"` format

### Phase 2 complete

- Every row of the audit's capabilities inventory marked `done` or
  `deliberately dropped` in the gap matrix
- Every row of the audit's legacy test inventory ported or explicitly waived
- PR #20 parity report + SDK migration guide merged, zero outstanding rows
- Issue #1231 checklist fully ticked
- `./testacc --resource vm --tier all` passes

### Out of scope (tracked separately)

- Graduation from "experimental" to stable
- `proxmox_cloned_vm` audit — separate epic
- SDK resource deprecation messaging — tied to v1.0 release planning

## Dependencies

- Proxmox VE minimum version == whatever SDK resource supports today
  (confirmed in PR #1 audit)
- Go 1.26+ (already in use)
- Plugin Framework: current `go.mod` version; no upgrades planned
- mitmproxy for `/bpg:debug-api` verification (standard practice)

## References

- [Issue #1231](https://github.com/bpg/terraform-provider-proxmox/issues/1231)
- [ADR-001: Use Plugin Framework](../docs/adr/001-use-plugin-framework.md)
- [ADR-003: Resource File Organization](../docs/adr/003-resource-file-organization.md)
- [ADR-004: Schema Design Conventions](../docs/adr/004-schema-design-conventions.md)
- [ADR-005: Error Handling](../docs/adr/005-error-handling.md)
- [ADR-006: Testing Requirements](../docs/adr/006-testing-requirements.md)
- [ADR-007: Resource Type Name Migration](../docs/adr/007-resource-type-name-migration.md)
- [CLAUDE.md](../CLAUDE.md) — project conventions
