# Development Tools

This directory contains tools and documentation for developing the Terraform Proxmox Provider.

## Files

### DEBUGGING.md

Comprehensive guide for debugging the provider, including:

- Using mitmproxy to intercept and analyze API calls
- Common debugging scenarios
- Log analysis techniques
- Troubleshooting tips

### SESSION_STATE_TEMPLATE.md

Template for maintaining context across multi-step tasks. Use this when:

- Working on large PRs or refactors
- Implementing features that span multiple sessions
- Debugging complex issues

**Usage:**

1. Copy to `.dev/<ISSUE_NUMBER>_SESSION_STATE.md` (e.g., `1234_SESSION_STATE.md`)
2. Files matching `*_SESSION_STATE.md` are auto-ignored by `.gitignore`
3. Update after each meaningful phase

> **Note:** The following patterns are auto-ignored by `.gitignore`:
> `*_SESSION_STATE.md`, `*_PROGRESS.md`, `*_REPORT.md`, `*_PLAN.md`

### proxmox_debug_script.py

Enhanced mitmproxy script for analyzing provider-API interactions. Features:

- **Categorizes API calls** — Automatically identifies Storage, VM, Container, Network, Cluster APIs
- **Filters output** — Shows only Proxmox API calls (excludes other traffic noise)
- **Highlights key parameters** — Marks important query params: `content`, `vmid`, `node`, `storage`, `type`
- **Pretty-prints JSON** — Formatted request/response bodies for readability
- **Shows data counts** — Displays number of items returned from list APIs
- **Visual indicators** — ✅ for success (2xx), ❌ for errors (4xx, 5xx)

---

## Working with LLM Agents

This section documents how to effectively use LLM agents (like Claude Code) for development work on this provider. The workflow is designed to maintain context across sessions, ensure quality, and produce verifiable results.

### Workflow Overview

```text
┌──────────────────────────────────────────────────────────────────────┐
│                        DEVELOPMENT WORKFLOW                          │
└──────────────────────────────────────────────────────────────────────┘

  ┌─────────────┐
  │ GitHub Issue│  ← All work starts with an issue
  └──────┬──────┘
         │
         ▼
  ┌─────────────┐     Creates:
  │/start-issue │────→ • Branch: fix/1234-description
  └──────┬──────┘      • Session state: .dev/1234_SESSION_STATE.md
         │             • Clears stale logs
         │
         ▼
  ┌─────────────┐
  │ Development │  ← Write code, create tests
  │             │
  │ (implement) │
  └──────┬──────┘
         │
         │ ┌─────────────┐
         ├─┤ /debug-api  │  ← Use during development to verify
         │ └─────────────┘    API calls are correct
         │
         ▼
  ┌─────────────┐     Runs:
  │   /ready    │────→ • make build, lint, test
  └──────┬──────┘      • Acceptance tests (with logging)
         │             • API verification prompt
         │             • Documentation check
         │
         ▼
  ┌─────────────┐     Creates:
  │/proof-report│────→ • .dev/1234_PROOF_REPORT.md
  └──────┬──────┘      • Test output evidence
         │             • API verification evidence
         │
         ▼
  ┌─────────────┐
  │  Submit PR  │  ← Use proof report content in PR description
  └─────────────┘


  ════════════════════════════════════════════════════════════════════

  RESUMING WORK (after break, context loss, or new session):

  ┌─────────────┐     Loads:
  │  /resume    │────→ • Session state context
  └──────┬──────┘      • Git state verification
         │             • Existing log files
         │             • Immediate next action
         ▼
  Continue from where you left off...
```

### Skills Reference

#### `/start-issue [issue-number]`

**When to use:** Beginning work on any GitHub issue.

**What it does:**

1. Verifies the issue exists on GitHub
2. Determines issue type (bug fix → `fix/`, feature → `feat/`)
3. Creates branch with proper naming: `{type}/{issue}-{description}`
4. Creates session state file from template
5. Populates session state with issue context
6. Clears stale log files from previous work

**Example:**

```text
You: /start-issue 1234
Agent: Creates fix/1234-vm-clone-timeout branch, session state, displays issue summary
```

#### `/resume [issue-number]`

**When to use:**

- Starting a new conversation to continue previous work
- After context loss or session timeout
- Returning to work after a break

**What it does:**

1. Lists available session state files (if no issue specified)
2. Loads session context and displays quick restore info
3. Verifies git state matches session (prompts to switch if needed)
4. Shows existing log files from previous runs
5. Displays immediate next action

**Example:**

```text
You: /resume
Agent: Shows available sessions, loads context, displays "Immediate Next Action: Verify test passes after fix"
```

#### `/debug-api [TestName] [parameter]`

**When to use:**

- Implementing new API parameters
- Debugging unexpected API behavior
- Verifying fix sends correct parameters
- When tests pass but behavior seems wrong

**What it does:**

1. Starts mitmproxy on port 8080
2. Runs acceptance test with proxy settings
3. Captures all API traffic to `/tmp/api_debug.log`
4. Analyzes traffic for specific parameters
5. Reports findings with recommendations

**Example:**

