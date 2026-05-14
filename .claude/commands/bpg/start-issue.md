---
name: start-issue
description: Use when starting work on a GitHub issue — sets up branch (or worktree), session state, and displays context. Hands off to `/bpg:investigate` for the investigation gate. Also use when user says "work on issue" or "fix #1234".
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
  - Skill
---

<objective>
Set up the environment to start work on a GitHub issue. Phase 1 of the issue workflow — the
investigation gate lives in `/bpg:investigate` and is invoked by the user after this skill
finishes.

Setup 1–6 (executed by this slash command):

1. Get and validate issue number
2. Verify the issue exists on GitHub; auto-create and continue if not
3. Determine issue type (fix/feat) from labels
4. Choose workspace (branch-in-place or `.claude/worktrees/`), pick branch base, create branch
5. Create a session state file from template
6. Display issue summary, update session state, and prompt the user to run `/bpg:investigate`

Phase 2 — investigation gate — is a separate skill: see `/bpg:investigate <ISSUE_NUM>`. That
skill drives the maintainer-talk → root-cause → pattern → hypothesis → approval flow before
TDD begins.

</objective>

<context>
Issue number: $ARGUMENTS

From [CLAUDE.md](../../../CLAUDE.md): "All work on fixes or features MUST have a corresponding GitHub issue."
</context>

<process>

### Convention: three kinds of variable

Each `Bash` tool call runs in a fresh shell — environment variables do NOT survive across
calls. This skill uses three notations so the requirement is unambiguous:

| Form      | Meaning                                                                                                                                  | Example use                                               |
| --------- | ---------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------- |
| `<UPPER>` | **Cross-block agent state.** The agent substitutes the entire token with a value carried from a prior step.                              | `git checkout -b "<BRANCH_NAME>" <BASE_REF>`              |
| `<lower>` | **User-judgment fill-in.** The agent picks a value based on context (search term, file path) — not from prior-step state.                | `gh pr list --search "is:merged <symptom>"`               |
| `$NAME`   | **Shell variable, local to a single bash block.** Defined and consumed within the same `Bash` call; the agent does NOT substitute these. | `ISSUE_JSON=$(gh issue view …)` then `echo "$ISSUE_JSON"` |

Concretely:

- **Cross-block (`<UPPER>`):** `<ISSUE_NUM>`, `<TITLE>`, `<BODY>`, `<LABELS>`, `<COMMENTS>`,
  `<ISSUE_TYPE>`, `<SHORT_DESC>`, `<BRANCH_NAME>`, `<BASE_REF>`, `<DRAFT_TITLE>`, `<DRAFT_BODY>`,
  `<DRAFT_LABEL>`.
- **User-judgment (`<lower>`):** `<symptom>`, `<keyword>`, `<branch>` — placeholders inside
  shell command examples that the agent fills based on what's being investigated.
- **Bash locals (`$NAME`):** `$ISSUE_JSON`, `$STATE`, `$ISSUE_URL`, `$AVAILABLE_LABELS`,
  `$TARGET_REPO`, `$DIRTY` exist for the duration of one `Bash` call.

Substitution mechanics: tools that validate input (git, gh) reject literal `<>` characters,
so a missed substitution generally surfaces as an error. Some tools (curl, printf, raw text
operations) accept the literal — be deliberate about substitution rather than relying on the
shell to fail. Quote substituted values as needed: `git commit -m "<TITLE>"` becomes
`git commit -m "the title"` after substitution.

## Setup 1: Get Issue Number

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

## Setup 2: Verify Issue Exists and Fetch Context

Fetch all needed issue fields in a single call so downstream setup steps reuse the same data:

```bash
ISSUE_JSON=$(gh issue view "<ISSUE_NUM>" --json title,body,labels,state,comments 2>/dev/null)
if [ -z "$ISSUE_JSON" ]; then
  ISSUE_EXISTS=false
else
  ISSUE_EXISTS=true
  TITLE=$(echo "$ISSUE_JSON" | jq -r '.title')
  BODY=$(echo "$ISSUE_JSON" | jq -r '.body')
  STATE=$(echo "$ISSUE_JSON" | jq -r '.state')
  LABELS=$(echo "$ISSUE_JSON" | jq -r '.labels[].name')
  COMMENTS=$(echo "$ISSUE_JSON" | jq -r '.comments[] | "--- @\(.author.login):\n\(.body)"')
fi
```

