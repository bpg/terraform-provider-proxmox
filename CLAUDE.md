# CLAUDE.md

Instructions for Claude Code working on this Terraform Provider for Proxmox VE.

---

## Critical Rules

**Never violate these — they cause bugs, test failures, or provider misbehavior.**

| Never Do                              | Reason                                   |
| ------------------------------------- | ---------------------------------------- |
| Start work without a GitHub issue     | All work must be tracked                 |
| Make assumptions without verification | Always verify with code/tests/mitmproxy  |
| Skip acceptance tests                 | Tests reproduce and verify fixes         |
| Commit without running linter         | Always `make lint` first                 |
| Commit without explicit user request  | User controls git operations             |
| Add changes beyond what's requested   | Only implement what's asked              |
| Post comments to GitHub issues/PRs    | Provide text for user to post themselves |
| Add Co-Authored-By lines to commits   | Use only `-s` flag for DCO sign-off      |

| Always Do                                    | Reason                                              |
| -------------------------------------------- | --------------------------------------------------- |
| Verify GitHub issue exists first             | No issue = flag deficiency, offer to help           |
| Talk to maintainer before code investigation | Leverage domain knowledge, avoid wasted exploration |
| Ask questions when uncertain                 | Never assume; clarify before proceeding             |
| Create acceptance test BEFORE fixing         | Proves issue exists, proves fix works               |
| Verify API calls with mitmproxy              | Tests passing ≠ correct API calls                   |
| Maintain session state for multi-step work   | Enables context recovery across sessions            |
| Run full checklist before completion         | See Production Readiness Checklist                  |

---

## GitHub Issue Requirements

**All work on fixes or features MUST have a corresponding GitHub issue.**

### Before Starting Work

1. **Verify issue exists** — Search for an existing issue
2. **If no issue exists** — Flag deficiency, do NOT proceed

### When No Issue Exists

Flag this to the user:

> "No GitHub issue found for this work. All fixes and features must be tracked with an issue before implementation begins."

Then offer to help create one:

1. Ask: "Would you like me to help draft a GitHub issue?"
2. Determine type: **Bug** or **Feature/Enhancement**
3. Draft content following the template structure
4. Provide draft for user to submit at: `https://github.com/bpg/terraform-provider-proxmox/issues/new/choose`
5. Wait for issue number before proceeding

### Naming Conventions

| Artifact      | Format                          | Example                         |
| ------------- | ------------------------------- | ------------------------------- |
| Branch        | `{type}/{issue}-{desc}`         | `fix/1234-clone-timeout`        |
| Plans         | `.dev/{issue}_PLAN.md`          | `.dev/1234_PLAN.md`             |
| PR body       | `.dev/{issue}_PR_BODY.md`       | `.dev/1234_PR_BODY.md`          |
| Session state | `.dev/{issue}_SESSION_STATE.md` | `.dev/1234_SESSION_STATE.md`    |
| Test names    | Descriptive, NO issue numbers   | `TestAccResourceVMClone`        |
| VM names      | Descriptive, NO issue numbers   | `test-vm-clone`                 |
| Commits       | Conventional, NO issue numbers  | `fix(vm): handle clone timeout` |

---

## Quick Reference

### Essential Commands

```bash
make build              # Build provider binary
make lint               # Run Go linter (auto-fixes formatting and most issues)
make test               # Run unit tests
make docs               # Regenerate Framework resource/datasource docs (not SDK)
./testacc TestName      # Run specific acceptance test
npx --yes prettier --write "path/to/*.md"         # Format markdown (tables, whitespace)
npx --yes markdownlint-cli2 --fix "path/to/*.md"  # Lint markdown (rules prettier doesn't cover)
```

### Linting Rules

**Never manually format or lint code. Always use the appropriate linter tool.**

| File type | Linter command                                                                        | When to run                  |
| --------- | ------------------------------------------------------------------------------------- | ---------------------------- |
| Go `.go`  | `make lint`                                                                           | After editing any `.go` file |
| Markdown  | `npx --yes prettier --write "file.md" && npx --yes markdownlint-cli2 --fix "file.md"` | After editing any `.md` file |

