# Session State Template

Use this template to maintain context across multi-step tasks, PR reviews, or complex refactors.

**Instructions:**

1. Copy this template to `.dev/<TASK>_SESSION_STATE.md` (e.g., `.dev/CLONED_VM_SESSION_STATE.md`)
2. Files matching `*_SESSION_STATE.md` are auto-ignored by `.gitignore`
3. Update after each meaningful phase of work

---

# [Task Name] - Session State

**Last Updated:** YYYY-MM-DD
**Status:** [In Progress | Blocked | Ready for Review | Complete]
**Branch:** `branch-name`

## Git State

**Current Branch:** `branch-name`

**Uncommitted Changes:**
```
M  path/to/modified/file.go
A  path/to/new/file.go
```

**Action Required:** [Describe any git actions needed before continuing]

## Session Summary

[1-2 paragraph summary of what was accomplished and current state]

## What Was Done

### Phase 1: [Name] ✅
- [Completed task 1]
- [Completed task 2]

### Phase 2: [Name] 🚧
- [In-progress task]

## Implementation Details

### Files Modified

| File | Lines | Purpose |
|------|-------|---------|
| `path/to/file.go` | 100-150 | [Brief description] |

### Key Code Locations

```
package/
├── file.go
│   ├── FunctionName() - Line XXX - [Purpose]
│   └── AnotherFunc() - Line YYY - [Purpose]
```

### API Endpoints Used

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api2/json/nodes/{node}/qemu/{vmid}` | GET | Fetch VM config |

### Mitmproxy Verification

**Verified API calls:**
```
GET /api2/json/path?param=value  ✅
POST /api2/json/path             ✅
```

**Logged to:** `/tmp/session_debug.log`

## Critical Warnings

### ⚠️ Gotchas Discovered

1. **[Issue Name]**
   - Description of the issue
   - How it manifests
   - Workaround or solution

### 🚨 Do NOT

- [Thing to avoid and why]

### 💡 Do

- [Recommended approach and why]

## Next Steps

### Immediate (Do Next)

1. **[Task Name]** — [Brief description]
   - Estimated effort: X hours/days
   - Dependencies: [None | List them]

### Pending

2. **[Task Name]** — [Brief description]
3. **[Task Name]** — [Brief description]

### Blocked

- **[Task Name]** — Blocked by: [Reason]

## Questions to Resolve

- [ ] [Question that needs answering]
- [ ] [Another question]

## Quick Start Commands

```bash
# Resume work
cd /path/to/project
git checkout branch-name
cat .dev/TASK_SESSION_STATE.md

# Verify current state
git status
make build

# Continue development
# [specific commands for this task]

# Before committing
make lint
./testacc TestAccRelevantTest
```

## Reference Documents

- **Project Guidelines:** [AGENTS.md](../AGENTS.md)
- **Debugging Guide:** [DEBUGGING.md](DEBUGGING.md)
- **Related Issue:** #XXXX (external reference only, not in commits)

## Verification Checklist

Before marking complete:

- [ ] All acceptance tests pass
- [ ] Mitmproxy verification done
- [ ] `make lint` passes
- [ ] `make docs` run (if schema changed)
- [ ] `make example` passes
- [ ] No secrets in files
- [ ] Progress file excluded from commit

---

## Session Log

### YYYY-MM-DD - [Your Name/Agent]

**Duration:** X hours
**Summary:** [Brief summary of session work]

**Completed:**
- [Item 1]
- [Item 2]

**Issues Encountered:**
- [Issue and resolution]

**Next Session:**
- [What to do next]

---

*Template version: 1.0 | Based on AGENTS.md guidelines*
