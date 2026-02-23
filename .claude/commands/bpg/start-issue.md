---
name: start-issue
description: Start work on a GitHub issue with proper setup
argument-hint: \[issue-number\]
allowed-tools:
  - Read
  - Write
  - Bash
  - Grep
  - Glob
  - AskUserQuestion
  - WebFetch
---

<objective>
Set up the environment to start work on a GitHub issue:

1. Verify the issue exists and get context
2. Create a properly named branch
3. Create a session state file from template
4. Display issue summary for context

Use this skill when:

- Starting work on a new issue
- User says "work on issue #1234" or "fix #1234", or "investigate #1234"
- Beginning any fix or feature implementation
</objective>

<context>
Issue number: $ARGUMENTS

From [CLAUDE.md](../../CLAUDE.md): "All work on fixes or features MUST have a corresponding GitHub issue."
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

Validate it's a number:

```bash
ISSUE_NUM="${ARGUMENTS}"
if ! [[ "$ISSUE_NUM" =~ ^[0-9]+$ ]]; then
  echo "ERROR: Invalid issue number: $ISSUE_NUM"
  exit 1
fi
```

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

This ensures a clean slate for the new issue.

## Step 6: Display Summary

Present the setup summary to the user. Replace placeholders with actual values:

```text
=== ISSUE #{ISSUE_NUM} SETUP COMPLETE ===

Title: {TITLE}
Type: {ISSUE_TYPE}
Branch: {BRANCH_NAME}
Session: {SESSION_FILE}

Issue Summary:
{ISSUE_BODY_FIRST_PARAGRAPH}

Labels: {LABELS}

--- Next Steps ---
1. Review the full issue: gh issue view {ISSUE_NUM}
2. Update session state as you work
3. Create acceptance test BEFORE implementing fix
4. Run /bpg:ready before completing work
```

## Step 7: Update Session State

Update `.dev/${ISSUE_NUM}_SESSION_STATE.md` with:

- `Status:` → "In Progress"
- `Last Updated:` → current date
- `Current state:` → "Branch created, ready to implement"
- `Immediate next action:` → "Create acceptance test that reproduces the issue"

</process>

<success_criteria>

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

</success_criteria>

<tips>

- If `gh` CLI is not authenticated, the skill will prompt for manual verification
- Branch names are auto-truncated to avoid filesystem issues
- Session state file is gitignored, so it won't be committed
- **Acceptance tests first:** Always create the acceptance test BEFORE implementing the fix. The test must reproduce the bug with the current code and fail, proving the issue exists. Only then implement the fix and verify the test passes.
- Stale `/tmp/testacc.log` and `/tmp/api_debug.log` are cleared to ensure clean slate
- Use `/bpg:resume` to continue work if you need to come back later

</tips>
