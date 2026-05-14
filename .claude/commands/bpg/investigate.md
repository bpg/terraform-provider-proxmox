---
name: investigate
description: Use when ready to investigate a GitHub issue under the Iron Law gate — drives maintainer-talk, root-cause analysis, pattern matching, hypothesis writing, and explicit approval before TDD. Invoke after `/bpg:start-issue` has set up the workspace; may be re-invoked mid-session to re-anchor. Also use when user says "investigate #1234" or "start investigating".
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
Drive a hard-gated investigation protocol on a GitHub issue, ending with explicit maintainer
approval before TDD begins. Invoke this AFTER `/bpg:start-issue` has set up the branch,
workspace, and session state file. May be re-invoked mid-session to re-anchor — running it
again loads the protocol back into the agent's context without disturbing existing state.

Investigation Steps 0–6 (agent-enforced — no code, no tests, no commits until complete):

0. Talk to the maintainer (4 structured questions + free-text hunches)
1. Pick the analytical lens — `superpowers:systematic-debugging`,
   `superpowers:brainstorming`, `superpowers:dispatching-parallel-agents`, or a combination
2. Root cause via project tooling: gh comments, mitmproxy, Serena, git log
3. Pattern analysis against sibling resources, ADRs, past PRs (with evidence cited)
4. Single hypothesis written in X / Y / Z / T / M form
5. Parallel research dispatch when threads are independent
6. Present findings, get explicit approval via `AskUserQuestion`

After Step 6 approval, proceed to TDD (acceptance test first; waived only by maintainer for
non-functional changes). TDD itself is outside this skill's scope — see `/bpg:ready` when work
is complete.
</objective>

<context>
Issue number: $ARGUMENTS

**Prerequisites:** `/bpg:start-issue <ISSUE_NUM>` must have run first, producing:

- A branch (`fix/<num>-…` or `feat/<num>-…`) checked out, OR a worktree created under
  `.claude/worktrees/`
- A session state file at `.dev/<ISSUE_NUM>_SESSION_STATE.md`

If those don't exist, this skill will redirect to `/bpg:start-issue` rather than proceed.

From [CLAUDE.md](../../../CLAUDE.md): "All work on fixes or features MUST have a corresponding
GitHub issue."
</context>

<process>

### Convention: three kinds of variable

Each `Bash` tool call runs in a fresh shell — environment variables do NOT survive across
calls. This skill uses three notations so the requirement is unambiguous:

| Form      | Meaning                                                                                                                                  | Example use                                         |
| --------- | ---------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------- |
| `<UPPER>` | **Cross-block agent state.** The agent substitutes the entire token with a value carried from a prior step.                              | `gh issue view "<ISSUE_NUM>" --json comments`       |
| `<lower>` | **User-judgment fill-in.** The agent picks a value based on context (search term, file path) — not from prior-step state.                | `gh pr list --search "is:merged <symptom>"`         |
| `$NAME`   | **Shell variable, local to a single bash block.** Defined and consumed within the same `Bash` call; the agent does NOT substitute these. | `ISSUE_JSON=$(gh issue view …); echo "$ISSUE_JSON"` |

Cross-block state used here: `<ISSUE_NUM>`, `<TITLE>`, `<BODY>`, `<COMMENTS>`. Bash locals
include `$ISSUE_JSON`, `$SESSION_FILE`. Quote substituted values as needed.

## Prereq 1: Get Issue Number

If `$ARGUMENTS` is provided, use it as `<ISSUE_NUM>`. Otherwise:

```text
AskUserQuestion(
  header: "Issue Number",
  question: "Which GitHub issue are we investigating?",
  options: [
    { label: "Enter number", description: "I'll type the issue number" },
    { label: "List active sessions", description: "Show .dev/*_SESSION_STATE.md files" }
  ]
)
```

If "List active sessions", run `ls .dev/*_SESSION_STATE.md` and let the user pick.

