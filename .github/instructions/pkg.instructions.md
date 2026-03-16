---
applyTo: "pkg/**"
---

# pkg Layer Instructions

The `pkg/` directory contains **infrastructure helpers** that are injected into internal layers from `cmd/api/main.go`. They have no knowledge of business domain concepts.

## Package Responsibilities

| Package | File | Role |
|---|---|---|
| `pkg/config` | `config.go` | Loads `config/config.yaml` via `spf13/viper`. Exposes typed structs: `Config`, `ServerConfig`, `DatabaseConfig`, `RedisConfig`, `CacheConfig`, `KafkaConfig`, `LoggerConfig`, `PaginationConfig`, `MigrationConfig` |
| `pkg/database` | `postgres.go` | Opens and returns a `*database.DB` (wrapping `*sqlx.DB`) with `New()`, `Close()`, `HealthCheck()`. Configures `MaxOpenConns`, `MaxIdleConns`, `ConnMaxLifetime` |
| `pkg/cache` | `redis.go` | Wraps `go-redis/v9`. Provides a `*Cache` struct with typed helpers: `SetBot/GetBot/DeleteBot`, `SetBotByAPIKeyHash/GetBotByAPIKeyHash/DeleteBotByAPIKeyHash`, `SetMemories/GetMemories/InvalidateMemories`. `NewWithClient()` for test injection |
| `pkg/logger` | `logger.go` | Wraps `go.uber.org/zap`. Exposes a `*Logger` struct with `Info`, `Error`, `Warn`, `Debug`, `Fatal`. Provides `RequestLogger` middleware for Fiber (nil-safe) |
| `pkg/migration` | `migration.go` | Wraps `golang-migrate/migrate/v4`. Provides `Run()` to apply pending up-migrations and `Version()` to query current state |
| `pkg/kafka` | `producer.go` | Wraps `IBM/sarama` `SyncProducer`. Provides `Publish(topic string, payload any)` and `Close()`. JSON-serialises the payload |

## Rules for `pkg/config`

- Expose only **typed structs** — never return raw `viper.Get()` values to callers.
- Add new config sections as new nested structs.
- All config fields must have `mapstructure` tags to work with viper's unmarshalling.
- The `ServerConfig` struct must have an `Address() string` helper that returns `"host:port"`.

## Rules for `pkg/cache`

- Every domain entity that is cached must have dedicated typed helpers: `GetXxx(ctx, key, dest)`, `SetXxx(ctx, key, value)`, `DeleteXxx(ctx, key)`.
- **Cache key format**: `agent:<entity>:<identifier>` (e.g., `agent:bot:<uuid>`, `agent:apikey:<hash>`).
- TTLs must come from `CacheConfig` — never hardcode durations.
- All cache errors must be **silently swallowed** in the usecase (cache miss is not an application error).
- Use `json.Marshal` / `json.Unmarshal` for serialisation into Redis strings.
- Provide `NewWithClient()` constructor for unit tests to inject a mock Redis client.

## Rules for `pkg/database`

- Wrap `*sqlx.DB` in a `DB` struct that exposes `GetContext`, `SelectContext`, `ExecContext`, `BeginTxx`, and health check methods.
- Apply connection pool settings from config immediately after `sqlx.Open`.
- Ping the database in `New()` and return an error if it fails — fail fast at startup.

## Rules for `pkg/logger`

- Expose `Info`, `Error`, `Warn`, `Debug`, `Fatal` methods that delegate to the underlying `*zap.Logger`.
- Always use **structured fields** (`zap.String`, `zap.Error`, `zap.Int`, etc.) — never use `fmt.Sprintf` inside log calls.
- Expose `Sync() error` to flush buffered log entries (called via `defer log.Sync()` in `main.go`).
- `RequestLogger` must be nil-safe — if the logger is nil, skip logging and call `c.Next()`.

## Rules for `pkg/kafka`

- `Publish` serialises the payload to JSON and sends a `sarama.ProducerMessage`.
- The caller is responsible for checking `if producer != nil` before calling — Kafka is an optional dependency.
- `Close()` must be called on shutdown (deferred in `main.go`).

## Adding a New Infrastructure Helper

1. Create a new file inside the relevant sub-package (e.g., `pkg/email/smtp.go`).
2. Accept configuration via a dedicated config struct added to `pkg/config/config.go`.
3. Initialise the client in `cmd/api/main.go` and inject it wherever needed.
4. Do **not** import from `internal/` — `pkg/` must remain domain-agnostic.
