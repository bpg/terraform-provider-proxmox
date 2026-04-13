---
name: ready
description: Run production readiness checklist before declaring work complete
argument-hint: \[TestName\]
allowed-tools:
  - Read
  - Bash
  - Grep
  - Glob
  - AskUserQuestion
  - Write
  - Edit
  - Skill
---

<objective>
Execute the full production readiness checklist to verify work is complete and correct.

Use this skill when:

- About to declare a feature or fix complete
- Before creating a PR
- After implementing changes, to verify nothing was missed
- User asks "is this ready?" or "am I done?"

**This checklist is mandatory.** Never skip steps. **If any step fails, stop immediately and report the failure. Do not continue to the next step until the failure is resolved.**
</objective>

<context>
Test name (optional): $ARGUMENTS

Executes the Production Readiness Checklist from [CLAUDE.md](../../../CLAUDE.md#production-readiness-checklist).
</context>

<process>

## Step 0: Determine Issue Context

Detect issue number from branch name:

```bash
# Match number after fix/ or feat/ prefix to avoid matching unrelated numbers
ISSUE_NUM=$(git branch --show-current | grep -oE '(fix|feat)/[0-9]+' | grep -oE '[0-9]+')
```

If issue number was detected, confirm with user. If not detected, ask:

```text
AskUserQuestion(
  header: "Issue Number",
  question: "What GitHub issue number is this work for?",
  options: [
    { label: "Use #{ISSUE_NUM}", description: "Detected from branch name" },
    { label: "Enter manually", description: "I'll type the issue number" }
  ]
)
```

Note: Only show "Use #{ISSUE_NUM}" option if detection succeeded.

**Check for session state file:**

```bash
SESSION_STATE=".dev/${ISSUE_NUM}_SESSION_STATE.md"
if [ -f "$SESSION_STATE" ]; then
  echo "Found session state: $SESSION_STATE"
fi
```

If session state exists, read it for context (test names, previous results, etc.).

If no test name provided via `$ARGUMENTS`, detect from changes:

```bash
git diff --name-only main...HEAD | grep -E "\.go$" | head -10
```

Store issue number for PR preparation.

---

## Step 1: Build Check

```bash
echo "=== Step 1: Build ==="
make build
BUILD_EXIT=$?
```

**Result:**

- Exit 0: "BUILD PASSED"
- Exit non-zero: "BUILD FAILED — Fix build errors before continuing"

---

## Step 2: Lint Check

### Go Lint

```bash
echo "=== Step 2a: Go Lint ==="
make lint
LINT_EXIT=$?
```

**Result:**

- "0 issues": "GO LINT PASSED"
- Issues found: "GO LINT FAILED — Run `make lint` to auto-fix, then review changes"

### Markdown Lint

Find changed markdown files and lint them:

```bash
echo "=== Step 2b: Markdown Lint ==="
MD_FILES_TO_LINT=()
while IFS= read -r file; do MD_FILES_TO_LINT+=("$file"); done < <(git diff --name-only main...HEAD | grep -E '\.md$')

if [ ${#MD_FILES_TO_LINT[@]} -gt 0 ]; then
  npx --yes markdownlint-cli2 --fix "${MD_FILES_TO_LINT[@]}"
  MD_LINT_EXIT=$?
else
  echo "No markdown files changed"
  MD_LINT_EXIT=0
fi
```

**Result:**

- Exit 0: "MARKDOWN LINT PASSED"
- Issues found: "MARKDOWN LINT FAILED — Review and fix remaining issues"

### ADR-007: Resource Type Name Convention

Check if any **new** resources or data sources were added, and verify they follow [ADR-007](docs/adr/007-resource-type-name-migration.md) Phase 1: new resources must use the `proxmox_` prefix with a hardcoded `TypeName`, not `req.ProviderTypeName + "_suffix"`.

```bash
echo "=== Step 2c: ADR-007 Naming Convention ==="
# Find new Go files added in this branch that contain Metadata functions (resources/data sources)
NEW_FILES=$(git diff --name-only --diff-filter=A main...HEAD | grep -E '\.go$')
if [ -n "$NEW_FILES" ]; then
  # Check for the old pattern: req.ProviderTypeName + "_something"
  VIOLATIONS=$(grep -l 'req\.ProviderTypeName' $NEW_FILES 2>/dev/null)
  if [ -n "$VIOLATIONS" ]; then
    echo "ADR-007 VIOLATION: New files use req.ProviderTypeName instead of hardcoded proxmox_ prefix:"
    echo "$VIOLATIONS"
    grep -n 'req\.ProviderTypeName' $VIOLATIONS
  else
    echo "No ADR-007 violations in new files"
  fi
else
  echo "No new Go files added"
fi
```

**Result:**

- No violations or no new files: "ADR-007 PASSED"
- Violations found: "ADR-007 FAILED — New resources/data sources must hardcode `resp.TypeName = \"proxmox_...\"` per [ADR-007](docs/adr/007-resource-type-name-migration.md) Phase 1"

---

## Step 3: Unit Tests

```bash
echo "=== Step 3: Unit Tests ==="
make test
TEST_EXIT=$?
```

**Result:**

- Exit 0: "UNIT TESTS PASSED"
- Exit non-zero: "UNIT TESTS FAILED — Fix failing tests"

---

## Step 4: Acceptance Tests

Determine which tests to run:

If specific test provided via `$ARGUMENTS`:

```bash
TEST_PATTERN="${ARGUMENTS}"
```

Otherwise, try to detect from changes:

```bash
# Find test files changed compared to main branch
CHANGED_TESTS=$(git diff --name-only main...HEAD | grep "_test.go" | head -3)
if [ -n "$CHANGED_TESTS" ]; then
  # Extract TestAcc function names from changed files
  TEST_NAMES=$(grep -h "^func TestAcc" $CHANGED_TESTS 2>/dev/null | sed 's/func \(TestAcc[^(]*\).*/\1/')
  if [ -n "$TEST_NAMES" ]; then
    # Join with .* pattern for matching multiple tests
    TEST_PATTERN=$(echo "$TEST_NAMES" | head -1)
    echo "Detected test: $TEST_PATTERN"
  fi
fi
```

If no tests detected, ask user:

```text
AskUserQuestion(
  header: "Acceptance Test",
  question: "Which acceptance test(s) should I run?",
  options: [
    { label: "Skip", description: "No acceptance tests for this change" },
    { label: "Enter name", description: "I'll provide the test name" }
  ]
)
```

Run the tests with verbose output and capture to log:

```bash
echo "=== Step 4: Acceptance Tests ==="
if [ -n "$TEST_PATTERN" ]; then
  ./testacc "$TEST_PATTERN" -- -v 2>&1 | tee /tmp/testacc.log
  ACC_EXIT=$?
else
  echo "No acceptance tests specified"
  ACC_EXIT=0
fi
```

The test output is saved to `/tmp/testacc.log` for use in `/bpg:prepare-pr`.

**Result:**

- Exit 0: "ACCEPTANCE TESTS PASSED"
- Exit non-zero: "ACCEPTANCE TESTS FAILED — Fix failing tests"

---

## Step 5: API Verification

First, check if mitmproxy is available and configured:

```bash
MITM_AVAILABLE="no"
which mitmdump > /dev/null 2>&1 && MITM_AVAILABLE="yes"
PROXY_CONFIGURED="no"
grep -q '^HTTPS_PROXY' testacc.env 2>/dev/null && PROXY_CONFIGURED="yes"
echo "mitmproxy installed: $MITM_AVAILABLE, proxy configured in testacc.env: $PROXY_CONFIGURED"
```

Check if acceptance tests already include behavioral assertions (e.g., uptime-based reboot detection, status checks via the Proxmox API):

```bash
# Look for direct API verification in test code (uptime checks, status checks, etc.)
git diff --name-only HEAD | grep "_test.go" | xargs grep -l "GetVMStatus\|Uptime\|NodeClient" 2>/dev/null
```

Then ask the user:

```text
AskUserQuestion(
  header: "API Verify",
  question: "Does this change involve API calls that need verification?",
  options: [
    { label: "Yes, run mitmproxy", description: "Start mitmproxy and capture API traffic" },
    { label: "Tests verify behavior", description: "Acceptance tests already assert correct behavior" },
    { label: "No API changes", description: "Docs-only, refactor, etc." },
    { label: "Already done", description: "I verified API calls earlier" }
  ]
)
```

If "Yes, run mitmproxy":

- Check mitmproxy is installed and proxy is configured
- If not installed: "mitmproxy not found. Install with `brew install mitmproxy` or use behavioral assertions in tests instead."
- If not configured in `testacc.env`: "Proxy not configured. Uncomment `HTTP_PROXY` and `HTTPS_PROXY` in `testacc.env`, or run `/bpg:debug-api {test_name}` which sets them inline."
- Otherwise: "Run `/bpg:debug-api {test_name}` to verify API calls"

If "Tests verify behavior": Record as "Verified via test assertions" — this is acceptable when tests directly check the behavior (e.g., uptime-based reboot detection, API status checks) rather than just Terraform state attributes.

If "No API changes" or "Already done": Record as skipped/completed.

---

## Step 6: Documentation

Check if schema was changed:

```bash
SCHEMA_CHANGED=$(git diff --name-only main...HEAD | grep -E "(schema|resource|datasource).*\.go$" | wc -l)
```

If schema changed:

```bash
echo "=== Step 6: Documentation ==="
make docs
DOCS_EXIT=$?
git diff --name-only docs/
```

**Result:**

- Exit 0 and no diff: "DOCS PASSED (no changes needed)"
- Exit 0 with diff: "DOCS GENERATED — Review changes in docs/"
- Exit non-zero: "DOCS FAILED"

If schema not changed:

- "DOCS SKIPPED (no schema changes detected)"

### Breaking Changes → Upgrade Guide

Check if the branch contains breaking changes. If so, verify `docs/guides/upgrade.md` was updated:

```bash
echo "=== Step 6b: Upgrade Guide Check ==="
# Detect breaking changes from committed messages
BREAKING=$(git log main...HEAD --oneline 2>/dev/null | grep -iE '!:|BREAKING')
if [ -n "$BREAKING" ]; then
  echo "Breaking changes detected:"
  echo "$BREAKING"
  # Check if upgrade guide was modified (committed or uncommitted)
  UPGRADE_CHANGED=0
  git diff --name-only main...HEAD 2>/dev/null | grep -qc 'docs/guides/upgrade.md' && UPGRADE_CHANGED=1
  git diff --name-only 2>/dev/null | grep -qc 'docs/guides/upgrade.md' && UPGRADE_CHANGED=1
  if [ "$UPGRADE_CHANGED" -eq 0 ]; then
    echo "UPGRADE GUIDE NOT UPDATED"
  else
    echo "Upgrade guide updated"
  fi
else
  echo "No breaking changes detected in commits"
  echo "Note: if breaking changes are identified during PR preparation, /bpg:prepare-pr will check the upgrade guide"
fi
```

**Result:**

- No breaking changes: "UPGRADE GUIDE SKIPPED (no breaking changes)"
- Breaking changes + guide updated: "UPGRADE GUIDE PASSED"
- Breaking changes + guide NOT updated: "UPGRADE GUIDE FAILED — Breaking changes must be documented in `docs/guides/upgrade.md`. Add a section for the new version with: description of the change, before/after behavior, and action required."

---

## Step 7: Implementation Scrutiny

**This step is mandatory. Do NOT rubber-stamp the implementation.**

After all mechanical checks pass, critically examine the actual code changes. The goal is to catch logic errors, missed edge cases, incomplete fixes, and pattern violations that automated tools cannot detect.

### 7a: Read the full diff

```bash
git diff main...HEAD
# If no commits yet, fall back to unstaged diff:
git diff
```

### 7b: Self-review checklist

For each changed file, answer these questions. Write your answers out explicitly — do not skip any:

1. **Correctness:** Does this change actually fix the reported problem / implement the requested feature? Trace through the code path mentally. Could a caller still hit the original bug through a different path?
2. **Completeness:** Are there other locations in the codebase with the same pattern that also need fixing? Search proactively:
   ```bash
   # Search for the OLD pattern that was replaced — if it still exists elsewhere, flag it
   ```
3. **Edge cases:** What inputs, error conditions, or timing scenarios could break this? Consider: nil/empty values, concurrent access, error wrapping chains, API behavior differences across PVE versions.
4. **Consistency:** Does the fix follow the same pattern used by other resources in the codebase? Find at least one reference example and compare.
5. **Regression risk:** Could this change break existing behavior? Consider callers of the modified functions — do they all handle the new behavior correctly?
6. **Error handling:** Are errors properly propagated, wrapped, or classified? Check `errors.Is`/`errors.As` chains, sentinel error usage, and diagnostic messages.
7. **Test coverage gap:** If acceptance tests don't cover the specific changed behavior, is there a concrete reason (e.g., can't inject API errors), or is the test just missing?

