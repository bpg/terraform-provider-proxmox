# Architecture Decision Records

This directory contains Architecture Decision Records (ADRs) for terraform-provider-proxmox. ADRs document the key architectural patterns contributors must follow.

## ADRs

| ADR                                        | Title                        | Summary                                                                              |
|--------------------------------------------|------------------------------|--------------------------------------------------------------------------------------|
| [001](001-use-plugin-framework.md)         | Use Plugin Framework         | All new resources use Terraform Plugin Framework, not SDKv2                          |
| [002](002-api-client-structure.md)         | API Client Structure         | Layered domain clients in `proxmox/` with `ExpandPath()` pattern                     |
| [003](003-resource-file-organization.md)   | Resource File Organization   | Domain hierarchy, 3-file pattern, naming conventions                                 |
| [004](004-schema-design-conventions.md)    | Schema Design Conventions    | Attribute types, validators, model-API conversion, `CheckDelete`                     |
| [005](005-error-handling.md)               | Error Handling               | `"Unable to [Action] [Resource]"` format, 3-layer error architecture, retry policies |
| [006](006-testing-requirements.md)         | Testing Requirements         | Acceptance tests required, table-driven structure, test helpers                      |
| [007](007-resource-type-name-migration.md) | Resource Type Name Migration | Migrate from `proxmox_virtual_environment_` to `proxmox_` prefix in 3 phases         |

## Reading Order for New Resource Implementation

When implementing a new resource, read the ADRs in this order:

1. **[ADR-001](001-use-plugin-framework.md)** — Confirms you're using the right framework
2. **[ADR-003](003-resource-file-organization.md)** — Where to put files and what to name them
3. **[ADR-002](002-api-client-structure.md)** — How to access the Proxmox API
4. **[ADR-004](004-schema-design-conventions.md)** — Schema attributes, model conversion, field deletion
5. **[ADR-005](005-error-handling.md)** — Error message format and sentinel error handling
6. **[ADR-006](006-testing-requirements.md)** — Acceptance test structure and helpers
7. **[reference-examples.md](reference-examples.md)** — Copy-from code walkthroughs and checklist

## Reference Examples

The [reference-examples.md](reference-examples.md) document provides annotated walkthroughs of three real resources at increasing complexity:

1. **SDN VNet** — start here for any new resource (simplest clean 3-file pattern)
2. **Replication** — many optional fields, split create/update methods, perfect audit score
3. **Backup Job** — ConfigValidators, comma-separated-to-list, nested objects

It also includes a [checklist](reference-examples.md#checklist-for-new-resource-implementation) for new resource implementation.
