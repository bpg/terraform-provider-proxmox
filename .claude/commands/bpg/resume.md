---
name: resume
description: Resume work from a previous session
argument-hint: \[issue-number\]
allowed-tools:
  - Read
  - Edit
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

**Check `.dev/` directory exists:**

```bash
if [ ! -d ".dev" ]; then
  echo "No .dev/ directory found. No sessions to resume."
  echo "To start work on an issue, use: /bpg:start-issue <issue-number>"
  exit 0
fi
```

If `$ARGUMENTS` provided, look for that specific issue:

```bash
SESSION_FILE=".dev/${ARGUMENTS}_SESSION_STATE.md"
if [ -f "$SESSION_FILE" ]; then
  ISSUE_NUM="${ARGUMENTS}"
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

```text
No active session state files found.

To start work on an issue, use: /bpg:start-issue <issue-number>
```

## Step 2: Select Session

If multiple sessions exist, extract key info from each:

```bash
for file in .dev/*_SESSION_STATE.md; do
  # Extract issue number from filename (format: {issue}_SESSION_STATE.md)
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

Present options (dynamically built from found sessions):

```text
AskUserQuestion(
  header: "Select Session",
  question: "Which session would you like to resume?",
  options: [
    { label: "#{ISSUE_1}", description: "{STATUS_1} - {TITLE_1}" },
    { label: "#{ISSUE_2}", description: "{STATUS_2} - {TITLE_2}" }
  ]
)
```

Replace `{ISSUE_N}`, `{STATUS_N}`, `{TITLE_N}` with actual values from the session files.

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
echo "Current: $CURRENT_BRANCH, Expected: $EXPECTED_BRANCH"
```

If branches don't match, ask user:

```text
AskUserQuestion(
  header: "Branch Mismatch",
  question: "Current branch ({CURRENT_BRANCH}) differs from session ({EXPECTED_BRANCH}). Switch?",
  options: [
    { label: "Yes", description: "Checkout the session's branch" },
    { label: "No", description: "Stay on current branch" }
  ]
)
```

If "Yes", checkout the expected branch:

```bash
git checkout "$EXPECTED_BRANCH"
```

Check for uncommitted changes:

```bash
git status --porcelain
```

## Step 5: Display Context

Present the loaded context to the user. Replace placeholders with actual values from the session file:

```text
=== RESUMING ISSUE #{ISSUE_NUM} ===

{QUICK_CONTEXT_RESTORE_SECTION}

--- Git State ---
Branch: {CURRENT_BRANCH}
Uncommitted: {UNCOMMITTED_COUNT} files

--- User Decisions (already made) ---
{USER_DECISIONS_TABLE}

--- Unverified Assumptions ---
{UNVERIFIED_ASSUMPTIONS}

--- Immediate Next Action ---
{IMMEDIATE_NEXT_ACTION}

Session file: .dev/{ISSUE_NUM}_SESSION_STATE.md
```

**Check for existing log files from previous runs:**

```bash
[ -f /tmp/testacc.log ] && echo "Found test log: /tmp/testacc.log"
[ -f /tmp/api_debug.log ] && echo "Found API debug log: /tmp/api_debug.log"
```

These logs may contain useful output from previous `/bpg:ready` or `/bpg:debug-api` runs.

## Step 6: Offer Actions

```text
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

If "Continue work": Begin working on the immediate next action from the session state.
If "View full session": Display the entire session state file.
If "Check issue": Run `gh issue view {ISSUE_NUM}` to see any updates.

## Step 7: Update Session State

Update `.dev/${ISSUE_NUM}_SESSION_STATE.md` using Read and Edit tools:

- `Last Updated:` → current date
- Note session was resumed (append to context or current state)

</process>

<success_criteria>

- [ ] `.dev/` directory exists
- [ ] Session state files listed (or specific one found)
- [ ] Session selected (if multiple)
- [ ] Context loaded and displayed
- [ ] Git state verified (branch matches)
- [ ] Existing log files noted
- [ ] Immediate next action shown
- [ ] Ready to continue work
- [ ] Session state updated with resume timestamp

</success_criteria>

<tips>

- Session files are sorted by modification time (most recent first)
- Always verify git state matches session before continuing
- If assumptions are marked "Unverified", prioritize verifying them
- Update session state after completing the immediate next action
- Check for existing `/tmp/testacc.log` and `/tmp/api_debug.log` from previous runs
- These logs can save time if tests don't need to be re-run

</tips>
