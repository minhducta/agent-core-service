---
applyTo: "cmd/api/main.go"
---

# Entrypoint (main.go) Instructions

`cmd/api/main.go` is the **composition root** of the service. Its only job is to wire all dependencies together and start the server. It must not contain business logic.

## Startup Sequence (must follow this order)

```
1. config.Load(path)                          → cfg
2. logger.NewLogger(level)                    → log  (defer log.Sync())
3. database.New(cfg.Database)                 → db   (defer db.Close())
4. migration.Run(cfg.Database, cfg.Migration.Path) — auto-migrate on startup
5. cache.New(cfg.Redis, cfg.Cache)            → redisCache
6. kafka.NewProducer(cfg.Kafka)               → producer  (defer producer.Close())
7. repository.NewXxxRepository(db)  × N       → repos
8. usecase.NewXxxUsecase(repo, cache, log, producer, ...) × N → usecases
9. handler.NewXxxHandler(usecase) × N          → handlers
10. handler.NewRouter(handlers..., botResolver) → router
11. fiber.New(config) + middleware (CORS, Recovery, RequestLogger)
12. router.Setup(app)
13. Graceful shutdown via os.Signal (SIGTERM / SIGINT)
```

## Rules

1. **Manual DI only** — no DI framework. Each constructor is called explicitly in the order above.
2. When adding a new resource, add the repo → usecase → handler instantiation in their respective sections, following the existing pattern.
3. **Graceful shutdown**: wait for in-flight requests to complete before exiting. Use `app.ShutdownWithContext` and the signal handling already present.
4. The `BotResolver` (implemented by `BotUsecase`) is passed to `NewRouter` for use in the auth middleware — never pass a pre-built middleware function.
5. Always `defer db.Close()`, `defer log.Sync()`, and `defer producer.Close()` immediately after initialisation.
6. Log a startup banner with `log.Info(...)` including the server address.
7. **Never** call `os.Exit` directly — let the signal handler trigger `app.Shutdown`.
8. Infrastructure clients (Kafka producer) may fail to initialise — decide per-client whether to fail fast or log a warning and continue with `nil`.
