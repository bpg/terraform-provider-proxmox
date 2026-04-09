# ADR: Reference Examples for New Framework Resource Implementations

## Purpose

This document identifies the cleanest framework provider implementations in the codebase and annotates their key patterns. It serves as a copy-from guide for contributors adding new resources or datasources.

Three implementations were selected after auditing 47 resources and 36 datasources in `fwprovider/`:

| Rank      | Implementation   | Why Selected                                                                               |
|-----------|------------------|--------------------------------------------------------------------------------------------|
| Primary   | **SDN VNet**     | Simplest clean implementation, perfect 3-file pattern (resource, model, datasource)        |
| Secondary | **Replication**  | Many optional fields, split create/update methods, perfect audit score (66/66)             |
| Tertiary  | **Backup Job**   | `ConfigValidators`, comma-separated-to-list, nested objects, shared fillCommonFields       |

Start with SDN VNet for any new resource. Refer to Replication when you need many optional fields or split create/update shapes. Refer to Backup Job for `ConfigValidators`, comma-separated-to-list conversion, or nested objects.

---

## How to Use This Document

| I need to...                                 | Start here                                                                                                          |
|----------------------------------------------|---------------------------------------------------------------------------------------------------------------------|
| Create a basic resource                      | [SDN VNet walkthrough](#primary-reference-sdn-vnet)                                                                 |
| Handle many optional fields                  | [Replication walkthrough](#secondary-reference-replication)                                                         |
| Handle split create/update patterns          | [Replication walkthrough](#secondary-reference-replication)                                                         |
| Add cross-field validation                   | [Backup Job ConfigValidators](#configvalidators-for-mutually-exclusive-attributes)                                  |
| Convert comma-separated API values to lists  | [Backup Job comma-sep-to-list](#comma-separated-api-values-as-terraform-lists)                                      |
| Parse composite import IDs                   | [SDN Subnet ImportState](../../fwprovider/cluster/sdn/subnet/resource.go)                                           |
| Need nested attributes?                      | [SDN Subnet `SingleNestedAttribute`](../../fwprovider/cluster/sdn/subnet/resource.go)                               |
| Add a datasource                             | [SDN VNet datasource](#datasource) and [Replication datasource](../../fwprovider/cluster/replication/datasource.go) |
| Find advanced patterns (generics, factories) | [Advanced Patterns table](#advanced-patterns)                                                                       |

---

## Primary Reference: SDN VNet

**Files:**

- `fwprovider/cluster/sdn/vnet/resource.go` -- CRUD + schema
- `fwprovider/cluster/sdn/vnet/model.go` -- Terraform/API mapping
- `fwprovider/cluster/sdn/vnet/datasource.go` -- Read-only datasource
- `fwprovider/cluster/sdn/vnet/resource_test.go` -- Acceptance tests

### Interface Compliance Assertions

Every resource file starts with compile-time interface checks. These ensure your struct satisfies all required interfaces before runtime.

```go
var (
    _ resource.Resource                = &Resource{}
    _ resource.ResourceWithConfigure   = &Resource{}
    _ resource.ResourceWithImportState = &Resource{}
)
```

If your resource implements additional interfaces (e.g., `ResourceWithConfigValidators`), add them here. Using `&Resource{}` (pointer literal) is the standard style; the ACL example uses `(*aclResource)(nil)` which is equivalent.

### Resource Struct and Constructor

The struct holds the typed API client. The constructor returns the zero-value struct (the client is injected later via `Configure`).

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
            // User-provided ID: Required + RequiresReplace + validators.
            // For server-assigned IDs, use attribute.ResourceID() instead
            // (Computed + UseStateForUnknown).
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
- `Sensitive: true` on credential fields (tokens, passwords, API keys) to redact from plan output and state display.
- `attribute.ResourceID()` helper (from `fwprovider/attribute/`) generates a computed ID attribute with `UseStateForUnknown()`. Use this when the resource ID is assigned by the server, not provided by the user.

### Create

The full create flow: plan -> toAPI -> API call -> read back -> fromAPI -> set state.

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

> [!IMPORTANT]
> Always read back from the API after create. Never save plan data directly to state — the API may normalize values or populate computed fields.

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

The critical pattern: when the API returns `api.ErrResourceDoesNotExist`, call `resp.State.RemoveResource(ctx)` to tell Terraform the resource was deleted outside of Terraform. This triggers a plan to recreate it.

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

`attribute.CheckDelete` (defined in `fwprovider/attribute/attribute.go`) adds the API field name to the delete list when the plan value is null but the state value is not. This tells the Proxmox API to remove the field.

> [!NOTE]
> The third argument to `CheckDelete` is the **Proxmox API parameter name**, which often differs from the Terraform attribute name (e.g., `"isolate-ports"` vs `isolate_ports`, `"api-path-prefix"` vs `influx_api_path_prefix`).

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

Ignore `api.ErrResourceDoesNotExist` on delete -- the resource is already gone, which is the desired end state.

> [!CAUTION]
> If your resource's domain client uses `retry.NewTaskOperation` for delete, ensure the retry predicate excludes `ErrResourceDoesNotExist`. This error can arrive via HTTP 500 (not just 404), so `IsTransientAPIError` alone will incorrectly retry it. Use the combined predicate:
>
> ```go
> retry.WithRetryIf(func(err error) bool {
>     return retry.IsTransientAPIError(err) && !errors.Is(err, api.ErrResourceDoesNotExist)
> })
> ```
>
> See [ADR-005: Error Handling — Retry Policies](005-error-handling.md#retry-policies).

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

For simple resources where the import ID maps directly to the API lookup key, this pattern is sufficient. For composite IDs, see [SDN Subnet](../../fwprovider/cluster/sdn/subnet/resource.go) which parses `vnet-id/subnet-id` format.

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

    data.Zone = attribute.StringPtrFromValue(m.Zone)
    data.Alias = attribute.StringPtrFromValue(m.Alias)
    data.IsolatePorts = attribute.CustomBoolPtrFromValue(m.IsolatePorts)
    data.Tag = attribute.Int64PtrFromValue(m.Tag)
    data.VlanAware = attribute.CustomBoolPtrFromValue(m.VlanAware)

    return data
}
```

Pattern: use the `attribute` package helpers (`StringPtrFromValue`, `Int64PtrFromValue`, `CustomBoolPtrFromValue`) which return `nil` for both null and unknown values. This is safer than raw `Value*Pointer()` methods, which return `&""` / `&0` for unknown values. The API struct uses pointer fields to distinguish "not set" from zero values.

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

Pattern: use `types.StringPointerValue()`, `types.BoolPointerValue()`, etc. These return `types.StringNull()` when the pointer is `nil`. The SDN-specific `handleDeletedValue` and pending-state handling are domain-specific; most resources will not need them.

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

> [!IMPORTANT]
> Datasources must error on not-found — never call `RemoveResource`. A datasource that cannot find its target is a configuration error, not a drift signal.

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

Each entry is a named scenario with one or more steps. Steps within a scenario run sequentially (create, then update, then import).

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

Use `resource.ParallelTest` (not `resource.Test`) for acceptance tests that can run concurrently. See the [Replication tests](#acceptance-tests-paralleltest-and-table-driven-subtests) for the canonical `ParallelTest` example.

-> **Note:** The VNet test file currently uses `resource.Test` rather than `resource.ParallelTest`. This is a known deviation (D6=2 in the audit). New resources should follow the `ParallelTest` pattern shown in the Replication reference.

### Test: Verifying Unset Attributes

```go
test.NoResourceAttributesSet(
    "proxmox_virtual_environment_sdn_vnet.test_vnet",
    []string{
        "alias",
        "isolate_ports",
        "tag",
        // ... optional fields that should be null/unset
    },
),
```

Use `test.NoResourceAttributesSet` to verify optional fields that should be null/unset after a config change removes them. This complements `test.ResourceAttributes` for positive assertions.

---

## Secondary Reference: Replication

**Files:**

- `fwprovider/cluster/replication/resource.go` -- CRUD + schema
- `fwprovider/cluster/replication/model.go` -- Terraform/API mapping
- `fwprovider/cluster/replication/datasource.go` -- Read-only datasource
- `fwprovider/cluster/replication/resource_test.go` -- Acceptance tests

The Replication resource scored 66/66 (Grade A) in the framework audit -- the only perfect score. It demonstrates patterns beyond VNet. Only the differences are shown.

### Split Create/Update API Methods

The Proxmox API uses different request shapes for creating and updating a replication job. The create endpoint requires `target` and `type` (immutable), while the update endpoint does not accept them. The model reflects this with separate methods:

```go
// toAPICreate populates all fields including immutable ones (target, type).
func (m *model) toAPICreate() *replications.ReplicationCreate {
    data := &replications.ReplicationCreate{}
    data.ID = m.ID.ValueString()
    data.Target = m.Target.ValueString()  // immutable: only sent at creation
    data.Type = m.Type.ValueString()      // immutable: only sent at creation
    data.Comment = attribute.StringPtrFromValue(m.Comment)
    data.Disable = attribute.CustomBoolPtrFromValue(m.Disable)
    data.Rate = attribute.Float64PtrFromValue(m.Rate)
    data.Schedule = attribute.StringPtrFromValue(m.Schedule)
    return data
}

// toAPIUpdate omits immutable fields -- API rejects them on update.
func (m *model) toAPIUpdate() *replications.ReplicationUpdate {
    data := &replications.ReplicationUpdate{}
    data.ID = m.ID.ValueString()
    data.Comment = attribute.StringPtrFromValue(m.Comment)
    data.Disable = attribute.CustomBoolPtrFromValue(m.Disable)
    data.Rate = attribute.Float64PtrFromValue(m.Rate)
    data.Schedule = attribute.StringPtrFromValue(m.Schedule)
    return data
}
```

The CRUD methods call the matching model method:

```go
// Create uses toAPICreate (includes target + type).
repl := plan.toAPICreate()
err := r.client.Replication(plan.ID.ValueString()).CreateReplication(ctx, repl)

// Update uses toAPIUpdate (excludes target + type).
repl := plan.toAPIUpdate()
// ... CheckDelete, then set repl.Delete = toDelete ...
err := r.client.Replication(plan.ID.ValueString()).UpdateReplication(ctx, repl)
```

**When to split:** If the API's create and update endpoints accept different fields (common when some fields are immutable after creation), use `toAPICreate()` and `toAPIUpdate()`. If both endpoints accept the same shape, a single `toAPI()` is fine (see VNet).

### RequiresReplace on Immutable Fields

Three fields are immutable after creation. Each gets `RequiresReplace()` so Terraform will destroy and recreate the resource if these change:

```go
"id": schema.StringAttribute{
    Required: true,
    PlanModifiers: []planmodifier.String{
        stringplanmodifier.RequiresReplace(),
    },
    Validators: []validator.String{
        stringvalidator.RegexMatches(
            regexp.MustCompile(`^[0-9]+-[0-9]+$`),
            "id must be <GUEST>-<JOBNUM>",
        ),
    },
},
"target": schema.StringAttribute{
    Required: true,
    PlanModifiers: []planmodifier.String{
        stringplanmodifier.RequiresReplace(),
    },
},
```

Note how `RequiresReplace` aligns with the split API methods: the fields that require replace are exactly the fields excluded from `toAPIUpdate()`.

### Extensive CheckDelete in Update

Every mutable optional field needs a `CheckDelete` call:

```go
var toDelete []string

attribute.CheckDelete(plan.Comment, state.Comment, &toDelete, "comment")
attribute.CheckDelete(plan.Disable, state.Disable, &toDelete, "disable")
attribute.CheckDelete(plan.Rate, state.Rate, &toDelete, "rate")
attribute.CheckDelete(plan.Schedule, state.Schedule, &toDelete, "schedule")

repl.Delete = toDelete  // API uses `delete` param to remove fields
```

### BoolDefault for Optional+Computed Bool

The `disable` field has a default value, so null in config means `false` (not "unset"):

```go
"disable": schema.BoolAttribute{
    Optional:    true,
    Computed:    true,
    Default:     booldefault.StaticBool(false),
    Description: "Flag to disable/deactivate this replication.",
},
```

### fromAPI: Handling API Quirks

The API omits `disable` from the response when the value is `false`. The `fromAPI` method normalizes this:

```go
func (m *model) fromAPI(id string, data *replications.ReplicationData) {
    m.ID = types.StringValue(id)
    m.Target = types.StringValue(data.Target)
    m.Type = types.StringValue(data.Type)
    m.Comment = types.StringPointerValue(data.Comment)

    // API quirk: `disable` is omitted when false, not returned as false.
    if v := data.Disable.PointerBool(); v != nil {
        m.Disable = types.BoolValue(*v)
    } else {
        m.Disable = types.BoolValue(false)
    }

    m.Rate = types.Float64PointerValue(data.Rate)
    m.Schedule = types.StringPointerValue(data.Schedule)
}
```

### Short-Name Alias with MoveState

Per [ADR-007](007-resource-type-name-migration.md), existing Framework resources are being migrated from the verbose `proxmox_virtual_environment_*` prefix to the shorter `proxmox_*` prefix. Both names coexist during the transition. `MoveState` enables users to migrate their state using Terraform's `moved` block (Terraform >= 1.8) without destroying and recreating resources.

> [!NOTE]
> **Do not add short-name aliases or `MoveState` to new resources or datasources.** New resources should use the `proxmox_` prefix directly in their `Metadata()` method. This pattern is only for migrating existing `proxmox_virtual_environment_*` resources.

The alias embeds the original resource and overrides only `Metadata` and `Schema`:

```go
type resourceShort struct {
    Resource  // embed: inherits all CRUD methods
}

func (r *resourceShort) Metadata(...) {
    resp.TypeName = "proxmox_replication"  // short name
}

func (r *resourceShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    r.Resource.Schema(ctx, req, resp)
    resp.Schema.DeprecationMessage = ""  // remove deprecation from short name
}

func (r *resourceShort) MoveState(ctx context.Context) []resource.StateMover {
    schemaResp := &resource.SchemaResponse{}
    r.Schema(ctx, resource.SchemaRequest{}, schemaResp)
    return []resource.StateMover{
        migration.PrefixMoveState("proxmox_virtual_environment_replication", &schemaResp.Schema),
    }
}
```

### Acceptance Tests: ParallelTest and Table-Driven Subtests

The test file demonstrates the preferred pattern with `resource.ParallelTest`:

```go
tests := []struct {
    name  string
    steps []resource.TestStep
}{
    {"create and update minimal replication", func() []resource.TestStep {
        cid := newCID()  // random ID avoids collisions between parallel tests
        return []resource.TestStep{ /* create -> add fields -> import */ }
    }()},
    {"replication fields deletion", func() []resource.TestStep {
        // Step 1: create with all fields. Step 2: remove all optional fields.
        return []resource.TestStep{ /* ... */ }
    }()},
    {"replication fields deletion and re-addition", func() []resource.TestStep {
        // Step 1: all fields. Step 2: remove all. Step 3: re-add all.
        return []resource.TestStep{ /* ... */ }
    }()},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        resource.ParallelTest(t, resource.TestCase{
            ProtoV6ProviderFactories: te.AccProviders,
            Steps:                    tt.steps,
        })
    })
}
```

Key patterns:

- **Random IDs** (`newCID()`) allow parallel subtests to run without resource collisions.
- **Field deletion + re-addition** tests verify `CheckDelete` works end-to-end.
- **Import verification** combined with the last step of multi-step tests.

### Validation Tests (Unit, No Infrastructure)

Validation rules are tested separately with `resource.UnitTest`:

```go
func TestUnitResourceReplication_Validators(t *testing.T) {
    t.Parallel()
    te := test.InitEnvironment(t)

    resource.UnitTest(t, resource.TestCase{
        ProtoV6ProviderFactories: te.AccProviders,
        Steps: []resource.TestStep{
            {
                PlanOnly: true,
                Config:   `resource "proxmox_replication" "test" { id = "invalidid" ... }`,
                ExpectError: regexp.MustCompile(`id must be <GUEST>-<JOBNUM>`),
            },
            // ... regex edge cases, enum validators, missing required fields
        },
    })
}
```

Named `TestUnit*` (not `TestAcc*`) to signal it runs without acceptance infrastructure.

---

## Tertiary Reference: Backup Job

**Files:**

- `fwprovider/cluster/backup/resource.go` -- CRUD + schema + ConfigValidators
- `fwprovider/cluster/backup/model.go` -- Terraform/API mapping
- `fwprovider/cluster/backup/resource_test.go` -- Acceptance tests

The Backup Job resource scored 64/66 (Grade A) in the framework audit. It demonstrates patterns beyond VNet and Replication. Only the differences are shown.

### ConfigValidators for Mutually Exclusive Attributes

The Backup Job implements `ResourceWithConfigValidators` to enforce that `all`, `vmid`, and `pool` are mutually exclusive selection modes:

```go
var (
    _ resource.Resource                     = &backupJobResource{}
    _ resource.ResourceWithConfigure        = &backupJobResource{}
    _ resource.ResourceWithImportState      = &backupJobResource{}
    _ resource.ResourceWithConfigValidators = &backupJobResource{}
)
```

```go
func (r *backupJobResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
    return []resource.ConfigValidator{
        resourcevalidator.Conflicting(
            path.MatchRoot("all"),
            path.MatchRoot("vmid"),
        ),
        resourcevalidator.Conflicting(
            path.MatchRoot("all"),
            path.MatchRoot("pool"),
        ),
        resourcevalidator.Conflicting(
            path.MatchRoot("vmid"),
            path.MatchRoot("pool"),
        ),
    }
}
```

Import `resourcevalidator` from `terraform-plugin-framework-validators/resourcevalidator`. Pairwise `Conflicting` calls produce more precise error messages than a single 3-way call. Validated at plan time, before any API calls.

### Comma-Separated API Values as Terraform Lists

The Proxmox API represents `vmid` as a comma-separated string (`"100,101,102"`), but the Terraform schema exposes it as a `types.List`:

```go
"vmid": schema.ListAttribute{
    Description: "A list of guest VM/CT IDs to include in the backup job.",
    Optional:    true,
    ElementType: types.StringType,
},
```

**toAPI (join):** Extract list elements and join (use `attribute.IsDefined()` for list/set fields where typed helpers are not available):

```go
if attribute.IsDefined(m.VMIDs) {
    var vmids []string
    d := m.VMIDs.ElementsAs(ctx, &vmids, false)
    diags.Append(d...)
    if !d.HasError() && len(vmids) > 0 {
        vmidStr := strings.Join(vmids, ",")
        common.VMID = &vmidStr
    }
}
```

**fromAPI (split):** Parse comma-separated string back to `types.List`:

```go
if data.VMID != nil && *data.VMID != "" {
    ids := strings.Split(*data.VMID, ",")
    vmidValues := make([]attr.Value, len(ids))
    for i, id := range ids {
        vmidValues[i] = types.StringValue(strings.TrimSpace(id))
    }
    listVal, d := types.ListValue(types.StringType, vmidValues)
    diags.Append(d...)
    m.VMIDs = listVal
} else {
    m.VMIDs = types.ListNull(types.StringType)
}
```

> [!IMPORTANT]
> When the API value is nil or empty, always set the Terraform field to `types.ListNull(types.StringType)` (not an empty list). This prevents spurious diffs between "not set" and "set to empty."

See [ADR-004](004-schema-design-conventions.md#comma-separated-api-values--terraform-lists) for the design rationale.

### Shared fillCommonFields for Create and Update

When create and update API bodies share a common embedded struct, extract the shared field mapping into a single method to avoid duplicating 20+ field mappings:

```go
func (m *backupJobModel) toAPICreate(ctx context.Context, diags *diag.Diagnostics) *backup.CreateRequestBody {
    body := &backup.CreateRequestBody{}
    body.ID = m.ID.ValueString()
    body.Schedule = m.Schedule.ValueString()
    m.fillCommonFields(ctx, &body.RequestBodyCommon, diags)
    return body
}

func (m *backupJobModel) toAPIUpdate(ctx context.Context, state *backupJobModel, diags *diag.Diagnostics) *backup.UpdateRequestBody {
    body := &backup.UpdateRequestBody{}
    body.Schedule = attribute.StringPtrFromValue(m.Schedule)
    m.fillCommonFields(ctx, &body.RequestBodyCommon, diags)
    // ... CheckDelete calls, body.Delete = toDelete ...
    return body
}
```

Compare with Replication's fully separate `toAPICreate()`/`toAPIUpdate()` -- use `fillCommonFields` when the overlap is large, and fully separate methods when the overlap is small.

### Nested Object Attributes

Backup Job uses `SingleNestedAttribute` for structured sub-objects like `fleecing` and `performance`. Each nested object needs a separate model struct and an `attrTypes()` function:

```go
type fleecingModel struct {
    Enabled types.Bool   `tfsdk:"enabled"`
    Storage types.String `tfsdk:"storage"`
}

func fleecingAttrTypes() map[string]attr.Type {
    return map[string]attr.Type{
        "enabled": types.BoolType,
        "storage": types.StringType,
    }
}
```

In `fromAPI`, convert to `types.Object` using `types.ObjectValueFrom`:

```go
if data.Fleecing != nil {
    fleecingVal := fleecingModel{
        Enabled: types.BoolPointerValue(data.Fleecing.Enabled.PointerBool()),
        Storage: types.StringPointerValue(data.Fleecing.Storage),
    }
    obj, d := types.ObjectValueFrom(ctx, fleecingAttrTypes(), fleecingVal)
    diags.Append(d...)
    m.Fleecing = obj
} else {
    m.Fleecing = types.ObjectNull(fleecingAttrTypes())
}
```

In `toAPI`, extract the nested model from `types.Object`:

```go
if !m.Fleecing.IsNull() && !m.Fleecing.IsUnknown() {
    var fleecing fleecingModel
    d := m.Fleecing.As(ctx, &fleecing, basetypes.ObjectAsOptions{})
    diags.Append(d...)
    if !d.HasError() {
        // Inner fields are always known here (outer guard ensures the object is not null/unknown),
        // so raw Value*Pointer() is safe. Using helpers would also work.
        common.Fleecing = &backup.FleecingConfig{
            Enabled: attribute.CustomBoolPtrFromValue(fleecing.Enabled),
            Storage: attribute.StringPtrFromValue(fleecing.Storage),
        }
    }
}
```

> [!NOTE]
> For **custom composite import ID parsing** (e.g., `vnet-id/subnet-id`), see [SDN Subnet](../../fwprovider/cluster/sdn/subnet/resource.go) which splits `req.ID` on `/` and looks up the resource using both parts.

---

## Advanced Patterns

These patterns appear in more complex resources. Refer to them when your resource
outgrows the simple 3-file pattern.

| Pattern | Where to Find | When to Use |
| ------- | ------------- | ----------- |
| Generic factory resource | `fwprovider/cluster/sdn/zone/resource_generic.go` | Multiple resources sharing CRUD logic with different schemas |
| Go generics | `fwprovider/storage/resource_generic.go` with `storageResource[T, M]` | API structure allows generic handling across similar resource types |
| Modular sub-schemas | `fwprovider/nodes/vm/` with sub-packages `cdrom/`, `cpu/`, `memory/`, `rng/`, `vga/` | Very large resources that benefit from splitting into sub-packages |
| Opt-in management | `fwprovider/nodes/clonedvm/resource.go` with `optInManagedAttribute()` | Cloned resources where only explicitly listed fields are managed |
| `ValidateConfig` method | `fwprovider/cluster/sdn/subnet/resource.go` | Complex validation that requires reading multiple attributes |
| Shared model | `fwprovider/cluster/acme/plugin_model.go` shared by `resource_acme_dns_plugin.go` and `resource_acme_account.go` | Multiple resources sharing one model file with common types |
| Custom `stringset.Value` type | `fwprovider/types/stringset/` | Comma-separated list attributes (e.g., node lists) |
| Comma-separated to List/Map | fwprovider/cluster/backup/model.go with toAPICreate()/fromAPI() | API fields using comma-separated strings exposed as Terraform lists or maps |

---

## Compliance Scoring Rubric

Each Framework resource is scored on 7 dimensions (0-3), weighted by impact (max 66). Datasources skip D1/D4 (max 42). This is the canonical methodology for compliance assessment.

**Grades (resources):** A >=58, B >=48, C >=36, D <36. **Grades (datasources):** A >=37, B >=30, C >=23, D <23. **Critical failure:** D1=0 or D3=0 caps grade at D regardless of total.

### D1: CRUD Lifecycle (Weight 5)

Checkpoints: (a) read-back after Create, (b) read-back after Update, (c) Read handles `ErrResourceDoesNotExist` with `RemoveResource`, (d) Delete ignores not-found, (e) ImportState implemented, (f) `CheckDelete` for every optional field.

**0** = multiple critical violations. **1** = missing read-back or RemoveResource. **2** = all critical patterns present, minor gaps. **3** = fully canonical. Singletons (ClusterOptions, NodeFirewallOptions) may mark c/d N/A.

### D2: Error Messages (Weight 2)

Checkpoints: (a) consistent pattern within the resource, (b) summary includes resource name, (c) Title Case, (d) detail is `err.Error()` not double-wrapped, (e) distinct message per CRUD operation, (f) read-back errors distinguishable, (g) new code uses ADR-005 `"Unable to [Verb] [Resource]"`.

**0** = inconsistent or malformed. **1** = consistent but missing resource name. **2** = good with minor gaps. **3** = fully consistent.

### D3: Schema Design (Weight 4)

Checkpoints: (a) `types.*` for all fields, (b) correct Required/Optional/Computed, (c) `RequiresReplace` on immutable fields, (d) validators on constrained fields, (e) `Sensitive: true` on credentials, (f) all attributes have Description, (g) comma-separated API values as List/Set.

**0** = raw Go types, no validators. **1** = types.\* but missing validators/descriptions. **2** = mostly compliant. **3** = full compliance.

### D4: Model Conversion (Weight 3)

Checkpoints: (a) `toAPI()`/`fromAPI()` naming (or `toAPICreate()`/`toAPIUpdate()` for split shapes), (b) `attribute.*PtrFromValue()` helpers in toAPI (not raw `Value*Pointer()`), (c) `types.*PointerValue()` in fromAPI, (d) `attribute.CustomBoolPtrFromValue()` for bool-int64, (e) model in separate file.

**0** = inline conversion, no model file. **1** = non-standard names. **2** = correct patterns but legacy names. **3** = fully canonical.

### D5: File Organization (Weight 2)

Checkpoints: (a) 3-file pattern, (b) domain hierarchy, (c) interface assertions, (d) Configure nil guard, (e) zero-value constructor, (f) short-name alias with MoveState (for migrated resources; new `proxmox_` resources skip this).

**0** = major structural violations. **1** = model inlined or missing assertions. **2** = missing one element. **3** = full compliance.

### D6: Testing Quality (Weight 3)

Checkpoints: (a) test file exists, (b) `//go:build acceptance || all`, (c) `t.Parallel()`, (d) `test.InitEnvironment`, (e) table-driven with named scenarios, (f) `resource.ParallelTest`, (g) import test, (h) functional coverage (create + update + field removal), (i) validation test if applicable, (j) `test.ResourceAttributes` / `test.NoResourceAttributesSet`.

**0** = no tests. **1** = basic test, missing import/parallel. **2** = create + update + import with ParallelTest. **3** = comprehensive.

### D7: Advanced Correctness (Weight 3)

Checkpoints: (a) ImportState errors on not-found (not RemoveResource), (b) delete retry excludes `ErrResourceDoesNotExist`, (c) all `State.Set()` wrapped in `Diagnostics.Append()`, (d) domain client uses `%w`, (e) datasource uses `config.DataSource` and errors on not-found.

**0** = multiple violations. **1** = notable gaps. **2** = minor issues. **3** = all correct.

---

## Checklist for New Resource Implementation

### Merge Requirements (must pass before merge)

Use this checklist for the minimum viable implementation. All items must be complete.

#### Setup

- [ ] Create package directory under `fwprovider/` following domain hierarchy ([ADR-003](003-resource-file-organization.md))
- [ ] Create `resource.go`, `model.go`, and `resource_test.go` (3-file pattern)

#### resource.go

- [ ] Interface compliance assertions (`var _ resource.Resource = ...`)
- [ ] Resource struct with typed client field
- [ ] Constructor returning zero-value struct (`NewResource()`)
- [ ] `Metadata` returning type name with `proxmox_` prefix (per ADR-007)
- [ ] `Configure` with nil guard and `config.Resource` type assertion
- [ ] `Schema` with descriptions on all attributes, validators on constrained fields, and `RequiresReplace()` on immutable fields
- [ ] Create: plan → toAPI() (or toAPICreate()) → API call → **read back from API** → set state
- [ ] `Read`: handle `api.ErrResourceDoesNotExist` with `RemoveResource`
- [ ] `Update`: `CheckDelete` for every optional field, then update + **read back** → set state
- [ ] `Delete`: ignore `api.ErrResourceDoesNotExist`
- [ ] `ImportState`: implemented (simple or composite ID parsing)
- [ ] All `resp.State.Set()` calls wrapped in `resp.Diagnostics.Append()`

#### model.go

- [ ] Model struct with `tfsdk` tags matching schema attributes exactly
- [ ] `types.*` for all optional fields (never raw Go types)
- [ ] toAPI() (or toAPICreate()/toAPIUpdate()) using `attribute.*PtrFromValue()` helpers for optional fields
- [ ] `fromAPI()` using `types.*PointerValue()` for optional fields

#### resource_test.go

- [ ] Build tag: `//go:build acceptance || all`
- [ ] `t.Parallel()` at top of test function
- [ ] `test.InitEnvironment(t)` for provider factories
- [ ] At least create, update, and import test steps
- [ ] `resource.ParallelTest` (not `resource.Test`)

#### Registration & Build

- [ ] Resource registered in `fwprovider/provider.go`
- [ ] `make build` and `make lint` pass
- [ ] `make test` passes (unit tests)
- [ ] Acceptance tests pass
- [ ] `make docs` regenerates documentation

### Grade A Target (recommended, can be follow-up PR)

These items bring the resource to full compliance (D6=3, Grade A in the scoring rubric).

- [ ] `Sensitive: true` on credential fields (tokens, passwords, API keys)
- [ ] Table-driven test structure with named scenarios
- [ ] Validation test with `PlanOnly: true` and `ExpectError` (if applicable)
- [ ] Field removal test (verifies `CheckDelete` end-to-end)
- [ ] `test.ResourceAttributes` / `test.NoResourceAttributesSet` for bulk assertions
- [ ] Bool-to-Int64 conversions use `proxmoxtypes.CustomBoolPtr()` / `.PointerBool()`
- [ ] Domain client delete retry predicate excludes `ErrResourceDoesNotExist` (see [ADR-005](005-error-handling.md#retry-policies))
- [ ] API calls verified with mitmproxy
- [ ] Datasource added (if applicable) with `config.DataSource` and not-found error handling

---

## Known Deviations and Migration Status

As of 2026-03-27, a comprehensive audit scored all 47 Framework resources. Results:

**Grade Distribution:** 11 A (23%) -- 17 B (36%) -- 16 C (34%) -- 3 D (6%)

### Tier 1: Critical (Grade D -- need rework)

| Resource    | Score | Primary Issues                                              |
|-------------|-------|-------------------------------------------------------------|
| UserToken   | 31/66 | Bare State.Set calls, no read-back, missing RemoveResource  |
| ACMEAccount | 34/66 | No toAPI/fromAPI, no CheckDelete, credentials not Sensitive |
| ACMEPlugin  | 34/66 | Incomplete read-back, no CheckDelete, copy-paste bug        |

### Tier 2: High-ROI Fixes (fix shared code -- many resources improve)

| Fix                                                     | Affected Resources     |
|---------------------------------------------------------|------------------------|
| Storage generic base: add CheckDelete, fix error prefix | 7 storage resources    |
| HW Mapping shared.go: fix silent error swallowing       | 3 HW mapping resources |
| Error message standardization codemod                   | ~33 resources          |

### Weakest Dimension: D6 Testing (avg 1.66/3)

11 resources have zero acceptance tests: HAGroup, HAResource, HARule,
CIFS/LVM/LVMThin/ZFS Storage, Fabric OpenFabric/OSPF, FabricNode OpenFabric/OSPF.

### Strongest Dimension: D3 Schema Design (avg 2.87/3)

Nearly universal compliance with types.*, validators, and descriptions.

### Reference Implementation Scores

| Resource      | Score | Grade | Notes                                                         |
|---------------|-------|-------|---------------------------------------------------------------|
| Replication   | 66/66 | A     | Perfect score -- additional reference for split create/update |
| BackupJob     | 64/66 | A     | Excellent -- uses `proxmox_` prefix natively per ADR-007      |
| SDN VNet      | 63/66 | A     | Primary reference (gold standard)                             |
| SDN Subnet    | 61/66 | A     | Good reference for SingleNestedAttribute                      |
| SDN Zones (5) | 61/66 | A     | Good reference for generic base pattern                       |

Full audit report and per-resource scorecards are maintained internally.
