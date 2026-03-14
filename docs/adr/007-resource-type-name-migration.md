# ADR-007: Resource Type Name Migration

## Status

Accepted

## Date

2026-03-13

## Context

All resources in this provider use the prefix `proxmox_virtual_environment_` (29 characters), making resource type names verbose and hard to scan. For example, `proxmox_virtual_environment_hagroup`, `proxmox_virtual_environment_network_linux_bridge`. Users and IDE tooling must skip past this boilerplate to reach the meaningful portion of the name. The community has requested shorter prefixes (see [#2133](https://github.com/bpg/terraform-provider-proxmox/issues/2133), [#1231](https://github.com/bpg/terraform-provider-proxmox/issues/1231)).

### Constraints Investigated

Three candidate prefixes were evaluated: `pve_`, `proxmox_`, and `proxmox_virtual_environment_` (status quo).

**Terraform's implied provider resolution** is the binding constraint. Terraform extracts the first word before `_` from a resource type name to determine which provider manages it:

- `proxmox_virtual_environment_hagroup` -> implied provider = `proxmox` (matches `bpg/proxmox`)
- `proxmox_hagroup` -> implied provider = `proxmox` (matches `bpg/proxmox`)
- `pve_hagroup` -> implied provider = `pve` (no match — would require a `pve` local name)

For `pve_*` resources to auto-resolve, users would need to declare `pve = { source = "bpg/proxmox" }` in `required_providers`. However, Terraform warns against (and has bugs with) declaring the same source under two local names since v1.2 ([terraform#31196](https://github.com/hashicorp/terraform/issues/31196), [terraform#31218](https://github.com/hashicorp/terraform/pull/31218)). This makes `pve_` impractical without a separate provider binary or registry entry.

**Plugin Framework and Mux** impose no prefix constraints — resource TypeNames can be any non-empty unique string. The prefix convention is advisory only. Both `proxmox_virtual_environment_*` and `proxmox_*` resources can coexist in the same provider binary.

### Summary of Options

| Prefix                                      | Auto-resolves? | Single provider config?     | Breaking change? |
|---------------------------------------------|----------------|-----------------------------|------------------|
| `proxmox_virtual_environment_` (status quo) | Yes            | Yes                         | No               |
| `proxmox_`                                  | Yes            | Yes                         | Yes (v1.0)       |
| `pve_`                                      | No             | No (needs `pve` local name) | Yes + UX cost    |

## Decision

Adopt `proxmox_` as the resource type name prefix, migrated in three phases.

### Phase 1: New Resources (immediate, pre-v1.0)

All new Framework resources use the `proxmox_` prefix by hardcoding the TypeName:

```go
func (r *haRuleResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = "proxmox_harule"
}
```

The provider's `MetadataResponse.TypeName` remains `"proxmox_virtual_environment"` — it is not changed. Existing resources continue using `req.ProviderTypeName + "_suffix"` unchanged.

**No breaking changes in this phase.**

### Phase 2: Rename Framework Resources (pre-v1.0, transition)

For each existing **Framework** resource, register a second resource struct with the short name alongside the original. The original emits a deprecation warning.

This phase applies only to Framework resources (`fwprovider/`). SDK resources are handled separately in Phase 3.

```go
// Short-name alias wrapping the original implementation
type haGroupShort struct{ haGroupResource }

func (r *haGroupShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = "proxmox_hagroup"
}

func (r *haGroupShort) MoveState(ctx context.Context) []resource.StateMover {
    return []resource.StateMover{{
        StateMover: func(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
            if req.SourceTypeName != "proxmox_virtual_environment_hagroup" {
                return
            }
            // Schema is identical — copy state as-is
            resp.TargetState.Raw = req.SourceRawState
        },
    }}
}
```

Users migrate at their own pace using Terraform's `moved` block (requires Terraform >= 1.8):

```hcl
moved {
  from = proxmox_virtual_environment_hagroup.example
  to   = proxmox_hagroup.example
}
```

**No breaking changes in this phase** — both names work simultaneously.

### Phase 3: Migrate SDK Resources to Framework (v1.0)

SDK resources (`proxmoxtf/`) are feature-frozen. As part of v1.0, each SDK resource is rewritten in the Framework provider with the new `proxmox_*` name directly. The old `proxmox_virtual_environment_*` SDK registrations are removed.

Additionally, any remaining long-name Framework resources from Phase 2 that still have dual registrations have their old names removed.

After v1.0, only `proxmox_*` names remain across the entire provider.

**This is a breaking change**, requiring the v1.0 major version bump.

### Scope

This ADR covers **resource type name prefix migration only**. It does not cover SDK-to-Framework schema migration, which is a separate concern with its own state transformation requirements. SDK resources being rewritten in the Framework may have different schemas; those migrations require dedicated `MoveState` implementations that handle attribute mapping, not just prefix renaming.

### Naming Convention Update

This ADR supersedes the resource type name convention in [ADR-003](003-resource-file-organization.md):

| Element             | Old Convention                                 | New Convention            |
|---------------------|------------------------------------------------|---------------------------|
| Resource type names | `proxmox_virtual_environment_{domain}_{name}`  | `proxmox_{domain}_{name}` |

New resources **must** use the `proxmox_` prefix. Existing resources retain their current names until Phase 2.

## Consequences

### Positive

- Resource names are 20+ characters shorter (e.g., `proxmox_hagroup` vs `proxmox_virtual_environment_hagroup`)
- Better IDE autocomplete and Terraform registry browsing
- Both old and new names auto-resolve to the `proxmox` provider — no user configuration changes needed
- `moved` block support provides a smooth migration path
- No separate provider binary or registry entry required

### Negative

- Dual registration in Phase 2 increases the number of registered resources temporarily
- Phase 3 (v1.0) is a breaking change requiring user migration
- `MoveState` / cross-type `moved` blocks require Terraform >= 1.8; older versions need manual `terraform state mv`
- Provider TypeName stays as `proxmox_virtual_environment` which is inconsistent with the new resource prefix — acceptable since it's an internal detail not visible to users

### Migration Risk

For Framework resource renames (Phase 2), the state migration is low-risk: schemas are identical between old and new names, so `MoveState` simply copies the raw state. No attribute transformation is needed.

SDK-to-Framework migrations (Phase 3) carry higher risk due to potential schema differences and are scoped separately from this ADR.

## References

- [#2133: Reduce "proxmox_virtual_environment_" prefix](https://github.com/bpg/terraform-provider-proxmox/issues/2133)
- [#1231: Related discussion on prefix naming](https://github.com/bpg/terraform-provider-proxmox/issues/1231)
- [ADR-003: Resource File Organization](003-resource-file-organization.md) — naming conventions (superseded for type names)
- [Terraform Plugin Framework: State Move](https://developer.hashicorp.com/terraform/plugin/framework/resources/state-move)
- [Terraform: Moved Block](https://developer.hashicorp.com/terraform/language/modules/develop/refactoring)
- [Terraform: Provider Requirements](https://developer.hashicorp.com/terraform/language/providers/requirements) — implied provider resolution
- [terraform#31196](https://github.com/hashicorp/terraform/issues/31196) — duplicate source in required_providers bug
