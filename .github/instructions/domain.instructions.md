---
applyTo: "internal/domain/**"
---

# Domain Layer Instructions

The `internal/domain` package is the **innermost layer** of Clean Architecture. It has **zero dependencies on other internal packages** and must never import from `handler`, `usecase`, `repository`, or `pkg/`.

## File Responsibilities

| File | Contents |
|---|---|
| `bot.go` | Entity: `Bot` (ID, Name, Role, Vibe, Emoji, AvatarURL, APIKeyHash, LastSeenAt, Status); DTOs: `BotResponse`, `BootstrapResponse`; enum: `BotStatus` |
| `memory.go` | Entity: `BotMemory` (ID, BotID, Type, Content, Tags, Importance, ExpiresAt); DTOs: `CreateMemoryRequest`, `MemoryResponse` |
| `todo.go` | Entity: `Todo`, `TodoChecklistItem`; DTOs; enums: `TodoStatus`, `TodoPriority` |
| `heartbeat.go` | Entity: `Heartbeat`; DTOs: `HeartbeatRequest`, `HeartbeatResponse`, `PendingCommand` |
| `common.go` | Shared types: `PaginationMeta`, `RefLinks` |
| `repository.go` | Repository interfaces: `BotRepository`, `MemoryRepository`, `SkillRepository`, `PolicyRepository`, `TodoRepository`, `HeartbeatRepository` |
| `error.go` | Error code constants: `ErrCodeValidation`, `ErrCodeNotFound`, `ErrCodeUnauthorized`, `ErrCodeForbidden`, `ErrCodeInternal`, `ErrCodeConflict` |
| `event.go` | Kafka event constants: `EventBotOnline`, `EventBotOffline`, `EventBotDegraded`, `EventTodoCompleted`, `EventMemoryCreated` |

## Rules

1. **Entities use `db:` and `json:` struct tags** — `db:` for `sqlx`, `json:` for HTTP. Sensitive fields use `json:"-"`.
2. **All primary keys are `uuid.UUID`** — never `int` or plain `string`.
3. **Enums are typed string constants** — never raw strings in business logic.
4. **Repository interfaces defined here**, not in `repository/`.
5. **No business logic** in domain layer.
6. Mark sensitive fields (`APIKeyHash`) with `json:"-"`.
