---
name: code-review
description: Review a pull request for bugs and code compliance, displaying results locally
argument-hint: "<PR number>"
allowed-tools:
  - Read
  - Bash
  - Grep
  - Glob
  - Write
  - Edit
  - Task
  - TaskCreate
  - TaskUpdate
  - TaskList
  - AskUserQuestion
---

Review pull request #$ARGUMENTS.

Make a todo list of the steps below, then follow them precisely.

## Step 0: Setup

1. Determine the PR number from the argument. If not provided, ask the user.

## Step 1: Eligibility Check

Use a Haiku agent to check if the pull request is eligible for review. Use `gh pr view <number> --json state,isDraft,author` to check. Skip the review if:

- The PR is **closed** or **merged**
- The PR is a **draft**
- The PR author is a bot (e.g., dependabot, renovate, github-actions)

Also check if a session state file already exists at `.dev/review_PR_<number>_SESSION_STATE.md`. If it does, warn the user that a previous review exists and ask whether to re-review or resume from the existing state.

If the PR is not eligible, tell the user why and do not proceed.

## Step 2: Checkout PR in Worktree

Check out the PR in an isolated worktree to avoid disrupting the user's working tree.

First, clean up any leftover worktree from a previous interrupted review:

```bash
# Remove stale worktree if it exists
git worktree remove .claude/worktrees/review-<number> 2>/dev/null || true
```

Then create the worktree:

```bash
git fetch origin pull/<number>/head
git worktree add .claude/worktrees/review-<number> FETCH_HEAD --detach
```

> **Note:** `.claude/` should be in `.gitignore`. If it is not, warn the user that the worktree path may appear in `git status`.

Set `WORKTREE=.claude/worktrees/review-<number>` — all subsequent file reads in review agents must use this path prefix.

Get the full commit SHA for linking: `git -C .claude/worktrees/review-<number> rev-parse HEAD`

Determine the base branch from `gh pr view <number> --json baseRefName -q .baseRefName`.

## Step 3: Gather Project Guidelines

Collect the list of guideline files directly (no agent needed):

1. Get changed files: `gh pr diff <number> --name-only`
2. Extract the unique parent directories from those file paths.
3. Build the guideline file list:
   - `CONTRIBUTING.md` (primary contributor guidelines — always include)
   - The root `CLAUDE.md` (if it exists)
   - Any `CLAUDE.md` in the changed directories (use Glob: `$WORKTREE/<dir>/CLAUDE.md`)

Store this list for use in Steps 5 and 6.

## Step 4: Determine PR Size and Summarize

Use a Haiku agent to:

1. View the pull request diff (`gh pr diff <number>`)
2. Count the total lines changed (additions + deletions)
3. Return both a **summary** of the change and the **line count**

Classify the PR:

- **Small:** <50 lines changed
- **Medium:** 50–300 lines changed
- **Large:** >300 lines changed

## Step 5: Parallel Code Review

Launch review agents based on PR size. **Each agent's prompt must include the worktree path (`$WORKTREE`) and the list of guideline files from Step 3.** Agents should return a list of issues with the reason each was flagged (e.g., guidelines violation, bug, historical context).

**Small PRs (<50 lines) — 2 Sonnet agents:**

- Agent 1: Guidelines compliance
- Agent 2: Bug scan

**Medium PRs (50–300 lines) — 3 Sonnet agents:**

- Agent 1: Guidelines compliance
- Agent 2: Bug scan
- Agent 3: Historical context

**Large PRs (>300 lines) — 5 Sonnet agents:**

- Agent 1: Guidelines compliance
- Agent 2: Bug scan
- Agent 3: Historical context
- Agent 4: Prior PR comments
- Agent 5: Code comment compliance

### Agent Descriptions

a. **Guidelines compliance:** Audit the changes against `CONTRIBUTING.md` (the primary source of contributor guidelines) and any relevant `CLAUDE.md` files. Key areas to check from CONTRIBUTING.md: coding conventions, commit message format, PR scope (one change per PR, no mixed concerns), proof of work, DCO sign-off, Framework-only for new resources, documentation workflow, and test placement. Note that CLAUDE.md is guidance for Claude agents, so not all of its instructions apply during code review of human contributions.

b. **Bug scan:** Read the file changes in the pull request, then do a shallow scan for obvious bugs. Avoid reading extra context beyond the changes, focusing just on the changes themselves. Focus on large bugs, and avoid small issues and nitpicks. Ignore likely false positives.

