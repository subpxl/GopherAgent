---
description: edit policy — ask before modifying files outside the approved set
---

## Rule: Ask Before Editing

Only the following files may be edited without asking the user first:

- `main.go`
- `pkg/agent/agent.go`
- `pkg/session/session.go`

**For any other file in the project**, you MUST ask the user explicitly before making any changes, even if an issue was found during review. State clearly:
- Which file you want to change
- What the change is
- Why it is needed

Wait for explicit approval before proceeding.
