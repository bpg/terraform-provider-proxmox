# Session State Template

Use this template to maintain context across multi-step tasks, PR reviews, or complex refactors.

**Instructions:**

1. Copy this template to `.dev/<ISSUE_NUMBER>_SESSION_STATE.md` (e.g., `.dev/1234_SESSION_STATE.md`)
2. Files matching `*_SESSION_STATE.md` are auto-ignored by `.gitignore`
3. Update after each meaningful phase of work
4. **Critical:** Update before ending any session — the next agent depends on this

---

## Issue #[NUMBER] - [Title]

- **GitHub Issue:** `https://github.com/bpg/terraform-provider-proxmox/issues/NUMBER`
- **Last Updated:** YYYY-MM-DD HH:MM
- **Status:** In Progress | Blocked | Ready for Review | Complete
- **Branch:** `fix/NUMBER-short-description` or `feat/NUMBER-short-description`

---

## Quick Context Restore

> **For the next agent:** Read this section first to restore context quickly.

**What this issue is about:**
[1-2 sentences describing the problem/feature]

**Current state:**
[1-2 sentences describing where we are in the implementation]

**Immediate next action:**
[Exactly what to do next — be specific]

---

## Git State

**Current Branch:** `branch-name`
**Base Branch:** `main`

**Uncommitted Changes:**

```text
M  path/to/modified/file.go
A  path/to/new/file.go
```

**Recent Commits on Branch:**

```text
abc1234 fix(vm): description of change
def5678 test(vm): add acceptance test for issue
```

**Action Required:** [Describe any git actions needed before continuing]

---

## User Decisions

> Track explicit decisions made by the user to avoid re-asking.

| Decision | User Choice | Date |
| -------- | ----------- | ---- |
| Which provider to modify? | Framework only | YYYY-MM-DD |
| Test approach? | Unit + acceptance | YYYY-MM-DD |

---

## Assumptions Made

> Track assumptions the agent made. Mark as confirmed/rejected when verified.

| Assumption | Status | Notes |
| ---------- | ------ | ----- |
| Clone operation uses POST /api2/json/nodes/{node}/qemu/{vmid}/clone | Unverified | Need mitmproxy |
| Timeout is in seconds | Confirmed | Verified in API docs |
| SDK provider needs same fix | Rejected | User said Framework only |

---

## What Was Done

### Phase 1: [Name] - Complete

- [Completed task 1]
- [Completed task 2]

### Phase 2: [Name] - In Progress

- [In-progress task]

---

## Implementation Details

### Files Modified

| File | Lines | Purpose |
| ---- | ----- | ------- |
| `path/to/file.go` | 100-150 | [Brief description] |

### Key Code Locations

```text
package/
├── file.go
│   ├── FunctionName() - Line XXX - [Purpose]
│   └── AnotherFunc() - Line YYY - [Purpose]
```

### API Endpoints Used

| Endpoint | Method | Purpose |
| -------- | ------ | ------- |
| `/api2/json/nodes/{node}/qemu/{vmid}` | GET | Fetch VM config |

---

## Context Gathered

> Important findings from codebase exploration. Saves re-reading files.

### Relevant Code Patterns

```go
// Example of how similar features are implemented
// From: fwprovider/nodes/resource_example.go:150
func (r *Resource) Create(ctx context.Context, ...) {
    // Pattern used here
}
```

### API Behavior Notes

- [Observation about API behavior]
- [Edge case discovered]

### Related Code Locations

| What | File | Line | Notes |
| ---- | ---- | ---- | ----- |
| Similar feature | `fwprovider/nodes/vm/clone.go` | 200 | Uses same pattern |
| Validation logic | `fwprovider/validators/timeout.go` | 50 | Reuse this |

---

## Mitmproxy Verification

**Verified API calls:**

```text
GET /api2/json/path?param=value  OK
POST /api2/json/path             OK
```

**Logged to:** `/tmp/api_debug.log`

**Key findings:**

- [What was confirmed/discovered via mitmproxy]

---

## Hypotheses Tested

> For debugging: track what was tried and results.

| Hypothesis | Test | Result |
| ---------- | ---- | ------ |
| Timeout not being sent | Checked mitmproxy log | Wrong - Timeout IS sent |
| Wrong parameter name | Compared with API docs | Correct - Should be `timeout` not `wait` |

---

## Critical Warnings

### Gotchas Discovered

1. **[Issue Name]**
   - Description of the issue
   - How it manifests
   - Workaround or solution

### Do NOT

- [Thing to avoid and why]

### Do

- [Recommended approach and why]

---

## Next Steps

### Immediate (Do Next)

1. **[Task Name]** — [Brief description]
   - Dependencies: None | List them

### Pending

1. **[Task Name]** — [Brief description]
1. **[Task Name]** — [Brief description]

### Blocked

- **[Task Name]** — Blocked by: [Reason]

---

## Questions to Resolve

- [ ] [Question that needs answering]
- [ ] [Another question]

---

## Quick Start Commands

> **Using an LLM agent?** Run `/bpg:resume` to automatically load this session state.

```bash
# Resume work (manual)
cd /path/to/terraform-provider-proxmox
git checkout branch-name
cat .dev/NUMBER_SESSION_STATE.md

# Verify current state
git status
git log --oneline -5
make build

# Run relevant tests
./testacc TestAccRelevantTest

# Before committing
make lint
make docs  # if schema changed
```

---

## Verification Checklist

Before marking complete (or run `/bpg:ready` to automate):

- [ ] GitHub issue exists and is referenced
- [ ] All acceptance tests pass
- [ ] Mitmproxy verification done (use `/bpg:debug-api`)
- [ ] `make lint` passes
- [ ] `make docs` run (if schema changed)
- [ ] No secrets in files
- [ ] Session state file excluded from commit
- [ ] PR body prepared: `.dev/NUMBER_PR_BODY.md` (use `/bpg:prepare-pr`)

---

## Reference Documents

- **GitHub Issue:** `https://github.com/bpg/terraform-provider-proxmox/issues/NUMBER`
- **Project Guidelines:** [CLAUDE.md](../CLAUDE.md)
- **Debugging Guide:** [DEBUGGING.md](DEBUGGING.md)
- **Proxmox API Docs:** `https://pve.proxmox.com/pve-docs/api-viewer/`

---

## Session Log

### YYYY-MM-DD HH:MM - [Agent Name / Human]

**Summary:** [Brief summary of session work]

**Completed:**

- [Item 1]
- [Item 2]

**Issues Encountered:**

- [Issue and resolution]

**Context for Next Session:**

- [Critical context that must not be lost]
- [Key file to look at: `path/to/file.go:123`]

---

<!-- Template version: 2.0 | Aligned with CLAUDE.md guidelines -->