### 7c: Evaluate and decide

**If you find issues:** Stop, report each issue with file path and line number, and provide a concrete fix recommendation. Do NOT continue to the summary step. The issues must be fixed and this step re-run.

**If the change involves non-obvious design decisions** (e.g., choosing between approaches, new patterns, behavioral trade-offs): Invoke the `/grill-me` skill to interrogate the user about those decisions before proceeding. Examples of when to grill:
- The fix changes observable behavior (not just error messages)
- The approach differs from what the issue suggested
- There are multiple valid ways to solve it and the trade-offs aren't obvious
- The change touches shared/core code that many resources depend on

**If the implementation is clean:** State explicitly what you verified and why you're confident, then proceed.

**Result:**

- All questions answered satisfactorily: "SCRUTINY PASSED"
- Issues found: "SCRUTINY FAILED — {list of issues}"
- Design decisions need validation: invoke `/grill-me`, then re-evaluate

---

## Step 8: Summary and Proof of Work Report

Generate summary:

```text
=== PRODUCTION READINESS CHECKLIST ===

Issue: #${ISSUE_NUM}
Date: $(date +%Y-%m-%d)

| Step | Status |
|------|--------|
| Build | ${BUILD_STATUS} |
| Lint | ${LINT_STATUS} |
| ADR-007 Naming | ${ADR007_STATUS} |
| Unit Tests | ${TEST_STATUS} |
| Acceptance Tests | ${ACC_STATUS} |
| API Verification | ${API_STATUS} |
| Documentation | ${DOCS_STATUS} |
| Upgrade Guide | ${UPGRADE_GUIDE_STATUS} |
| Scrutiny | ${SCRUTINY_STATUS} |

Overall: ${OVERALL_STATUS}
```