c. **Historical context:** Read the git blame and history of the code modified, to identify any bugs in light of that historical context.

d. **Prior PR comments:** Find the last 10 merged PRs that touched the same files. Use `git log --pretty=format:"%H" -- <file>` in the worktree to get commits that touched each file, then `gh pr list --state merged --search "<commit SHA>"` to find associated PRs. Read comments on those PRs and check for any that may also apply to the current PR.

e. **Code comment compliance:** Read code comments in the modified files, and make sure the changes in the pull request comply with any guidance in the comments.

## Step 6: Score Issues

Batch all issues and launch **one Haiku agent per review agent that found issues** to score that agent's issues. Skip scoring for agents that reported no issues. Each scoring agent receives the PR diff, the issues from its corresponding review agent, and the list of guideline files (from Step 3). It returns a score for each issue.

For issues flagged due to guideline violations, the agent should double-check that `CONTRIBUTING.md` or the relevant `CLAUDE.md` actually calls out that issue specifically.

Scoring rubric (give this to the agents verbatim):

- 0: Not confident at all. This is a false positive that doesn't stand up to light scrutiny, or is a pre-existing issue.
- 25: Somewhat confident. This might be a real issue, but may also be a false positive. The agent wasn't able to verify that it's a real issue. If the issue is stylistic, it is one that was not explicitly called out in CONTRIBUTING.md or the relevant CLAUDE.md.
- 50: Moderately confident. The agent was able to verify this is a real issue, but it might be a nitpick or not happen very often in practice. Relative to the rest of the PR, it's not very important.
- 75: Highly confident. The agent double checked the issue, and verified that it is very likely it is a real issue that will be hit in practice. The existing approach in the PR is insufficient. The issue is very important and will directly impact the code's functionality, or it is an issue that is directly mentioned in CONTRIBUTING.md or the relevant CLAUDE.md.
- 100: Absolutely certain. The agent double checked the issue, and confirmed that it is definitely a real issue, that will happen frequently in practice. The evidence directly confirms this.

## Step 7: Filter

Filter out any issues with a score less than 50. If there are no issues that meet this criteria, report that no issues were found.

## Step 8: Display Results

Display the review results directly in the conversation. **Never post to GitHub.**

Reference files using full paths from the project root with line numbers.

Format for issues found:

```text
### Code review for PR #NUMBER

**Summary:** [1-2 sentence PR summary from Step 4]

Found N issues (M total identified, K filtered as low confidence):

1. (Score: 85) <brief description>
   File: `fwprovider/resource_vm.go:10-15`
   Category: <one of: Bug, Guidelines, Historical, Prior PR, Code comments>
   Details: <explanation and suggestion>

2. ...
```

Format for no issues:

```text
### Code review for PR #NUMBER

**Summary:** [1-2 sentence PR summary from Step 4]

No issues found. Checked for bugs and contributor guidelines compliance.
```

## Step 9: Save Session State

Write a session state file to `.dev/review_PR_NUMBER_SESSION_STATE.md`.

**Important:** Include ALL issues identified by the review agents, not just the ones that passed the score >= 50 filter. The session state must capture the complete picture so that filtered issues can be re-evaluated in follow-up sessions or after author changes.

**For clean reviews (no issues found), use this short format:**

```markdown
# Code Review Session State — PR #NUMBER

- **PR:** https://github.com/bpg/terraform-provider-proxmox/pull/NUMBER
- **Last Updated:** YYYY-MM-DD HH:MM
- **Status:** Review Complete — No Issues Found
- **PR Summary:** [1-2 sentence summary]
- **PR Size:** Small/Medium/Large (N lines changed)
- **Files Changed:** [list of files]
- **Agents Run:** [which agents were run based on PR size]
```

**For reviews with issues, use this full format:**