After the block runs, the agent captures the `TITLE`, `BODY`, `LABELS`, `COMMENTS`, `STATE`
shell values from the output and uses them as `<TITLE>`, `<BODY>`, `<LABELS>`, `<COMMENTS>` in
later steps. Do NOT re-run `gh issue view`.

### Closed-issue check

If `STATE == "CLOSED"`, warn the user before proceeding — re-opening work on a closed issue is
sometimes intentional (related fix, follow-up) and sometimes a mistake (wrong number):

```text
AskUserQuestion(
  header: "Closed issue",
  question: "Issue #<ISSUE_NUM> is CLOSED. Continue anyway?",
  options: [
    { label: "Continue",       description: "Intentional — follow-up work or related fix" },
    { label: "Pick a different issue", description: "Wrong number — back to Setup 1" },
    { label: "Stop",           description: "Cancel the skill" }
  ]
)
```

If "Pick a different issue", restart at Setup 1. If "Stop", exit cleanly. If "Continue",
proceed and record the rationale (one line) in session state under `Closed-issue rationale:`.

### Issue does not exist — draft, create, and continue

If `ISSUE_EXISTS=false`:

```text
No GitHub issue #{ISSUE_NUM} found.

All fixes and features must be tracked with an issue before implementation begins.
```

```text
AskUserQuestion(
  header: "Create Issue",
  question: "Would you like help creating a GitHub issue?",
  options: [
    { label: "Yes, bug report", description: "Draft a bug report and file it" },
    { label: "Yes, feature request", description: "Draft a feature request and file it" },
    { label: "No", description: "Stop — I'll create it manually" }
  ]
)
```

If "No": stop. Do not proceed without an issue.

If "Yes": determine type and draft using these templates as a structural guide (NOT verbatim — fill
real content):

- Bug: `.github/ISSUE_TEMPLATE/bug_report.md`
- Feature: `.github/ISSUE_TEMPLATE/feature_request.md`

Then gather missing details from the user via `AskUserQuestion` (one round of focused questions —
symptom, expected behavior, repro, environment) and assemble the body.

**Confirm the target repo before filing.** `gh issue create` files against the repo `gh` resolves
from the current directory. If the user is on a fork, that's almost never what they want:

```bash
TARGET_REPO=$(gh repo view --json nameWithOwner -q .nameWithOwner)
echo "Issue will be filed against: $TARGET_REPO"
```

Show this to the user and confirm. If `$TARGET_REPO` is a fork (not `bpg/terraform-provider-proxmox`),
ask whether to file upstream instead (`gh issue create --repo bpg/terraform-provider-proxmox …`).

**Constrain labels to ones that exist.** Free-text labels fail with `gh issue create`:

```bash
AVAILABLE_LABELS=$(gh label list --limit 100 --json name -q '.[].name')
```

Pick `<DRAFT_LABEL>` only from `$AVAILABLE_LABELS`. If the natural label doesn't exist, drop it
or pick the closest match (e.g., `bug` exists, `regression` may not).

Show the draft to the user and ask for confirmation:

```text
AskUserQuestion(
  header: "File Issue",
  question: "Ready to file this issue on GitHub?",
  options: [
    { label: "File it",    description: "Create the issue with this draft and continue setup" },
    { label: "Edit first", description: "Iterate on the draft before filing" },
    { label: "Cancel",     description: "Don't file — stop the skill" }
  ]
)
```

If "File it":

```bash
ISSUE_URL=$(gh issue create \
  --title "<DRAFT_TITLE>" \
  --body "<DRAFT_BODY>" \
  --label "<DRAFT_LABEL>")
ISSUE_NUM=$(basename "$ISSUE_URL")
```

