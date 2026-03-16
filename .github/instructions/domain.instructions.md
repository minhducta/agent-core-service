---
applyTo: "internal/domain/**"
---

# Domain Layer Instructions

The `internal/domain` package is the **innermost layer** of Clean Architecture. It has **zero dependencies on other internal packages** and must never import from `handler`, `usecase`, `repository`, or `pkg/`.

## File Responsibilities

| File | Contents |
|---|---|
| `bot.go` | Entity: `Bot` (ID, Name, Role, Vibe, Emoji, AvatarURL, APIKeyHash, LastSeenAt, Status); DTOs: `BotResponse`, `BotProfileResponse`, `BotIdentityResponse`, `BootstrapResponse`; enum: `BotStatus` |
| `memory.go` | Entity: `BotMemory` (ID, BotID, Type, Content, Tags, Importance, ExpiresAt); DTOs: `CreateMemoryRequest`; enum: `MemoryType` |
| `skill.go` | Entity: `BotSkill` (ID, BotID, Name, Description, UsageGuide) |
| `policy.go` | Entity: `BotPolicy` (ID, BotID, Action, Effect, Conditions); enum: `PolicyEffect` |
| `todo.go` | Entity: `Todo`, `TodoChecklistItem`; DTOs: `UpdateTodoRequest`, `UpdateChecklistItemRequest`, `TodoListResponse`; enums: `TodoStatus`, `TodoPriority` |
| `heartbeat.go` | Entity: `Heartbeat`; DTOs: `HeartbeatRequest`, `HeartbeatStatusResponse`; enum: `HeartbeatStatus` |
| `common.go` | Shared types: `PaginationMeta`, `RefLinks` |
| `repository.go` | Repository interfaces: `BotRepository`, `BotMemoryRepository`, `BotSkillRepository`, `BotPolicyRepository`, `TodoRepository`, `TodoChecklistRepository`, `HeartbeatRepository` |
| `error.go` | Error code constants: `ErrCodeValidation`, `ErrCodeNotFound`, `ErrCodeUnauthorized`, `ErrCodeForbidden`, `ErrCodeInternal`, `ErrCodeConflict`; sentinel errors |
| `event.go` | Kafka event constants: `EventBotOnline`, `EventBotOffline`, `EventBotDegraded`, `EventTodoCompleted`, `EventMemoryCreated`, etc. |

## Rules

1. **Entities use `db:` and `json:` struct tags** — `db:` for `sqlx`, `json:` for HTTP. Sensitive fields use `json:"-"`.
2. **All primary keys are `uuid.UUID`** — never `int` or plain `string`.
3. **Enums are typed string constants** — never raw strings in business logic.
4. **Repository interfaces defined here**, not in `repository/`.
5. **Request/response DTOs** are defined alongside their entity.
6. **No business logic** in the domain layer. DTOs may have simple helper methods (e.g., `ToResponse()`), but orchestration belongs in `usecase/`.
7. Mark sensitive fields (`APIKeyHash`) with `json:"-"`.
8. **Kafka event types** are string constants defined in `event.go` — they represent domain events published to the message broker.

## Adding a New Entity

1. Define the struct with `db:` and `json:` tags (mark sensitive fields `json:"-"`).
2. Add enums (status, type, etc.) as typed string constants.
3. Add request/response DTOs.
4. Add the repository interface to `repository.go`.