```markdown
# Code Review Session State — PR #NUMBER

- **PR:** https://github.com/bpg/terraform-provider-proxmox/pull/NUMBER
- **Last Updated:** YYYY-MM-DD HH:MM
- **Status:** Review Complete | Pending Author Changes

---

## Quick Context Restore

**What this PR is about:**
[1-2 sentence summary of the PR]

**Review status:**
[Summary of review outcome — how many issues found, key concerns]

**Immediate next action:**
[What to do next — e.g. "Fix issues locally" or "Wait for author to address issues"]

---

## Review Summary

**PR Summary:**
[Summary from Step 4]

**PR Size:** Small/Medium/Large (N lines changed)

**Files Changed:**
[List of files modified in the PR]

**Agents Run:**
[Which agents were run based on PR size]

---

## Issues Found (Score >= 50)

| # | Score | Description | File | Status |
|---|-------|-------------|------|--------|
| 1 | 85 | Brief description | `path/to/file.go:10` | Open |

### Issue Details

#### Issue 1: Title

- **Score:** 85
- **File:** `path/to/file.go:10-15`
- **Category:** Bug / Guidelines / Historical context
- **Description:** Detailed description
- **Suggestion:** What should be changed

---

## All Identified Issues (Including Filtered)

Record every issue from all review agents here, regardless of score.

| # | Score | Agent | Description | File | Reason Filtered |
|---|-------|-------|-------------|------|-----------------|
| 1 | 85 | Bug scan | Brief description | `path/to/file.go:10` | — (included) |
| 2 | 40 | Guidelines | Brief description | `path/to/other.go:5` | Low confidence |
| 3 | 25 | Historical | Brief description | `path/to/file.go:20` | Likely false positive |

---

## Session Log

### YYYY-MM-DD HH:MM - Code Review Agent

**Completed:**

- Reviewed PR #NUMBER
- Found N issues (M filtered out as low confidence)

**Next Steps:**

- User to review findings and choose: fix locally or end review
```

## Step 10: Clean Up Worktree

Remove the isolated worktree:

```bash
git worktree remove .claude/worktrees/review-<number>
```

If removal fails, warn the user and provide the manual cleanup command:
`git worktree remove --force .claude/worktrees/review-<number>`

## Step 11: Next Steps

Tell the user where the session state file was saved, then ask what they'd like to do next:

- **Fix issues** — Check out the PR branch locally and work through the identified issues
- **Done** — End the review (they can resume later with `/bpg:resume`)

If there are no issues (clean review), skip the prompt and just report the result.

## Step 12: Fix Mode (if selected)

If the user chose to fix issues:

1. **Check for uncommitted changes** before switching branches:

   ```bash
   git status --porcelain
   ```

   If there are uncommitted changes, ask the user how to proceed (stash, commit, or abort).

2. **Check out the PR branch** for local editing:

   ```bash
   gh pr checkout <number>
   ```

3. **Create a todo list** from the issues found in Step 8 (using TaskCreate), one task per issue, ordered by score (highest first). Each task should include the file path, line numbers, and a brief description of what to fix.

4. **Work through the issues** — fix each one, marking tasks complete as you go.

5. When all issues are addressed, run `make lint` and `make test` to verify, then tell the user they can commit and push when ready.

6. **Update the session state file** — mark each fixed issue's Status as "Fixed" in the issues table, and update the Session Log with a new entry describing the fixes made.

---

## False Positive Guidance

Examples of false positives, for Steps 5 and 6:

- Pre-existing issues
- Something that looks like a bug but is not actually a bug
- Pedantic nitpicks that a senior engineer wouldn't call out
- Issues that a linter, typechecker, or compiler would catch (eg. missing or incorrect imports, type errors, broken tests, formatting issues, pedantic style issues like newlines). No need to run these build steps yourself — it is safe to assume that they will be run separately as part of CI.
- General code quality issues (eg. lack of test coverage, general security issues, poor documentation), unless explicitly required in CONTRIBUTING.md or CLAUDE.md
- Issues that are called out in project guidelines, but explicitly silenced in the code (eg. due to a lint ignore comment)
- Changes in functionality that are likely intentional or are directly related to the broader change
- Real issues, but on lines that the user did not modify in their pull request

## Notes

- Do not check build signal or attempt to build or typecheck the app. These will run separately, and are not relevant to your code review.
- Use `gh` to interact with Github (eg. to fetch a pull request), rather than web fetch.
- You must cite each issue with the full file path from the project root and line numbers (e.g., `fwprovider/resource_vm.go:10-15`). When referencing a guideline violation, cite the specific section of `CONTRIBUTING.md` or `CLAUDE.md`.
- When linking to code on GitHub (e.g., in session state), use the full git SHA (not HEAD or abbreviated). Line range format is `L[start]-L[end]`, provide at least 1 line of context before and after.
- **Never post reviews to GitHub.** All results are displayed locally in the conversation only.
