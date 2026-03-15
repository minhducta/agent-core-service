---
applyTo: "internal/repository/**"
---

# Repository Layer Instructions

Implements `domain` repository interfaces using raw SQL via `sqlx`.

## Rules

1. **Never import from `usecase/` or `handler/`** — only `domain/` and `pkg/`.
2. **Raw SQL only** — no ORM. Use `sqlx.NamedExec`, `sqlx.Get`, `sqlx.Select`.
3. **Parameterised queries always** — never string interpolation.
4. **Wrap errors** with `fmt.Errorf("repo.X: %w", err)` for traceability.
5. **Multi-tenancy enforced at query level** — always filter by `bot_id` where applicable.

## Testing

- Use `github.com/DATA-DOG/go-sqlmock` — never a real database in unit tests.
- Table-driven tests with `t.Run`.
- Test happy path + error path for every public method.