(Note: `$ISSUE_URL` here is a bash local within this same block, so use `$ISSUE_URL`, not the
placeholder form. The agent captures the resulting `ISSUE_NUM` value and uses it as
`<ISSUE_NUM>` in subsequent steps.)

Then re-enter Setup 2 with the new `<ISSUE_NUM>`. Reuse the `<DRAFT_TITLE>`/`<DRAFT_BODY>`
already in the agent's context — there's no need to re-fetch via `gh issue view` since you just
wrote them. Set `<LABELS>` from `<DRAFT_LABEL>`. Then continue to Setup 3.

If "Edit first": prompt the user for the changes inline ("What would you like to change?"),
update the draft in the agent's context, and re-display. Repeat until the user picks "File it"
or "Cancel". Do NOT exit the skill — iteration happens here, not in another invocation.

If "Cancel": stop the skill cleanly.

## Setup 3: Determine Issue Type

Reuse `LABELS` from Setup 2 (do NOT re-run `gh issue view`):

```bash
if echo "<LABELS>" | grep -qi "bug"; then
  ISSUE_TYPE="fix"
elif echo "<LABELS>" | grep -qi "enhancement\|feature"; then
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

## Setup 4: Create Branch or Worktree

Reuse `TITLE` from Setup 2 (do NOT re-run `gh issue view`). Normalize for use in a branch name:

```bash
SHORT_DESC=$(echo "<TITLE>" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | sed 's/--*/-/g' | cut -c1-40 | sed 's/-$//')
BRANCH_NAME="<ISSUE_TYPE>/<ISSUE_NUM>-<SHORT_DESC>"
```

### Determine branch base

`origin/main` is the default for ~99% of work, but backports to a release branch or hotfixes
against a tag need a different base. Ask before choosing the workspace so both paths use the
same `<BASE_REF>`:

```bash
git fetch origin  # ensure remote refs are current before the user picks
```

```text
AskUserQuestion(
  header: "Branch base",
  question: "Base the new branch on which ref?",
  options: [
    { label: "origin/main",    description: "Default for new fixes and features" },
    { label: "Different base", description: "Backport, hotfix, or stacked branch — I'll specify the ref next" }
  ]
)
```

If "origin/main" is chosen, set `<BASE_REF>` = `origin/main` and skip to the workspace choice.

If "Different base", prompt the user in plain prose for the exact ref (e.g.
`origin/release-0.105`, `v0.105.0`, `feature/foo`). Then verify it resolves before continuing:

```bash
if ! git rev-parse --verify "<BASE_REF>" >/dev/null 2>&1; then
  echo "ERROR: ref '<BASE_REF>' does not resolve. Did you mean origin/<branch> for a remote?"
  exit 1
fi
```

If the verify fails, re-prompt the user for the ref and try again — do NOT proceed with an
unresolved ref. Common forms: `origin/<branch>` for remote branches, plain `<branch>` for local
branches, `v<x.y.z>` for tags. Loop until the user provides a ref that resolves or chooses to
abort.

### Choose: branch in place vs isolated worktree

Worktrees are useful when multiple PRs are in flight or when the current checkout has
uncommitted work.

```text
AskUserQuestion(
  header: "Workspace",
  question: "Where should the work happen?",
  options: [
    { label: "Branch here",       description: "Create a branch in this checkout" },
    { label: "Isolated worktree", description: "Create a git worktree under .claude/worktrees/<ISSUE_NUM>-<SHORT_DESC>" }
  ]
)
```

If "Isolated worktree": invoke `superpowers:using-git-worktrees` via the `Skill` tool and let it
handle the worktree mechanics. Pass `<BRANCH_NAME>`, `<BASE_REF>`, and target path
`.claude/worktrees/<ISSUE_NUM>-<SHORT_DESC>`. After the worktree is created, all subsequent
commands in this skill must operate inside the worktree directory (use absolute paths or
`cd $WORKTREE_PATH && …` per Bash call — shell state does not persist across calls). Skip the
rest of this section.

**Worktree upstream check still required.** After the worktree skill creates the branch,
verify `git -C "$WORKTREE_PATH" branch -vv | grep "^\* "` does NOT show `[origin/main]`.
If it does, run `git -C "$WORKTREE_PATH" branch --unset-upstream` before any push attempt.
The same `feedback-never-push-to-main.md` failure mode applies regardless of how the branch
was created.

### Branch-in-place path

Pre-flight: refuse to clobber any local work (modified, staged, OR untracked):

```bash
DIRTY=$(git status --porcelain)
if [ -n "$DIRTY" ]; then
  echo "ERROR: Working tree is not clean. Stash, commit, or use the worktree option."
  echo "$DIRTY"
  exit 1
