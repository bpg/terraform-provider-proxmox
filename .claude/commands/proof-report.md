---
name: proof-report
description: Generate proof of work document for PR submission
argument-hint: \[issue-number\]
allowed-tools:
  - Read
  - Write
  - Bash
  - Grep
  - Glob
  - AskUserQuestion
---

<objective>
Generate a `.dev/<ISSUE_NUMBER>_PROOF_REPORT.md` document that serves as proof of work for PR submission.

Use this skill when:

- Completing a feature or bug fix
- Preparing to submit a PR
- After running `/ready` checklist
- User asks to "document the work" or "create proof"

The proof report captures test results, API verification, and implementation summary for the PR description.
</objective>

<context>
Issue number: $ARGUMENTS

Report location: `.dev/<ISSUE_NUMBER>_PROOF_REPORT.md`
Pattern is gitignored, so reports won't be committed.

From [CLAUDE.md](../../CLAUDE.md): "After completing work, create `.dev/{issue}_PROOF_REPORT.md`"
</context>

<process>

## Step 1: Determine Issue Number

If `$ARGUMENTS` provided, use it. Otherwise detect from branch name:

```bash
# Try to detect from branch name (e.g., fix/1234-description or feat/1234-description)
# Match number after fix/ or feat/ prefix to avoid matching unrelated numbers
ISSUE_NUM=$(git branch --show-current | grep -oE '(fix|feat)/[0-9]+' | grep -oE '[0-9]+')
```

If still unclear, ask:

```text
AskUserQuestion(
  header: "Issue Number",
  question: "What is the GitHub issue number for this work?",
  options: [
    { label: "From branch", description: "Use: #${ISSUE_NUM}" },
    { label: "Enter manually", description: "I'll type the issue number" }
  ]
)
```

Set the report path:

```bash
REPORT_PATH=".dev/${ISSUE_NUM}_PROOF_REPORT.md"
```

## Step 2: Check Prerequisites

**Ensure `.dev/` directory exists:**

```bash
mkdir -p .dev
```

**Check for session state file:**

```bash
SESSION_STATE=".dev/${ISSUE_NUM}_SESSION_STATE.md"
if [ -f "$SESSION_STATE" ]; then
  echo "Found session state: $SESSION_STATE"
fi
```

If session state exists, read it to extract:

- User decisions made during the work
- Hypotheses tested and their results
- Key implementation notes

## Step 3: Gather Test Results

**Find related acceptance tests:**

```bash
# Search for tests related to the changes (compared to main branch)
TEST_FILES=$(git diff --name-only main...HEAD | grep "_test.go" | head -5)
```

**Determine test name to run:**

Extract test function names from changed test files:

```bash
# Find TestAcc functions in changed test files
grep -h "^func TestAcc" ${TEST_FILES} | sed 's/func \(TestAcc[^(]*\).*/\1/' | head -5
```

If multiple tests found, ask user which to run. If session state exists, check if a test name was recorded there.

**Run the test with verbose output and capture to log:**

If tests haven't been run yet, or user wants fresh results:

```bash
./testacc ${TEST_NAME} -- -v 2>&1 | tee /tmp/testacc.log
```

This captures detailed test output for the report. Record the log file location in the session state if it exists.

**Get recent test output (if available):**

```bash
# Check if tests were run recently
if [ -f /tmp/testacc.log ]; then
  RECENT_OUTPUT=$(tail -50 /tmp/testacc.log)
  echo "Found test log: /tmp/testacc.log"
fi
```

**Ask about test status:**

```text
AskUserQuestion(
  header: "Test Results",
  question: "What were the acceptance test results?",
  options: [
    { label: "All passed", description: "All tests passed successfully" },
    { label: "Partial", description: "Some tests passed, need to note exceptions" },
    { label: "Run now", description: "Run tests and capture output" }
  ]
)
```

If "Partial", ask for details about which tests and why.
If "Run now", execute the test command above and capture the output.

## Step 4: Gather API Verification

**Check for mitmproxy logs:**

