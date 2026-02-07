# CLAUDE.md

Instructions for Claude Code working on this Terraform Provider for Proxmox VE.

---

## Critical Rules

**Never violate these — they cause bugs, test failures, or provider misbehavior.**

| Never Do | Reason |
| -------- | ------ |
| Start work without a GitHub issue | All work must be tracked |
| Make assumptions without verification | Always verify with code/tests/mitmproxy |
| Skip acceptance tests | Tests reproduce and verify fixes |
| Commit without running linter | Always `make lint` first |
| Commit without explicit user request | User controls git operations |
| Add changes beyond what's requested | Only implement what's asked |

| Always Do | Reason |
| --------- | ------ |
| Verify GitHub issue exists first | No issue = flag deficiency, offer to help |
| Ask questions when uncertain | Never assume; clarify before proceeding |
| Create acceptance test BEFORE fixing | Proves issue exists, proves fix works |
| Verify API calls with mitmproxy | Tests passing ≠ correct API calls |
| Maintain session state for multi-step work | Enables context recovery across sessions |
| Run full checklist before completion | See Production Readiness Checklist |

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

| Artifact | Format | Example |
| -------- | ------ | ------- |
| Branch | `{type}/{issue}-{desc}` | `fix/1234-clone-timeout` |
| Plans | `.dev/YYYY-MM-DD-{feature}.md` | `.dev/2026-02-03-reference-examples.md` |
| PR body | `.dev/{issue}_PR_BODY.md` | `.dev/1234_PR_BODY.md` |
| Session state | `.dev/{issue}_SESSION_STATE.md` | `.dev/1234_SESSION_STATE.md` |
| Test names | Descriptive, NO issue numbers | `TestAccResourceVMClone` |
| VM names | Descriptive, NO issue numbers | `test-vm-clone` |
| Commits | Conventional, NO issue numbers | `fix(vm): handle clone timeout` |

---

## Quick Reference

### Essential Commands

```bash
make build              # Build provider binary
make lint               # Run Go linter (auto-fixes formatting and most issues)
make test               # Run unit tests
make docs               # Generate documentation
./testacc TestName      # Run specific acceptance test
npx --yes markdownlint-cli2 --fix "path/to/*.md"  # Lint markdown files
```

### Linting Rules

**Never manually format or lint code. Always use the appropriate linter tool.**

| File type  | Linter command                                  | When to run                     |
|------------|------------------------------------------------ |---------------------------------|
| Go (`.go`) | `make lint`                                     | After editing any `.go` file    |
| Markdown   | `npx --yes markdownlint-cli2 --fix "file.md"`  | After editing any `.md` file    |

### Acceptance Test Script (`./testacc`)

```bash
./testacc TestAccResourceVM           # Run single test
./testacc "TestAccResource.*"         # Run tests matching pattern
./testacc --no-proxy TestName         # Run without mitmproxy
./testacc --verbose TestName          # Verbose output
./testacc TestName -- -v              # Pass flags through to go test
```

Requires `testacc.env` with:

```bash
TF_ACC=1
PROXMOX_VE_API_TOKEN="root@pam!<token>=<value>"
PROXMOX_VE_ENDPOINT="https://<host>:8006/"
PROXMOX_VE_SSH_AGENT="true"
PROXMOX_VE_SSH_USERNAME="root"
# Optional: PROXMOX_VE_ACC_NODE_NAME, PROXMOX_VE_ACC_NODE_SSH_ADDRESS, etc.
```

### Production Readiness Checklist

**Run `/ready` to execute automatically.**

1. `make build` — Must pass
2. `make lint` — Must show 0 issues
3. `make test` — All unit tests pass
4. `./testacc TestAccYourFeature` — Acceptance tests pass
5. `/debug-api` — Verify API calls with mitmproxy
6. `make docs` — Regenerate if schema changed
7. `/prepare-pr` — Generate PR body from template

### Commit Guidelines

