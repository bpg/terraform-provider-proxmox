# Architecture Decision Records

This directory contains Architecture Decision Records (ADRs) for terraform-provider-proxmox. ADRs document the key architectural patterns contributors must follow.

## ADRs

| ADR                                      | Title                      | Summary                                                                              |
|------------------------------------------|----------------------------|--------------------------------------------------------------------------------------|
| [001](001-use-plugin-framework.md)       | Use Plugin Framework       | All new resources use Terraform Plugin Framework, not SDKv2                          |
| [002](002-api-client-structure.md)       | API Client Structure       | Layered domain clients in `proxmox/` with `ExpandPath()` pattern                     |
| [003](003-resource-file-organization.md) | Resource File Organization | Domain hierarchy, 3-file pattern, naming conventions                                 |
| [004](004-schema-design-conventions.md)  | Schema Design Conventions  | Attribute types, validators, model-API conversion, `CheckDelete`                     |
| [005](005-error-handling.md)             | Error Handling             | `"Unable to [Action] [Resource]"` format, 3-layer error architecture, retry policies |
| [006](006-testing-requirements.md)       | Testing Requirements       | Acceptance tests required, table-driven structure, test helpers                      |

## Reference Examples

The [reference-examples.md](reference-examples.md) document provides annotated walkthroughs of three real resources at increasing complexity:

1. **SDN VNet** — start here for any new resource
2. **Metrics Server** — many optional fields, sensitive attributes
3. **ACL** — cross-field validation, custom import parsing

It also includes a [checklist](reference-examples.md#checklist-for-new-resource-implementation) for new resource implementation.
