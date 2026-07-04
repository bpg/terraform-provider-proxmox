# ADR-008: Sub-block Contract

## Status

Accepted

## Date

2026-04-20

## Context

VM-shaped resources (`proxmox_vm`, `proxmox_cloned_vm`) compose dozens of independent configuration blocks: `cpu`, `vga`, `rng`, `cdrom`, `memory`, `disk`, `network_device`, `usb`, `hostpci`, `numa`, `serial_device`, `parallel`, `virtiofs`, `audio_device`, `efi_disk`, `tpm_state`, `agent`, `watchdog`, `amd_sev`, `initialization`, `operating_system`, `smbios`, `startup`. Each block needs its own schema, model, and conversion logic; without a shared contract, every sub-package invents slightly different shapes and the top-level resource has to special-case each one.

The existing five sub-packages (`cpu`, `vga`, `rng`, `cdrom`, `memory` under `fwprovider/nodes/vm/`) already converged on a function-based, Value-centric pattern: each sub-package exports an opaque `Value` type alias, a `NewValue` constructor that maps a PVE GET response into the Value, and `FillCreateBody` / `FillUpdateBody` functions that mutate the API request body from a Value. The top-level resource holds these Values in its model struct and calls the sub-package functions during CRUD.

> **Note.** This ADR codifies the target contract. The existing five sub-packages pre-date it and deviate in places: `memory` exports only `NewValue` + `FillUpdateBody` (no `FillCreateBody`) and lacks `datasource_schema.go` / `resource_test.go`; `cpu`/`vga`/`rng` schemas carry `Optional+Computed` on fields that should drop `Computed` per [ADR-004 §Provider Defaults vs PVE Defaults](004-schema-design-conventions.md#provider-defaults-vs-pve-defaults); the body-taking helper described in the [`FillUpdateBody`](#fillcreatebody-and-fillupdatebody) example ships as `attribute.CheckDeleteBody` (Go does not allow overloading, so the slice-taking universal form keeps the shorter `attribute.CheckDelete` name) but the existing sub-package call-sites still hand-roll deletion with `ShouldBeRemoved` / `IsDefined` cascades; `cpu` still hosts the `numa` / `hotplugged` misnomers covered in [Block Name Maps to a Single PVE Concept](#block-name-maps-to-a-single-pve-concept). Conformance with this contract — including the `numa.enabled` / `vcpus` rehomes — is tracked under [#1231](https://github.com/bpg/terraform-provider-proxmox/issues/1231).

Two factors make this contract worth codifying now rather than after the next 15 sub-packages land:

1. **Joint ownership.** Five of the existing sub-packages are also consumed by `proxmox_cloned_vm`. Any contract change must keep both consumers green. Codification raises the cost of accidental drift.
2. **PVE Read behavior.** Sub-block schema choices (Optional vs Optional+Computed, NewValue null-Object handling) follow PVE's actual Read semantics — see [ADR-004 §Provider Defaults vs PVE Defaults](004-schema-design-conventions.md#provider-defaults-vs-pve-defaults). The two ADRs are interlocking: the contract here describes _how_ sub-packages are shaped; ADR-004 describes _which_ shape applies to a given attribute.

## Decision

### Signature Families

Sub-packages fall into one of two families based on the underlying PVE shape.

**Single-nested** — used when PVE exposes the block as a single property (`cpu=...`, `vga=...`, `memory=...`). The sub-package's `Value` aliases `types.Object`.

```go
// Value is the framework type the parent resource holds in its model.
type Value = types.Object

func ResourceSchema() schema.Attribute       // returns schema.SingleNestedAttribute
func DataSourceSchema() dsschema.Attribute   // returns dsschema.SingleNestedAttribute
func NullValue() Value                       // returns types.ObjectNull(attributeTypes())

func NewValue(ctx context.Context, config *vms.GetResponseData, diags *diag.Diagnostics) Value
func FillCreateBody(ctx context.Context, planValue Value, body *vms.CreateRequestBody, diags *diag.Diagnostics)
func FillUpdateBody(ctx context.Context, planValue, stateValue Value, body *vms.UpdateRequestBody, diags *diag.Diagnostics)
```

Reference implementation: `fwprovider/nodes/vm/cpu/`.

**Map-keyed** — used when PVE exposes the block as a slot-keyed family (`net0=…`, `net1=…`; `scsi0=…`, `virtio0=…`; `ide2=…`). The sub-package's `Value` aliases `types.Map`. Function signatures are identical to the single-nested family except for the underlying alias.

```go
type Value = types.Map  // remaining signatures unchanged
```

Reference implementation: `fwprovider/nodes/vm/cdrom/`.

### Sub-package `Model` Struct

Each sub-package defines a `Model` struct holding the framework field types (`types.String`, `types.Int64`, etc.) with `tfsdk:` tags matching the schema attribute names, plus an `attributeTypes()` helper:

```go
type Model struct {
    Affinity types.String `tfsdk:"affinity"`
    Cores    types.Int64  `tfsdk:"cores"`
    // ...
}

func attributeTypes() map[string]attr.Type { /* ... */ }
```

`Model` is exported by Go visibility (capital `M`) but is treated as an implementation detail — the parent resource's model holds the sub-package's `Value` alias, not `Model`. Outside the sub-package, only `Value`, `NewValue`, `FillCreateBody`, `FillUpdateBody`, `NullValue`, `ResourceSchema`, and `DataSourceSchema` are referenced. `Model` lives in `model.go` because `attributeTypes()` and `(de)serialization` to the framework's `basetypes.ObjectAsOptions{}` need it; do not import it from the parent or another sub-package.

### `NewValue` (FromAPI Direction)

`NewValue` maps a PVE GET response to the sub-package's `Value`. Two rules apply uniformly:

1. **Use `types.*PointerValue()` for nil-as-null mapping.** Provider does _not_ substitute PVE's documented internal defaults for absent fields — see [ADR-004 §Provider Defaults vs PVE Defaults](004-schema-design-conventions.md#provider-defaults-vs-pve-defaults). Empirical mitmproxy traces are the source of truth for which fields PVE auto-populates.

2. **Return `NullValue()` (i.e. `types.ObjectNull(attributeTypes())` / `types.MapNull(...)`) when the underlying API device pointer is nil.** Returning a non-null Object with null inner fields creates a permanent plan-vs-state diff for users who don't define the block in HCL. The block-level null guard reflects "user has no `vga` block in HCL → state has null `vga`".

   **Carve-out (hypothetical):** if PVE ever auto-populates specific inner fields on Read whenever the block is set, those fields would keep `Optional+Computed` and the Object would stay non-null. Verify per field with mitmproxy before assuming it applies — the `cpu.cores`/`cpu.sockets` candidates for this carve-out were empirically disproved (see `.dev/1231_AUDIT.md` §4) and the `cpu` sub-package went Optional-only on all attributes. Document any real carve-out per field with a short code comment citing the trace evidence.

### `FillCreateBody` and `FillUpdateBody`

`FillCreateBody` reads the plan and writes the create request body. Use the [`attribute` package helpers from ADR-004](004-schema-design-conventions.md#model-api-conversion) to convert framework values to API pointers — they handle null and unknown identically:

```go
body.CPUAffinity = attribute.StringPtrFromValue(plan.Affinity)
body.CPUCores    = attribute.Int64PtrFromValue(plan.Cores)
body.CPULimit    = attribute.Float64PtrFromValue(plan.Limit)
body.NUMAEnabled = attribute.CustomBoolPtrFromValue(plan.Numa)
```

Avoid hand-rolled `IsUnknown()` cascades — every field would need its own check, and `ValueXxxPointer()` returns `&""` / `&0` / `&false` for unknown values, which sends bogus zeros to the API.

`FillUpdateBody` reads the plan and prior state, then writes the update request body. Use the `attribute` helpers to populate fields (nil-safe for null/unknown), and `attribute.CheckDeleteBody` to record removals directly on the body — the helper appends to `body.Delete` internally via `body.AppendDelete`:

```go
attribute.CheckDeleteBody(plan.Affinity, state.Affinity, body, "affinity")
attribute.CheckDeleteBody(plan.Limit,    state.Limit,    body, "cpulimit")

body.CPUAffinity = attribute.StringPtrFromValue(plan.Affinity)
body.CPULimit    = attribute.Float64PtrFromValue(plan.Limit)
```

The CheckDelete string is the **PVE API parameter name** (lowercase, matching the `url:` struct tag on the request body), not the Go field name. The same string ends up in PVE's wire request as `delete=affinity,cpulimit,…`. The body-taking form requires the body to implement `attribute.DeleteAppender` — a single-method interface satisfied by `*vms.UpdateRequestBody.AppendDelete(apiName)` — so the helper can record the removal in place.

The body-taking form coexists with the universal `(plan, state, *[]string, "api-name")` signature documented in [ADR-004 §Field Deletion on Update](004-schema-design-conventions.md#field-deletion-on-update) — non-VM Framework resources (SDN, metrics, Replication, etc.) use the universal form when their body type does not implement `DeleteAppender`. VM sub-packages prefer `CheckDeleteBody`, eliminating the local `[]string` plumbing.

Hand-rolled `ShouldBeRemoved` + `IsDefined` cascades that branch on `plan.Field.Equal(state.Field)` are not used in either form — `attribute.StringPtrFromValue` already returns nil for null/unknown, so the explicit `IsDefined` guard is redundant in `FillCreateBody` _and_ `FillUpdateBody`. A reflection-based `body.ToDelete(GoFieldName)` helper on `vms.UpdateRequestBody` exists for legacy callers, but new sub-package code must not use it — the canonical pattern is `attribute.CheckDeleteBody` with the API name. See [Common Mistakes](#common-mistakes).

### Guard Ordering in `FillUpdateBody`

When a single-nested sub-package normalizes null/unknown plan Objects into an all-null `Model` for per-field diffing (an `unpackOrEmpty`-style helper), the top-of-function guards must come in a specific order:

```go
// FIRST: skip when the plan is unknown or identical to state.
if planValue.IsUnknown() || planValue.Equal(stateValue) {
    return
}

plan := unpackOrEmpty(ctx, planValue, diags)
state := unpackOrEmpty(ctx, stateValue, diags)
if diags.HasError() {
    return
}

// Per-field CheckDeleteBody calls, then compound/atomic field deletion.
// Gate compound deletion on plan.Type.IsNull() (not !IsDefined) so an
// inner-field plan that is unknown is not mis-read as a removal.

// THEN: skip population on full block removal.
if planValue.IsNull() {
    return
}

plan.toAPI(ctx, body, diags)
```

Two failure modes make the ordering load-bearing. First, the unpack helper collapses a block-level unknown into all-null fields — without the top `IsUnknown()` guard, every `CheckDeleteBody(plan.X, state.X, ...)` fires a spurious `delete=`. Second, `!attribute.IsDefined(plan.Type)` matches both null and unknown, so an inner field set from an unresolved reference (e.g. `cpu.type = other_resource.x.computed_attr`) would mis-trigger a compound `delete=cpu`; gate compound deletion on `plan.Type.IsNull()` so only explicit user removal triggers it.

### Conversion Pitfalls

**Zero is a value.** A user-set `0` on an `Int64` attribute must reach PVE when `0` is a meaningful PVE value (e.g. "no limit"). The pattern `if v.ValueInt64() != 0` silently drops it — use `attribute.Int64PtrFromValue`, which returns a pointer to `0` when the user set 0 and nil when null/unknown. When the API struct field is `*int` rather than `*int64`, cast explicitly:

```go
if v := attribute.Int64PtrFromValue(m.MaxBytes); v != nil {
    n := int(*v)
    dev.MaxBytes = &n
}
```

Verify with mitmproxy that PVE accepts `0` as meaningful for the specific field; if PVE rejects it, raise the validator's lower bound instead of silently dropping the value.

**Empty string vs null on non-pointer subfields.** Some PVE wire types have `string` (not `*string`) subfields (e.g. `CustomRNGDevice.Source`). In `NewValue`, map an empty string from the API to `types.StringNull()` — otherwise plan=null vs state=`""` produces a permanent diff.

### Map-keyed Update Diff

Map-keyed sub-packages diff plan against state with `utils.MapDiff`, which returns three maps (`toCreate`, `toUpdate`, `toDelete`):

```go
toCreate, toUpdate, toDelete := utils.MapDiff(plan, state)

for slot, dev := range toCreate { body.AddDevice(slot, dev.toAPI()) }   // sub-package-specific add helper
for slot, dev := range toUpdate { body.AddDevice(slot, dev.toAPI()) }
for slot := range toDelete      { body.Delete = append(body.Delete, slot) }
```

Slot-level deletion appends the slot key (e.g. `"net1"`, `"ide2"`) directly to `body.Delete`. This is the same `body.Delete []string` (`url:"delete,omitempty,comma"`) that scalar field deletion populates via `attribute.CheckDeleteBody` — slot keys and scalar API names share one comma-separated `delete=…` parameter on the wire. The per-device add helper (`AddCustomStorageDevice`, `AddNetworkDevice`, etc.) is sub-package-specific; cdrom is the reference, calling `body.AddCustomStorageDevice(iface, dev.toAPI())` — the `toAPI()` / `fromAPI()` names mandated by [ADR-004 §Model-API Conversion](004-schema-design-conventions.md#model-api-conversion).

Do not add an `if planValue.IsNull() { return }` early-return ahead of the diff. When the user removes the whole block after having slots, the plan is null and the slot deletes come precisely from `MapDiff(nil, state)`'s `toDelete` — Go ranges over nil maps as no-ops, so the null plan needs no special-casing, and an early return silently leaks every slot on block removal.

### File Layout

Per [ADR-003 §File Pattern](003-resource-file-organization.md#file-pattern), each sub-package contains:

| File                          | Purpose                                                                                                                                                                                                          |
| ----------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `resource_schema.go`          | `ResourceSchema()` returning `schema.SingleNestedAttribute` or `schema.MapNestedAttribute`                                                                                                                       |
| `datasource_schema.go`        | `DataSourceSchema()` — the datasource counterpart. Output attributes are `Computed: true` only (never `Optional`); see [`CLAUDE.md` §Datasource Schema Attributes](../../CLAUDE.md#datasource-schema-attributes) |
| `model.go`                    | `Model` struct (treated as sub-package-internal — see [Sub-package `Model` Struct](#sub-package-model-struct)), `attributeTypes()`, `NullValue()`, optional `Model` (de)serialization helpers                    |
| `resource.go`                 | `NewValue`, `FillCreateBody`, `FillUpdateBody` — the three exported workhorses                                                                                                                                   |
| `resource_test.go`            | Acceptance tests covering the block's CRUD + import scenarios                                                                                                                                                    |
| `model_test.go` (when needed) | Unit tests on `NewValue` nil-paths and `FillUpdateBody` diff matrices                                                                                                                                            |

### Single-vs-Map Rule

Use `SingleNestedAttribute` when PVE currently exposes the device as one slot only:

- **Architecturally single** (VM hardware model precludes multiples): `efi_disk` (PVE `efidisk0` — one EFI variable store per VM, firmware-defined), `tpm_state` (PVE `tpmstate0` — TPM spec is single-instance per system).
- **Conventionally single** (PVE currently allows one slot but isn't architecturally constrained): `audio_device` (PVE `audio0` — qemu-server source `$id //= 0`). The trade-off is explicit: if PVE later adds additional slots (e.g. `audio1+`), migrating the schema to map-keyed is a breaking change. Accepted to keep HCL ergonomic — PVE has had `audio0` only for years and the additional indentation is not justified by speculative forward-compat.

Use `MapNestedAttribute` (keyed by PVE slot name) for everything else: `disk`, `network_device`, `cdrom`, `usb`, `hostpci`, `numa`, `serial_device`, `parallel`, `virtiofs`. Even when PVE exposes a single-slot device today, prefer map-keyed if PVE could plausibly add slots later — `cdrom` is the reference (PVE has a pending feature request to support multiple cdroms; the addition is _additive_ because cdrom is map-keyed).

The map's `Validators` must include a `mapvalidator.KeysAre(stringvalidator.RegexMatches(...))` bounded to PVE's current per-family slot count (verified from `qemu-server.git src/PVE/QemuServer.pm` constants — `MAX_NETS`, `MAX_USB_DEVICES`, `MAX_HOSTPCI_DEVICES`, etc.). Relax in a future additive PR if PVE expands the bounds.

### Block Name Maps to a Single PVE Concept

Sub-block names normally correspond to one PVE concept (`cpu` → CPU emulation, `vga` → VGA device, `rng` → RNG device). Avoid the "virtual sub-block" pattern — a block that contains attributes mapping to unrelated PVE concepts.

The `cpu` block is the historical exception: `cpu.numa` (PVE: `numa=1` — a top-level VM toggle) and `cpu.hotplugged` (PVE: `vcpus=N` — a top-level vCPU count) are SDK-inherited misnomers scheduled for relocation — `numa.enabled` under a dedicated top-level `numa` block, `vcpus` as a top-level scalar — tracked under [#1231](https://github.com/bpg/terraform-provider-proxmox/issues/1231). Future sub-packages should not repeat the pattern.

### Per-Commit Verification Gate

Each sub-package port commit passes all of:

1. `make lint` — 0 issues.
2. `go build ./...` — clean.
3. `./testacc TestAccResourceVM2<Sub>` — all scenarios pass.
4. `./testacc TestAccResourceClonedVM` — the joint-ownership gate: sub-packages consumed by `proxmox_cloned_vm` break it silently otherwise.
5. Mitmproxy trace reviewed — Create and Update send the expected wire parameters; block removal emits `delete=<key>` (compound) or `delete=<slot>` per removed slot (map-keyed), not silence; field removal emits per-field `delete=<api-name>` for multi-key sub-blocks.
6. `make docs` — regenerate when the schema shape changed visibly (`Computed` flag, description text).

Skipping step 3 masks schema-shape bugs; skipping step 4 breaks `proxmox_cloned_vm` silently; skipping step 5 misses behavioral regressions that state assertions cannot catch (a leaked block-removal delete leaves state clean while PVE still has the device).

### Mandatory Acceptance Scenarios

Per [ADR-006 §Functional Coverage](006-testing-requirements.md#functional-coverage-requirement), every sub-package covers:

- No block in HCL → state has no block attributes set (verifies the `NullValue` guard).
- Block set with a subset of fields → set fields in state, others null (plan stability).
- Update changes the block → state matches the new config (compound replacement / map update).
- Block removed entirely after being set → state has no block attributes (verifies block-level delete).
- Map-keyed only: add slot, update slot, remove slot — the three `MapDiff` paths; plus apply → empty re-plan (reorder immunity).
- Fields where zero is meaningful: set to `0` → state has `0` (see [Conversion Pitfalls](#conversion-pitfalls)).

## Consequences

### Positive

- New sub-packages follow a recipe: copy `cpu/` for single-nested or `cdrom/` for map-keyed, swap field names. No design decisions per sub-package.
- Joint-ownership invariant is enforceable: `proxmox_cloned_vm` and `proxmox_vm` consume sub-packages via the same contract. `proxmox_cloned_vm` calls `FillCreateBody` to pre-seed the cloned VM's options (cdrom/cpu/rng/vga today), using `FillUpdateBody` only where the sub-package omits `FillCreateBody` (memory is the current case — see the Note above); every method it does call has the same signature as the `proxmox_vm` call.
- Top-level resource shrinks: each sub-block adds one schema entry, one `NewValue` call in Read, and one `Fill…Body` call in Create/Update. Per [ADR-005 §Read-Back After Create and Update](005-error-handling.md#read-back-after-create-and-update), Read is the single point that materializes API → state.
- Reorder-immunity for map-keyed devices: inserting `net2` does not churn `net0` or `net1`, removing `net1` does not renumber the rest. Slot key is identity.

### Negative

- The `Value = types.Object` / `Value = types.Map` aliases hide the fact that they are framework primitive types — readers must follow the alias to know what they're holding. Acceptable given the consistency benefit.
- The function-based shape can't enforce signatures via Go interfaces — every sub-package re-declares `NewValue`, `FillCreateBody`, `FillUpdateBody`. Drift is caught only by review and the joint-ownership tests.
- `attributeTypes()` and `NullValue()` are hand-written boilerplate per sub-package. The framework does not generate either from the schema, so manual maintenance is unavoidable.

### Common Mistakes

- Returning a non-null Object with null inner fields from `NewValue` when the API device is nil — produces a permanent plan-vs-state diff after the parent resource's schema is `Optional`-only. See [ADR-004 §Provider Defaults vs PVE Defaults](004-schema-design-conventions.md#provider-defaults-vs-pve-defaults).
- Substituting provider-invented defaults (e.g. `Type → "kvm64"`, `Cores → 1`) in `NewValue`. Empirical mitmproxy traces determine whether PVE auto-populates a field; a `$confdesc` default or an existing provider sentinel is suggestive evidence, not proof. See [ADR-004 §Provider Defaults vs PVE Defaults](004-schema-design-conventions.md#provider-defaults-vs-pve-defaults) for the empirical rule, and the abandoned `cpu.cores` / `cpu.sockets` carve-out there as a worked cautionary example.
- Calling `body.ToDelete(GoFieldName)` directly in `FillUpdateBody`. The canonical pattern is `attribute.CheckDeleteBody` with an API name; the body's reflection helper stays out of caller code.
- Passing a Go struct field name (e.g. `"CPUAffinity"`) where the PVE API name (`"affinity"`) is expected — applies to both `attribute.CheckDeleteBody`'s fourth argument and any direct `body.Delete = append(...)` on slot keys. PVE rejects the literal Go-style name. Use the `url:` tag's first part.
- Hand-rolling `if !plan.Field.Equal(state.Field) { if ShouldBeRemoved(...) { … } else if IsDefined(...) { … } }` cascades. Use `attribute.CheckDeleteBody` (deletion tracking) and `attribute.*PtrFromValue` (population) — the helpers handle null/unknown so the cascade is unnecessary.
- Importing a sub-package's `Model` from the parent resource or another sub-package. Even though `Model` is Go-visible (capital `M`), it is sub-package-internal by contract — the parent resource composes via `Value` / `NewValue` / `Fill…Body` only.
- Exposing the block as `proxmox_vm.attribute_name` when the block contains attributes from unrelated PVE concepts — split into separate blocks instead. The cpu→numa→vcpus relocation is the cautionary tale.
- Map-keyed slot regex too loose (`^[a-z]+[0-9]+$`) or too tight (only the slots used in tests). Bound it from the PVE Perl source's `MAX_*` constants.
- An `if planValue.IsNull() { return }` early-return before the map diff — leaks every slot delete on whole-block removal. See [Map-keyed Update Diff](#map-keyed-update-diff).
- Running per-field `CheckDeleteBody` calls before the block-level `IsUnknown()` guard — a block-unknown plan collapses to all-null fields and fires spurious deletes. See [Guard Ordering in `FillUpdateBody`](#guard-ordering-in-fillupdatebody).
- Using `!= 0` as an unset check on integer fields where `0` is a valid PVE value — use `attribute.Int64PtrFromValue`. See [Conversion Pitfalls](#conversion-pitfalls).
- Guarding body assignment with `reflect.DeepEqual(dev, &vms.Custom<X>{})` zero-checks — assign `body.<Device> = plan.toAPI()` directly; the compound serializer already skips zero fields via `omitempty`.

## References

- [ADR-003: Resource File Organization](003-resource-file-organization.md) — file pattern this contract uses
- [ADR-004: Schema Design Conventions](004-schema-design-conventions.md) — attribute helpers, `CheckDelete`, PVE-defaults rule
- [ADR-005: Error Handling](005-error-handling.md) — Read-back-after-Create/Update flow that calls `NewValue`
- [ADR-006: Testing Requirements](006-testing-requirements.md) — acceptance test scenarios per sub-block
- `fwprovider/nodes/vm/cpu/` — single-nested reference implementation
- `fwprovider/nodes/vm/cdrom/` — map-keyed reference implementation
- `fwprovider/attribute/attribute.go` — `StringPtrFromValue`, `Int64PtrFromValue`, `Float64PtrFromValue`, `CustomBoolPtrFromValue`, `CheckDelete`, `CheckDeleteBody`, `DeleteAppender`, `IsDefined`
- `utils/maps.go` — `MapDiff` for map-keyed update diff
- [#1231](https://github.com/bpg/terraform-provider-proxmox/issues/1231) — VM resource Plugin Framework migration epic
