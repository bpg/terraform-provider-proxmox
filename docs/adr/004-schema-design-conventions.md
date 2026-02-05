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

Try to avoid setting `Computed: true` with Default values in the schema. This can lead to unexpected behavior if the default changes on the server side.

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

**`toAPI()`** — Terraform model to API request struct. Use `Value*Pointer()` methods so null values map to `nil`:

```go
func (m *model) toAPI() *vnets.VNet {
    data := &vnets.VNet{}
    data.Zone = m.Zone.ValueStringPointer()
    data.Alias = m.Alias.ValueStringPointer()
    data.Tag = m.Tag.ValueInt64Pointer()
    data.IsolatePorts = proxmoxtypes.CustomBoolPtr(m.IsolatePorts.ValueBoolPointer())
    return data
}
```

For `CustomBool` fields (API uses `0`/`1` integers for booleans), use `proxmoxtypes.CustomBoolPtr()` which converts `*bool` → `*CustomBool`.

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
- Setting `Computed: true` with `Default` — leads to unexpected behavior when server defaults change.

## References

- [Reference Examples](reference-examples.md) — annotated code for all patterns above
- [ADR-005: Error Handling](005-error-handling.md) — error patterns in CRUD methods
- [Terraform Plugin Framework: Schemas](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas)
- [Terraform Plugin Framework: Attributes](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes)
