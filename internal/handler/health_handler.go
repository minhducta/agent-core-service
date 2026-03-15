package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhducta/agent-core-service/pkg/cache"
	"github.com/minhducta/agent-core-service/pkg/database"
)

// HealthHandler handles liveness/readiness probes
type HealthHandler struct {
	db    *database.DB
	cache *cache.Cache
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(db *database.DB, cache *cache.Cache) *HealthHandler {
	return &HealthHandler{db: db, cache: cache}
}

// HealthCheck checks DB and Redis connectivity
// GET /health
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	dbStatus := "ok"
	if h.db != nil {
		if err := h.db.HealthCheck(); err != nil {
			dbStatus = "error: " + err.Error()
		}
	}

	redisStatus := "ok"
	if h.cache != nil {
		if err := h.cache.HealthCheck(c.Context()); err != nil {
			redisStatus = "error: " + err.Error()
		}
	}

	httpStatus := fiber.StatusOK
	if dbStatus != "ok" || redisStatus != "ok" {
		httpStatus = fiber.StatusServiceUnavailable
	}

	return c.Status(httpStatus).JSON(fiber.Map{
		"service":  "agent-core-service",
		"status":   "running",
		"database": dbStatus,
		"redis":    redisStatus,
	})
}

// Ready checks if the service is ready to accept traffic
// GET /ready
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	if h.db != nil {
		if err := h.db.HealthCheck(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not ready",
				"error":  err.Error(),
			})
		}
	}

	return c.JSON(fiber.Map{"status": "ready"})
}