fi
```

Check if branch exists:

```bash
git show-ref --verify --quiet "refs/heads/<BRANCH_NAME>" && echo "Branch exists" || echo "Branch available"
```

If branch already exists, ask user:

```text
AskUserQuestion(
  header: "Branch Exists",
  question: "Branch {BRANCH_NAME} already exists. What would you like to do?",
  options: [
    { label: "Switch to it", description: "Checkout existing branch (no rebase)" },
    { label: "Create new", description: "Use a different branch name" }
  ]
)
```

- "Switch to it": `git checkout "<BRANCH_NAME>"` — then verify upstream (see below).
- "Create new": ask for a new name or append a numeric suffix.

If the branch doesn't exist, create it from the chosen base:

```bash
git checkout -b "<BRANCH_NAME>" --no-track <BASE_REF>
```

**`--no-track` is non-negotiable.** Without it, `git checkout -b NEW_BRANCH origin/main`
implicitly sets the new branch's upstream to `origin/main`. A subsequent bare `git push`,
`git push origin HEAD`, or `git push -u origin HEAD` can then resolve the remote ref via
that tracking config and land the feature branch's commits **directly on `main`**. This has
happened in this repo (post-mortem in `feedback-never-push-to-main.md`). `--no-track`
ensures the branch starts with no upstream; the first push must explicitly name the target
remote ref via `-u origin "<BRANCH_NAME>:<BRANCH_NAME>"`, which is the safe form.

**Verify immediately after creating the branch:**

```bash
git branch -vv | grep "^\* "
```

The output line for the current branch must NOT contain `[origin/main]`. If it does,
something failed — STOP, run `git branch --unset-upstream`, and notify the user before
proceeding.

**Switch-to-it path:** if the user picked "Switch to it" on an existing branch above, run
the same verify step. If the existing branch was created without `--no-track` previously,
it may still have `origin/main` as upstream — unset it before continuing:

```bash
git checkout "<BRANCH_NAME>"
git branch -vv | grep "^\* "
# If the line shows [origin/main], unset:
# git branch --unset-upstream
```

## Setup 5: Create Session State File

Ensure `.dev/` directory exists and create session file:

```bash
mkdir -p .dev
SESSION_FILE=".dev/<ISSUE_NUM>_SESSION_STATE.md"
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

## Setup 6: Display Summary and Complete Setup

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

Update `.dev/<ISSUE_NUM>_SESSION_STATE.md` with:

- `Status:` → "Setup complete, awaiting investigation"
- `Last Updated:` → current date
- `Current state:` → "Workspace ready; investigation not started"
- `Immediate next action:` → "Run `/bpg:investigate <ISSUE_NUM>` to begin the investigation
  gate"

Then prompt the user with their next move:

```text
AskUserQuestion(
  header: "Next step",
  question: "Setup is complete. Start investigation now?",
  options: [
    { label: "Run /bpg:investigate now", description: "Continue immediately into the investigation gate" },
    { label: "Pause here",                description: "I'll run /bpg:investigate later when ready" }
  ]
)
```

If "Run /bpg:investigate now": tell the user to invoke `/bpg:investigate <ISSUE_NUM>` (this
slash command cannot invoke another slash command directly — only the user can). If "Pause
here": exit cleanly. Either way, do NOT begin investigation work in this skill.

## Push hygiene (read once, apply on every push later in the workflow)

