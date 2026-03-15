package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhducta/agent-core-service/internal/middleware"
)

// Router holds all handlers and registers routes
type Router struct {
	healthHandler    *HealthHandler
	botHandler       *BotHandler
	memoryHandler    *MemoryHandler
	skillHandler     *SkillHandler
	policyHandler    *PolicyHandler
	todoHandler      *TodoHandler
	heartbeatHandler *HeartbeatHandler
	botResolver      middleware.BotResolver
}

// NewRouter creates a new Router
func NewRouter(
	healthHandler *HealthHandler,
	botHandler *BotHandler,
	memoryHandler *MemoryHandler,
	skillHandler *SkillHandler,
	policyHandler *PolicyHandler,
	todoHandler *TodoHandler,
	heartbeatHandler *HeartbeatHandler,
	botResolver middleware.BotResolver,
) *Router {
	return &Router{
		healthHandler:    healthHandler,
		botHandler:       botHandler,
		memoryHandler:    memoryHandler,
		skillHandler:     skillHandler,
		policyHandler:    policyHandler,
		todoHandler:      todoHandler,
		heartbeatHandler: heartbeatHandler,
		botResolver:      botResolver,
	}
}

// Setup registers all routes on the fiber app
func (r *Router) Setup(app *fiber.App) {
	// Liveness / readiness probes
	app.Get("/health", r.healthHandler.HealthCheck)
	app.Get("/ready", r.healthHandler.Ready)

	// All v1 routes require a valid API key
	v1 := app.Group("/v1", middleware.APIKeyAuth(r.botResolver))

	// Bot identity routes
	me := v1.Group("/me")
	me.Get("", r.botHandler.GetProfile)
	me.Get("/identity", r.botHandler.GetIdentity)
	me.Get("/bootstrap", r.botHandler.GetBootstrap)
	me.Get("/memories", r.memoryHandler.ListMemories)
	me.Post("/memories", r.memoryHandler.CreateMemory)
	me.Delete("/memories/:id", r.memoryHandler.DeleteMemory)
	me.Get("/skills", r.skillHandler.ListSkills)
	me.Get("/policies", r.policyHandler.ListPolicies)

	// Todo routes
	todos := v1.Group("/todos")
	todos.Get("", r.todoHandler.ListTodos)
	todos.Patch("/:id", r.todoHandler.UpdateTodo)
	todos.Get("/:id/checklist", r.todoHandler.GetChecklist)
	todos.Patch("/:id/checklist/:item_id", r.todoHandler.UpdateChecklistItem)

	// Heartbeat routes
	heartbeat := v1.Group("/heartbeat")
	heartbeat.Post("", r.heartbeatHandler.RecordHeartbeat)
	heartbeat.Get("/status", r.heartbeatHandler.GetHeartbeatStatus)
}
