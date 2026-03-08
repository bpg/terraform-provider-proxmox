---
name: start-issue
description: Use when starting work on a GitHub issue — sets up branch, session state, and displays context. Also use when user says "work on issue", "fix #1234", or "investigate #1234".
argument-hint: \[issue-number\]
allowed-tools:
  - Read
  - Edit
  - Write
  - Bash
  - Grep
  - Glob
  - Agent
  - AskUserQuestion
  - WebFetch
---

<objective>
Set up the environment to start work on a GitHub issue, then guide the agent through investigation and TDD.

Phase 1 (setup — executed immediately):

1. Verify the issue exists and get context
2. Create a properly named branch
3. Create a session state file from template
4. Display issue summary and update session state

Phase 2 (after setup — guidance for next steps):

5. Investigate the issue before writing any code
6. Fix with TDD (acceptance tests first)
</objective>

<context>
Issue number: $ARGUMENTS

From [CLAUDE.md](../../../CLAUDE.md): "All work on fixes or features MUST have a corresponding GitHub issue."
</context>

<process>

## Step 1: Get Issue Number

If `$ARGUMENTS` provided, use it. Otherwise ask:

```text
AskUserQuestion(
  header: "Issue Number",
  question: "What GitHub issue number should we work on?",
  options: [
    { label: "Enter number", description: "I'll type the issue number" },
    { label: "Browse issues", description: "Help me find an issue to work on" }
  ]
)
```

If "Browse issues", run `gh issue list --limit 10` to show recent issues.

Validate the input is a number. If not numeric, ask the user to provide a valid issue number.

## Step 2: Verify Issue Exists

Fetch issue details from GitHub:

```bash
gh issue view "$ISSUE_NUM" --json title,body,labels,state 2>/dev/null
```

If issue doesn't exist or gh fails:

```text
No GitHub issue #{ISSUE_NUM} found.

All fixes and features must be tracked with an issue before implementation begins.

Would you like me to help draft a GitHub issue?
```

```text
AskUserQuestion(
  header: "Create Issue",
  question: "Would you like help creating a GitHub issue?",
  options: [
    { label: "Yes, bug report", description: "Draft a bug report" },
    { label: "Yes, feature request", description: "Draft a feature request" },
    { label: "No", description: "I'll create it manually" }
  ]
)
```

If user wants help, determine type and draft using templates from:

- Bug: `.github/ISSUE_TEMPLATE/bug_report.md`
- Feature: `.github/ISSUE_TEMPLATE/feature_request.md`

## Step 3: Determine Issue Type

From labels or title, determine if this is a bug fix or feature:

```bash
LABELS=$(gh issue view "$ISSUE_NUM" --json labels -q '.labels[].name' 2>/dev/null)
if echo "$LABELS" | grep -qi "bug"; then
  ISSUE_TYPE="fix"
elif echo "$LABELS" | grep -qi "enhancement\|feature"; then
  ISSUE_TYPE="feat"
fi
```

If type cannot be determined from labels, ask user:

```text
AskUserQuestion(
  header: "Issue Type",
  question: "Is this a bug fix or a new feature?",
  options: [
    { label: "Bug fix", description: "Use 'fix/' branch prefix" },
    { label: "Feature", description: "Use 'feat/' branch prefix" }
  ]
)
```

## Step 4: Create Branch

Generate branch name from issue title:

```bash
TITLE=$(gh issue view "$ISSUE_NUM" --json title -q '.title' 2>/dev/null)
# Normalize: lowercase, replace spaces/special chars with hyphens, truncate
SHORT_DESC=$(echo "$TITLE" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | sed 's/--*/-/g' | cut -c1-40 | sed 's/-$//')
BRANCH_NAME="${ISSUE_TYPE}/${ISSUE_NUM}-${SHORT_DESC}"
```

Check if branch exists:

```bash
git show-ref --verify --quiet "refs/heads/${BRANCH_NAME}" && echo "Branch exists" || echo "Branch available"
```

If branch already exists, ask user:

```text
AskUserQuestion(
  header: "Branch Exists",
  question: "Branch {BRANCH_NAME} already exists. What would you like to do?",
  options: [
    { label: "Switch to it", description: "Checkout existing branch" },
    { label: "Create new", description: "Use a different branch name" }
  ]
)
```

If "Switch to it": `git checkout "${BRANCH_NAME}"`
If "Create new": Ask for new name or append suffix.

If branch doesn't exist, create it:

```bash
git checkout -b "${BRANCH_NAME}"
```

## Step 5: Create Session State File

Ensure `.dev/` directory exists and create session file:

```bash
mkdir -p .dev
SESSION_FILE=".dev/${ISSUE_NUM}_SESSION_STATE.md"
```

Verify template exists:

```bash
if [ ! -f ".dev/SESSION_STATE_TEMPLATE.md" ]; then
  echo "ERROR: Template not found at .dev/SESSION_STATE_TEMPLATE.md"
fi
```

If template is missing, stop and alert user. The template should exist in the repository.

Check if session file already exists:

