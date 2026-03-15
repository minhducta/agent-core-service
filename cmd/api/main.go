package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minhducta/agent-core-service/internal/handler"
	"github.com/minhducta/agent-core-service/internal/middleware"
	"github.com/minhducta/agent-core-service/internal/repository"
	"github.com/minhducta/agent-core-service/internal/usecase"
	"github.com/minhducta/agent-core-service/pkg/cache"
	"github.com/minhducta/agent-core-service/pkg/config"
	"github.com/minhducta/agent-core-service/pkg/database"
	"github.com/minhducta/agent-core-service/pkg/kafka"
	"github.com/minhducta/agent-core-service/pkg/logger"
	"github.com/minhducta/agent-core-service/pkg/migration"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize logger
	appLogger := logger.NewLogger(cfg.Logger.Level)
	defer appLogger.Sync()

	// Initialize database
	db, err := database.New(cfg.Database)
	if err != nil {
		appLogger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Run migrations
	if cfg.Migration.Enabled {
		applied, err := migration.Run(cfg.Database, cfg.Migration.Path)
		if err != nil {
			appLogger.Fatal("failed to run migrations", zap.Error(err))
		}
		appLogger.Info("migrations applied", zap.Int("count", applied))
	}

	// Initialize cache
	appCache, err := cache.New(cfg.Redis, cfg.Cache)
	if err != nil {
		appLogger.Fatal("failed to connect to Redis", zap.Error(err))
	}

	// Initialize Kafka producer
	producer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		appLogger.Fatal("failed to create Kafka producer", zap.Error(err))
	}
	defer producer.Close()

	// Initialize repositories
	botRepo := repository.NewBotRepository(db)
	memoryRepo := repository.NewBotMemoryRepository(db)
	skillRepo := repository.NewBotSkillRepository(db)
	policyRepo := repository.NewBotPolicyRepository(db)
	todoRepo := repository.NewTodoRepository(db)
	checklistRepo := repository.NewTodoChecklistRepository(db)
	heartbeatRepo := repository.NewHeartbeatRepository(db)

	// Initialize usecases
	botUC := usecase.NewBotUsecase(
		botRepo, memoryRepo, skillRepo, policyRepo,
		appCache, producer, appLogger.Logger,
	)
	memoryUC := usecase.NewMemoryUsecase(memoryRepo, appCache, producer, appLogger.Logger)
	skillUC := usecase.NewSkillUsecase(skillRepo, appLogger.Logger)
	policyUC := usecase.NewPolicyUsecase(policyRepo, appLogger.Logger)
	todoUC := usecase.NewTodoUsecase(todoRepo, checklistRepo, producer, appLogger.Logger)
	heartbeatUC := usecase.NewHeartbeatUsecase(heartbeatRepo, botRepo, producer, appLogger.Logger)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler(db, appCache)
	botHandler := handler.NewBotHandler(botUC)
	memoryHandler := handler.NewMemoryHandler(memoryUC)
	skillHandler := handler.NewSkillHandler(skillUC)
	policyHandler := handler.NewPolicyHandler(policyUC)
	todoHandler := handler.NewTodoHandler(todoUC)
	heartbeatHandler := handler.NewHeartbeatHandler(heartbeatUC)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		Prefork:      cfg.Server.Prefork,
		AppName:      "agent-core-service",
	})

	// Apply global middleware
	app.Use(middleware.CORS())
	app.Use(middleware.Recovery())
	app.Use(middleware.RequestLogger(appLogger.Logger))

	// Setup routes
	router := handler.NewRouter(
		healthHandler,
		botHandler,
		memoryHandler,
		skillHandler,
		policyHandler,
		todoHandler,
		heartbeatHandler,
		botUC,
	)
	router.Setup(app)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		addr := cfg.Server.Address()
		appLogger.Info("starting agent-core-service", zap.String("addr", addr))
		if err := app.Listen(addr); err != nil {
			appLogger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	<-quit
	appLogger.Info("shutting down agent-core-service...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		appLogger.Fatal("server forced to shutdown", zap.Error(err))
	}

	appLogger.Info("server stopped gracefully")
}