```bash
if [ -f /tmp/api_debug.log ]; then
  API_LOG_EXISTS=true
  # Extract relevant API calls
  API_CALLS=$(grep -E "GET|POST|PUT|DELETE" /tmp/api_debug.log | grep "api2/json" | head -10)
fi
```

**Ask about API verification:**

```text
AskUserQuestion(
  header: "API Verification",
  question: "Was API behavior verified with mitmproxy?",
  options: [
    { label: "Yes, verified", description: "API calls confirmed correct" },
    { label: "Not applicable", description: "No API changes in this work" },
    { label: "Not done", description: "Need to run /debug-api first" }
  ]
)
```

If "Not done", suggest running `/debug-api` first.

## Step 5: Gather Implementation Summary

**Get changed files:**

```bash
CHANGED_FILES=$(git diff --name-only main...HEAD | grep -E "\.go$" | head -20)
```

**Get commit history:**

```bash
COMMITS=$(git log --oneline main..HEAD | head -10)
```

**Ask for summary:**

```text
AskUserQuestion(
  header: "Summary",
  question: "Briefly describe what was implemented/fixed:",
  options: [
    { label: "Bug fix", description: "Fixed an issue or defect" },
    { label: "New feature", description: "Added new functionality" },
    { label: "Enhancement", description: "Improved existing functionality" },
    { label: "Refactor", description: "Code restructuring without behavior change" }
  ]
)
```

After user selects the type, ask for additional details if needed. If session state exists, pull the summary from there.

## Step 6: Verify Checklist Items

**Ask if checks were already run:**

```text
AskUserQuestion(
  header: "Verification",
  question: "Did you already run /ready or verify build/lint/test?",
  options: [
    { label: "Yes, all passed", description: "Skip re-running checks" },
    { label: "Run now", description: "Run make build, lint, and test" }
  ]
)
```

If "Run now", verify each checklist item:

```bash
# Verify make build
make build && BUILD_STATUS="x" || BUILD_STATUS=" "

# Verify make lint
make lint && LINT_STATUS="x" || LINT_STATUS=" "

# Verify make test
make test && TEST_STATUS="x" || TEST_STATUS=" "
```

If "Yes, all passed", set all status variables to "x".

For acceptance tests and API verification, use the results gathered in previous steps.

**Ask about documentation:**

```text
AskUserQuestion(
  header: "Documentation",
  question: "Did you regenerate documentation (make docs)?",
  options: [
    { label: "Yes", description: "Docs regenerated after schema changes" },
    { label: "Not needed", description: "No schema changes in this work" },
    { label: "Not done", description: "Need to run make docs" }
  ]
)
```

## Step 7: Generate Report

Create the report file using the Write tool with the following template.

Replace placeholders with actual values gathered in previous steps:

- `{ISSUE_NUM}` - The issue number
- `{DATE}` - Current date (YYYY-MM-DD format)
- `{USER_SUMMARY}` - Summary from user/session state
- `{CHANGED_FILES_LIST}` - List of changed .go files
- `{COMMITS_LIST}` - Recent commits
- `{TEST_COMMAND_AND_OUTPUT}` - The full test command and terminal output from /tmp/testacc.log
- `{TEST_SUMMARY_TABLE}` - Table rows with: Test Name | Duration | Status | Description

**Example test command and output format:**

```text
$ ./testacc TestAccResourceVM.*Migration -- -v
=== RUN   TestAccResourceVMMigrationStopped
=== PAUSE TestAccResourceVMMigrationStopped
=== RUN   TestAccResourceVMMigrationRunning
=== PAUSE TestAccResourceVMMigrationRunning
=== CONT  TestAccResourceVMMigrationStopped
=== CONT  TestAccResourceVMMigrationRunning
--- PASS: TestAccResourceVMMigrationStopped (11.08s)
--- PASS: TestAccResourceVMMigrationRunning (24.02s)
PASS
ok   github.com/bpg/terraform-provider-proxmox/fwprovider 25.123s
```

**Example test summary table rows:**