## Prereq 2: Verify Session State Exists

```bash
SESSION_FILE=".dev/<ISSUE_NUM>_SESSION_STATE.md"
if [ ! -f "$SESSION_FILE" ]; then
  echo "ERROR: $SESSION_FILE not found"
  exit 1
fi
```

If the session file is missing, tell the user:

> No session state found for issue #`<ISSUE_NUM>`. Run `/bpg:start-issue <ISSUE_NUM>` first to
> set up the branch and session state, then return here.

Then exit cleanly. Do NOT proceed without the session state file.

## Prereq 3: Reload Issue Context

Re-fetch the issue body and comments — they may have changed since `/bpg:start-issue` ran:

```bash
ISSUE_JSON=$(gh issue view "<ISSUE_NUM>" --json title,body,comments,state 2>/dev/null)
TITLE=$(echo "$ISSUE_JSON" | jq -r '.title')
BODY=$(echo "$ISSUE_JSON" | jq -r '.body')
COMMENTS=$(echo "$ISSUE_JSON" | jq -r '.comments[] | "--- @\(.author.login):\n\(.body)"')
STATE=$(echo "$ISSUE_JSON" | jq -r '.state')
```

The agent captures `TITLE`, `BODY`, `COMMENTS` from the output and uses them as `<TITLE>`,
`<BODY>`, `<COMMENTS>` in the steps below.

If `STATE == "CLOSED"`, warn the user (the issue may have been closed since setup) and ask
whether to continue. Do not auto-proceed on a closed issue.

Read the session state file to recover branch name, prior maintainer context (if Step 0 has
run before), prior hypothesis (if re-anchoring after a partial investigation), etc.

## Prereq 4: Re-anchor or fresh start

```text
AskUserQuestion(
  header: "Investigation state",
  question: "Is this a fresh investigation or a re-anchor of work already in progress?",
  options: [
    { label: "Fresh start",       description: "First time running investigate for this issue — start at Step 0" },
    { label: "Re-anchor",         description: "Investigation in progress; resume from the last completed step" },
    { label: "Restart from Step 0", description: "Discard prior investigation and start over" }
  ]
)
```

For "Re-anchor": read the session state file, identify the last completed step from sections
present (`Maintainer Context`, `Root Cause`, `Pattern Analysis`, `Hypothesis`, `Step 6 concerns`),
and resume at the next step.

For "Restart from Step 0": archive existing investigation sections in session state under a
`Prior investigation (archived <date>)` heading, then proceed to Step 0.

</process>

---

## The Iron Law

```text
NO FIXES WITHOUT ROOT-CAUSE INVESTIGATION FIRST
```

Mirrors `superpowers:systematic-debugging`. Violating the letter of this process is violating
the spirit. Trivial fixes (single-character typos, comment-only edits, version bumps) still go
through every step — the steps just collapse to one-line answers when there's nothing to
investigate, and the test requirement at Step 6 may be waived by the maintainer when the change
is non-functional. Skipping the gate altogether is faster only in fantasy.

## Step 0 — Talk to the maintainer (mandatory)

Before reading any code, gather maintainer context. Per [CLAUDE.md](../../../CLAUDE.md), the
maintainer's domain knowledge often shortcuts hours of code spelunking.