Run prettier first (formats tables and whitespace), then markdownlint (validates rules prettier
doesn't cover). Config lives in `.prettierrc` / `.prettierignore` and `.markdownlint.json` /
`.markdownlintignore`.

### Acceptance Test Script (`./testacc`)

```bash
./testacc TestAccResourceVM           # Run single test
./testacc "TestAccResource.*"         # Run tests matching pattern
./testacc --tier light                # Light tests only (~30s)
./testacc --tier medium               # Medium tests only (~3 min)
./testacc --tier heavy                # Heavy tests only (~15 min)
./testacc --tier light,medium         # Combine tiers
./testacc --tier all                  # All tiers with smart parallelism (~15 min)
./testacc --resource vm               # All VM-related tests
./testacc --resource sdn              # All SDN tests
./testacc --no-proxy TestName         # Run without mitmproxy
./testacc TestName -- -count 2        # Pass flags through to go test
```

### Test Tiers

Tests are classified via `//testacc:tier=X` annotations in test files:

| Tier   | Description                    | Parallelism | Time    |
| ------ | ------------------------------ | ----------- | ------- |
| light  | API-only, no VMs or containers | -p 8        | ~30s    |
| medium | Simple VMs with unique IDs     | -p 4        | ~3 min  |
| heavy  | Cloud images, shared state     | -p 1        | ~15 min |

Resource targeting via `//testacc:resource=X` annotations: `vm`, `container`, `firewall`, `sdn`, `file`, `pool`, `acme`, `access`, `backup`, `ha`, `hardwaremapping`, `metrics`, `options`, `replication`, `apt`, `datastores`, `storage`, `network`, `misc`

Requires `testacc.env` — see [CONTRIBUTING.md](CONTRIBUTING.md#acceptance-tests) for setup.

### Production Readiness Checklist

**Run `/bpg:ready` to execute automatically.**

1. `make build` — Must pass
2. `make lint` — Must show 0 issues
3. `make test` — All unit tests pass
4. `./testacc TestAccYourFeature` — Acceptance tests pass
5. `/bpg:debug-api` — Verify API calls with mitmproxy
6. `make docs` — Regenerate Framework docs if schema changed
7. `/bpg:prepare-pr` — Generate PR body from template

### Commit Guidelines

See [CONTRIBUTING.md](CONTRIBUTING.md#commit-message-conventions). Key rules:

- Format: `{type}({scope}): {description}`
- **Types:** `feat`, `fix`, `chore`
- **Scopes:** `vm`, `lxc`, `provider`, `core`, `docs`, `ci`
- Lowercase, no period, under 72 chars, NO issue numbers
- **DCO sign-off required:** use `git commit -s` (adds `Signed-off-by` line)
- **No Claude attribution:** never add `Co-Authored-By` trailers

---

## Agent Development Practices

See [.dev/README.md](.dev/README.md#agent-development-practices) for detailed guidance on state persistence, debugging, context management, error recovery, and session handoff.

Key principles:

- **Proof over trust** — "Tests pass" ≠ correct behavior. Verify with mitmproxy or behavioral assertions. Scrutinize implementation against plan before declaring done.
- **Atomic commits** — Each commit = working, resumable state
- **Parallel agents** — Use for independent research tasks; avoid for file modifications or sequential workflows
- **Session state** — Update before any context switch; write "next action" for a stranger

---

## Project Architecture

### Prerequisites

- **Go 1.26+** required
- **golangci-lint 2.11.4** — installed automatically by `make lint`
- **Line length limit:** 150 characters (enforced by linter)
- **Comment line wrap:** ~120 characters (not 70–80; the linter allows 150, so narrow wrapping wastes vertical space)

### Overview

- **Dual-provider:** SDK v2 (`proxmoxtf/`) and Plugin Framework (`fwprovider/`)
- **New features:** Framework only; SDK is feature-frozen

### Directory Structure

```text
├── proxmox/           # Shared API client
│   └── retry/         # Unified retry logic (TaskOperation, APICallOperation, PollOperation)
├── fwprovider/        # Framework provider ← NEW CODE HERE
│   ├── test/          # Shared test utilities and acceptance tests
│   ├── config/        # Provider configuration types (Resource, DataSource)
│   ├── attribute/     # Attribute helpers (ResourceID, CheckDelete, IsDefined)
│   ├── types/         # Custom attribute types (stringset, etc.)
│   └── validators/    # Custom validators
├── proxmoxtf/         # Legacy SDK provider (feature-frozen)
├── utils/             # Shared utilities (maps, sets, strings, IP)
├── .dev/              # Development tools, plans, and session files
├── example/           # Example Terraform configurations
├── templates/         # Doc templates for Framework resources/datasources
└── docs/              # Provider documentation (mixed: see Documentation section)
```

### API Client

```text
proxmox.Client
├── Node(name) → nodes.Client
├── Cluster() → cluster.Client
├── Access() → access.Client
├── Pool() → pools.Client
├── Storage() → storage.Client
├── Version() → version.Client
├── API() → api.Client (raw HTTP)
└── SSH() → ssh.Client
```

---

## Development Workflow

### Fixing Issues

1. **Verify GitHub issue exists** — Flag deficiency if not
2. **Create branch:** `fix/{issue}-description`
3. **Create session state:** `.dev/{issue}_SESSION_STATE.md`
4. **Talk to maintainer** — Before diving into code, ask about: known regression status, scope (isolated or systemic), which code path (Framework/SDK/both), initial hunches, fix scope (narrow or broad). Assess the user's workflow/config yourself.
5. **Create acceptance test** that reproduces the issue
6. **Verify test fails** with current code
7. **Implement fix**
8. **Verify test passes**
9. **Run linter:** `make lint`
10. **Verify with mitmproxy**
11. **Complete checklist**

### Adding Features

1. **Verify GitHub issue exists** — Flag deficiency if not
2. **Create branch:** `feat/{issue}-description`
3. **Create session state:** `.dev/{issue}_SESSION_STATE.md`
4. **Talk to maintainer** — Discuss scope, design choices, and any constraints before implementation.
5. Implement in Framework provider only (`fwprovider/`)
6. Add validation, acceptance tests, documentation
7. **Complete checklist**

---

## Code Patterns

### Framework (fwprovider/)

Each resource has 3 files: `resource_*.go` (CRUD), `*_model.go` (API mapping), `resource_*_test.go` (acceptance tests). Client access flows through `config.Resource` → `cfg.Client.Domain().SubClient()`.

```go
schema.StringAttribute{
    Required: true,
    Validators: []validator.String{
        stringvalidator.OneOf("a", "b"),
    },
}
resp.Diagnostics.AddError("Unable to Create Resource", err.Error())
```

**Error diagnostic conventions:** New code should use `"Unable to [Action] [Resource]"` format (see [ADR-005](docs/adr/005-error-handling.md)). Include the resource name/ID in the summary (e.g., `fmt.Sprintf("Unable to Read VM %q", name)`) — domain clients do not reliably include it in `err.Error()`. No trailing period. Pass `err.Error()` as the detail string — never double-wrap. Legacy prefixes ("Could not", "Error") are acceptable in existing code.

### Attribute Helper Functions (fwprovider/attribute/)

When exporting model fields to API request bodies, **always use the `attribute` package helpers** instead of manual `IsDefined` + `ValueXxxPointer()` patterns:

```go
// GOOD — use helpers for model → API body assignments
body.Comment = attribute.StringPtrFromValue(m.Comment)
body.MTU = attribute.Int64PtrFromValue(m.MTU)
body.Disable = attribute.CustomBoolPtrFromValue(m.Disable)
body.Rate = attribute.Float64PtrFromValue(m.Rate)

// BAD — don't use manual IsDefined + pointer extraction
if attribute.IsDefined(m.MTU) {
    body.MTU = m.MTU.ValueInt64Pointer()
}
```

Available helpers (all return nil for null/unknown, pointer to value otherwise):

| Helper                             | Input type      | Output type                |
| ---------------------------------- | --------------- | -------------------------- |
| `attribute.StringPtrFromValue`     | `types.String`  | `*string`                  |
| `attribute.Int64PtrFromValue`      | `types.Int64`   | `*int64`                   |
| `attribute.Float64PtrFromValue`    | `types.Float64` | `*float64`                 |
| `attribute.CustomBoolPtrFromValue` | `types.Bool`    | `*proxmoxtypes.CustomBool` |

Use `attribute.IsDefined()` only when you need to branch on whether a field has a value (e.g., conditional logic), not for simple pointer extraction.

**Custom types** (`customtypes.IPCIDRValue`, `customtypes.IPAddrValue`, etc.) cannot use these helpers since they accept only `types.String`/`types.Int64`/`types.Bool`. For custom types, continue using `.ValueStringPointer()` directly — these are Optional-only fields so null/unknown handling is safe.

### Datasource Schema Attributes

In a **datasource**, attributes that are purely output (populated by the provider during Read) must be `Computed: true` only — never `Optional`. This applies to all attributes except lookup keys (which are `Required`).

| Attribute role   | Schema flags     | Example                               |
| ---------------- | ---------------- | ------------------------------------- |
| Lookup key       | `Required: true` | `id`, `node_name`                     |
| Read-only output | `Computed: true` | `name`, `status`, `tags`, `cpu` block |

**Why not `Optional` on outputs?** `Optional` on a datasource output lets users write values in config that are silently ignored — misleading UX and confusing docs (attributes appear under "Optional" instead of "Read-Only").

**Nil API values in Computed fields:** After Read, Computed attributes must have a known value — null means "unknown" which is only valid during planning. Convert nil API pointers to sensible defaults: `""` for strings, `false` for bools, empty collections for sets/maps. Use `types.StringValue("")` instead of `types.StringPointerValue(nil)`.

**Nested blocks in datasources** (e.g., `cpu`, `vga`, `rng`): The datasource should have its own `DataSourceSchema()` with `Computed: true` on the block and all inner attributes. Do not reuse `ResourceSchema()` which has `Optional: true, Computed: true` for resource write semantics.

### Comma-Separated API Values

When the Proxmox API uses comma-separated strings (e.g., `vmid=100,101,102`), **always expose them as Terraform list or set attributes** — never as raw comma-separated strings. Convert in `toAPI()` (join) and `fromAPI()` (split). See [ADR-004](docs/adr/004-schema-design-conventions.md#comma-separated-api-values--terraform-lists) for details and code examples.

### Retry Patterns (proxmox/retry/)

Three operation types — choose based on the API call pattern:

```go
// Async UPID tasks (create, clone, delete, start):
op := retry.NewTaskOperation("name", retry.WithRetryIf(retry.IsTransientAPIError))
op.DoTask(ctx, dispatchFn, waitFn)

// Synchronous blocking calls (PUT /config):
op := retry.NewAPICallOperation("name", retry.WithRetryIf(retry.ErrorContains("got timeout")))
op.Do(ctx, fn)

// Polling loops (wait for status, config unlock):
op := retry.NewPollOperation("name", retry.WithRetryIf(func(err error) bool { ... }))
op.DoPoll(ctx, fn)
```

**Delete predicate trap:** `ErrResourceDoesNotExist` can arrive via HTTP 500, so `IsTransientAPIError` alone will match it. Delete operations must combine predicates:

```go
retry.WithRetryIf(func(err error) bool {
    return retry.IsTransientAPIError(err) && !errors.Is(err, api.ErrResourceDoesNotExist)
})
```

See [ADR-005: Error Handling](docs/adr/005-error-handling.md#retry-policies) for full details.

### SDK (proxmoxtf/) — Legacy Only

```go
"key": {
    Type:     schema.TypeString,
    Required: true,
    ValidateDiagFunc: validation.ToDiagFunc(
        validation.StringInSlice([]string{"a", "b"}, false)),
}
```

When fixing validation issues, update BOTH providers where applicable.

---

## Testing Notes

- **VMs with `started = true`** need boot disk with cloud image; use `stop_on_destroy = true`
- **Naming:** Descriptive names only, NO issue numbers
- **API verification:** Use `/bpg:debug-api` for mitmproxy workflow
- **TDD with behavioral assertions:** Tests MUST actually fail without the fix — if a test passes both with and without the fix, it doesn't prove anything. Don't rely only on Terraform state attributes; use direct API checks (e.g., `te.NodeClient().VM(vmID).GetVMStatus(ctx)` to check uptime before/after to detect reboots). See `resource_vm_hotplug_test.go` and `resource_vm_disks_test.go` for patterns.
- **Connection issues:** If acceptance tests fail due to Proxmox host unreachable or similar, ask the user — don't work around it with unit tests or other substitutes
- **Functional coverage:** Tests must cover ALL major use cases for the resource — not just one happy path. Different input modes (e.g., `all` vs `vmid` vs `pool`), list attributes with multiple elements, compound fields, nested objects, and import round-trips must each have test scenarios. PRs with insufficient functional coverage will be rejected. See [ADR-006](docs/adr/006-testing-requirements.md#functional-coverage-requirement).

---

## Documentation

Docs under `docs/` are a **mix** of auto-generated and manually maintained files.

| Provider                  | Docs generation                                                             | Edit where                                                                                                                                                                               |
| ------------------------- | --------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Framework (`fwprovider/`) | Auto-generated by `make docs` from schema + optional `templates/` overrides | Edit `templates/resources/<name>.md.tmpl` (or `templates/data-sources/<name>.md.tmpl`). If no custom template exists, docs come from the schema `MarkdownDescription` fields in Go code. |
| SDK (`proxmoxtf/`)        | **Manually maintained**                                                     | Edit `docs/` files directly                                                                                                                                                              |

**Key rules:**

- `make docs` only regenerates Framework resource/datasource docs and guides with templates; SDK docs are untouched
- Manual edits to `docs/` files for Framework resources **will be lost** on `make docs` — always edit the template or schema description instead
- Manual edits to `docs/` files for SDK resources are safe — they are the source of truth
- Custom templates in `templates/` override default `tfplugindocs` generation for specific Framework resources

**Guides** use two patterns:

| Pattern                 | Source of truth                                                                                           | Examples                                                   |
| ----------------------- | --------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------- |
| **A (template-driven)** | `templates/guides/<name>.md.tmpl` with `{{ codefile }}` directives; examples in `examples/guides/<name>/` | `clone-vm`, `vm-lifecycle`                                 |
| **B (direct markdown)** | `docs/guides/<name>.md` edited directly; inline HCL blocks                                                | `multi-node`, `upgrade`, `migration-vm-clone`, `cloned-vm` |

For Pattern A guides, edit the template — `docs/guides/<name>.md` is auto-generated by `make docs` and will be overwritten.

---

## Session Management

For multi-step work, maintain session state in `.dev/{issue}_SESSION_STATE.md` using [the template](.dev/SESSION_STATE_TEMPLATE.md). Update before ending session, before context-heavy operations, after completing a phase, or when blocked.

---

## Communication Style

| Do                                                | Don't                            |
| ------------------------------------------------- | -------------------------------- |
| Be concise and direct                             | Apologize                        |
| Use technical terminology                         | Summarize changes made           |
| Explain reasoning                                 | Make up information              |
| Admit uncertainty                                 | Show implementation unless asked |
| Use friendly, conversational tone in ticket notes | Use formal/corporate language    |
| Lead with key finding, skip preamble              | Write verbose preambles          |

**Code comments:** Minimal — explain "why", not "what". One concise comment per block, skip when code is self-explanatory.

- Good: `// Reboot before resize: pending changes include old disk size`
- Bad: `// Update the size in the update body to match the plan size when the disk is growing. Without this, the UpdateVM API call creates a pending change with the OLD size, and a subsequent reboot would revert the resize done by ResizeVMDisk.`

---

## Skills

| Skill              | Purpose                                                            |
| ------------------ | ------------------------------------------------------------------ |
| `/bpg:start-issue` | Start work on a GitHub issue (branch + session state)              |
| `/bpg:resume`      | Resume work from a previous session                                |
| `/bpg:ready`       | Run production readiness checklist                                 |
| `/bpg:debug-api`   | Debug API calls with mitmproxy                                     |
| `/bpg:prepare-pr`  | Prepare PR body from template with proof of work                   |
| `/bpg:done`        | Wrap up session — extract learnings, finalize state, archive files |

See [.dev/README.md](.dev/README.md#working-with-llm-agents) for detailed workflow documentation and how skills connect together.

---

## References

- [CONTRIBUTING.md](CONTRIBUTING.md) — Contributing guide
- [docs/adr/](docs/adr/README.md) — Architecture Decision Records and reference examples
- [.dev/DEBUGGING.md](.dev/DEBUGGING.md) — Debugging guide
- [.dev/SESSION_STATE_TEMPLATE.md](.dev/SESSION_STATE_TEMPLATE.md) — Session template
- [Proxmox API](https://pve.proxmox.com/pve-docs/api-viewer/)
- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