### Write Proof of Work Report

Write a `.dev/${ISSUE_NUM}_REPORT.md` file containing the proof of work evidence. This file is required by CONTRIBUTING.md for AI-assisted contributions and is used by `/bpg:prepare-pr` to fill the PR body.

**First, assess test coverage quality.** Read the test code and the implementation diff, then classify:

- **Strong** — Tests precisely replicate and validate the behavior from the PR. A test would fail without the change and pass with it.
- **Partial** — Tests exercise the resource but don't specifically target the changed behavior.
- **None** — No acceptance tests, or tests don't exercise the changed code paths.

**Then write the report, scaling evidence depth to coverage:**

```markdown
# Proof of Work Report — Issue #${ISSUE_NUM}

**Date:** YYYY-MM-DD
**Branch:** ${BRANCH_NAME}
**Test Coverage:** Strong / Partial / None

## Checklist Results

| Step | Status |
|------|--------|
| Build | PASSED/FAILED |
| Lint | PASSED/FAILED |
| ADR-007 Naming | PASSED/SKIPPED |
| Unit Tests | PASSED/FAILED |
| Acceptance Tests | PASSED/FAILED/SKIPPED |
| API Verification | PASSED/SKIPPED |
| Documentation | PASSED/SKIPPED |
| Upgrade Guide | PASSED/SKIPPED |

## Acceptance Test Output

\`\`\`
<command used>
<trimmed test output: RUN/PASS/FAIL lines + summary>
\`\`\`

## Additional Evidence

<Include this section when test coverage is Partial or None>

### Terraform Configuration

\`\`\`hcl
<the HCL config used to exercise the change>
\`\`\`

### terraform plan Output

\`\`\`
<plan output showing expected diff — especially for bug fixes>
\`\`\`

### terraform apply Output

\`\`\`
<apply output showing successful creation/modification>
\`\`\`

### Before / After

<For bug fixes: what happened before the fix vs after.
E.g. "Before: plan showed unexpected diff on every run. After: clean plan.">

## API Verification

<mitmproxy logs, behavioral assertion details, or "N/A — no API changes">

## Notes

<any additional context, decisions, or caveats>
```

