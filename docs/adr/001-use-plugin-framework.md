# ADR-001: Use Terraform Plugin Framework for All New Resources

## Status

Accepted

## Date

2026-02-01 (retroactive documentation)

## Context

The terraform-provider-proxmox codebase contains resources implemented using two different Terraform provider development frameworks:

1. **Terraform Plugin SDKv2** (legacy) - Resources in `proxmoxtf/resource/`
2. **Terraform Plugin Framework** (current) - Resources in `fwprovider/`

HashiCorp has declared SDKv2 to be in maintenance mode, with the Plugin Framework being the recommended approach for new provider development. The Plugin Framework offers several advantages:

- Better type safety with Go generics
- Cleaner separation of schema, model, and CRUD logic
- Native support for Protocol 6
- Improved plan modification capabilities
- Better support for complex nested attributes

The provider currently uses a multiplexer (`tf6muxserver`) to serve both SDKv2 and Framework resources from a single binary, allowing gradual migration.

## Decision

**All new resources and data sources MUST be implemented using the Terraform Plugin Framework.**

This includes:

- New Proxmox functionality (e.g., new SDN features, new API endpoints)
- Resources requested by the community
- Any resource not currently implemented

SDKv2 may still be used for:

- Bug fixes to existing SDKv2 resources (until they are migrated)
- Enhancements to existing SDKv2 resources (until they are migrated)
- Backports of critical fixes

### Migration Priority

The following SDKv2 resources are prioritized for migration to Plugin Framework:

1. VM (`proxmoxtf/resource/vm`)
2. Container (`proxmoxtf/resource/container`)

## Consequences

### Positive

- New code follows modern patterns and is easier to maintain
- Contributors learn one framework, not two
- Prepares codebase for v1.0 release
- Better alignment with Terraform ecosystem direction

### Negative

- Contributors familiar only with SDKv2 need to learn Framework
- Some patterns differ between Framework and SDK, causing potential confusion
- Migration of existing resources requires effort

### Implementation Notes

1. **File Location**: New resources go in `fwprovider/` directory, organized by domain
2. **File Pattern**: 3-file structure per resource — see [ADR-003](003-resource-file-organization.md)
3. **Tests**: Acceptance tests colocated with resource; shared test utilities in `fwprovider/test/`
4. **Registration**: Resources are registered in `fwprovider/provider.go` via `Resources()` and `DataSources()` methods
5. **Client Access**: Use `config.Resource` or `config.DataSource` from configure methods
6. **Reference Examples**: See [reference-examples.md](reference-examples.md) for annotated walkthroughs

### Common Mistakes

- Adding new resources in `proxmoxtf/` (SDKv2) instead of `fwprovider/` (Framework).
- Importing `github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema` in new code — use `github.com/hashicorp/terraform-plugin-framework/resource/schema` instead.

## References

- [Terraform Plugin Framework documentation](https://developer.hashicorp.com/terraform/plugin/framework)
- [Migration guide](https://developer.hashicorp.com/terraform/plugin/framework/migrating)
- [ADR-002: API Client Structure](002-api-client-structure.md) — domain client hierarchy
- [ADR-003: Resource File Organization](003-resource-file-organization.md) — file naming and placement
- [ADR-004: Schema Design Conventions](004-schema-design-conventions.md) — attributes, validators, model conversion
- [ADR-005: Error Handling](005-error-handling.md) — error message format and sentinel errors
- [ADR-006: Testing Requirements](006-testing-requirements.md) — acceptance test structure
- [Reference Examples](reference-examples.md) — annotated walkthroughs
