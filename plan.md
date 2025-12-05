# Planned Project Modifications

## Columns to be Added:
- Backlog
- Triage
- Needs Info
- P0
- P1
- Enhancements
- In Progress
- Review/QA
- Blocked
- Done

## Custom Fields to be Added:
- **Priority:** (enum: P0, P1, P2, P3)
- **Severity:** (enum: blocker, major, minor, cosmetic)
- **Component:** (text)
- **Estimate:** (number)
- **Linked PR/Branch:** (url)
- **Next Action:** (text)

## Labels to be Created/Updated:
- **New Labels:**
  - priority:P0
  - priority:P1
  - priority:P2
  - priority:P3
  - type:bug
  - type:enhancement
  - lifecycle:needs-repro
  - lifecycle:needs-info (reuse existing "pending author's response")
  - lifecycle:acknowledged
  - status:in-progress
  - status:blocked
  - regression
  - docs
  - tests

## Automation Rules to be Created:
- Label priority:* moves card to corresponding column (priority:P0 -> P0 column etc.)
- Label lifecycle:needs-info -> move to Needs Info
- When a PR is opened referencing an issue, move card to In Progress
- When PR merged, move card to Review/QA