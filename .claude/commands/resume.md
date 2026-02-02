---
name: resume
description: Resume work from a previous session
argument-hint: [issue-number]
allowed-tools:
  - Read
  - Bash
  - Grep
  - Glob
  - AskUserQuestion
---

<objective>
Resume work from a previous session by:

1. Listing available session state files
2. Loading the selected session's context
3. Displaying the "Quick Context Restore" section
4. Showing immediate next action

Use this skill when:

- Starting a new conversation to continue previous work
- User says "resume", "continue", or "pick up where I left off"
- Returning to work after a break
</objective>

<context>
Issue number (optional): $ARGUMENTS

Session files location: `.dev/*_SESSION_STATE.md`
</context>

<process>

## Step 1: Find Session State Files

If `$ARGUMENTS` provided, look for that specific issue:

```bash
SESSION_FILE=".dev/${ARGUMENTS}_SESSION_STATE.md"
if [ -f "$SESSION_FILE" ]; then
  # Found it, proceed to Step 3
else
  echo "No session file found for issue #${ARGUMENTS}"
  # Fall through to list all
fi
```

List all available session files:

```bash
ls -t .dev/*_SESSION_STATE.md 2>/dev/null | head -10
```

If no session files found:

```
No active session state files found.

To start work on an issue, use: /start-issue <issue-number>
```

## Step 2: Select Session

If multiple sessions exist, extract key info from each:

```bash
for file in .dev/*_SESSION_STATE.md; do
  # Extract issue number from filename
  ISSUE=$(basename "$file" | grep -oE '^[0-9]+')
  # Extract status
  STATUS=$(grep -m1 "Status:" "$file" | sed 's/.*Status:[[:space:]]*//')
  # Extract title
  TITLE=$(grep -m1 "^## Issue" "$file" | sed 's/.*- //')
  # Extract last updated
  UPDATED=$(grep -m1 "Last Updated:" "$file" | sed 's/.*Updated:[[:space:]]*//')
  echo "#${ISSUE} | ${STATUS} | ${UPDATED} | ${TITLE}"
done
```

Present options:

```
AskUserQuestion(
  header: "Select Session",
  question: "Which session would you like to resume?",
  options: [
    { label: "#1234", description: "In Progress - fix/1234-clone-timeout" },
    { label: "#5678", description: "Blocked - feat/5678-new-feature" }
  ]
)
```

## Step 3: Load Session Context

Read the selected session file:

```bash
cat ".dev/${ISSUE_NUM}_SESSION_STATE.md"
```

Extract key sections:

1. **Quick Context Restore** — Primary context for fast bootstrap
2. **Git State** — Current branch and uncommitted changes
3. **User Decisions** — Decisions already made
4. **Assumptions Made** — What's verified vs unverified
5. **Next Steps > Immediate** — What to do next

## Step 4: Verify Git State

Check current git state matches session:

```bash
CURRENT_BRANCH=$(git branch --show-current)
EXPECTED_BRANCH=$(grep "Current Branch:" "$SESSION_FILE" | sed 's/.*`\(.*\)`.*/\1/')

if [ "$CURRENT_BRANCH" != "$EXPECTED_BRANCH" ]; then
  echo "WARNING: Current branch ($CURRENT_BRANCH) differs from session ($EXPECTED_BRANCH)"
  AskUserQuestion(
    header: "Branch Mismatch",
    question: "Switch to the session's branch?",
    options: [
      { label: "Yes", description: "Checkout $EXPECTED_BRANCH" },
      { label: "No", description: "Stay on current branch" }
    ]
  )
fi
```

Check for uncommitted changes:

```bash
git status --porcelain
```

## Step 5: Display Context

```
=== RESUMING ISSUE #${ISSUE_NUM} ===

${QUICK_CONTEXT_RESTORE_SECTION}

--- Git State ---
Branch: ${CURRENT_BRANCH}
Uncommitted: ${UNCOMMITTED_COUNT} files

--- User Decisions (already made) ---
${USER_DECISIONS_TABLE}

--- Unverified Assumptions ---
${UNVERIFIED_ASSUMPTIONS}

--- Immediate Next Action ---
${IMMEDIATE_NEXT_ACTION}

Session file: .dev/${ISSUE_NUM}_SESSION_STATE.md
```

## Step 6: Offer Actions

```
AskUserQuestion(
  header: "Ready",
  question: "How would you like to proceed?",
  options: [
    { label: "Continue work", description: "Start on the immediate next action" },
    { label: "View full session", description: "Display complete session state" },
    { label: "Check issue", description: "View GitHub issue for updates" }
  ]
)
```

If "View full session": Display the entire session state file.
If "Check issue": Run `gh issue view ${ISSUE_NUM}`.

</process>

<success_criteria>

- [ ] Session state files listed (or specific one found)
- [ ] Session selected (if multiple)
- [ ] Context loaded and displayed
- [ ] Git state verified
- [ ] Immediate next action shown
- [ ] Ready to continue work

</success_criteria>

<tips>
- Session files are sorted by modification time (most recent first)
- Always verify git state matches session before continuing
- If assumptions are marked "Unverified", prioritize verifying them
- Update session state after completing the immediate next action
</tips>
