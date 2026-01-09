# AGENTS.md

Universal instructions for AI agents working on this Terraform Provider for Proxmox VE.

> **Navigation:** This document uses a tiered structure. Start with [Quick Reference](#quick-reference) for common tasks. Dive into detailed sections as needed. Use `Ctrl+F` to search for specific topics.

---

## Table of Contents

1. [Quick Reference](#quick-reference) — Essential commands and rules
2. [Critical Rules](#critical-rules) — Never violate these
3. [Project Architecture](#project-architecture) — Codebase structure
4. [Development Workflow](#development-workflow) — How to implement changes
5. [Testing](#testing) — Running and writing tests
6. [Debugging](#debugging) — API verification with mitmproxy
7. [Code Patterns](#code-patterns) — Framework vs SDK patterns
8. [Session Management](#session-management) — For multi-step tasks
9. [Agent Behavior](#agent-behavior) — Communication and knowledge management

---

## Quick Reference

### Essential Commands

```bash
make build              # Build provider binary
make lint               # Run linter (it will auto-fix most issues)
make test               # Run unit tests
make docs               # Generate documentation
./testacc TestName      # Run specific acceptance test
./testacc --all         # Run all acceptance tests (only for massive changes, as old tests are fragile and prone to breaking)
make example            # Validate with example configs (mostly for SDKv2 provider changes)
```

### Production Readiness Checklist

**Complete ALL before declaring work done:**

1. `make build` — Must pass
2. `make lint` — Must show 0 issues
3. `make test` — All unit tests pass
4. `./testacc TestAccYourFeature` — Acceptance tests pass
5. Verify API calls with mitmproxy (see [Debugging](#debugging))
6. `make docs` — Regenerate if schema changed
7. `make example` — Must complete without errors

### Commit & PR Guidelines

See [CONTRIBUTING.md](CONTRIBUTING.md#commit-message-conventions) for commit format and [Submitting changes](CONTRIBUTING.md#submitting-changes) for PR requirements.

**Key rules:** Lowercase, no period, under 72 chars, NO issue numbers. Include proof of work.

---

## Critical Rules

> **⚠️ Violations cause bugs, test failures, or provider misbehavior.**

### Never Do

| Rule                                  | Reason                                    |
|---------------------------------------|-------------------------------------------|
| Make assumptions without verification | Always verify with code/tests/mitmproxy   |
| Remove unrelated code                 | Preserve existing patterns                |
| Add changes beyond what's requested   | Only implement what's asked               |
| Leave unused variables                | Causes linter failures                    |
| Skip acceptance tests                 | Tests reproduce and verify fixes          |
| Use issue numbers in test/VM names    | Use descriptive names only                |
| Include issue numbers in commits      | Project convention                        |
| Commit without running linter         | Always `make lint` first                  |
| Commit without explicit user request  | User controls git operations              |

### Always Do

| Rule                                      | Reason                                            |
|-------------------------------------------|---------------------------------------------------|
| Ask questions when uncertain              | Never assume; clarify before proceeding           |
| Verify information before presenting      | Never speculate or invent                         |
| Create acceptance test BEFORE fixing      | Proves issue exists, proves fix works             |
| Update BOTH providers for validation      | SDK (`proxmoxtf/`) and Framework (`fwprovider/`)  |
| Run full checklist before completion      | See [Quick Reference](#quick-reference)           |
| Use mitmproxy to verify API calls         | Tests passing ≠ correct API calls                 |

---

## Project Architecture

### Overview

Terraform/OpenTofu provider for Proxmox VE 9.x using dual-provider architecture: legacy SDK v2 and modern Plugin Framework, multiplexed into a single binary.

**Key facts:**

- Go 1.25+ required
- Mozilla Public License v2.0
- Active migration from SDK to Framework (targeting v1.0)

### Directory Structure

```text
.
├── main.go                 # Entry point with mux server
├── proxmox/                # Shared API client (used by both providers)
│   ├── api/               # REST client with auth, TLS, retry
│   ├── nodes/             # VM, container, storage operations
│   ├── cluster/           # ACME, HA, SDN, metrics
│   └── access/            # ACL, users, roles, tokens
├── fwprovider/            # Framework provider (modern) ← NEW CODE HERE
│   ├── nodes/vm/          # VM2 resource (modular sub-packages)
│   ├── nodes/clonedvm/    # Dedicated cloned VM resource
│   ├── test/              # Acceptance tests
│   ├── types/             # Custom attribute types
│   └── validators/        # Custom validators
├── proxmoxtf/             # Legacy SDK provider
│   ├── resource/          # Resource implementations
│   └── datasource/        # DataSource implementations
├── .dev/                  # Development tools and guides
│   ├── DEBUGGING.md       # Comprehensive debugging guide
│   └── proxmox_debug_script.py  # Mitmproxy helper
└── docs/                  # Auto-generated documentation
```

### Dual Provider System

| Aspect      | Framework (`fwprovider/`)      | SDK (`proxmoxtf/`)          |
|-------------|--------------------------------|-----------------------------|
| Status      | Modern, preferred              | Legacy, maintained          |
| Structure   | Modular sub-packages           | Monolithic files            |
| VM Resource | `vm2` (no clone)               | `vm` (with clone)           |
| Validation  | `int64validator.Between()`     | `validation.IntBetween()`   |
| Errors      | `resp.Diagnostics.AddError()`  | `diag.FromErr()`            |

**Important:** When fixing validation issues, update BOTH providers where applicable to maintain consistency.

### API Client

Both providers share `proxmox.Client`:

```text
proxmox.Client
├── Node(name) → nodes.Client (VMs, containers, network)
├── Cluster() → cluster.Client (ACME, HA, SDN)
├── Access() → access.Client (ACL, users, tokens)
├── Storage() → storage.Client
├── API() → api.Client (raw HTTP)
└── SSH() → ssh.Client (node access)
```

---

## Development Workflow

### Fixing Issues

1. **Create acceptance test** that reproduces the issue
2. **Verify test fails** with current code
3. **Implement fix** in appropriate provider(s)
4. **Verify test passes**
5. **Run related tests** to check for regressions
6. **Run linter:** `make lint`
7. **Verify with mitmproxy** (see [Debugging](#debugging))
8. **Complete checklist** from [Quick Reference](#quick-reference)

### Adding Features

1. New features go in Framework provider only (`fwprovider/`). SDK is feature-frozen.
2. Follow existing patterns; consider adding matching datasource(s).
3. Add validation, acceptance tests, and documentation.
4. Complete production readiness checklist.

See [CONTRIBUTING.md](CONTRIBUTING.md#provider-implementation-guidance) for reference implementations, documentation workflow, and best practices.

---

## Testing

See [CONTRIBUTING.md](CONTRIBUTING.md#testing) for full testing documentation.

### Quick Commands

```bash
./testacc TestAccResourceVM2CPU              # Single test
./testacc TestAccResourceVM2CPU -- -count 1  # With flags
./testacc --all                              # All tests (use sparingly)
./testacc --no-proxy TestName                # Without proxy
```

### Agent-Specific Notes

- **Test requirements:** VMs with `started = true` need boot disk with cloud image; use `stop_on_destroy = true` without qemu-guest-agent
- **Naming:** Do NOT use issue numbers in test names or VM names
- **Sandbox:** Tests may fail due to restricted filesystem access. Request "all" permissions or set `GOCACHE`/`GOMODCACHE` inside workspace.

### Stuck VM Cleanup

```bash
qm set <vmid> --onboot 0 --skiplock
kill -9 $(cat /var/run/qemu-server/<vmid>.pid)
rm -f /var/lock/qemu-server/lock-<vmid>.conf
qm destroy <vmid> --purge --skiplock
```

---

## Debugging

> **Full guide:** [.dev/DEBUGGING.md](.dev/DEBUGGING.md)

### Why Mitmproxy is Critical

**Tests passing ≠ API calls are correct.** Always verify with mitmproxy when:

- Implementing new API parameters
- Modifying existing API calls
- Before declaring a feature complete
- Debugging API errors

### Quick Start

```bash
# 1. Start proxy
mitmdump --flow-detail 2 > /tmp/api.log 2>&1 &

# 2. Run test
./testacc TestAccYourFeature

# 3. Verify your parameter is sent
grep "your_param=" /tmp/api.log

# 4. Stop proxy
pkill -f mitmdump
```

### Flow Detail Levels

| Level | Shows                                |
|-------|--------------------------------------|
| 2     | Full URL + headers (recommended)     |
| 4     | Everything untruncated (deep debug)  |

### Enhanced Debug Script

```bash
mitmdump -s .dev/proxmox_debug_script.py > /tmp/debug.log 2>&1 &
```

Features: API categorization, filtered output, highlighted params, JSON formatting.

---

## Code Patterns

### Schema Definition

**Framework:**

```go
schema.StringAttribute{
    Required: true,
    Validators: []validator.String{
        stringvalidator.OneOf("a", "b"),
    },
}
```

**SDK:**

```go
"key": {
    Type:     schema.TypeString,
    Required: true,
    ValidateDiagFunc: validation.ToDiagFunc(
        validation.StringInSlice([]string{"a", "b"}, false)),
}
```

### Error Handling

**Framework:**

```go
resp.Diagnostics.AddError("Summary", "Detail")
```

**SDK:**

```go
return diag.FromErr(err)
```

### Validation Consistency

Both providers must use identical validation ranges:

- Framework: `int64validator.Between(min, max)`
- SDK: `validation.IntBetween(min, max)`

---

## Session Management

For large efforts (PRs, multi-step refactors), maintain a session state file.

### Template

Use `.dev/SESSION_STATE_TEMPLATE.md` as a starting point. Key sections:

1. **Git State** — Branch, uncommitted changes
2. **Session Summary** — What was accomplished
3. **Implementation Details** — File locations, line numbers, API endpoints
4. **Critical Warnings** — Gotchas discovered
5. **Next Steps** — Prioritized action items
6. **Quick Start** — Commands to resume work

### Naming and Location

- **Location:** `.dev/` directory
- **Pattern:** `<TASK>_SESSION_STATE.md` (e.g., `CLONED_VM_SESSION_STATE.md`)
- **Gitignore:** `*_SESSION_STATE.md` is already in `.gitignore`

### Rules

- Never commit session state files (auto-ignored)
- Never include secrets or credentials
- Update after each meaningful phase
- Include mitmproxy verification results (sanitized)

---

## Agent Behavior

### Knowledge Management

**Update your knowledge/memory when you:**

- Discover important patterns
- Make repeated mistakes
- Learn project-specific conventions
- Find information that contradicts stored knowledge

**Store learnings about:**

- Proxmox VE API behavior
- Terraform provider patterns
- Testing practices
- Common pitfalls

**Never store:**

- Credentials, tokens, or secrets
- Transient task-specific details
- Implementation plans (use session files instead)

### Communication Style

| Do                            | Don't                                     |
|-------------------------------|-------------------------------------------|
| Be concise and direct         | Apologize                                 |
| Use technical terminology     | Summarize changes made                    |
| Explain reasoning             | Ask for confirmation of provided info     |
| Admit uncertainty             | Make up information                       |
| Provide file paths            | Show current implementation unless asked  |

### Code Changes

- Make changes file by file
- Provide all edits for a file in one chunk
- Use explicit, descriptive variable names
- Follow existing code style
- Consider edge cases
- Minimal comments (only essential ones)
- No capitalized comments

### Comments and Style

- **No excessive comments** — Only explain non-obvious logic
- **No redundant comments** — Don't describe what code clearly shows
- **No emojis in code** — Keep code professional and clean
- **Minimal emojis in docs** — Use sparingly, only where they add clarity

### Markdown Files

After updating any markdown file:

1. Run linter to check formatting
2. Fix spelling, structural issues
3. Verify links work
4. Maintain consistent heading levels

---

## Proxmox VE Compatibility

| Version | Support                       |
|---------|-------------------------------|
| 9.x     | Full (target)                 |
| 8.x     | Limited, not testing priority |
| 7.x     | Not supported                 |

### Known Limitations

- Serial device required for Debian 12/Ubuntu VMs to prevent kernel panic on disk resize
- Snippets/backups require PAM account for SFTP upload
- Cluster hardware mappings require `root` PAM account
- Lock errors when creating multiple VMs simultaneously (use `parallelism=1`)

---

## References

- **Contributing Guide:** [CONTRIBUTING.md](CONTRIBUTING.md) — Single source of truth for development workflow
- **Debugging Guide:** [.dev/DEBUGGING.md](.dev/DEBUGGING.md)
- **Session Template:** [.dev/SESSION_STATE_TEMPLATE.md](.dev/SESSION_STATE_TEMPLATE.md)
- **Proxmox API:** <https://pve.proxmox.com/pve-docs/api-viewer/>
- **Terraform Plugin Framework:** <https://developer.hashicorp.com/terraform/plugin/framework>

---

*This document contains agent-specific instructions. For general development workflow, see [CONTRIBUTING.md](CONTRIBUTING.md).*
