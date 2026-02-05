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
|----------|------------------------------------------------|-------------------------------------------------|
| Resource | CRUD operations, schema, constructor           | `resource.go` or `resource_{name}.go`           |
| Model    | Terraform model struct, `toAPI()`, `fromAPI()` | `model.go` or `model_{name}.go`                 |
| Tests    | Acceptance tests                               | `resource_test.go` or `resource_{name}_test.go` |

Use the short form (`resource.go`, `model.go`) when the package contains a single resource. Use the qualified form (`resource_{name}.go`) when multiple resources share a package.

Optional additional files:

| File                                      | When to Add                                              |
|-------------------------------------------|----------------------------------------------------------|
| `datasource.go` or `datasource_{name}.go` | Read-only data source for the same API object            |
| `resource_schema.go`                      | Schema definition is large enough to warrant separation  |
| `datasource_schema.go`                    | Data source schema is large enough to warrant separation |

### Naming Conventions

| Element                   | Convention                                    | Example                                                 |
|---------------------------|-----------------------------------------------|---------------------------------------------------------|
| Package names             | Lowercase, singular noun matching API domain  | `zone`, `vnet`, `subnet`, multi-word: `hardwaremapping` |
| Terraform attribute names | `snake_case`                                  | `isolate_ports`, `vlan_aware`                           |
| Go struct fields          | PascalCase with `tfsdk` tag                   | `IsolatePorts types.Bool \`tfsdk:"isolate_ports"\``     |
| Resource type names       | `proxmox_virtual_environment_{domain}_{name}` | `proxmox_virtual_environment_sdn_vnet`                  |
| Test function names       | `TestAcc{Resource}{Scenario}`                 | `TestAccResourceSDNVNet`                                |
| Constructor functions     | `NewResource()`, `NewDataSource()`            | —                                                       |

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
- [Reference Examples](reference-examples.md) — annotated walkthrough of the 3-file pattern
