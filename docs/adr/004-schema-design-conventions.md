# ADR-004: Schema Design Conventions

## Status

Accepted

## Date

2026-02-04 (retroactive documentation)

## Context

Terraform resource schemas define the user-facing contract: which attributes exist, whether they are required or optional, how they are validated, and when changes force recreation. Without consistent conventions, schema designs diverge across resources, confusing users and increasing review burden.

## Decision

### Attribute Types

Use `types.*` from the Terraform Plugin Framework for all model fields:

```go
type model struct {
    ID           types.String `tfsdk:"id"`
    Zone         types.String `tfsdk:"zone"`
    IsolatePorts types.Bool   `tfsdk:"isolate_ports"`
    Tag          types.Int64  `tfsdk:"tag"`
}
```

Never use raw Go types (`string`, `bool`, `int`) for optional fields. The `types.*` wrappers distinguish between null, unknown, and set values.

### Required vs Optional vs Computed

| Scenario                                  | Schema Setting                   |
| ----------------------------------------- | -------------------------------- |
| User must provide                         | `Required: true`                 |
| User may provide, no server default       | `Optional: true`                 |
| User may provide, server provides default | `Optional: true, Computed: true` |
| Server-only value, user cannot set        | `Computed: true`                 |

Use `Optional + Computed` when the Proxmox API actively supplies a value for an omitted field on Read. This allows Terraform to show the server-assigned value in state without requiring the user to specify it. _Documented_ PVE defaults are not enough — see [Provider Defaults vs PVE Defaults](#provider-defaults-vs-pve-defaults) below for the empirical rule (PVE often documents a default that it does _not_ surface on Read).

Use `Computed: true` with `Default` only for boolean fields where the default is a fixed value that matches the API's omission behavior (e.g., `booldefault.StaticBool(false)` when the API treats omission as `false`). Avoid this pattern for string or numeric fields where server-side defaults may change independently of the provider. See the Replication resource's `disable` field in [reference-examples.md](reference-examples.md#booldefault-for-optionalcomputed-bool) for the canonical example. This bool-omission pattern is the documented carve-out from [Provider Defaults vs PVE Defaults](#provider-defaults-vs-pve-defaults) below; do not generalize it to non-boolean fields.

### Provider Defaults vs PVE Defaults

The provider does not duplicate PVE's documented defaults via schema `Default(...)`. Schema choice is driven by **PVE Read behavior** — what the API actually returns on GET — not by the value PVE would apply at QEMU launch time. This rule is empirically grounded: `qemu-server.git src/PVE/QemuServer.pm`'s `$confdesc` documents internal defaults applied at launch, but `parse_vm_config()` returns only literal config-file content, so most documented defaults never appear in the GET response.

| PVE Read behavior                            | Schema target                      | Examples                                                                         |
| -------------------------------------------- | ---------------------------------- | -------------------------------------------------------------------------------- |
| Auto-populates a value on GET                | `Optional + Computed` (no Default) | `boot` (always present on PVE GET response)                                      |
| Returns null/absent when unset               | `Optional` only                    | `cpu.affinity`, `cpu.limit`, `cpu.type`, `vga.type`, `rng.source`, `description` |
| Provider-only attribute (no PVE counterpart) | `Optional + Default`               | `purge_on_destroy`, `stop_on_destroy`, `delete_unreferenced_disks_on_destroy`    |

> **Note.** The `cpu`/`vga`/`rng` examples above describe the target classification. The current `fwprovider/nodes/vm/{cpu,vga,rng}` schemas carry `Optional+Computed` inherited from the SDK and the initial Framework scaffold; conformance — dropping `Computed` from all `cpu`/`vga`/`rng`/`memory` attributes — is tracked under [#1231](https://github.com/bpg/terraform-provider-proxmox/issues/1231). The originally-proposed `Optional+Computed` carve-out on `cpu.cores`/`cpu.sockets` was abandoned after empirical testing (see `.dev/1231_AUDIT.md` §4) showed PVE does not auto-populate those keys on Read.

PVE Read behavior must be verified empirically per attribute. The prescribed method:

1. Run a focused acceptance test that exercises the block with no fields set (`mitmdump --mode regular@8082 --flow-detail 4`), capture GET `/config` responses, inspect what PVE returned.
2. Run a second acceptance test with the block set, to capture any auto-populate behavior PVE adds.
3. Cross-reference with `qemu-server.git`'s `$confdesc` for the documented default; **do not assume the documented default is what PVE returns**.
4. Cross-reference with any existing provider sentinels (`if Field == nil { Field = X }`) — they are evidence the original author observed auto-population, but verify before keeping or dropping.

See `/bpg:debug-api` for the mitmproxy workflow.

**Per-field carve-outs are possible but rare.** A sub-block can have a mix of auto-populated and null-absent fields, in which case only the auto-populated fields keep `Optional+Computed`. Verify per field via mitmproxy rather than trusting `$confdesc` or prior sentinels: the `cpu` block was originally expected to carve out `cores` and `sockets` as `Optional+Computed` because PVE was thought to auto-populate them, but empirical traces (see `.dev/1231_AUDIT.md` §4) showed it does not, and the carve-out was abandoned before landing.

**`NewValue` must coordinate.** When changing a sub-block's schema from `Optional+Computed` to `Optional` only, the sub-package's `NewValue` (FromAPI) function MUST return `NullValue()` (i.e. `types.ObjectNull(attributeTypes())`) when the underlying API device pointer is nil. Returning a non-null Object with null inner fields creates a permanent plan-vs-state diff for users without the block in HCL. See [ADR-008 §`NewValue` (FromAPI Direction)](008-sub-block-contract.md#newvalue-fromapi-direction).

### Immutable Fields

Fields that cannot be changed after resource creation must use `RequiresReplace()`:

```go
"id": schema.StringAttribute{
    Required: true,
    PlanModifiers: []planmodifier.String{
        stringplanmodifier.RequiresReplace(),
    },
},
```

### Resource ID Attribute

Use the `attribute.ResourceID()` helper from `fwprovider/attribute/` to define the `id` attribute for resources where the ID is server-assigned or derived:

```go
"id": attribute.ResourceID(),
```

This helper returns a `schema.StringAttribute` with `Computed: true`, `UseStateForUnknown()`, and `RequiresReplace()` plan modifiers and a standard description.

For resources where the user provides the ID (e.g., SDN VNet's `id` is user-specified), define the attribute manually with `Required: true` and `RequiresReplace()`.

### Validators

Use validators from the `terraform-plugin-framework-validators` module for standard rules. Use project-specific validators from `fwprovider/validators/` for reusable domain rules (e.g., `validators.SDNID()`).

```go
// Standard validator
"type": schema.StringAttribute{
    Validators: []validator.String{
        stringvalidator.OneOf("graphite", "influxdb"),
    },
},

// Range validator
"mtu": schema.Int64Attribute{
    Validators: []validator.Int64{
        int64validator.Between(512, 65536),
    },
},

// Regex validator
"alias": schema.StringAttribute{
    Validators: []validator.String{
        stringvalidator.RegexMatches(
            regexp.MustCompile(`^[a-zA-Z0-9-]+$`),
            "must contain only alphanumeric characters and dashes",
        ),
    },
},
```

### Enum Validators

Use `OneOf` only for **short, stable** PVE enums — typically ≤5 values, with no growth pressure. Examples that qualify: `cpu.architecture` (`"aarch64"`, `"x86_64"`), `memory.hugepages` (`"2"`, `"1024"`, `"any"`), `vga.clipboard` (`"vnc"`).

For **long or version-evolving** enums, drop client-side validation entirely and defer to PVE. Examples that fail the rule: `cpu.type` (~75 CPU model strings, growing as Intel/AMD release new architectures), `vga.type` (~14 VGA types, with new `qxl` variants added over time), `machine`, `bios`, `scsi_hardware`, `audio_device.driver`. Hard-coding these in a `OneOf` validator means the provider has to ship a release every time PVE adds a value, and the validator silently lies until then.

When the constraint is format-only (regex, length, range), keep the validator regardless of enum length — `cpu.affinity` (CPU-list regex), `cpu.cores` range `[1, 1024]`, `name` DNS regex are all stable formats worth checking at plan time.

**Slot regex for map-keyed devices.** Per [ADR-008 §Single-vs-Map Rule](008-sub-block-contract.md#single-vs-map-rule), every map-keyed device sub-package's `mapvalidator.KeysAre(stringvalidator.RegexMatches(...))` must be bounded by the PVE source's `MAX_*` constants (`MAX_NETS`, `MAX_USB_DEVICES`, etc.) so out-of-range slot keys fail at plan time. Relax in a future additive PR if PVE expands the bounds.

### Cross-Field Validation

When validation depends on multiple attributes, implement `ResourceWithConfigValidators`:

```go
func (r *myResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
    return []resource.ConfigValidator{
        resourcevalidator.Conflicting(
            path.MatchRoot("group_id"),
            path.MatchRoot("user_id"),
        ),
    }
}
```

For more complex validation that requires parsing attribute values, implement `ResourceWithValidateConfig` and add logic in the `ValidateConfig` method.

**When to add a validator vs document the constraint.** Cross-attribute constraints are documented in `MarkdownDescription` by default — let PVE reject invalid combinations at apply time. Promote to a plan-time validator only when (a) the constraint is hit _frequently_ in support issues, or (b) PVE's apply-time error is unhelpful (cryptic, mislocated, or only surfaces after partial side-effects). Every plan-time validator is a maintenance liability: it duplicates PVE's logic, drifts as PVE evolves, and adds a failure point distinct from the API's own. Keep the bar high.

### Sensitive Attributes

Mark secret fields (tokens, passwords) as sensitive so Terraform redacts them:

```go
"token": schema.StringAttribute{
    Sensitive: true,
},
```

### Attribute Descriptions

Every schema attribute must have a non-empty `Description`. This text appears in Terraform CLI output (`terraform show`, `terraform plan`) and in auto-generated documentation.

| Field                 | When to Use                                                     | Format                   |
| --------------------- | --------------------------------------------------------------- | ------------------------ |
| `Description`         | **Always**                                                      | Plain text, one sentence |
| `MarkdownDescription` | Only when the description needs formatting (links, code, lists) | Markdown syntax          |

When both are set, `MarkdownDescription` is used for documentation generation and `Description` is used for CLI output. When only `Description` is set, it is used for both. For most attributes, `Description` alone is sufficient.

### Model-API Conversion

Every model implements conversion methods for mapping between Terraform state and API request/response structs.

**Method Naming Convention:**

| Method                            | Purpose                                                                                              |
| --------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `toAPI()`                         | Convert Terraform model to a single API request struct (when create and update share the same shape) |
| `toAPICreate()` / `toAPIUpdate()` | Convert to separate create and update request structs (when they differ)                             |
| `fromAPI()`                       | Convert API response to Terraform model                                                              |

When create and update request types differ (e.g., create includes immutable fields while update does not), use `toAPICreate()` and `toAPIUpdate()` rather than a single `toAPI()`. A shared helper (e.g., `fillCommonFields()`) can reduce duplication when the overlap is large. See the [Replication reference](reference-examples.md#split-createupdate-api-methods) for the canonical example.

Avoid alternative naming patterns such as `importFromAPI()`, `toAPIRequestBody()`, `toOptionsRequestBody()`, `toCreateRequest()`, `toCreateAPIRequest()`, or `intoUpdateBody()`. While functionally equivalent, consistent naming makes patterns discoverable across resources.

> **Legacy code note:** Existing resources may use older naming patterns. New code must use the standard names above. Existing resources will be migrated over time.

**`toAPI()`** — Terraform model to API request struct. Use the `attribute` package helpers so null and unknown values both map to `nil`:

```go
func (m *model) toAPI() *vnets.VNet {
    data := &vnets.VNet{}
    data.Zone = attribute.StringPtrFromValue(m.Zone)
    data.Alias = attribute.StringPtrFromValue(m.Alias)
    data.Tag = attribute.Int64PtrFromValue(m.Tag)
    data.IsolatePorts = attribute.CustomBoolPtrFromValue(m.IsolatePorts)
    return data
}
```

The helpers `StringPtrFromValue`, `Int64PtrFromValue`, `Float64PtrFromValue`, and `CustomBoolPtrFromValue` (all in `fwprovider/attribute/`) return nil for null **and** unknown values, making them safe for Optional+Computed fields. Prefer these over raw `Value*Pointer()` methods, which return `&""` / `&0` / `&false` for unknown values — a common source of bugs.

> **Note:** Custom attribute types (`customtypes.IPCIDRValue`, etc.) cannot use these helpers. For those, continue using `.ValueStringPointer()` directly.

> **Write-through variant — when the API has no nested struct:** Some PVE endpoints scatter a logical group's fields across the top-level request body rather than nesting them under a dedicated struct (e.g. `memory`, `balloon`, `shares`, `hugepages`, `keephugepages` all live directly on the VM config endpoint). When there is no single struct to return, `toAPI` may take the request body and populate the relevant fields in place: `func (m *Model) toAPI(body *SomeRequestBody)`. Keep the verb (`toAPI`); the write-through shape is a side effect of the API's flatness, not a new pattern. See `fwprovider/nodes/vm/memory/model.go` for the reference.

**`fromAPI()`** — API response to Terraform model. Use `types.*PointerValue()` so `nil` maps to null:

```go
func (m *model) fromAPI(id string, data *vnets.VNetData) {
    m.ID = types.StringValue(id)
    m.Zone = types.StringPointerValue(data.Zone)
    m.Alias = types.StringPointerValue(data.Alias)
    m.Tag = types.Int64PointerValue(data.Tag)
    m.IsolatePorts = types.BoolPointerValue(data.IsolatePorts.PointerBool())
}
```

For `CustomBool` fields, use `.PointerBool()` to convert `*CustomBool` → `*bool`, then `types.BoolPointerValue()` handles `nil` naturally — consistent with the other pointer value conversions.

> **When the API type uses `*int64` instead of `*CustomBool`:** The preferred approach is to update the API struct to use `*proxmoxtypes.CustomBool`. If that is not feasible, use the project-wide `CustomBoolPtr()` and `.PointerBool()` methods rather than defining local conversion helpers (e.g., `boolToInt64Ptr()` / `int64ToBoolPtr()`).

### Field Deletion on Update

When an optional field is removed from configuration, the Proxmox API requires explicit deletion via a `delete` parameter. Use `attribute.CheckDelete()` to detect these transitions:

```go
var toDelete []string
attribute.CheckDelete(plan.Alias, state.Alias, &toDelete, "alias")
attribute.CheckDelete(plan.Tag, state.Tag, &toDelete, "tag")

update := &vnets.VNetUpdate{
    VNet:   *plan.toAPI(),
    Delete: toDelete,
}
```

The third argument to `CheckDelete` is the **Proxmox API parameter name**, which may differ from the Terraform attribute name (e.g., `"api-path-prefix"` vs `influx_api_path_prefix`).

### Comma-Separated API Values → Terraform Lists

When the Proxmox API accepts or returns a comma-separated string (e.g., `vmid=100,101,102`, `exclude-path=/tmp,/var`), **always expose it as a Terraform list or set attribute** — never as a raw comma-separated string. This gives users proper HCL list syntax, element-level validation, and `for_each`/`dynamic` block compatibility.

In the schema:

```go
"vmid": schema.ListAttribute{
    Description: "A list of guest VM/CT IDs to include in the backup job.",
    Optional:    true,
    ElementType: types.StringType,
},
```

In `toAPI()` — join the list into a comma-separated string for the API:

```go
if !m.VMIDs.IsNull() && !m.VMIDs.IsUnknown() {
    var ids []string
    diags.Append(m.VMIDs.ElementsAs(ctx, &ids, false)...)
    if len(ids) > 0 {
        joined := strings.Join(ids, ",")
        common.VMID = &joined
    }
}
```

In `fromAPI()` — split the comma-separated string into a list:

```go
if data.VMID != nil && *data.VMID != "" {
    ids := strings.Split(*data.VMID, ",")
    values := make([]attr.Value, len(ids))
    for i, id := range ids {
        values[i] = types.StringValue(strings.TrimSpace(id))
    }
    m.VMIDs, _ = types.ListValue(types.StringType, values)
} else {
    m.VMIDs = types.ListNull(types.StringType)
}
```

For unordered values (e.g., tags, node lists), use `stringset.Value` (a custom set type) instead of `types.List`.

### Custom Types

The project provides custom attribute types in `fwprovider/types/`:

| Type                      | Package                       | Use Case                                           |
| ------------------------- | ----------------------------- | -------------------------------------------------- |
| `stringset.Value`         | `fwprovider/types/stringset/` | Comma-separated list attributes (e.g., node lists) |
| `customtypes.IPAddrValue` | `fwprovider/types/`           | IP address validation                              |
| `customtypes.IPCIDRValue` | `fwprovider/types/`           | CIDR block validation                              |

## Consequences

### Positive

- Consistent user experience across resources
- Null/unknown handling is correct by construction
- Validators catch errors at plan time, before API calls
- Field deletion works correctly with the Proxmox API

### Negative

- Boilerplate for `toAPI`/`fromAPI` and `CheckDelete` on every optional field
- Custom types add a learning curve for new contributors

### Common Mistakes

- Using raw Go types (`string`, `bool`, `int`) for optional model fields — use `types.*` wrappers.
- Using `types.StringValue("")` instead of `types.StringNull()` for absent values — empty string and null are different in Terraform.
- Forgetting `CheckDelete` calls in Update for optional fields — the Proxmox API won't clear the field.
- Using the Terraform attribute name instead of the Proxmox API parameter name in `CheckDelete`.
- Setting `Computed: true` with `Default` on string or numeric fields — leads to unexpected behavior when server defaults change. This combination is acceptable for boolean fields with fixed defaults (see guidance above).
- Adding a schema `Default(...)` for an attribute that mirrors a PVE default. Provider should not duplicate PVE's defaults — use `Optional+Computed` when PVE auto-populates the value on Read, `Optional` only when PVE returns absent. See [Provider Defaults vs PVE Defaults](#provider-defaults-vs-pve-defaults).
- Hard-coding a long, version-evolving PVE enum in a `OneOf` validator (e.g. CPU types, VGA types, machine types, BIOS modes, `scsi_hardware`, `audio_device.driver`). Defer to PVE for these — the validator silently lies between PVE releases. See [Enum Validators](#enum-validators).
- Promoting a cross-attribute constraint to a plan-time validator without evidence the constraint is hit frequently or that PVE's apply-time error is unhelpful. Document in `MarkdownDescription` first; promote only when needed.
- Exposing comma-separated API values as a single `types.String` instead of `types.List` or `stringset.Value` — use proper Terraform list/set types so users get HCL list syntax and element-level operations.
- Using non-standard model method names (`importFromAPI`, `toAPIRequestBody`, `toCreateRequest`, etc.) instead of the canonical `toAPI()` / `toAPICreate()` / `toAPIUpdate()` / `fromAPI()`. See [Model-API Conversion](#model-api-conversion).
- Omitting `Description` on schema attributes — every attribute must have a non-empty description.
- Defining local bool-to-int64 conversion helpers instead of updating the API type to use `*proxmoxtypes.CustomBool`.

## References

- [Reference Examples](reference-examples.md) — annotated code for all patterns above
- [ADR-005: Error Handling](005-error-handling.md) — error patterns in CRUD methods
- [ADR-008: Sub-block Contract](008-sub-block-contract.md) — sub-package shape for VM-style composite resources, including `NewValue` null-Object coordination
- [Terraform Plugin Framework: Schemas](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas)
- [Terraform Plugin Framework: Attributes](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes)
