## Context: Option C (“cloned VM” resource) for terraform-provider-proxmox

This repository contains two Terraform provider implementations:

- **Legacy SDK provider**: `proxmoxtf/`
- **Framework provider (FWK)**: `fwprovider/`

The goal is to implement **Option C** from the cloning semantics discussion: introduce a **dedicated “cloned VM” resource** (FWK-first) so that cloning is treated as an **imperative create-time function**, and subsequent management is explicit and unambiguous.

This avoids mixing Proxmox clone inheritance (unknown state) with Terraform’s declarative “desired state” semantics inside a single VM resource.

## What was already done in this repo (important)

### 1) FWK `vm2` clone support removed

We removed clone-related complexity from the FWK `vm2` resource and its datasource.

- `proxmox_virtual_environment_vm2` no longer supports `clone`.
- Clone-specific update behavior was removed (no “keep inherited values because clone” branching).
- Clone-focused acceptance tests for `vm2` were removed.
- Non-`vm2` clone tests were intentionally kept as-is.

### 2) `./testacc` improved

The `./testacc` script was upgraded so acceptance tests are easier to run:

- Sources env file (`testacc.env`) instead of using `xargs` (avoids ARG_MAX issues).
- Supports:
  - `--pkg <pkg>` to override the package
  - `--all` to run all FWK acceptance tests
  - `--env <file>` to source a different env file
  - `--no-proxy` to disable proxy env for the run
  - `--` passthrough for `go test` flags, e.g. `./testacc TestName -- -timeout 10m -count 1`
- When proxy env vars are set, `./testacc` auto-adds `NO_PROXY/no_proxy` entries for `127.0.0.1,localhost,::1` to avoid Terraform/provider local reattach failures.

### 3) Acceptance tests execution caveat (Cursor)

When running acceptance tests inside Cursor, Terraform/provider can fail with:

- `Error: Request cancelled` (plugin ApplyResourceChange cancelled)

Cause: **sandbox write restrictions** preventing Terraform from writing to system temp dirs.

Resolution:

- Run acceptance tests **without sandbox restrictions** when needed.

## Desired behavior of the new “cloned VM” resource (Option C)

### High-level contract

- **Create**:
  1) Call Proxmox clone API to create the VM from a template/source VM.
  2) Apply only explicitly-managed configuration.
  3) Apply explicit deletes if requested.
  4) Read back the managed subset and store state.

- **Read**:
  - Track identity (`node_name`, `vm_id`) and **only managed fields**.
  - Do not attempt to mirror the full inherited remote VM state into Terraform state.

- **Update**:
  - Update only managed fields.
  - Removing configuration from Terraform should usually mean **stop managing**, not delete remote config.

- **Delete**:
  - Delete the VM (standard Terraform lifecycle).

## Schema strategy (critical)

The current clone pain is largely caused by list-based device modeling and inherited “unknown” state.

For the new resource:

- Use **maps keyed by slot identifiers**, not lists.
  - network: `net0`, `net1`, ...
  - disks: `scsi0`, `virtio0`, `ide2`, ...

This makes diffs addressable and removes ambiguity.

### Explicit delete mechanism

To delete inherited devices, require explicit intent.

Example concept:

- `delete = { network = ["net1"], disk = ["scsi2"] }`

Rule:

- **Omission ≠ deletion**.
- Only explicit entries in `delete` cause deletion.

## Recommended implementation approach (FWK-first)

Implement the new resource in **FWK only** first:

- No breaking changes for existing SDK users.
- Avoids forcing clone semantics into `vm2`.
- Provides a clean path forward; SDK parity can be considered later.

## Step-by-step implementation plan

### Step 0 — Confirm naming

Pick the final resource type name.

Suggested:

- `proxmox_virtual_environment_cloned_vm` (FWK)

### Step 1 — Add resource skeleton (FWK)

Create a new resource under `fwprovider/` and register it in the FWK provider.

Deliverables:

- resource file(s): Create/Read/Update/Delete implemented
- model + schema
- registration in FWK provider resources

### Step 2 — Design schema and model

Include:

- Identity:
  - `node_name` (required)
  - `vm_id` (optional or computed+optional, similar pattern to `vm2`)

- `clone` block (ForceNew):
  - `id` / `source_vmid` (required)
  - optional: full/linked, target storage, retries, etc.

- Managed configuration blocks:
  - Reuse existing FWK vm blocks where possible (`cpu`, `rng`, `vga`, `cdrom`, later disks/network/etc.).

- Device maps:
  - `network` map
  - `disk` map

- `delete` block:
  - `network` list of slot keys
  - `disk` list of slot keys

### Step 3 — Implement Create

Pseudo-flow:

1) Determine `vm_id` (generate if empty)
2) Call clone API
3) Apply managed config (blocks + device maps)
4) Apply deletes (only explicit)
5) Read back and set state (managed subset)

### Step 4 — Implement Read

Read only:

- Identity
- Managed blocks
- Managed map keys

Do not import every inherited remote device into state.

### Step 5 — Implement Update

Map diff semantics:

- Added key: create/apply device
- Changed key: update device
- Removed key: stop managing (no remote delete)
- Deleted key: only if in `delete` list

Non-map blocks:

- Update only explicit changes; avoid trying to “clean up” inherited config on omission.

### Step 6 — Implement Delete

- Stop/shutdown behavior (as per existing patterns)
- Delete VM

### Step 7 — Add acceptance tests (must cover Option C semantics)

Add tests (in `fwprovider/test/`, or another acceptance location consistent with existing patterns) that prove:

1) **Inheritance preserved**:
   - Template has `net0` and `net1`
   - cloned resource manages only `net0`
   - Expect `net1` still exists

2) **Explicit delete**:
   - Same template
   - `delete.network=["net1"]`
   - Expect `net1` removed

3) **Stop managing ≠ delete**:
   - Manage `net0` then remove it from config without listing it in delete
   - Expect `net0` still exists remotely

4) **Map-key stability**:
   - Update `net0` settings without affecting other NICs

Execution:

- Use `./testacc TestName` and when running from Cursor, disable sandbox restrictions if needed.

### Step 8 — Documentation/migration note (minimal)

Update docs/templates as needed:

- `vm2` no longer supports clone
- new cloned-vm resource is recommended for clone workflows
- explain explicit delete mechanism

## Operational notes for future sessions

- Run lint: `make lint`
- Run unit tests: `go test ./...`
- Run acceptance tests: `./testacc TestName` (consider sandbox constraints in Cursor)

## Status checklist

- [x] Remove clone from FWK `vm2` resource and datasource
- [x] Remove vm2 clone acceptance tests
- [x] Improve `./testacc`
- [x] Validate vm2 acceptance tests pass
- [ ] Implement new FWK cloned-vm resource (Option C)
- [ ] Add acceptance tests that validate Option C semantics
- [ ] Wire docs/migration guidance
