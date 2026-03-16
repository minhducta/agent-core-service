---
applyTo: "internal/usecase/**"
---

# Usecase Layer Instructions

The `internal/usecase` package contains **all business logic**. It orchestrates repositories, cache, and the Kafka producer, and is the only layer that decides *what* to do with data.

## File Responsibilities

| File | Contents |
|---|---|
| `bot_usecase.go` | Bot identity: GetProfile, GetIdentity, GetBootstrap, ResolveByAPIKey; cache-aside for bot lookups; API key hash resolution |
| `memory_usecase.go` | Memory CRUD: ListMemories, CreateMemory, DeleteMemory; cache invalidation; Kafka event publishing |
| `skill_usecase.go` | Skill read operations: ListSkills |
| `policy_usecase.go` | Policy read operations: ListPolicies |
| `todo_usecase.go` | Todo management: ListTodos, UpdateTodo, GetChecklist, UpdateChecklistItem; Kafka events for state changes |
| `heartbeat_usecase.go` | Heartbeat recording and status retrieval; bot last_seen_at update; Kafka event publishing |

## Usecase Struct Pattern

```go
type XxxUsecase struct {
    xxxRepo  domain.XxxRepository
    // ... other repos as needed
    cache    *cache.Cache
    log      *logger.Logger
    producer *kafka.Producer   // may be nil
}

func NewXxxUsecase(
    xxxRepo domain.XxxRepository,
    cache   *cache.Cache,
    log     *logger.Logger,
    producer *kafka.Producer,
) *XxxUsecase {
    return &XxxUsecase{...}
}
```

> Inject only the repositories and infrastructure that *this* usecase actually needs.

## Rules

1. **Cache-aside pattern**: on reads, check cache first; on miss, query the repo, then populate the cache. On any write (create/update/delete), invalidate relevant cache keys. Always check `if uc.cache != nil` before cache operations — cache may be nil in tests.

2. **Log at usecase level** using `uc.logger.Error(...)` / `uc.logger.Info(...)` with `zap` fields. Pass structured error context, not raw error strings. Do not log in the repository.

3. **Return domain errors**, not raw database errors. When repo returns `nil, nil` for a missing record, return the appropriate sentinel error (e.g., `domain.ErrBotNotFound`). Handlers map these to HTTP status codes.

4. **Never call `c.BodyParser` or reference `fiber.Ctx`** inside a usecase. Usecases receive plain Go types and `context.Context`.

5. **Business rule validation** lives here (not in handlers or repos).

6. **Kafka event publishing** is done via `uc.producer.Publish(eventType, payload)`. Always check `if uc.producer != nil` before publishing — Kafka may be unavailable or nil in tests.

7. Keep usecase methods **focused and single-purpose**. Complex flows should be broken into private helper methods within the same file.

## Adding a New Usecase Method

1. Add the method to the appropriate `*_usecase.go` file.
2. If it requires a new repository operation, define the method in the `domain.XxxRepository` interface first.
3. Implement cache invalidation for all write paths.
4. Publish a Kafka event if the action is domain-significant.
5. Wire the method to a handler in `internal/handler/`.
