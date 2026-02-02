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

**Usage:**

```bash
mitmdump --mode regular --listen-port 8080 -s .dev/proxmox_debug_script.py > /tmp/debug.log 2>&1 &
./testacc TestAccYourTest
cat /tmp/debug.log
pkill -f mitmdump
```

## Quick Start

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

## Related Documentation

- **Agent Instructions:** [CLAUDE.md](../CLAUDE.md) — Primary guidelines for AI-assisted development
- **Debugging Details:** [DEBUGGING.md](DEBUGGING.md) — In-depth debugging guide

## Available Skills

Skills automate common workflows. Use them with `/skill-name`:

| Skill | Purpose |
| ----- | ------- |
| `/start-issue` | Start work on a GitHub issue (branch + session state) |
| `/resume` | Resume work from a previous session |
| `/ready` | Run production readiness checklist |
| `/debug-api` | Debug API calls with mitmproxy |
| `/proof-report` | Generate proof of work document |

## Contributing

When adding new tools or scripts to this directory:

1. Document the tool in this README
2. Reference it from DEBUGGING.md if applicable
3. Update CLAUDE.md if it's a commonly-used workflow
