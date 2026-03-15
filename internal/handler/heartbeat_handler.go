package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/minhducta/agent-core-service/internal/domain"
	"github.com/minhducta/agent-core-service/internal/usecase"
)

// HeartbeatHandler handles /v1/heartbeat routes
type HeartbeatHandler struct {
	heartbeatUC *usecase.HeartbeatUsecase
}

// NewHeartbeatHandler creates a new HeartbeatHandler
func NewHeartbeatHandler(heartbeatUC *usecase.HeartbeatUsecase) *HeartbeatHandler {
	return &HeartbeatHandler{heartbeatUC: heartbeatUC}
}

// RecordHeartbeat processes a heartbeat from the calling bot
// POST /v1/heartbeat
func (h *HeartbeatHandler) RecordHeartbeat(c *fiber.Ctx) error {
	botID, err := parseBotID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errResponse(domain.ErrCodeUnauthorized, "unauthorized"))
	}

	var req domain.HeartbeatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errResponse(domain.ErrCodeValidation, "invalid request body"))
	}

	if req.Status == "" {
		req.Status = domain.HeartbeatStatusOK
	}

	hb, err := h.heartbeatUC.RecordHeartbeat(c.Context(), botID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse(domain.ErrCodeInternal, "failed to record heartbeat"))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": hb})
}

// GetHeartbeatStatus returns the latest heartbeat status
// GET /v1/heartbeat/status
func (h *HeartbeatHandler) GetHeartbeatStatus(c *fiber.Ctx) error {
	botID, err := parseBotID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errResponse(domain.ErrCodeUnauthorized, "unauthorized"))
	}

	resp, err := h.heartbeatUC.GetStatus(c.Context(), botID)
	if err != nil {
		if errors.Is(err, domain.ErrHeartbeatNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(errResponse(domain.ErrCodeNotFound, "no heartbeat recorded"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse(domain.ErrCodeInternal, "failed to get heartbeat status"))
	}

	return c.JSON(fiber.Map{"data": resp})
}
