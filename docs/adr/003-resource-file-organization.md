# ADR-003: Resource File Organization

## Status

Accepted

## Date

2026-02-04 (retroactive documentation)

## Context

Contributors need clear guidance on where to place files and how to name them when implementing new resources. Without documented conventions, file placement and naming varies across the codebase, increasing review burden and making patterns harder to discover.

## Decision

### Domain Hierarchy

Resources are organized by Proxmox API domain under `fwprovider/`:

```text
fwprovider/
├── access/                  # /access/* endpoints
├── cluster/                 # /cluster/* endpoints
│   ├── acme/                # /cluster/acme/*
│   ├── ha/                  # /cluster/ha/*
│   ├── metrics/             # /cluster/metrics/*
│   ├── options/             # /cluster/config
│   └── sdn/                 # /cluster/sdn/*
│       ├── zone/            # /cluster/sdn/zones
│       ├── vnet/            # /cluster/sdn/vnets
│       └── subnet/          # /cluster/sdn/vnets/{vnet}/subnets
├── nodes/                   # /nodes/{node}/* endpoints
│   ├── apt/                 # /nodes/{node}/apt
│   ├── firewall/            # /nodes/{node}/firewall
│   ├── network/             # /nodes/{node}/network
│   └── vm/                  # /nodes/{node}/qemu/{vmid}
├── pools/                   # /pools/*
└── storage/                 # /storage/*
```

The directory hierarchy mirrors the Proxmox API path structure. This makes it straightforward to locate the code for any given API endpoint.

### File Pattern

Each resource package contains a minimum of three files:

| File     | Purpose                                        | Naming                                          |
| -------- | ---------------------------------------------- | ----------------------------------------------- |
| Resource | CRUD operations, schema, constructor           | `resource.go` or `resource_{name}.go`           |
| Model    | Terraform model struct, `toAPI()`, `fromAPI()` | `model.go` or `model_{name}.go`                 |
| Tests    | Acceptance tests                               | `resource_test.go` or `resource_{name}_test.go` |

Use the short form (`resource.go`, `model.go`) when the package contains a single resource. Use the qualified form (`resource_{name}.go`) when multiple resources share a package.

Optional additional files:

| File                                      | When to Add                                                              |
| ----------------------------------------- | ------------------------------------------------------------------------ |
| `datasource.go` or `datasource_{name}.go` | Read-only data source (see [ADR-005: Datasource Error Handling][ds-err]) |
| `resource_schema.go`                      | Schema definition is large enough to warrant separation                  |
| `datasource_schema.go`                    | Data source schema is large enough to warrant separation                 |

[ds-err]: 005-error-handling.md#datasource-error-handling

### Singleton Resources

Some Proxmox resources represent pre-existing cluster or node configuration that always exists (e.g., cluster options, node firewall options). These "singleton" resources have modified lifecycle semantics:

| Operation | Singleton Behavior                                                    |
| --------- | --------------------------------------------------------------------- |
| Create    | Applies settings via PUT (the resource already exists)                |
| Read      | Fetches current settings (not-found handling is N/A)                  |
| Update    | Applies changed settings via PUT                                      |
| Delete    | Resets all managed fields to defaults (does not destroy the resource) |
| Import    | Reads current settings (always succeeds for valid node/cluster)       |

Singleton resources follow the same 3-file pattern and naming conventions. The Delete method should enumerate all managed fields and send them in a `delete` list to reset them to Proxmox defaults.

Singleton resources use the same error message format as regular resources. Since the underlying resource always exists, not-found handling in Read and Delete is typically N/A.

**Common Mistakes:**

- Implementing a no-op Delete that leaves managed fields in their configured state.
- Checking for `ErrResourceDoesNotExist` in Read when the resource always exists.
- Not reading back after Create/Update (this requirement still applies to singletons). See [ADR-005](005-error-handling.md#read-back-after-create-and-update).

### Naming Conventions

| Element                   | Convention                                                                 | Example                                                 |
| ------------------------- | -------------------------------------------------------------------------- | ------------------------------------------------------- |
| Package names             | Lowercase, singular noun matching API domain                               | `zone`, `vnet`, `subnet`, multi-word: `hardwaremapping` |
| Terraform attribute names | `snake_case`                                                               | `isolate_ports`, `vlan_aware`                           |
| Go struct fields          | PascalCase with `tfsdk` tag                                                | `IsolatePorts types.Bool \`tfsdk:"isolate_ports"\``     |
| Resource type names (new) | `proxmox_{domain}_{name}` ([ADR-007](007-resource-type-name-migration.md)) | `proxmox_sdn_vnet`                                      |
| Resource type names (old) | `proxmox_virtual_environment_{domain}_{name}` (legacy, pre-ADR-007)        | `proxmox_virtual_environment_sdn_vnet`                  |
| Test function names       | `TestAcc{Resource}{Scenario}`                                              | `TestAccResourceSDNVNet`                                |
| Constructor functions     | `NewResource()`, `NewDataSource()`                                         | —                                                       |

### Test Colocation

Test files are colocated with their resource files in the same package directory. Shared test utilities (provider factories, config rendering, assertion helpers) live in `fwprovider/test/`.

### Registration

New resources and data sources must be registered in `fwprovider/provider.go`:

- Resources: add to the `Resources()` method
- Data Sources: add to the `DataSources()` method

### Corresponding API Client

Each resource package in `fwprovider/` typically has a corresponding client package in `proxmox/` following the same domain hierarchy. See [ADR-002](002-api-client-structure.md) for API client organization.

## Consequences

### Positive

- Contributors can locate code by API path
- Consistent structure reduces review friction
- New resources follow a predictable pattern
- Tests are easy to find alongside their resources

### Negative

- Deeply nested domains (e.g., SDN subnets) create deep directory trees
- Single-resource packages may feel like overhead for very simple resources

### Common Mistakes

- Placing new resources in `proxmoxtf/` instead of `fwprovider/`.
- Using `resource_{name}_model.go` or `{name}_model.go` instead of `model_{name}.go` for model files.
- Forgetting to register the resource in `fwprovider/provider.go`.
- Creating flat directory structures when the API path has nesting (e.g., `sdn_vnet.go` instead of `sdn/vnet/resource.go`).

## References

- [ADR-001: Use Plugin Framework](001-use-plugin-framework.md)
- [ADR-002: API Client Structure](002-api-client-structure.md)
- [ADR-004: Schema Design Conventions](004-schema-design-conventions.md) — schema and model patterns
- [ADR-006: Testing Requirements](006-testing-requirements.md) — test file placement and structure
- [ADR-007: Resource Type Name Migration](007-resource-type-name-migration.md) — resource type naming and `moved` block support
- [ADR-008: Sub-block Contract](008-sub-block-contract.md) — sub-package file layout for VM-style composite resources (uses the same 3-file pattern with `resource_schema.go` separated out and `datasource_schema.go` added)
- [Reference Examples](reference-examples.md) — annotated walkthrough of the 3-file pattern
