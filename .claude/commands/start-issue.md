---
name: start-issue
description: Start work on a GitHub issue with proper setup
argument-hint: <issue-number>
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

```
AskUserQuestion(
  header: "Issue Number",
  question: "What GitHub issue number should we work on?",
  options: [
    { label: "I'll provide it", description: "Enter the issue number" }
  ]
)
```

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

```
"No GitHub issue #${ISSUE_NUM} found.

All fixes and features must be tracked with an issue before implementation begins.

Would you like me to help draft a GitHub issue?"
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
else
  # Ask user
  AskUserQuestion(
    header: "Issue Type",
    question: "Is this a bug fix or a new feature?",
    options: [
      { label: "Bug fix", description: "Use 'fix/' branch prefix" },
      { label: "Feature", description: "Use 'feat/' branch prefix" }
    ]
  )
fi
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
if git show-ref --verify --quiet "refs/heads/${BRANCH_NAME}"; then
  echo "Branch ${BRANCH_NAME} already exists."
  AskUserQuestion(
    header: "Branch Exists",
    question: "Branch already exists. What would you like to do?",
    options: [
      { label: "Switch to it", description: "Checkout existing branch" },
      { label: "Create new", description: "Use a different branch name" }
    ]
  )
else
  git checkout -b "${BRANCH_NAME}"
fi
```

## Step 5: Create Session State File

Copy template and populate:

```bash
SESSION_FILE=".dev/${ISSUE_NUM}_SESSION_STATE.md"
if [ -f "$SESSION_FILE" ]; then
  echo "Session state file already exists: $SESSION_FILE"
else
  cp .dev/SESSION_STATE_TEMPLATE.md "$SESSION_FILE"
  # Populate with issue details
fi
```

Update the session state with issue context:

- Replace `[NUMBER]` with actual issue number
- Replace `[Title]` with issue title
- Set status to "In Progress"
- Set branch name
- Fill "What this issue is about" from issue body

## Step 6: Display Summary

```
=== ISSUE #${ISSUE_NUM} SETUP COMPLETE ===

Title: ${TITLE}
Type: ${ISSUE_TYPE}
Branch: ${BRANCH_NAME}
Session: ${SESSION_FILE}

Issue Summary:
${ISSUE_BODY_FIRST_PARAGRAPH}

Labels: ${LABELS}

--- Next Steps ---
1. Review the full issue: gh issue view ${ISSUE_NUM}
2. Update session state as you work
3. Create acceptance test BEFORE implementing fix
4. Run /ready before completing work
```

</process>

<success_criteria>

- [ ] Issue number provided and validated
- [ ] Issue verified to exist on GitHub
- [ ] Branch created with correct naming: `{type}/{issue}-{description}`
- [ ] Session state file created: `.dev/{issue}_SESSION_STATE.md`
- [ ] Issue context displayed to user

</success_criteria>

<tips>
- If `gh` CLI is not authenticated, the skill will prompt for manual verification
- Branch names are auto-truncated to avoid filesystem issues
- Session state file is gitignored, so it won't be committed
- Always create the acceptance test before implementing the fix
</tips>
