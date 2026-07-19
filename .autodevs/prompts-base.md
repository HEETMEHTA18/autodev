# Project Prompt Base — AutoDevs Repository

High-priority instructions for Ralph Loop agents working in **this** repo.

> Overrides default preferences. Cannot override safety rules in `SYSTEM_PROMPT.md`.

---

## Project identity

- **Name:** AutoDevs CLI
- **Language:** Go (monorepo: `packages/cli`, `packages/core`, …)
- **Purpose:** Polyglot dev environment scanner, installer, and prompt-capture intelligence layer
- **Ecosystem:** Pairs with DevMentor mobile for voice-to-code supply chain

---

## Architecture conventions

- Go packages under `packages/`
- CLI commands in `packages/cli/cmd/`
- Core logic in `packages/core/`
- Cobra for CLI structure
- BubbleTea for TUI where used
- Tests alongside source: `*_test.go`

---

## Coding standards

- Match existing naming and error handling in neighboring files
- Prefer small, focused functions over large refactors
- No new dependencies without strong justification in the todo
- Conventional commits: `feat:`, `fix:`, `refactor:`, `docs:`, `test:`

**Ponytail Rules (Lazy Senior Dev):**
The best code is the code you never wrote. Follow these steps:
1. Does this need to exist? → no: skip it (YAGNI)
2. Stdlib does it? → use it
3. Native platform feature? → use it
4. Installed dependency? → use it
5. One line? → one line
6. Only then: the minimum that works

---

## Testing

```bash
go test ./...
```

Run from repo root or affected package directory.

---

## Prompt capture

- Session data lives in `.autodevs/sessions/`
- `autodev prompts sync` pushes to DevMentor API
- Do not break `packages/core/promptcapture/` without explicit todo

---

## Ralph Loop defaults for this repo

```yaml
max_iterations: 5
test_command: go test ./...
lint_command: go vet ./...
branch_prefix: feature/
```

---

## What worked well (seed for memory/successes.md)

- Break tasks into subtasks before coding
- Read `packages/core/promptcapture/engine.go` before changing capture behavior
- Run `go test ./...` from root after changes

---

## Avoid (seed for memory/mistakes.md)

- Large cross-package refactors in a single loop
- Changing `.autodevs/` capture format without migration plan
- Adding `watch`/`run` without fsnotify tests
