package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/minhducta/agent-core-service/internal/domain"
)

// MemoryHandler handles /v1/me/memories routes
type MemoryHandler struct {
	memoryUC MemoryUsecase
}

// NewMemoryHandler creates a new MemoryHandler
func NewMemoryHandler(memoryUC MemoryUsecase) *MemoryHandler {
	return &MemoryHandler{memoryUC: memoryUC}
}

// ListMemories returns memories for the calling bot
// GET /v1/me/memories
func (h *MemoryHandler) ListMemories(c *fiber.Ctx) error {
	botID, err := parseBotID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errResponse(domain.ErrCodeUnauthorized, "unauthorized"))
	}

	memories, err := h.memoryUC.ListMemories(c.Context(), botID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse(domain.ErrCodeInternal, "failed to list memories"))
	}

	return c.JSON(fiber.Map{"data": memories})
}

// CreateMemory creates a new memory
// POST /v1/me/memories
func (h *MemoryHandler) CreateMemory(c *fiber.Ctx) error {
	botID, err := parseBotID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errResponse(domain.ErrCodeUnauthorized, "unauthorized"))
	}

	var req domain.CreateMemoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errResponse(domain.ErrCodeValidation, "invalid request body"))
	}

	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(errResponse(domain.ErrCodeValidation, "content is required"))
	}

	if req.Tags == nil {
		req.Tags = []string{}
	}

	memory, err := h.memoryUC.CreateMemory(c.Context(), botID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse(domain.ErrCodeInternal, "failed to create memory"))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": memory})
}

// DeleteMemory removes a memory
// DELETE /v1/me/memories/:id
func (h *MemoryHandler) DeleteMemory(c *fiber.Ctx) error {
	botID, err := parseBotID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errResponse(domain.ErrCodeUnauthorized, "unauthorized"))
	}

	memoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errResponse(domain.ErrCodeValidation, "invalid memory id"))
	}

	if err := h.memoryUC.DeleteMemory(c.Context(), botID, memoryID); err != nil {
		if errors.Is(err, domain.ErrMemoryNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(errResponse(domain.ErrCodeNotFound, "memory not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse(domain.ErrCodeInternal, "failed to delete memory"))
	}

	return c.SendStatus(fiber.StatusNoContent)
}