```text
You: /debug-api TestAccResourceVM content
Agent: Starts proxy, runs test, shows API calls containing "content" parameter
```

**Key insight:** Tests passing ≠ correct API calls. Always verify with mitmproxy for API changes.

#### `/ready [TestName]`

**When to use:**

- Before declaring work complete
- Before creating a PR
- After implementing changes

**What it does:**

1. Runs `make build` — Must pass
2. Runs `make lint` — Must show 0 issues
3. Runs `make test` — All unit tests pass
4. Runs acceptance tests with verbose output → `/tmp/testacc.log`
5. Prompts for API verification status
6. Checks if documentation needs regeneration
7. Updates session state with results

**Example:**

```text
You: /ready TestAccResourceVMClone
Agent: Runs all checks, reports status, suggests /proof-report if all pass
```

#### `/proof-report [issue-number]`

**When to use:**

- After `/ready` passes all checks
- Before submitting a PR
- When asked to document the work

**What it does:**

1. Gathers test results from `/tmp/testacc.log`
2. Gathers API verification from `/tmp/api_debug.log`
3. Pulls context from session state
4. Verifies checklist items
5. Creates `.dev/{issue}_PROOF_REPORT.md`

**Output format matches PR template:** Summary, test output, test coverage table, API verification, checklist.

**Example:**

```text
You: /proof-report 1234
Agent: Creates .dev/1234_PROOF_REPORT.md with full evidence, ready for PR description
```

### Common Scenarios

#### Scenario 1: Fix a Bug

```text
1. /start-issue 1234           ← Setup branch and session
2. [Investigate and implement fix]
3. /debug-api TestAccBugFix    ← Verify API calls
4. /ready TestAccBugFix        ← Run full checklist
5. /proof-report 1234          ← Generate evidence
6. [Create PR using proof report content]
```

#### Scenario 2: Resume After Break

```text
1. /resume 1234                ← Load context
2. [Agent shows: "Next action: Run tests after implementing fix"]
3. [Continue from where you left off]
```

#### Scenario 3: Context Window Full

When the agent's context fills up during long work:

```text
1. Agent updates session state before context loss
2. [New conversation]
3. /resume 1234                ← Restore full context
4. [Continue seamlessly]
```

#### Scenario 4: Multiple Issues

```text
1. /start-issue 1234           ← Work on first issue
2. [Complete work, /ready, /proof-report]
3. /start-issue 5678           ← Start second issue (clears logs)
4. [Work on second issue]
```

### Shared State Between Skills

The skills share state through files:

| File | Written By | Read By |
| ---- | ---------- | ------- |
| `.dev/{issue}_SESSION_STATE.md` | `/start-issue`, all skills update | `/resume`, `/ready`, `/proof-report` |
| `/tmp/testacc.log` | `/ready`, `/debug-api` | `/proof-report`, `/resume` |
| `/tmp/api_debug.log` | `/debug-api` | `/proof-report`, `/resume` |

This allows:

- `/proof-report` to use test results from `/ready` without re-running
- `/resume` to note existing logs from previous runs
- Session state to accumulate context across the workflow

### Tips for Effective Agent-Assisted Development

1. **Always start with `/start-issue`** — Sets up proper branch naming and session tracking

2. **Update session state frequently** — The agent will do this, but remind it before long operations

3. **Use `/debug-api` liberally** — API verification catches bugs that tests miss

4. **Don't skip `/ready`** — The checklist exists because each item has caught real bugs

5. **Keep proof reports** — They're gitignored but useful for PR descriptions and future reference

6. **Resume, don't restart** — After breaks, use `/resume` instead of re-explaining context

7. **Trust but verify** — Review the agent's work, especially for complex logic

---

## Quick Start (Manual Debugging)

### Debug API Calls

```bash
# Start proxy with enhanced script
mitmdump -s .dev/proxmox_debug_script.py --flow-detail 2 > /tmp/debug.log 2>&1 &

# Run your test
./testacc TestAccDatasourceFile

# View the output (shows categorized API calls)
cat /tmp/debug.log | grep "API"

# Stop the proxy
pkill -f mitmdump
```

### Verify New API Parameter

```bash
# Basic proxy with full URLs
mitmdump --flow-detail 2 > /tmp/test.log 2>&1 &

# Run test
./testacc TestAccYourNewFeature

# Check if parameter is sent
grep "your_param=" /tmp/test.log

# Cleanup
pkill -f mitmdump
```

---

## Related Documentation

- **Agent Instructions:** [CLAUDE.md](../CLAUDE.md) — Primary guidelines for AI-assisted development
- **Debugging Details:** [DEBUGGING.md](DEBUGGING.md) — In-depth debugging guide
- **Session Template:** [SESSION_STATE_TEMPLATE.md](SESSION_STATE_TEMPLATE.md) — Template for session files

---

## Contributing

When adding new tools or scripts to this directory:

1. Document the tool in this README
2. Reference it from DEBUGGING.md if applicable
3. Update CLAUDE.md if it's a commonly-used workflow
4. If adding a new skill, follow the patterns in `.claude/commands/`