```bash
if [ -f "$SESSION_FILE" ]; then
  echo "Session state file already exists: $SESSION_FILE"
else
  cp .dev/SESSION_STATE_TEMPLATE.md "$SESSION_FILE"
fi
```

If file already exists, ask user:

```text
AskUserQuestion(
  header: "Session Exists",
  question: "Session state file already exists. What would you like to do?",
  options: [
    { label: "Use existing", description: "Continue with existing session state" },
    { label: "Overwrite", description: "Start fresh with new session state" }
  ]
)
```

**Populate the session state** using the Read and Edit tools:

1. Read the session file
2. Replace `[NUMBER]` with `{ISSUE_NUM}`
3. Replace `[Title]` with issue title from GitHub
4. Set `Status:` to "In Progress"
5. Set `Current Branch:` to `{BRANCH_NAME}`
6. Fill "What this issue is about" from issue body (first paragraph)
7. Set `Last Updated:` to current date

**Clear any stale log files from previous work:**

```bash
rm -f /tmp/testacc.log /tmp/api_debug.log
```

## Step 6: Display Summary and Complete Setup

Present the setup summary to the user:

```text
=== ISSUE #{ISSUE_NUM} SETUP COMPLETE ===

Title: {TITLE}
Type: {ISSUE_TYPE}
Branch: {BRANCH_NAME}
Session: {SESSION_FILE}

Issue Summary:
{ISSUE_BODY_FIRST_PARAGRAPH}

Labels: {LABELS}
```

Update `.dev/${ISSUE_NUM}_SESSION_STATE.md` with:

- `Status:` → "In Progress"
- `Last Updated:` → current date
- `Current state:` → "Setup complete, ready to investigate"
- `Immediate next action:` → "Investigate the issue — explore relevant code and identify root cause"

**Then ask the user whether to continue with investigation or wait for further instructions.**

</process>

---

## After Setup: Investigation and TDD

**These steps are mandatory before any fix is implemented. Do NOT skip investigation.**

### Investigate First

**Do NOT jump to writing code or tests.** Investigate first, then discuss with the user.

Invoke the `/superpowers:systematic-debugging` skill to guide the investigation. Then:

1. **Explore the relevant code** — Use Serena's `get_symbols_overview` to understand file structure, then `find_symbol` to drill into specific functions. Use `find_referencing_symbols` to trace call chains. Fall back to Grep/Glob for cross-file pattern searches.
2. **Look up API docs if needed** — Use Context7 with `/websites/pve_proxmox_pve-docs` to check Proxmox API endpoint parameters and behavior. Use `/hashicorp/terraform-plugin-framework` for Framework API questions.
3. **Identify the root cause** — Form a hypothesis about what's wrong and why. Trace the code path from the user's reported behavior to the underlying bug.
4. **Check for related patterns** — Look for similar attributes/resources that may have the same issue or that already handle the case correctly.
5. **Present findings to the user** — Summarize root cause, proposed fix, and open questions.
6. **Ask the user** if you should continue with the fix or discuss further.

Update session state with investigation findings before proceeding.

### Fix with TDD

Only proceed after investigation is complete and the user confirms.

Follow TDD (Red-Green-Refactor):

1. **RED — Write a failing acceptance test first**
   - Create an acceptance test that reproduces the bug
   - Run it with `./testacc` and **verify it fails for the expected reason**
   - If the test fails for a different reason (e.g. connection issues, missing infrastructure), **ask the user** — do NOT work around it
   - If the bug cannot be reproduced with acceptance tests, ask the user how to proceed
2. **GREEN — Implement the minimal fix**
   - Write the simplest code that makes the failing test pass
   - Run the test again and verify it passes
3. **Verify — No regressions**
   - Run related existing acceptance tests to confirm no regressions
   - Run `make lint`
4. Run `/bpg:ready` before completing work

### What NOT to Do

- Do NOT create unit tests, extract helper functions, or refactor production code as workarounds when acceptance tests can't reproduce the bug. Ask the user instead.
- Do NOT skip investigation and jump straight to writing code or tests.
- Do NOT work around infrastructure problems (Proxmox unreachable, missing services). Ask the user.

<success_criteria>

Setup phase (executed by this skill):

- [ ] Issue number provided and validated
- [ ] Issue verified to exist on GitHub
- [ ] Issue type determined (fix/feat)
- [ ] Branch created with correct naming: `{type}/{issue}-{description}`
- [ ] `.dev/` directory exists
- [ ] Session state template exists
- [ ] Session state file created: `.dev/{issue}_SESSION_STATE.md`
- [ ] Session state populated with issue context
- [ ] Stale log files cleared
- [ ] Issue context displayed to user
- [ ] Session state updated with current status
- [ ] User asked whether to continue or wait

</success_criteria>

<tips>

- If `gh` CLI is not authenticated, the skill will prompt for manual verification
- Branch names are auto-truncated to avoid filesystem issues
- Session state file is gitignored, so it won't be committed
- Stale `/tmp/testacc.log` and `/tmp/api_debug.log` are cleared to ensure clean slate
- Use `/bpg:resume` to continue work if you need to come back later

</tips>
