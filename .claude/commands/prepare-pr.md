---
name: prepare-pr
description: Prepare PR body from template with proof of work
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
Generate a filled-out PR body based on `.github/PULL_REQUEST_TEMPLATE.md` and save it to `.dev/<ISSUE>_PR_BODY.md`.

Use this skill when:

- Preparing to submit a PR
- After running `/ready` checklist
- User asks to "prepare PR" or "create PR body"

The output file can be used directly with `gh pr create --body-file`.
</objective>

<context>
Issue number: $ARGUMENTS

Output location: `.dev/<ISSUE>_PR_BODY.md`
Pattern is gitignored, so the file won't be committed.

PR template location: `.github/PULL_REQUEST_TEMPLATE.md`
</context>

<process>

## Step 1: Determine Issue Number

If `$ARGUMENTS` provided, use it. Otherwise detect from branch name:

```bash
ISSUE_NUM=$(git branch --show-current | grep -oE '(fix|feat)/[0-9]+' | grep -oE '[0-9]+')
```

If still unclear, ask the user.

Set paths:

```bash
PR_BODY=".dev/${ISSUE_NUM}_PR_BODY.md"
SESSION_STATE=".dev/${ISSUE_NUM}_SESSION_STATE.md"
```

## Step 2: Read PR Template

Read `.github/PULL_REQUEST_TEMPLATE.md` to understand the exact structure that needs to be filled.

## Step 3: Gather Context

Run these in parallel where possible:

**Changed files and diff summary:**

```bash
git diff --stat main...HEAD
```

**Commit history:**

```bash
git log --oneline main..HEAD
```

**Detect change type (bug fix, feature, etc.):**

Infer from branch prefix (`fix/` or `feat/`) and commit messages.

**Check for breaking changes:**

Look at commits for the bang character in conventional commit prefix or significant schema changes.

**Check session state** (if exists) for previously gathered context, user decisions, and test results.

## Step 4: Compose PR Title (Squash Commit Message)

The PR title becomes the squash commit message on merge. It must follow conventional commits exactly.

**Format:** `{type}({scope}): {description}`

**Rules (from CONTRIBUTING.md):**

- **Types:** `feat` (new features), `fix` (bug fixes), `chore` (maintenance)
- **Scopes:** `vm`, `lxc`, `provider`, `core`, `docs`, `ci`
- Lowercase description, no period at the end, under 72 characters
- No issue numbers in the title
- For breaking changes, add bang before colon: feat(vm)!: remove legacy clone

Infer type from branch prefix (`fix/` or `feat/`), scope from changed files, and description from the commit history and diff.

If breaking changes were detected in Step 3, include the bang character in the title.

Store the title — it will be written as the first line of the PR body file.

## Step 5: Write "What does this PR do?" Section

Compose a clear, concise summary:

- State the problem (reference the issue)
- Describe what changed and why
- Keep it to 2-5 sentences

Base this on: commit messages, changed files, session state, and the issue context.

## Step 6: Fill Contributor's Note Checklist

Check each item by inspecting actual state — do not assume:

- **make lint** — Check if /ready was run (session state) or run make lint now
- **Documentation updated** — Check git diff for docs/ changes; if schema changed, verify make docs was run
- **Acceptance tests added/updated** — Check git diff for _test.go files
- **Backward compatibility** — Check for schema field removals or type changes in diff
- **Reference examples followed** — Only for new resources; check if resource follows patterns
- **make example run** — Ask user if applicable (SDK/provider config changes)

Mark items `[x]` only when verified. Leave `[ ]` for items not done or not applicable.

## Step 7: Build Proof of Work Section

This is the most important section. Follow the PR template guidelines:

> REQUIRED for code changes. Include at minimum:
>
> - Acceptance test output (`./testacc TestAccYourResource`)
> - For bug fixes: test output showing the fix works
> - For API changes: either mitmproxy logs showing correct API calls,
>   or terraform/tofu output showing successful resource creation/update
>   together with the test resource configuration used

**Gather test output:**

```bash
# Check for saved test log
if [ -f /tmp/testacc.log ]; then
  echo "Found test log"
fi
```

Read `/tmp/testacc.log` if it exists. Extract:

- The test command that was run
- Test names, durations, and pass/fail status
- The final summary line (`ok` / `FAIL`)

**Format the proof of work as:**

1. The exact command used (e.g., `./testacc TestAccResourceVMUpdate -- -v`)
2. Trimmed test output in a code block — include RUN/PASS/FAIL lines and the summary, skip verbose terraform plan/apply noise unless it shows the fix
3. For API changes: mitmproxy excerpt or terraform output if available (check `/tmp/api_debug.log`)

If no test log exists, ask:

