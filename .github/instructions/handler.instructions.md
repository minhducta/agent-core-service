---
applyTo: "internal/handler/**"
---

# Handler Layer Instructions

The `internal/handler` package is the **delivery layer**. Handlers only parse HTTP input, delegate to a usecase, and return a JSON response. They contain **no business logic**.

## File Responsibilities

| File | Contents |
|---|---|
| `router.go` | `Router` struct holding all handlers + `BotResolver`; `Setup(*fiber.App)` registers every route under `/v1/` with public and protected (API Key) groups |
| `bot_handler.go` | `GetProfile`, `GetIdentity`, `GetBootstrap` |
| `memory_handler.go` | `ListMemories`, `CreateMemory`, `DeleteMemory` |
| `skill_handler.go` | `ListSkills` |
| `policy_handler.go` | `ListPolicies` |
| `todo_handler.go` | `ListTodos`, `UpdateTodo`, `GetChecklist`, `UpdateChecklistItem` |
| `heartbeat_handler.go` | `RecordHeartbeat`, `GetHeartbeatStatus` |
| `health_handler.go` | `GET /health` and `GET /ready` probes |
| `helpers.go` | Shared utilities: `parseBotID(c *fiber.Ctx)`, `parsePathUUID(c *fiber.Ctx, key string)`, `errResponse(code, message string)` |
| `interfaces.go` | Usecase interfaces consumed by handlers: `BotUsecase`, `MemoryUsecase`, `SkillUsecase`, `PolicyUsecase`, `TodoUsecase`, `HeartbeatUsecase` |

## Handler Struct Pattern

```go
type XxxHandler struct {
    xxxUC XxxUsecase
}

func NewXxxHandler(xxxUC XxxUsecase) *XxxHandler {
    return &XxxHandler{xxxUC: xxxUC}
}
```

## Request Handling Rules

1. **Parse body** with `c.BodyParser(&req)`. On error, return `400` with `domain.ErrCodeValidation`.
2. **Parse UUID path params** with `parsePathUUID(c, "id")`. On error, return `400`.
3. **Extract `botID`** from `c.Locals("botId")` via `parseBotID(c)` — set by the API Key auth middleware for protected routes.
4. **Never** read `botID` from the request body — always from `Locals`.
5. Call **exactly one usecase method** per handler. Do not call repo or cache directly.
6. Map usecase errors to HTTP status codes:
   - `ErrCodeNotFound` → `404`
   - `ErrCodeValidation` / `ErrCodeConflict` → `400` / `409`
   - `ErrCodeUnauthorized` → `401`
   - `ErrCodeForbidden` → `403`
   - `ErrCodeInternal` → `500`
7. **Success response shape:** `{"data": <entity>}` for single items; `{"data": [...], "meta": {...}}` for paginated lists.
8. **Error response shape:** `{"error": {"code": "...", "message": "..."}}`.
9. Never log inside a handler — logging is the usecase's responsibility.

## Route Registration (router.go)

- **Public** routes: `/health`, `/ready`.
- **Protected** routes (wrapped in `middleware.APIKeyAuth(botResolver)`): all `/v1/*` endpoints.
- Every new handler must be:
  1. Added as a field on `Router`.
  2. Accepted as a parameter in `NewRouter(...)`.
  3. Wired in `Setup()`.

## Adding a New Handler Method

1. Add the method to the handler file for the relevant resource.
2. Register the route in `router.go` under the appropriate group.
3. Wire the handler in `cmd/api/main.go` (pass to `NewRouter`).