```markdown
| TestAccResourceVMMigrationStopped | 11.08s | ✅ PASS | Stopped VM migration |
| TestAccResourceVMMigrationRunning | 24.02s | ✅ PASS | Running VM migration |
```

- `{API_VERIFICATION_SECTION}` - Description of API verification status
- `{API_CALLS_SNIPPET}` - Relevant API calls from /tmp/api_debug.log
- `{BUILD_STATUS}`, `{LINT_STATUS}`, etc. - "x" if passed, " " if not

**Report template:**

```markdown
# Issue #{ISSUE_NUM} - Proof of Work Report

**Date:** {DATE}
**GitHub Issue:** https://github.com/bpg/terraform-provider-proxmox/issues/{ISSUE_NUM}
**Author:** Claude Code

---

## Summary

{USER_SUMMARY}

## Files Changed

{CHANGED_FILES_LIST}

## Commits

{COMMITS_LIST}

## Test Results

\`\`\`
{TEST_COMMAND_AND_OUTPUT}
\`\`\`

### Test Coverage Summary

| Test Name | Duration | Status | Description |
|-----------|----------|--------|-------------|
{TEST_SUMMARY_TABLE}

## API Verification

{API_VERIFICATION_SECTION}

### Verified API Calls

\`\`\`
{API_CALLS_SNIPPET}
\`\`\`

## Checklist

- [{BUILD_STATUS}] `make build` passes
- [{LINT_STATUS}] `make lint` shows 0 issues
- [{TEST_STATUS}] `make test` passes
- [{ACC_TEST_STATUS}] Acceptance tests pass
- [{API_CHECK}] API calls verified with mitmproxy
- [{DOCS_CHECK}] Documentation regenerated (if schema changed)

---

*Generated by `/proof-report` skill. Use this content in the PR description.*
```

**Write the report:**

Use the Write tool to create the file at `.dev/{ISSUE_NUM}_PROOF_REPORT.md` with the populated template content.

## Step 8: Present Result

```markdown
=== PROOF REPORT GENERATED ===

Location: ${REPORT_PATH}

This file is gitignored and won't be committed.

Use the content in your PR description:
1. Copy relevant sections to the PR template
2. Include test output and API verification
3. Reference this report for detailed evidence

To view: cat ${REPORT_PATH}
```

**Offer to display:**

```text
AskUserQuestion(
  header: "View Report",
  question: "Would you like to see the generated report?",
  options: [
    { label: "Yes", description: "Display the full report" },
    { label: "No", description: "Done" }
  ]
)
```

If "Yes", display the report content.

</process>

<report_template>
The report follows this structure:

1. **Summary** — Brief description of what was done (from user input or session state)
2. **Files Changed** — List of modified files (compared to main branch)
3. **Commits** — Recent commit history (from branch point)
4. **Test Results** — Raw command + terminal output, then summary table (Name, Duration, Status, Description)
5. **API Verification** — Mitmproxy output showing correct parameters
6. **Checklist** — Production readiness confirmation (verified, not assumed)

This format aligns with the PR template requirements in [CONTRIBUTING.md](../../CONTRIBUTING.md).
</report_template>

<success_criteria>

- [ ] Issue number determined
- [ ] `.dev/` directory exists
- [ ] Session state checked for existing context
- [ ] Test results gathered (with log at /tmp/testacc.log)
- [ ] API verification status captured
- [ ] Implementation summary written
- [ ] Checklist items verified (not assumed)
- [ ] Report file created at `.dev/{ISSUE_NUM}_PROOF_REPORT.md`
- [ ] Report displayed or path provided to user
</success_criteria>

<tips>
- Run `/ready` first to ensure all checks pass before generating the report
- The report is gitignored, so you can include detailed output without worrying about committing it
- Copy the most relevant sections to the PR description
- Keep API verification snippets focused on the specific parameters being tested
- If session state exists (`.dev/{ISSUE_NUM}_SESSION_STATE.md`), pull context from there to save time
- Test output is saved to `/tmp/testacc.log` - record this location in session state for future reference
- The test log persists across agent sessions and can be reused if tests don't need re-running
</tips>
