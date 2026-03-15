---
applyTo: "internal/usecase/**"
---

# Usecase Layer Instructions

Orchestrates repositories, cache (Redis), and Kafka event publishing.

## Rules

1. **Inject dependencies via constructor** — `NewXxxUsecase(repo, cache, kafka, logger)`.
2. **Never import `handler/`** — only `domain/`, `pkg/`, and `repository/` interfaces.
3. **Publish Kafka events** for significant state changes (bot online/offline, todo completed).
4. **Cache-aside pattern**: check Redis first → on miss, query DB → store in Redis.
5. **Return domain errors** (`domain.ErrCodeNotFound`, etc.) — handlers map to HTTP status.

## Testing

- Mock repository via `testify/mock` implementing domain interfaces.
- Mock Redis via `go-redis/redismock/v9`.
- Table-driven tests required.
