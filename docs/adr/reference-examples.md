# ADR: Reference Examples for New Framework Resource Implementations

## Purpose

This document identifies the cleanest framework provider implementations in the codebase
and annotates their key patterns. It serves as a copy-from guide for contributors adding
new resources or datasources.

Three implementations were selected after reviewing all ~48 resources and ~35 datasources
in `fwprovider/`:

| Rank      | Implementation     | Why Selected                                                                             |
|-----------|--------------------|------------------------------------------------------------------------------------------|
| Primary   | **SDN VNet**       | Simplest clean implementation, perfect 3-file pattern (resource, model, datasource)      |
| Secondary | **Metrics Server** | Richer schema with many optional fields, `attribute.ResourceID()`, separate `name`/`id`  |
| Tertiary  | **ACL**            | Cross-field validation (`ConfigValidators`), custom import ID parsing, non-standard CRUD |

Start with SDN VNet for any new resource. Refer to Metrics Server or ACL only when your
resource needs the additional patterns they demonstrate.

---

## How to Use This Document

| I need to...                                 | Start here                                                                                                                        |
|----------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------|
| Create a basic resource                      | [SDN VNet walkthrough](#primary-reference-sdn-vnet)                                                                               |
| Handle many optional fields                  | [Metrics Server](#secondary-reference-metrics-server)                                                                             |
| Add cross-field validation                   | [ACL ConfigValidators](#tertiary-reference-acl)                                                                                   |
| Need nested attributes?                      | [SDN Subnet `SingleNestedAttribute`](../../fwprovider/cluster/sdn/subnet/resource.go)                                             |
| Add a datasource                             | [SDN VNet datasource](#datasource) and [Metrics Server datasource](../../fwprovider/cluster/metrics/datasource_metrics_server.go) |
| Find advanced patterns (generics, factories) | [Advanced Patterns table](#advanced-patterns)                                                                                     |

---

## Primary Reference: SDN VNet

**Files:**

- `fwprovider/cluster/sdn/vnet/resource.go` -- CRUD + schema
- `fwprovider/cluster/sdn/vnet/model.go` -- Terraform/API mapping
- `fwprovider/cluster/sdn/vnet/datasource.go` -- Read-only datasource
- `fwprovider/cluster/sdn/vnet/resource_test.go` -- Acceptance tests

### Interface Compliance Assertions

Every resource file starts with compile-time interface checks. These ensure your struct
satisfies all required interfaces before runtime.

```go
var (
    _ resource.Resource                = &Resource{}
    _ resource.ResourceWithConfigure   = &Resource{}
    _ resource.ResourceWithImportState = &Resource{}
)
```

If your resource implements additional interfaces (e.g., `ResourceWithConfigValidators`),
add them here. Using `&Resource{}` (pointer literal) is the standard style; the ACL
example uses `(*aclResource)(nil)` which is equivalent.

### Resource Struct and Constructor

The struct holds the typed API client. The constructor returns the zero-value struct
(the client is injected later via `Configure`).

```go
type Resource struct {
    client *cluster.Client
}

func NewResource() resource.Resource {
    return &Resource{}
}
```

### Configure Method

Extracts `config.Resource` from provider data and stores the typed client.

```go
func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    cfg, ok := req.ProviderData.(config.Resource)
    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected config.Resource, got: %T", req.ProviderData),
        )
        return
    }

    r.client = cfg.Client.Cluster()
}
```

Key points:

- Guard against `nil` provider data (happens during plan-only operations).
- Type-assert to `config.Resource` (resources) or `config.DataSource` (datasources).
- Navigate the client hierarchy: `cfg.Client.Cluster()`, `cfg.Client.Access()`, etc.

### Schema Definition

```go
func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages Proxmox VE SDN VNet.",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Description: "The unique identifier of the SDN VNet.",
                Required:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
                Validators: validators.SDNID(),
            },
            "alias": schema.StringAttribute{
                Optional:    true,
                Description: "An optional alias for this VNet.",
                Validators: []validator.String{
                    stringvalidator.RegexMatches(
                        regexp.MustCompile(`^[()._a-zA-Z0-9\s-]+$`),
                        "alias must contain only alphanumeric characters...",
                    ),
                    stringvalidator.LengthAtMost(256),
                },
            },
            // ... more attributes
        },
    }
}
```

Patterns to note:

- `RequiresReplace()` on immutable fields (forces destroy+recreate on change).
- Validators from `fwprovider/validators` package for reusable rules.
- Inline validators from `terraform-plugin-framework-validators` for one-off rules.

### Create

The full create flow: plan -> toAPI -> API call -> read back -> set state.

```go
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan model

    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    vnet := plan.toAPI()

    err := r.client.SDNVnets(plan.ID.ValueString()).CreateVnet(ctx, vnet)
    if err != nil {
        resp.Diagnostics.AddError("Unable to Create SDN VNet", err.Error())
        return
    }

    // Read back to get actual state (important for computed fields)
    data, err := r.client.SDNVnets(plan.ID.ValueString()).GetVnet(ctx)
    if err != nil {
        resp.Diagnostics.AddError("Unable to Read SDN VNet After Creation", err.Error())
        return
    }

    readModel := &model{}
    readModel.fromAPI(plan.ID.ValueString(), data)

    resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}
```

Always read back after create. The API may normalize values or populate computed fields.

### Read (with not-found handling)

```go
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state model

    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    data, err := r.client.SDNVnets(state.ID.ValueString()).GetVnet(ctx)
    if err != nil {
        if errors.Is(err, api.ErrResourceDoesNotExist) {
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError("Unable to Read SDN VNet", err.Error())
        return
    }

    readModel := &model{}
    readModel.fromAPI(state.ID.ValueString(), data)

    resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}
```

The critical pattern: when the API returns `api.ErrResourceDoesNotExist`, call
`resp.State.RemoveResource(ctx)` to tell Terraform the resource was deleted outside of
Terraform. This triggers a plan to recreate it.

### Update (with CheckDelete for optional fields)

```go
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan model
    var state model

    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    var toDelete []string
    attribute.CheckDelete(plan.Alias, state.Alias, &toDelete, "alias")
    attribute.CheckDelete(plan.IsolatePorts, state.IsolatePorts, &toDelete, "isolate-ports")
    attribute.CheckDelete(plan.Tag, state.Tag, &toDelete, "tag")
    attribute.CheckDelete(plan.VlanAware, state.VlanAware, &toDelete, "vlanaware")

    vnet := plan.toAPI()
    reqData := &vnets.VNetUpdate{
        VNet:   *vnet,
        Delete: toDelete,
    }

    err := r.client.SDNVnets(plan.ID.ValueString()).UpdateVnet(ctx, reqData)
    // ... error handling, read back, set state
}
```

`attribute.CheckDelete` (defined in `fwprovider/attribute/attribute.go:47-59`) adds
the API field name to the delete list when the plan value is null but the state value
is not. This tells the Proxmox API to remove the field. The third argument is the
**API parameter name** (not the Terraform attribute name).

### Delete (ignore not-found)

```go
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state model

    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    err := r.client.SDNVnets(state.ID.ValueString()).DeleteVnet(ctx)
    if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
        resp.Diagnostics.AddError("Unable to Delete SDN VNet", err.Error())
    }
}
```

Ignore `api.ErrResourceDoesNotExist` on delete -- the resource is already gone, which
is the desired end state.

### ImportState

```go
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    data, err := r.client.SDNVnets(req.ID).GetVnet(ctx)
    if err != nil {
        if errors.Is(err, api.ErrResourceDoesNotExist) {
            resp.Diagnostics.AddError("SDN VNet Not Found",
                fmt.Sprintf("SDN VNet with ID '%s' was not found", req.ID))
            return
        }
        resp.Diagnostics.AddError("Unable to Import SDN VNet", err.Error())
        return
    }

    readModel := &model{}
    readModel.fromAPI(req.ID, data)

    resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}
```

For simple resources where the import ID maps directly to the API lookup key, this
pattern is sufficient. For composite IDs, see [ACL's custom parsing](#custom-import-id-parsing).

### Model Struct (model.go)

```go
type model struct {
    ID           types.String `tfsdk:"id"`
    Zone         types.String `tfsdk:"zone"`
    Alias        types.String `tfsdk:"alias"`
    IsolatePorts types.Bool   `tfsdk:"isolate_ports"`
    Tag          types.Int64  `tfsdk:"tag"`
    VlanAware    types.Bool   `tfsdk:"vlan_aware"`
}
```

Rules for model structs:

- Every field gets a `tfsdk` tag matching the schema attribute name.
- Use `types.String`, `types.Bool`, `types.Int64` -- never raw Go types for optional fields.
- Required fields that are never null can use plain Go types (see ACL model for this pattern).

### toAPI Method

```go
func (m *model) toAPI() *vnets.VNet {
    data := &vnets.VNet{}

    data.Zone = m.Zone.ValueStringPointer()
    data.Alias = m.Alias.ValueStringPointer()
    data.IsolatePorts = proxmoxtypes.CustomBoolPtr(m.IsolatePorts.ValueBoolPointer())
    data.Tag = m.Tag.ValueInt64Pointer()
    data.VlanAware = proxmoxtypes.CustomBoolPtr(m.VlanAware.ValueBoolPointer())

    return data
}
```

Pattern: optional fields use `ValueStringPointer()` / `ValueBoolPointer()` /
`ValueInt64Pointer()` which return `nil` when the Terraform value is null. The API
struct uses pointer fields to distinguish "not set" from zero values.

### fromAPI Method

```go
func (m *model) fromAPI(id string, data *vnets.VNetData) {
    m.ID = types.StringValue(id)

    m.Zone = m.handleDeletedValue(data.Zone)
    m.Alias = m.handleDeletedValue(data.Alias)
    m.IsolatePorts = types.BoolPointerValue(data.IsolatePorts.PointerBool())
    m.Tag = types.Int64PointerValue(data.Tag)
    m.VlanAware = types.BoolPointerValue(data.VlanAware.PointerBool())

    // Handle pending changes from SDN (not yet applied)
    if data.Pending != nil {
        // ... override fields with pending values
    }
}
```

Pattern: use `types.StringPointerValue()`, `types.BoolPointerValue()`, etc. These
return `types.StringNull()` when the pointer is `nil`. The SDN-specific `handleDeletedValue`
and pending-state handling are domain-specific; most resources will not need them.

### Datasource

The datasource (`datasource.go`) follows the same structure but uses:

- `config.DataSource` instead of `config.Resource` in Configure.
- `datasource.Schema` with all non-ID attributes set to `Computed: true`.
- Only a `Read` method, reading from `req.Config` instead of `req.State`.

```go
var _ datasource.DataSource = &DataSource{}
var _ datasource.DataSourceWithConfigure = &DataSource{}
```

```go
func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var readModel model

    resp.Diagnostics.Append(req.Config.Get(ctx, &readModel)...)
    if resp.Diagnostics.HasError() {
        return
    }

    vnet, err := d.client.SDNVnets(readModel.ID.ValueString()).GetVnet(ctx)
    if err != nil {
        if errors.Is(err, api.ErrResourceDoesNotExist) {
            resp.Diagnostics.AddError("SDN VNet Not Found", ...)
            return
        }
        resp.Diagnostics.AddError("Unable to Read SDN VNet", err.Error())
        return
    }

    state := model{}
    state.fromAPI(readModel.ID.ValueString(), vnet)

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
```

Key difference from resource Read: datasource errors on not-found instead of removing
from state, because a datasource that cannot find its target is a configuration error.

### Acceptance Tests

```go
func TestAccResourceSDNVNet(t *testing.T) {
    t.Parallel()

    te := test.InitEnvironment(t)
```

The test environment (`fwprovider/test/test_environment.go`) provides:

- `te.AccProviders` -- pre-configured provider factories.
- `te.RenderConfig()` -- template rendering with `{{.NodeName}}` and custom vars.
- `test.ResourceAttributes()` -- bulk attribute assertion helper.
- `test.NoResourceAttributesSet()` -- assert attributes are not set.

#### Table-driven test structure

```go
tests := []struct {
    name  string
    steps []resource.TestStep
}{
```

Each entry is a named scenario with one or more steps. Steps within a scenario
run sequentially (create, then update, then import).

#### Multi-step test with import verification

```go
{"create and update vnet with simple zone", []resource.TestStep{
    {
        Config: te.RenderConfig(`
        resource "proxmox_virtual_environment_sdn_vnet" "test_vnet" {
            id    = "testv"
            zone  = proxmox_virtual_environment_sdn_zone_simple.test_zone.id
            alias = "Test VNet"
        }`),
        Check: resource.ComposeTestCheckFunc(
            test.ResourceAttributes("proxmox_virtual_environment_sdn_vnet.test_vnet", map[string]string{
                "id":    "testv",
                "zone":  "testz",
                "alias": "Test VNet",
            }),
        ),
    },
    {
        // Second step: update + import
        Config: te.RenderConfig(`...updated config...`),
        Check:            resource.ComposeTestCheckFunc(...),
        ResourceName:      "proxmox_virtual_environment_sdn_vnet.test_vnet",
        ImportStateId:     "testv",
        ImportState:       true,
        ImportStateVerify: true,
    },
}},
```

#### Running tests with ParallelTest

```go
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        resource.ParallelTest(t, resource.TestCase{
            ProtoV6ProviderFactories: te.AccProviders,
            Steps:                    tt.steps,
        })
    })
}
```

Use `resource.ParallelTest` (not `resource.Test`) for acceptance tests that can run
concurrently.

---

## Secondary Reference: Metrics Server

**Files:**

- `fwprovider/cluster/metrics/resource_metrics_server.go`
- `fwprovider/cluster/metrics/metrics_server_model.go`
- `fwprovider/cluster/metrics/resource_metrics_server_test.go`

The Metrics Server demonstrates patterns beyond VNet. Only the differences are shown.

### Sensitive Attributes

The Metrics Server schema marks secret fields as sensitive so Terraform redacts them
from plan output and state display:

```go
"influx_token": schema.StringAttribute{
    Sensitive: true,
    // ...
},
```

### Validator Combinations

The Metrics Server schema shows 18 attributes including fields specific to different
backend types (InfluxDB, Graphite, OpenTelemetry).

Notable patterns:

- `int64validator.Between(512, 65536)` for numeric range validation.
- `stringvalidator.OneOf("graphite", "influxdb", "opentelemetry")` for enum fields.
- Multiple `RequiresReplace()` fields (on `name` and `type`).

### Extensive CheckDelete in Update

When a resource has many optional fields, the Update method needs a `CheckDelete` call
for each one:

```go
var toDelete []string

attribute.CheckDelete(plan.Disable, state.Disable, &toDelete, "disable")
attribute.CheckDelete(plan.MTU, state.MTU, &toDelete, "mtu")
attribute.CheckDelete(plan.Timeout, state.Timeout, &toDelete, "timeout")
attribute.CheckDelete(plan.InfluxAPIPathPrefix, state.InfluxAPIPathPrefix, &toDelete, "api-path-prefix")
attribute.CheckDelete(plan.InfluxBucket, state.InfluxBucket, &toDelete, "bucket")
// ... 10 more CheckDelete calls
```

The third argument is always the **Proxmox API parameter name**, which often differs
from the Terraform attribute name (e.g., `"api-path-prefix"` vs `influx_api_path_prefix`).

### Bool to Int64 Conversion

The Proxmox API uses `int64` values (0/1) to represent booleans for some fields. The
Metrics Server model handles this conversion with helper functions:

```go
func boolToInt64Ptr(boolPtr *bool) *int64 {
    if boolPtr != nil {
        var result int64

        if *boolPtr {
            result = int64(1)
        } else {
            result = int64(0)
        }

        return &result
    }

    return nil
}

func int64ToBoolPtr(int64ptr *int64) *bool {
    if int64ptr != nil {
        var result bool

        if *int64ptr == 0 {
            result = false
        } else {
            result = true
        }

        return &result
    }

    return nil
}
```

These are used in the model's `importFromAPI` and `toAPIRequestBody` methods to convert
between Terraform's `types.Bool` and the API's `int64` representation:

```go
// toAPI
data.Disable = boolToInt64Ptr(m.Disable.ValueBoolPointer())

// fromAPI
m.Disable = types.BoolPointerValue(int64ToBoolPtr(data.Disable))
```

Use this pattern when the Proxmox API represents a boolean concept as 0/1 integers.

### Test: Verifying Unset Attributes

```go
test.NoResourceAttributesSet(
    "proxmox_virtual_environment_metrics_server.acc_influxdb_server",
    []string{
        "disable",
        "timeout",
        "influx_api_path_prefix",
        // ... more fields
    },
),
```

Use `test.NoResourceAttributesSet` to verify optional fields that should be null/unset.

---

## Tertiary Reference: ACL

**Files:**

- `fwprovider/access/resource_acl.go`
- `fwprovider/access/resource_acl_model.go`
- `fwprovider/access/resource_acl_test.go`

The ACL resource demonstrates patterns beyond VNet. Only the differences are shown.

### ConfigValidators for Cross-Field Validation

The ACL resource implements `ResourceWithConfigValidators` to enforce that `group_id`,
`token_id`, and `user_id` are mutually exclusive:

```go
var (
    _ resource.Resource                     = (*aclResource)(nil)
    _ resource.ResourceWithConfigure        = (*aclResource)(nil)
    _ resource.ResourceWithImportState      = (*aclResource)(nil)
    _ resource.ResourceWithConfigValidators = (*aclResource)(nil)
)
```

```go
func (r *aclResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
    return []resource.ConfigValidator{
        resourcevalidator.Conflicting(
            path.MatchRoot("group_id"),
            path.MatchRoot("token_id"),
            path.MatchRoot("user_id"),
        ),
    }
}
```

The `resourcevalidator.Conflicting` validator (from `terraform-plugin-framework-validators`)
produces a clear error when more than one of the listed attributes is set. This is
validated at plan time, before any API calls.

### Custom Import ID Parsing

When the import ID is a composite string (not a simple API identifier), parse it in
`ImportState`:

```go
func (r *aclResource) ImportState(
    ctx context.Context,
    req resource.ImportStateRequest,
    resp *resource.ImportStateResponse,
) {
    model, err := parseACLResourceModelFromID(req.ID)
    if err != nil {
        resp.Diagnostics.AddError("Unable to import ACL",
            "failed to parse ID: "+err.Error())
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
```

The parsing function splits the composite ID and populates the model:

```go
const aclIDFormat = "{path}?{group|user@realm|user@realm!token}?{role}"

func parseACLResourceModelFromID(id string) (*aclResourceModel, error) {
    parts := strings.Split(id, "?")
    if len(parts) != 3 {
        return nil, fmt.Errorf("invalid ACL resource ID format %#v, expected %v", id, aclIDFormat)
    }

    path := parts[0]
    entityID := parts[1]
    roleID := parts[2]

    model := &aclResourceModel{
        ID:      types.StringValue(id),
        GroupID: types.StringNull(),
        Path:    path,
        RoleID:  roleID,
        TokenID: types.StringNull(),
        UserID:  types.StringNull(),
    }

    switch {
    case strings.Contains(entityID, "!"):
        model.TokenID = types.StringValue(entityID)
    case strings.Contains(entityID, "@"):
        model.UserID = types.StringValue(entityID)
    default:
        model.GroupID = types.StringValue(entityID)
    }

    return model, nil
}
```

Define the format as a constant so error messages are self-documenting. Initialize
mutually exclusive fields to `types.StringNull()` to avoid nil-vs-empty confusion.

### Validation Tests

The ACL test file includes a dedicated validation test that verifies cross-field
constraints without making API calls:

```go
func TestAccAcl_Validators(t *testing.T) {
    t.Parallel()

    te := test.InitEnvironment(t)

    resource.UnitTest(t, resource.TestCase{
        ProtoV6ProviderFactories: te.AccProviders,
        Steps: []resource.TestStep{
            {
                PlanOnly: true,
                Config: `resource "proxmox_virtual_environment_acl" "test" {
                    group_id = "test"
                    path = "/"
                    role_id = "test"
                    token_id = "test"
                }`,
                ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Combination`),
            },
            // ... more validation cases
        },
    })
}
```

Key patterns:

- `resource.UnitTest` instead of `resource.Test` -- no real infrastructure needed.
- `PlanOnly: true` -- validates during plan, no apply.
- `ExpectError` with a regex -- asserts the expected validation error.

---

## Advanced Patterns

These patterns appear in more complex resources. Refer to them when your resource
outgrows the simple 3-file pattern.

| Pattern | Where to Find | When to Use |
| ------- | ------------- | ----------- |
| Generic factory resource | `fwprovider/cluster/sdn/zone/resource_generic.go` | Multiple resources sharing CRUD logic with different schemas |
| Go generics | `fwprovider/storage/resource_generic.go` with `storageResource[T, M]` | API structure allows generic handling across similar resource types |
| Modular sub-schemas | `fwprovider/nodes/vm/` with sub-packages `cdrom/`, `cpu/`, `memory/`, `rng/`, `vga/` | Very large resources that benefit from splitting into sub-packages |
| `ValidateConfig` method | `fwprovider/cluster/sdn/subnet/resource.go` | Complex validation that requires reading multiple attributes |
| Shared model | `fwprovider/cluster/acme/plugin_model.go` shared by `resource_acme_dns_plugin.go` and `resource_acme_account.go` | Multiple resources sharing one model file with common types |
| Custom `stringset.Value` type | `fwprovider/types/stringset/` | Comma-separated list attributes (e.g., node lists) |

---

## Checklist for New Resource Implementation

Use this checklist when implementing a new framework resource. Each item links to the
relevant pattern in this document.

### Setup

- [ ] Create package directory under `fwprovider/` following domain hierarchy
- [ ] Create `resource.go`, `model.go`, and `resource_test.go` (3-file pattern)
- [ ] Add `datasource.go` if a read-only datasource is also needed

### resource.go

- [ ] Interface compliance assertions (`var _ resource.Resource = ...`)
- [ ] Resource struct with typed client field
- [ ] Constructor returning zero-value struct (`NewResource()`)
- [ ] `Metadata` returning `req.ProviderTypeName + "_your_resource"`
- [ ] `Configure` with nil guard, `config.Resource` type assertion, client extraction
- [ ] `Schema` with descriptions, validators, and plan modifiers on immutable fields
- [ ] `Create`: plan -> toAPI -> API call -> read back -> set state
- [ ] `Read`: handle `api.ErrResourceDoesNotExist` with `RemoveResource`
- [ ] `Update`: `CheckDelete` for every optional field, then update + read back
- [ ] `Delete`: ignore `api.ErrResourceDoesNotExist`
- [ ] `ImportState`: simple ID pass-through or custom parsing

### model.go

- [ ] Model struct with `tfsdk` tags matching schema attributes exactly
- [ ] `types.*` for optional fields, plain Go types only for always-present fields
- [ ] `toAPI()` method using `Value*Pointer()` for optional fields
- [ ] `fromAPI()` method using `types.*PointerValue()` for optional fields

### resource_test.go

- [ ] `//go:build acceptance || all` build tag
- [ ] `t.Parallel()` at top of test function
- [ ] `test.InitEnvironment(t)` for provider factories and config rendering
- [ ] Table-driven test structure with named scenarios
- [ ] Create, update, and import steps
- [ ] `test.ResourceAttributes` for bulk assertions
- [ ] `resource.ParallelTest` for concurrent execution
- [ ] Validation test with `PlanOnly: true` and `ExpectError` if applicable

### Registration

- [ ] Register resource in provider's resource list
- [ ] Register datasource in provider's datasource list (if applicable)
- [ ] Run `make docs` to regenerate documentation if schema changed

### Verification

- [ ] `make build` passes
- [ ] `make lint` shows 0 issues
- [ ] `make test` passes (unit tests)
- [ ] `./testacc TestAccYourResource` passes (acceptance tests)
- [ ] API calls verified with mitmproxy