**What to include based on coverage level:**

- **Strong:** Checklist + test output is sufficient. Omit the "Additional Evidence" section.
- **Partial:** Checklist + test output + at least one item from Additional Evidence (terraform config + plan/apply output, or before/after comparison).
- **None:** Checklist + full Additional Evidence section (terraform config, plan output, apply output, and before/after if applicable). This is the minimum to prove the change works without tests.

If all passed:

```text
All checks passed. Proof of work saved to: .dev/${ISSUE_NUM}_REPORT.md

Next step: Prepare PR body.
Run: /prepare-pr ${ISSUE_NUM}
```

If any failed:

```text
CHECKLIST INCOMPLETE

Failed steps: ${FAILED_STEPS}

Fix the issues above and run /bpg:ready again.
```

## Step 9: Update Session State

Detect issue number and update `.dev/${ISSUE_NUM}_SESSION_STATE.md` using Read and Edit tools:

- `Last Updated:` → current date
- `Current state:` → checklist results summary (e.g., "All checks passed" or "Blocked on failing lint")
- `Immediate next action:` → if all passed: "Prepare PR with `/bpg:prepare-pr`"; if failed: description of what to fix
- Add/update a **Verification Results** section with the checklist table from Step 7
- `Status:` → "Ready for PR" if all passed, keep "In Progress" if any failed

</process>

<success_criteria>

- [ ] Issue number determined
- [ ] Session state checked for context
- [ ] Build passes
- [ ] Go lint shows 0 issues
- [ ] Markdown lint passes on changed `.md` files
- [ ] ADR-007 naming convention check passes (new resources use `proxmox_` prefix)
- [ ] Unit tests pass
- [ ] Acceptance tests pass (or explicitly skipped)
- [ ] API verification done (or explicitly skipped for non-API changes)
- [ ] Documentation regenerated (if schema changed)
- [ ] Upgrade guide updated (if breaking changes detected)
- [ ] Implementation scrutiny passed (all 7 questions answered, no issues found)
- [ ] `/grill-me` invoked if non-obvious design decisions detected
- [ ] Summary presented to user
- [ ] Session state updated with results
</success_criteria>

<tips>
- If you're unsure which acceptance tests to run, look for tests matching the resource/datasource name
- Schema changes = any modification to attribute definitions, validators, or type definitions
- The checklist is designed to catch issues before PR review, saving time for everyone
- Test output is saved to `/tmp/testacc.log` - this persists for `/bpg:prepare-pr` to use
- If session state exists, update it with results so `/bpg:prepare-pr` can skip re-verification
</tips>
