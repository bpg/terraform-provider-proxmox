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
|-------------------------------------------|----------------------------------|
| User must provide                         | `Required: true`                 |
| User may provide, no server default       | `Optional: true`                 |
| User may provide, server provides default | `Optional: true, Computed: true` |
| Server-only value, user cannot set        | `Computed: true`                 |

Use `Optional + Computed` when the Proxmox API supplies a default value for an omitted field. This allows Terraform to show the server-assigned value in state without requiring the user to specify it.

Use `Computed: true` with `Default` only for boolean fields where the default is a fixed value that matches the API's omission behavior (e.g., `booldefault.StaticBool(false)` when the API treats omission as `false`). Avoid this pattern for string or numeric fields where server-side defaults may change independently of the provider. See the Replication resource's `disable` field in [reference-examples.md](reference-examples.md#booldefault-for-optionalcomputed-bool) for the canonical example.

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

### Sensitive Attributes

Mark secret fields (tokens, passwords) as sensitive so Terraform redacts them:

```go
"token": schema.StringAttribute{
    Sensitive: true,
},
```

### Model-API Conversion

Every model implements two conversion methods:

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

When the API's create and update endpoints accept different field sets (e.g., immutable fields only accepted at creation), use separate methods named `toAPICreate()` and `toAPIUpdate()`. See the [Replication reference](reference-examples.md#split-createupdate-api-methods) for the canonical example.

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
|---------------------------|-------------------------------|----------------------------------------------------|
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
- Exposing comma-separated API values as a single `types.String` instead of `types.List` or `stringset.Value` — use proper Terraform list/set types so users get HCL list syntax and element-level operations.

## References

- [Reference Examples](reference-examples.md) — annotated code for all patterns above
- [ADR-005: Error Handling](005-error-handling.md) — error patterns in CRUD methods
- [Terraform Plugin Framework: Schemas](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas)
- [Terraform Plugin Framework: Attributes](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes)
