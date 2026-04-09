---
name: done
description: Wrap up a session — extract learnings to memory, finalize session state, archive .dev files
argument-hint: \[issue-number\]
allowed-tools:
  - Read
  - Edit
  - Write
  - Bash
  - Grep
  - Glob
  - AskUserQuestion
---

<objective>
Wrap up work on an issue by:

1. Reviewing the session for learnings worth preserving
2. Writing new or updating existing memory files
3. Finalizing the session state file
4. Archiving all issue-related files from `.dev/` to `.dev/archive/`

Use this skill when:

- Work on an issue is complete (PR created or ready to create)
- User says "done", "wrap up", "finish up", or "archive this"
- Before switching to a different issue
</objective>

<context>
Issue number (optional): $ARGUMENTS

Memory location: `~/.claude/projects/-Users-pasha-code-terraform-provider-proxmox/memory/`
Memory index: `~/.claude/projects/-Users-pasha-code-terraform-provider-proxmox/memory/MEMORY.md`
Session files: `.dev/{issue}_*`
Archive: `.dev/archive/`
</context>

<process>

## Step 0: Determine Issue Number

If `$ARGUMENTS` provided, use it. Otherwise detect from branch name:

```bash
ISSUE_NUM=$(git branch --show-current | grep -oE '(fix|feat)/[0-9]+' | grep -oE '[0-9]+')
```

If still unclear, ask the user.

## Step 1: Extract Learnings

Review the conversation history and identify learnings that would help future sessions on this project. Look for:

- **Patterns discovered** — code patterns, API behaviors, Terraform Framework idioms that weren't obvious
- **Traps and gotchas** — things that looked right but were wrong, subtle bugs, incorrect assumptions
- **Infrastructure knowledge** — test environment setup, PVE host configuration, tooling quirks
- **Process improvements** — workflow steps that worked well or poorly, better approaches discovered mid-session
- **Codebase conventions** — implicit rules or conventions found during review that aren't documented in CLAUDE.md

For each learning, decide:

- Does it fit in an **existing** memory file? → Update that file
- Is it a new topic? → Create a new memory file
- Is it too specific to this issue? → Skip (session state captures it already)

**Quality bar:** Only persist learnings that would save time or prevent mistakes in future sessions. Don't persist obvious things or one-off details.

### Read existing memory files

Read the memory index to understand what's already captured:

```bash
cat ~/.claude/projects/-Users-pasha-code-terraform-provider-proxmox/memory/MEMORY.md
```

Then read any memory files that might overlap with the session's learnings to avoid duplication.

### Write or update memory files

For new files:
- Use descriptive kebab-case names
- Keep focused — one topic per file
- Write for a future agent who has no context

For updates to existing files:
- Add new sections or bullet points
- Don't rewrite what's already there unless it's wrong

After writing/updating, update `MEMORY.md` index if new files were created.

### Present learnings to user

Show the user what you're capturing and ask for confirmation:

```text
I identified these learnings from this session:

1. **[Topic]** — [Brief description] → [New file / Update to existing file]
2. **[Topic]** — [Brief description] → [New file / Update to existing file]

Should I save these? Any corrections or additions?
```

Wait for user confirmation before writing.

## Step 2: Finalize Session State

Read the session state file `.dev/{issue}_SESSION_STATE.md` and update it with final details:

- `Last Updated:` → current date
- `Status:` → "Complete" (or "PR Created" / "Merged" if applicable)
- `Current state:` → final summary of what was accomplished
- `Immediate next action:` → "Archived. No further action needed." (or link to open PR if applicable)
- Ensure the **What Was Done** section is complete and accurate
- Ensure **Verification Results** reflect the final test run
- Add a **Session Log** entry with date and summary

## Step 3: Archive Files

Move all issue-related files from `.dev/` to `.dev/archive/`:

```bash
ISSUE_NUM=<detected>
mkdir -p .dev/archive

# Find all files for this issue
for f in .dev/${ISSUE_NUM}_*; do
  [ -f "$f" ] && mv "$f" .dev/archive/
done
```

Verify the move:

```bash
ls .dev/${ISSUE_NUM}_* 2>/dev/null && echo "WARNING: files still in .dev/" || echo "All files archived"
ls .dev/archive/${ISSUE_NUM}_*
```

## Step 4: Switch to Main Branch

Return to the main branch so the workspace is clean for the next task:

```bash
git checkout main
```

## Step 5: Summary

Present a brief summary:

```text
Session wrapped up for issue #${ISSUE_NUM}:

Memory:
  - [Created/Updated]: [file] — [description]
  - [Created/Updated]: [file] — [description]

Archived:
  - .dev/archive/${ISSUE_NUM}_SESSION_STATE.md
  - .dev/archive/${ISSUE_NUM}_PR_BODY.md
  - ...

Switched to: main
```

</process>

<success_criteria>

- [ ] Learnings identified and presented to user for confirmation
- [ ] Memory files created or updated (with index updated)
- [ ] Session state finalized with complete summary
- [ ] All `.dev/{issue}_*` files moved to `.dev/archive/`
- [ ] Switched to `main` branch
- [ ] Summary presented to user
</success_criteria>

<tips>
- Don't over-persist. If a learning is already in CLAUDE.md or an existing memory, skip it.
- Session state should be self-contained — a future agent reading only that file should understand what happened.
- The archive keeps files accessible but out of the way. They're gitignored so they never get committed.
- If the user hasn't created a PR yet, mention it in the summary as a remaining action.
</tips>
