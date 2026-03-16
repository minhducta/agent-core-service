---
applyTo: "internal/repository/**"
---

# Repository Layer Instructions

The `internal/repository` package implements the `domain.XxxRepository` interfaces using **raw SQL via `sqlx`**. It is the only layer that communicates with PostgreSQL.

## File Responsibilities

| File | Contents |
|---|---|
| `bot_repository.go` | PostgreSQL implementation of `domain.BotRepository` — GetByID, GetByAPIKeyHash, UpdateLastSeen |
| `memory_repository.go` | PostgreSQL implementation of `domain.BotMemoryRepository` — ListByBotID, Create, Delete |
| `skill_policy_repository.go` | PostgreSQL implementation of `domain.BotSkillRepository` and `domain.BotPolicyRepository` — ListByBotID |
| `todo_repository.go` | PostgreSQL implementation of `domain.TodoRepository` and `domain.TodoChecklistRepository` — ListByBotID, GetByID, Update, ListByTodoID |
| `heartbeat_repository.go` | PostgreSQL implementation of `domain.HeartbeatRepository` — Create, GetStatusByBotID |

## Repository Struct Pattern

```go
type XxxRepository struct {
    db *database.DB
}

func NewXxxRepository(db *database.DB) *XxxRepository {
    return &XxxRepository{db: db}
}
```

The returned type is the concrete `*XxxRepository`, assigned to a `domain.XxxRepository` variable in `main.go`.

## SQL Rules

1. **Always use parameterised queries** — `$1, $2, …` placeholders with `sqlx`. Never concatenate user input into SQL strings.

2. **Use `db.GetContext` / `db.SelectContext`** for reads (they scan directly into structs via `db:` tags). Use `db.ExecContext` for writes.

3. **Multi-tenancy enforced at query level** — always filter by `bot_id` where applicable.

4. **Pagination**: `List` methods accept `page int, limit int` and produce `LIMIT $N OFFSET $M`. Always run a separate `COUNT(*)` query and return `([]Entity, totalCount int, error)`.

5. **Transactions**: use `db.BeginTxx(ctx, nil)` for operations that must be atomic. Defer `tx.Rollback()` immediately after `Begin`, and call `tx.Commit()` only on success.

6. **No business logic**: repositories only translate between Go structs and SQL. Validation and orchestration belong in `usecase/`.

7. **No logging inside repository**: the usecase layer logs errors received from the repository.

8. **Return `nil, nil`** when a record is not found (checking `sql.ErrNoRows`) — the usecase decides how to handle a missing record.

## Adding a New Repository Method

1. Declare the method signature in `domain/repository.go` under the appropriate interface.
2. Implement it in the corresponding `*_repository.go` file.
3. Write parameterised SQL — reference the migration files in `migrations/` for exact column names.
4. If the operation requires a new column or table, add a migration first.
