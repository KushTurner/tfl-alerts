# CLAUDE.md

## Error handling

Wrap errors with context at every layer using `%w`:
```go
fmt.Errorf("finding users for train %s: %w", train.Line, err)
```
Never log and return — pick one. Return errors up the stack and log once at the top. Use sentinel errors for known failure cases:
```go
var ErrNotFound = errors.New("not found")
```

## Interfaces

Keep interfaces small and defined at the point of use. If mocking an interface requires implementing methods irrelevant to the test, split it into smaller interfaces by responsibility.

## Layer boundaries

DB types and external API response types must not leak past their own layer. The DB layer maps to domain structs, the TFL client maps to domain structs — the service layer only ever sees domain types.

## No ORMs

Raw SQL only.

## Testing

Test observable behaviour and business rules, not wiring. If a test can only break due to an internal rename or structural change — not a logic bug — it's not worth writing. Almost every new feature should have tests. The DB layer gets integration tests; everything else gets unit tests.

## PRs and features

Before raising a PR:
- Run tests via `make test` — all must pass
- Branch must be off the latest `main` with no merge conflicts
- PR description must explain what changed and why — enough context for someone reviewing cold

## Commits and branches

Commit messages are short imperative sentences, no prefix or tag (e.g. `Fix notification windows ignoring BST/GMT offset`). Branch names carry the type: `feature/`, `fix/`.

## Simplicity

Prefer the simplest solution. Only extract duplication after 3+ occurrences, and only if it's genuinely the same concept. Don't add abstractions for hypothetical future requirements.
