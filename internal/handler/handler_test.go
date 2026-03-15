package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewHealthHandler(t *testing.T) {
	h := NewHealthHandler(nil, nil)
	assert.NotNil(t, h)
}

func TestNewBotHandler(t *testing.T) {
	h := NewBotHandler(nil)
	assert.NotNil(t, h)
}

func TestNewMemoryHandler(t *testing.T) {
	h := NewMemoryHandler(nil)
	assert.NotNil(t, h)
}

func TestNewTodoHandler(t *testing.T) {
	h := NewTodoHandler(nil)
	assert.NotNil(t, h)
}

func TestNewHeartbeatHandler(t *testing.T) {
	h := NewHeartbeatHandler(nil)
	assert.NotNil(t, h)
}

// TestHealthCheck_NilDeps ensures /health returns 200 when no real deps are wired
func TestHealthCheck_NilDeps(t *testing.T) {
	app := fiber.New()
	h := NewHealthHandler(nil, nil)
	app.Get("/health", h.HealthCheck)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

// TestReady_NilDeps ensures /ready returns 200 when db is nil
func TestReady_NilDeps(t *testing.T) {
	app := fiber.New()
	h := NewHealthHandler(nil, nil)
	app.Get("/ready", h.Ready)

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

// TestListMemories_Unauthorized ensures missing botId Local → 401
func TestListMemories_Unauthorized(t *testing.T) {
	app := fiber.New()
	h := NewMemoryHandler(nil)
	app.Get("/memories", h.ListMemories)

	req := httptest.NewRequest(http.MethodGet, "/memories", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// TestCreateMemory_InvalidBody ensures body parse error → 400
func TestCreateMemory_InvalidBody(t *testing.T) {
	app := fiber.New()
	h := NewMemoryHandler(nil)
	// Inject a fake botId to pass auth check
	app.Post("/memories", func(c *fiber.Ctx) error {
		c.Locals("botId", "550e8400-e29b-41d4-a716-446655440000")
		return c.Next()
	}, h.CreateMemory)

	req := httptest.NewRequest(http.MethodPost, "/memories", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// TestUpdateTodo_InvalidID ensures bad UUID → 400
func TestUpdateTodo_InvalidID(t *testing.T) {
	app := fiber.New()
	h := NewTodoHandler(nil)
	app.Patch("/todos/:id", func(c *fiber.Ctx) error {
		c.Locals("botId", "550e8400-e29b-41d4-a716-446655440000")
		return c.Next()
	}, h.UpdateTodo)

	req := httptest.NewRequest(http.MethodPatch, "/todos/not-a-uuid", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// TestRecordHeartbeat_Unauthorized ensures missing botId → 401
func TestRecordHeartbeat_Unauthorized(t *testing.T) {
	app := fiber.New()
	h := NewHeartbeatHandler(nil)
	app.Post("/heartbeat", h.RecordHeartbeat)

	req := httptest.NewRequest(http.MethodPost, "/heartbeat", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}