Use one `AskUserQuestion` call with these four structured questions (the tool's max is 4):

```text
AskUserQuestion(
  questions: [
    {
      header: "Regression",
      question: "Is this a known regression?",
      options: [
        { label: "Yes — known regression",       description: "Was working before, broke recently" },
        { label: "No — never worked / new",      description: "First-time bug or new behavior" },
        { label: "Unsure",                       description: "Haven't checked git history" }
      ]
    },
    {
      header: "Scope",
      question: "Is the issue isolated or systemic?",
      options: [
        { label: "Isolated (single resource/attribute)", description: "One place" },
        { label: "Systemic (multiple)",                  description: "Pattern repeats across resources" },
        { label: "Unsure",                               description: "" }
      ]
    },
    {
      header: "Code path",
      question: "Which provider is affected?",
      options: [
        { label: "Framework (fwprovider/)", description: "New code lives here" },
        { label: "SDK (proxmoxtf/)",        description: "Legacy, feature-frozen" },
        { label: "Both",                    description: "Joint-ownership resource" },
        { label: "Unsure",                  description: "" }
      ]
    },
    {
      header: "Fix scope",
      question: "Acceptable fix scope?",
      options: [
        { label: "Narrow patch",                  description: "Smallest change to make tests pass" },
        { label: "Larger refactor OK",            description: "Maintainer accepts collateral cleanup" },
        { label: "Decide after Step 4 hypothesis", description: "" }
      ]
    }
  ]
)
```

Then ask one free-text follow-up in plain prose (NOT via `AskUserQuestion`, which doesn't
support free text):

> Any hunches about specific files, attributes, recent commits, or sibling resources we should
> compare to? Reply "none" if you have nothing to add.

Record all five answers under a "Maintainer Context" section in
`.dev/<ISSUE_NUM>_SESSION_STATE.md`. Proceed to Step 1 only after the user has answered (or
explicitly waived this step with words like "I have nothing — start digging").

### Step content scales with issue size — but every step cites evidence

Trivial changes (single-character typo, comment-only edit, version bump, lint-only change,
non-functional doc rewording) still go through every step. A step may collapse to a one-line
answer when there's genuinely nothing to find — but **the answer must cite what was looked at**.
Bare assertions don't count.

Bad (assertion):

- Step 3 pattern: "no working twin"
- Step 5 parallel: "no independent threads"

Good (assertion + evidence):

- Step 3 pattern: "ran `rg -n 'pve_version' fwprovider/`, only one declaration site, no twin"
- Step 5 parallel: "two threads — sibling-attr search and past-PR search — both feed Step 3
  output, so sequential, no parallelism applies"

The discipline is "I checked X and the result was nothing," not "I assume nothing." When you
catch yourself writing a one-line step answer with no evidence cited, return to the step and
do the actual check.

The acceptance-test requirement (the TDD gate at Step 6) is the only thing the maintainer may
explicitly waive, and only when the change is non-functional. Step 4's hypothesis is always
written with a real test plan — the agent does not pre-classify a change as non-functional.
The maintainer makes that call at Step 6.

## Step 1 — Pick the analytical lens

Different issue shapes need different superpowers. Pick before digging:

| Issue shape                                                          | Primary superpower                                                   | Why                                                                                     |
| -------------------------------------------------------------------- | -------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| Bug, test failure, "doesn't work as documented"                      | `superpowers:systematic-debugging`                                   | 4-phase root-cause gate; symptom fixes are failure                                      |
| Feature request, ambiguous "should we do X?", design unclear         | `superpowers:brainstorming`                                          | Stress-tests intent and requirements before investigation narrows on the wrong question |
| Hand-wavy issue body, or maintainer flagged design discussion needed | `superpowers:brainstorming` then `systematic-debugging`              | Lock down scope first, then debug                                                       |
| 2+ independent research threads surface during Steps 2–3             | `superpowers:dispatching-parallel-agents` (in addition to the above) | Parallel agents save context and wall time when threads don't share state               |

Invoke the chosen skill explicitly via the `Skill` tool — don't just hold it in mind.

**Lenses can combine and re-pick.** A bug that turned out to be a design problem warrants
re-picking the lens (debug → brainstorm) when that becomes clear during Step 2 or 3. A feature
where the design is settled but implementation has bugs benefits from both lenses interleaved.
The table is a starting point, not a commitment — record any lens change in session state with
a one-line reason.

## Step 2 — Root cause (Phase 1)

Map systematic-debugging's Phase 1 onto this codebase:

| Phase 1 activity              | Project tooling                                                                                                                                                                        |
| ----------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Read errors carefully         | Reuse `<BODY>` and `<COMMENTS>` from Prereq 3 — full thread, not just the body. Do NOT re-fetch.                                                                                       |
| Reproduce consistently        | Acceptance test (`./testacc TestName`) or mitmproxy via `/bpg:debug-api`                                                                                                               |
| Check recent changes          | `git log --grep "<keyword>" --since=6.months -- path/` and `gh pr list --search "<symptom>"`                                                                                           |
| Gather evidence at boundaries | Terraform plan ↔ provider model ↔ HTTP request ↔ PVE API ↔ PVE backend. mitmproxy reveals which boundary breaks (memory: `vm-disk-pve-config-traps`, `acceptance-test-raw-pve-config`) |
| Trace data flow backward      | Serena `find_symbol` → `find_referencing_symbols` from the failing attribute outward, until you hit the source (model → `toAPI` → client → HTTP)                                       |

**Output of Step 2:** a written statement of WHAT goes wrong and WHERE in the code path.
Record in session state under `Root Cause:`.

## Step 3 — Pattern analysis (Phase 2)

Find what already works that's similar to what's broken:

- **Sibling resources/attributes** — Most provider bugs have a working twin. Use Serena/Grep to
  find an attribute or resource that handles the same case correctly. Compare line by line.
- **ADRs** — `docs/adr/` documents design conventions that often explain the right behavior:
  ADR-004 (schema design), ADR-005 (error handling), ADR-006 (testing), ADR-008 (sub-package
  ports).
- **Past similar bugs** — `gh pr list --search "is:merged <symptom>"` and
  `gh issue list --state closed --search "<symptom>"`. Yours is rarely the first occurrence.
- **Proxmox API source** when docs are insufficient — see memory `proxmox-api-source-code` for
  paths into the Perl source.

**Output of Step 3:** a list of differences between the working twin and the broken case, no
matter how trivial. Record in session state under `Pattern Analysis:`.

## Step 4 — Hypothesis (Phase 3)

Write a single hypothesis before touching code. **The form is unconditional** — the agent
always drafts a real acceptance test plan. Whether to waive the test is the maintainer's call
at Step 6, not a Step 4 classification:

```text
I think the root cause is X, because Y.
The fix is Z.
The failing acceptance test that proves the fix is T; it currently fails with message M
and will pass once Z is applied.
```

Variables: X = root cause, Y = supporting evidence, Z = fix, T = acceptance test plan,
M = current failure message of T. If any variable is unfilled, you are not ready to write code
— return to Step 2 or 3.

For changes that look non-functional (typo, comment, version bump): still draft T as a real
test plan — even if "the test" is `verify the rendered doc shows the corrected text via
`make docs && grep …`". The maintainer decides at Step 6 whether to take that path or waive.
Do NOT pre-classify the change as non-functional; that's the maintainer's judgment.

**Hypothesis sets** are allowed during investigation, not after Step 4 is committed: if
multiple plausible root causes survive Steps 2–3, write them all down as candidates with the
discriminating evidence each predicts. Then run that evidence-gathering, narrow to ONE, and
only then commit to the form above. Do not enter Step 6 with multiple hypotheses live — pick
one or return to Step 2.

Record the committed hypothesis in session state under `Hypothesis:`.

## Step 5 — Dispatch parallel research if it pays off

When Steps 2–3 surface 2+ independent investigative threads, invoke
`superpowers:dispatching-parallel-agents` and run them concurrently. Examples of independent
threads:

- Thread A — "Find every Framework attribute of type X and report which ones handle nil-back
  correctly"
- Thread B — "Search merged PRs and closed issues for prior fixes touching `<symptom>` and
  summarize patterns"
- Thread C — "Read PVE Perl source for endpoint Z and report the parameter shape"

Do NOT dispatch parallel agents for sequential threads (where C needs A's output) or for
threads that share state.

## Step 6 — Present and get approval

Update `.dev/<ISSUE_NUM>_SESSION_STATE.md` with sections for: Maintainer Context (Step 0), Root
Cause (Step 2), Pattern Analysis (Step 3), Hypothesis (Step 4). Then summarize for the user:

1. Root cause (one sentence)
2. Proposed fix (one sentence)
3. Acceptance test plan
4. Open questions or risks

**Approval is required before moving to TDD, and approval is always machine-readable.** Do not
parse natural-language replies for approval signals — soft phrases like "looks fine" or "ok"
are ambiguous and the agent will rationalize them as approval under pressure. Always use
`AskUserQuestion`:

```text
AskUserQuestion(
  header: "Step 6 approval",
  question: "Approve this hypothesis and proceed to TDD?",
  options: [
    { label: "Proceed to TDD",        description: "Hypothesis stands; write the failing acceptance test" },
    { label: "Proceed (waive test)",  description: "Non-functional change; skip the acceptance test, run V verification only" },
    { label: "Refine first",          description: "Hypothesis needs more work — go back to Step 2/3/4" },
    { label: "Block — concerns",      description: "Don't proceed; I have concerns to discuss first" }
  ]
)
```

Only "Proceed to TDD" or "Proceed (waive test)" unblocks the next phase. Any natural-language
reply outside this question — including "approved", "looks good", silence — is NOT approval;
re-issue the `AskUserQuestion` and wait. The "waive test" option is the only place a maintainer
waives TDD, and it must be explicit.

**On "Refine first":** prompt the user with one focused `AskUserQuestion` asking _which_ part
needs revision — root cause (Step 2), pattern analysis (Step 3), hypothesis form (Step 4), or
test plan (within Step 4's T). Re-do that step and only that step, then return to Step 6 with
an updated summary. Do NOT silently re-do all of Steps 2–4.

**On "Block — concerns":** ask the user to articulate the concerns in plain prose. Record them
in session state under `Step 6 concerns:` with timestamp. Then ask which step they want to
revisit (or whether they want to abandon the issue). Do NOT proceed without a follow-up
direction; "Block" means "stop and clarify," not "stop forever."

## Fix with TDD

Only proceed after Step 6 approval. Follow TDD (Red-Green-Refactor):

1. **RED — Write a failing acceptance test first**
   - Create an acceptance test that reproduces the bug
   - Run it with `./testacc` and **verify it fails for the expected reason** (matches M from
     Step 4's hypothesis)
   - If it fails for a different reason (connection issues, missing infrastructure), **ask the
     user** — do NOT work around it
   - If the bug cannot be reproduced with acceptance tests, ask the user how to proceed
2. **GREEN — Implement the minimal fix**
   - Write the simplest code that makes the failing test pass
   - Run the test again and verify it passes
3. **Verify — No regressions**
   - Run related existing acceptance tests to confirm no regressions
   - Run `make lint`
4. Run `/bpg:ready` before completing work

## Rationalizations — STOP if you catch yourself thinking…

| Excuse                                                | Reality                                                                                        |
| ----------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| "Issue is small, skip Step 0"                         | Maintainer's hunch is the highest-bandwidth signal you have. Asking takes two minutes.         |
| "I already know the area, skip pattern analysis"      | Confirmation bias is highest when you "know the area." Compare to a working twin anyway.       |
| "Just write the test, hypothesis will fall out"       | Test-first without hypothesis = fishing. You'll write tests that pass for the wrong reason.    |
| "Investigation is theatre, the user wants the fix"    | The user wants the _right_ fix that doesn't regress. Investigation is the cheapest path there. |
| "Parallel agents make me look thorough"               | Parallelism is for independent threads only. Faking it wastes context.                         |
| "Step 6 approval is a formality"                      | Approval is where the maintainer catches your wrong framing before code is written.            |
| "It's a feature request, debugging skill doesn't fit" | True — use `superpowers:brainstorming` instead. Step 1's table tells you which skill applies.  |
| "Test feels like overkill for this fix"               | Only the maintainer waives the test, and only for non-functional changes. Don't pre-decide.    |

## What NOT to Do

- Do NOT write unit tests, helper functions, or refactors as workarounds when acceptance tests
  can't reproduce the bug. Ask the user.
- Do NOT skip investigation and jump straight to code or tests. Iron Law.
- Do NOT work around infrastructure problems (Proxmox unreachable, missing services). Ask the
  user.
- Do NOT propose multiple hypotheses simultaneously. One at a time (Phase 3 discipline).
- Do NOT dispatch parallel agents for sequential threads or threads that share state.

<success_criteria>

**Prereq criteria — enforced by this slash command.** The skill's `<process>` block does not
complete until all of these are true:

- [ ] Issue number provided (from `$ARGUMENTS` or `AskUserQuestion`)
- [ ] Session state file `.dev/<ISSUE_NUM>_SESSION_STATE.md` exists; if not, redirected to
      `/bpg:start-issue` and exited
- [ ] Issue context (`<TITLE>`, `<BODY>`, `<COMMENTS>`) reloaded via `gh issue view`
- [ ] Closed-state warning issued if `STATE == "CLOSED"`
- [ ] Re-anchor decision recorded (fresh / re-anchor / restart)

**Investigation criteria — enforced by the agent after the slash command exits.** These are NOT
checked by the slash command. The Iron Law lives in the agent's discipline; if any item is
unchecked when TDD begins, the agent has violated the contract:

- [ ] Step 0 — Maintainer context (4 structured answers + free-text hunches) gathered or
      explicitly waived; recorded in session state
- [ ] Step 1 — Analytical lens picked (or "none — mechanical change"); chosen `superpowers:*`
      skill invoked via the `Skill` tool when one applies. Lens may be re-picked if
      investigation reveals it was wrong.
- [ ] Step 2 — Root cause identified (WHAT and WHERE) and written down
- [ ] Step 3 — Working twin / sibling pattern identified and differences enumerated, OR a
      one-line "no working twin applies" with **evidence cited** (e.g., the grep that came back
      empty)
- [ ] Step 4 — Single hypothesis written in the X / Y / Z / T / M form, always with a real
      test plan. Hypothesis sets allowed during exploration but narrowed to ONE before Step 6.
      The agent never pre-classifies a change as non-functional; only Step 6 can waive the test.
- [ ] Step 5 — Parallel agents dispatched if and only if 2+ independent threads exist (else
      one-line "no parallel threads" with **why** — e.g., threads are sequential, share state)
- [ ] Step 6 — Findings summarized to user; explicit approval received via `AskUserQuestion`
      before TDD begins. Test waiver, if any, must be explicit at this step.

</success_criteria>

<tips>

- This skill is re-invokable. If the agent loses context mid-investigation, re-run
  `/bpg:investigate <ISSUE_NUM>` to reload the protocol; pick "Re-anchor" at Prereq 4 to resume
  from the last completed step (read from session state)
- Trivial changes (typo, comment, version bump) still run every step — each step may be a
  one-line answer when there's nothing to dig into, but the answer must cite what was looked at.
  The maintainer may waive the acceptance-test requirement at Step 6, but only for genuinely
  non-functional changes
- `<TITLE>`, `<BODY>`, `<COMMENTS>` are reloaded fresh on each invocation. Session state holds
  the durable artifacts (Maintainer Context, Root Cause, Pattern Analysis, Hypothesis,
  concerns)
- Memory references (`vm-disk-pve-config-traps`, `acceptance-test-raw-pve-config`,
  `proxmox-api-source-code`) point to entries in the agent's persistent memory. If a memory has
  been removed or renamed, fall back to grepping the codebase or reading the relevant ADR
  rather than treating the reference as authoritative
- After Step 6 approval, TDD is outside this skill's scope. Run `/bpg:ready` when work is
  complete

</tips>