See [CONTRIBUTING.md](CONTRIBUTING.md#commit-message-conventions). Key rules:

- Format: `{type}({scope}): {description}`
- **Types:** `feat`, `fix`, `chore`
- **Scopes:** `vm`, `lxc`, `provider`, `core`, `docs`, `ci`
- Lowercase, no period, under 72 chars, NO issue numbers
- **DCO sign-off required:** use `git commit -s` (adds `Signed-off-by` line)

---

## Agent Development Practices

### Parallel Agents

Use parallel agents for independent tasks to speed up work:

**Good candidates for parallel execution:**

- Research tasks (explore different parts of codebase simultaneously)
- Running independent test suites
- Searching for patterns across different directories
- Gathering context from multiple unrelated files

**Not suitable for parallel execution:**

- Tasks with dependencies (B needs output of A)
- File modifications (risk of conflicts)
- Sequential workflows (test → fix → verify)

**How to request:** Ask for agents to run "in parallel" explicitly.

### State Persistence

LLMs have no memory between sessions. Externalize state to files:

- **Session state file** — The agent's memory across context resets
- **Update before ANY context switch** — End of session, new task, long operation
- **Write "next action" for a stranger** — Assume no prior context

### Track Decisions, Not Just Actions

- **User decisions** — Never re-ask; record in session state
- **Agent assumptions** — Make explicit; mark verified/rejected
- **Reasoning** — "Why" matters more than "what"

### Hypothesis-Driven Debugging

- Form hypothesis → test → record result
- Prevents circular debugging across sessions
- Use "Hypotheses Tested" table in session state

### Minimize Re-exploration

- Cache code patterns and file locations in session state
- Record dead ends so they're not re-explored
- Note key file:line references for quick restoration

### Atomic Commits

- Each commit = working, resumable state
- If session dies mid-work, resume from last commit

### Proof Over Trust

- "Tests pass" ≠ correct behavior
- Always verify with mitmproxy
- Include evidence in PR proof of work section

### Context Window Management

For long-running tasks:

- **Checkpoint frequently** — Update session state after every successful test run
- **Summarize completed work** — Don't keep raw exploration in context; distill findings
- **Chunk large changes** — Break into atomic commits to create resume points
- **Use `/resume`** — Start new sessions by loading session state, not from memory

### Error Recovery

When things go wrong:

- **Test failures** — Record in session state, add to "Hypotheses Tested", don't mark complete
- **API errors** — Capture in mitmproxy log, document in session state
- **Context loss** — Always resume from session state file using `/resume`
- **Blocked work** — Update session status to "Blocked", document blocker, move to next task

### Session Handoff

When handing off work:

- **To another agent** — Ensure "Quick Context Restore" is complete and current
- **To human** — Create PR using `/prepare-pr`, reference session state location
- **From human** — Use `/resume`, ask about any "Unverified" assumptions

---

## Project Architecture

### Prerequisites

- **Go 1.25+** required
- **golangci-lint 2.8.0** — installed automatically by `make lint`
- **Line length limit:** 150 characters (enforced by linter)

### Overview

- **Dual-provider:** SDK v2 (`proxmoxtf/`) and Plugin Framework (`fwprovider/`)
- **New features:** Framework only; SDK is feature-frozen

### Directory Structure

```text
├── proxmox/           # Shared API client
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
└── docs/              # Auto-generated documentation
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
4. **Create acceptance test** that reproduces the issue
5. **Verify test fails** with current code
6. **Implement fix**
7. **Verify test passes**
8. **Run linter:** `make lint`
9. **Verify with mitmproxy**
10. **Complete checklist**

### Adding Features

1. **Verify GitHub issue exists** — Flag deficiency if not
2. **Create branch:** `feat/{issue}-description`
3. **Create session state:** `.dev/{issue}_SESSION_STATE.md`
4. Implement in Framework provider only (`fwprovider/`)
5. Add validation, acceptance tests, documentation
6. **Complete checklist**

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
- **API verification:** Use `/debug-api` for mitmproxy workflow

---

## Session Management

For multi-step work, maintain session state using [.dev/SESSION_STATE_TEMPLATE.md](.dev/SESSION_STATE_TEMPLATE.md).

**Location:** `.dev/{issue}_SESSION_STATE.md`

**Key sections to maintain:**

- Quick Context Restore — For fast agent bootstrap
- User Decisions — Prevent re-asking
- Assumptions Made — Track verification status
- Context Gathered — Save re-reading files
- Hypotheses Tested — For debugging sessions

**Update triggers:**

- Before ending session
- Before context-heavy operations
- After completing a phase
- When blocked or switching tasks

---

## Communication Style

| Do | Don't |
| -- | ----- |
| Be concise and direct | Apologize |
| Use technical terminology | Summarize changes made |
| Explain reasoning | Make up information |
| Admit uncertainty | Show implementation unless asked |

---

## Skills

| Skill | Purpose |
| ----- | ------- |
| `/start-issue` | Start work on a GitHub issue (branch + session state) |
| `/resume` | Resume work from a previous session |
| `/ready` | Run production readiness checklist |
| `/debug-api` | Debug API calls with mitmproxy |
| `/prepare-pr` | Prepare PR body from template with proof of work |

See [.dev/README.md](.dev/README.md#working-with-llm-agents) for detailed workflow documentation and how skills connect together.

---

## References

- [CONTRIBUTING.md](CONTRIBUTING.md) — Contributing guide
- [docs/adr/](docs/adr/README.md) — Architecture Decision Records and reference examples
- [.dev/DEBUGGING.md](.dev/DEBUGGING.md) — Debugging guide
- [.dev/SESSION_STATE_TEMPLATE.md](.dev/SESSION_STATE_TEMPLATE.md) — Session template
- [Proxmox API](https://pve.proxmox.com/pve-docs/api-viewer/)
- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