```text
AskUserQuestion(
  header: "Test Output",
  question: "No test log found at /tmp/testacc.log. How should proof of work be filled?",
  options: [
    { label: "Run tests now", description: "Run acceptance tests and capture output" },
    { label: "Skip for now", description: "Leave proof of work empty (will delay review)" },
    { label: "Already have output", description: "I'll paste the test output" }
  ]
)
```

If "Run tests now": detect test name from changed test files and run:

```bash
./testacc ${TEST_NAME} -- -v 2>&1 | tee /tmp/testacc.log
```

## Step 8: Set Issue Link

Determine `Closes` vs `Relates`:

- Bug fixes typically `Closes #ISSUE`
- Features typically `Closes #ISSUE`
- Partial work uses `Relates #ISSUE`

Ask if ambiguous:

```text
AskUserQuestion(
  header: "Issue Link",
  question: "Should this PR close or relate to #${ISSUE_NUM}?",
  options: [
    { label: "Closes", description: "PR fully resolves the issue" },
    { label: "Relates", description: "PR is partial or related work" }
  ]
)
```

## Step 9: Generate PR Body

Ensure `.dev/` directory exists:

```bash
mkdir -p .dev
```

Write the filled template to `.dev/${ISSUE_NUM}_PR_BODY.md`.

The output file has two parts:

1. **PR title** — on the first line, prefixed with `# PR Title:` so it's easy to find and copy
2. **PR body** — the filled template, matching `.github/PULL_REQUEST_TEMPLATE.md` structure exactly

Keep all HTML comments from the template. Keep the Community Note section verbatim.

**Template with placeholders:**

```markdown
# PR Title: {PR_TITLE}

---

### What does this PR do?

{PR_SUMMARY}

### Contributor's Note

- [{LINT}] I have run `make lint` and fixed any issues.
- [{DOCS}] I have updated documentation (FWK: schema descriptions + `make docs`; SDK: manual `/docs/` edits).
- [{TESTS}] I have added / updated acceptance tests (**required** for new resources and bug fixes — see [ADR-006](docs/adr/006-testing-requirements.md)).
- [{COMPAT}] I have considered backward compatibility (no breaking schema changes without `!` in PR title).
- [{REFERENCE}] For new resources: I followed the [reference examples](docs/adr/reference-examples.md).
- [{EXAMPLE}] I have run `make example` to verify the change works (mainly for SDK / provider config changes).

{BREAKING_CHANGES_SECTION_IF_APPLICABLE}

### Proof of Work

{PROOF_OF_WORK}

### Community Note

- Please vote on this pull request by adding a reaction to the original pull request comment to help the community and maintainers prioritize this request
- Please do not leave "+1" or other comments that do not add relevant new information or questions, they generate extra noise for pull request followers and do not help prioritize the request

{CLOSES_OR_RELATES} #{ISSUE_NUM}
```

Replace each placeholder with actual gathered values.

For checklist items: use `x` if verified, space if not done/not applicable.

For breaking changes: only include the section if there are breaking changes.

For proof of work: include trimmed test output in a fenced code block. Keep it focused — show the command, test results, and summary. No need to include full terraform plan output unless it demonstrates the fix.

## Step 10: Present Result

Display:

```text
PR body saved to: .dev/${ISSUE_NUM}_PR_BODY.md
PR title: ${PR_TITLE}

To create the PR:
  gh pr create --title "${PR_TITLE}" --body-file .dev/${ISSUE_NUM}_PR_BODY.md

To preview:
  cat .dev/${ISSUE_NUM}_PR_BODY.md
```

Ask if the user wants to view the generated body.

</process>

<success_criteria>

- [ ] Issue number determined
- [ ] PR template read from `.github/PULL_REQUEST_TEMPLATE.md`
- [ ] PR title composed (conventional commits format, under 72 chars)
- [ ] "What does this PR do?" section filled with clear summary
- [ ] Contributor's Note checklist verified (not assumed)
- [ ] Proof of Work section filled with test output
- [ ] Issue link set (Closes/Relates)
- [ ] Breaking changes section included if applicable
- [ ] PR body written to `.dev/{ISSUE}_PR_BODY.md` with title on first line
- [ ] PR creation command provided to user (with `--title` and `--body-file`)
</success_criteria>

<tips>
- Run `/ready` first to ensure all checks pass and test output is captured
- The proof of work section is what reviewers check first — make it thorough
- Keep test output trimmed: RUN/PASS/FAIL lines + summary, not full terraform noise
- The output file works directly with `gh pr create --body-file`
- If session state exists, pull results from there to avoid redundant work
- PR title becomes the squash commit message — make it follow conventional commits exactly
</tips>