This skill creates the branch but does NOT push. Pushes happen later in `/bpg:prepare-pr` or
manually. Even though `--no-track` was used at creation, follow these rules on every future
push to make accidents impossible:

1. **Always verify tracking immediately before pushing:**

   ```bash
   git branch -vv | grep "^\* "
   ```

   The current-branch line must NOT show `[origin/main]`. If it does, STOP and run
   `git branch --unset-upstream` before pushing.

2. **Only push with an explicit `<local>:<remote>` refspec.** The safe forms (in priority
   order):

   ```bash
   # Best: explicit refspec, both sides named
   git push -u origin "<BRANCH_NAME>:<BRANCH_NAME>"

   # Acceptable: HEAD with explicit refs/heads/ on the remote side
   git push -u origin "HEAD:refs/heads/<BRANCH_NAME>"
   ```

3. **Forbidden push forms** (all carry the risk of landing on `main`):

   ```bash
   git push                    # forbidden — relies on push.default
   git push origin             # forbidden
   git push origin HEAD        # forbidden — remote ref resolved via tracking config
   git push -u origin HEAD     # forbidden — same risk surface
   ```

4. **Verify after pushing:**

   ```bash
   git ls-remote origin "<BRANCH_NAME>"
   ```

   This MUST print a line for the branch. If `main` advanced unexpectedly, stop and notify
   the user immediately — do NOT attempt to fix without explicit instruction.

5. **PR creation:** `gh pr create --head "<BRANCH_NAME>"` — pass `--head` explicitly; do not
   rely on the inferred branch.

</process>

<success_criteria>

**Setup criteria — enforced by this slash command.** The skill's `<process>` block does not
complete until all of these are true:

- [ ] Issue number provided and validated
- [ ] Issue verified to exist on GitHub (or freshly created via the auto-continue path)
- [ ] Closed-issue check passed (open issue, OR user explicitly chose to continue with rationale
      recorded)
- [ ] Issue type determined (fix/feat)
- [ ] Workspace chosen (branch-in-place vs `.claude/worktrees/`)
- [ ] Working tree clean (branch-in-place path) OR worktree created
- [ ] Branch based on `origin/main` (or alternate base verified with `git rev-parse --verify`)
- [ ] Branch created with correct naming: `{type}/{issue}-{description}`
- [ ] Branch created with `--no-track` (or upstream explicitly unset post-checkout) — verified
      via `git branch -vv | grep "^\* "` showing NO `[origin/main]` on the current branch line
- [ ] `.dev/` directory exists
- [ ] Session state template exists
- [ ] Session state file created: `.dev/{issue}_SESSION_STATE.md`
- [ ] Session state populated with issue context
- [ ] Stale log files cleared
- [ ] Issue context displayed to user
- [ ] User prompted to invoke `/bpg:investigate <ISSUE_NUM>` next

The investigation gate is a separate skill — see `/bpg:investigate <ISSUE_NUM>` for its success
criteria. This skill's contract ends with "workspace ready, hand-off prompt issued."

</success_criteria>

<tips>

- If `gh` CLI is not authenticated, the skill will prompt for manual verification
- Branch names are auto-truncated to 40 chars after the issue number to avoid filesystem issues
- Worktrees go under `.claude/worktrees/` (already gitignored). When using a worktree, prefix
  every subsequent command with `cd $WORKTREE_PATH && …` or use absolute paths — shell state
  does not persist across `Bash` tool calls
- Session state file is gitignored, so it won't be committed
- Stale `/tmp/testacc.log` and `/tmp/api_debug.log` are cleared at setup. **Caveat:** these
  paths are global, so two parallel sessions (e.g., two worktrees on different issues) will
  stomp each other's logs. Run `/bpg:debug-api` per-issue to avoid confusion
- After this skill completes, run `/bpg:investigate <ISSUE_NUM>` to start the investigation
  gate. Use `/bpg:resume` to come back to a paused session later
- The "Push hygiene" section near the end of `<process>` is **mandatory reading** before any
  `git push` later in the workflow. It is the only thing standing between your feature
  branch's commits and an accidental landing on `origin/main`

</tips>
