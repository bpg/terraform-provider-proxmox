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
---

<objective>
Execute the full production readiness checklist to verify work is complete and correct.

Use this skill when:

- About to declare a feature or fix complete
- Before creating a PR
- After implementing changes, to verify nothing was missed
- User asks "is this ready?" or "am I done?"

**This checklist is mandatory.** Never skip steps. If a step fails, stop and fix before continuing.
</objective>

<context>
Test name (optional): $ARGUMENTS

Executes the Production Readiness Checklist from [CLAUDE.md](../../CLAUDE.md#production-readiness-checklist).
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

Store issue number for proof report generation.

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

If failed, stop and report. Do not continue.

---

## Step 2: Lint Check

```bash
echo "=== Step 2: Lint ==="
make lint
LINT_EXIT=$?
```

**Result:**

- "0 issues": "LINT PASSED"
- Issues found: "LINT FAILED — Run `make lint` to auto-fix, then review changes"

If failed, stop and report. Do not continue.

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

If failed, stop and report. Do not continue.

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

The test output is saved to `/tmp/testacc.log` for use in the proof report.

**Result:**

- Exit 0: "ACCEPTANCE TESTS PASSED"
- Exit non-zero: "ACCEPTANCE TESTS FAILED — Fix failing tests"

If failed, stop and report. Do not continue.

---

## Step 5: API Verification (Mitmproxy)

Ask if API verification is needed:

```text
AskUserQuestion(
  header: "API Verification",
  question: "Does this change involve API calls that need mitmproxy verification?",
  options: [
    { label: "Yes", description: "Run mitmproxy verification" },
    { label: "No", description: "No API changes (docs-only, refactor, etc.)" },
    { label: "Already done", description: "I verified API calls earlier" }
  ]
)
```

If "Yes":

- Suggest: "Run `/debug-api {test_name}` to verify API calls"
- Or run inline mitmproxy check (abbreviated version)

If "No" or "Already done": Record as skipped/completed.

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

---

## Step 7: Summary and Proof Report Prompt

Generate summary:

```text
=== PRODUCTION READINESS CHECKLIST ===

Issue: #${ISSUE_NUM}
Date: $(date +%Y-%m-%d)

| Step | Status |
|------|--------|
| Build | ${BUILD_STATUS} |
| Lint | ${LINT_STATUS} |
| Unit Tests | ${TEST_STATUS} |
| Acceptance Tests | ${ACC_STATUS} |
| API Verification | ${API_STATUS} |
| Documentation | ${DOCS_STATUS} |

Overall: ${OVERALL_STATUS}
```

If all passed:

```text
All checks passed.

Next step: Create proof of work report.
Run: /proof-report ${ISSUE_NUM}
```

If any failed:

```text
CHECKLIST INCOMPLETE

Failed steps: ${FAILED_STEPS}

Fix the issues above and run /ready again.
```

**Update session state:**

Ensure `.dev/` directory exists:

```bash
mkdir -p .dev
```

If `.dev/${ISSUE_NUM}_SESSION_STATE.md` exists, update it with:

- Checklist results (which steps passed/failed)
- Test log location (`/tmp/testacc.log`)
- Timestamp of this run

This allows `/proof-report` to pull results without re-running checks.

</process>

<success_criteria>

- [ ] Issue number determined
- [ ] Session state checked for context
- [ ] Build passes
- [ ] Lint shows 0 issues
- [ ] Unit tests pass
- [ ] Acceptance tests pass (or explicitly skipped)
- [ ] API verification done (or explicitly skipped for non-API changes)
- [ ] Documentation regenerated (if schema changed)
- [ ] Summary presented to user
- [ ] Session state updated with results
</success_criteria>

<tips>
- If you're unsure which acceptance tests to run, look for tests matching the resource/datasource name
- Schema changes = any modification to attribute definitions, validators, or type definitions
- The checklist is designed to catch issues before PR review, saving time for everyone
- Test output is saved to `/tmp/testacc.log` - this persists for `/proof-report` to use
- If session state exists, update it with results so `/proof-report` can skip re-verification
</tips>
