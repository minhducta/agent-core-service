package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/minhducta/agent-core-service/internal/domain"
)

// BotHandler handles /v1/me routes
type BotHandler struct {
	botUC BotUsecase
}

// NewBotHandler creates a new BotHandler
func NewBotHandler(botUC BotUsecase) *BotHandler {
	return &BotHandler{botUC: botUC}
}

// GetProfile returns the bot profile + ref_links
// GET /v1/me
func (h *BotHandler) GetProfile(c *fiber.Ctx) error {
	botID, err := parseBotID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errResponse(domain.ErrCodeUnauthorized, "unauthorized"))
	}

	resp, err := h.botUC.GetProfile(c.Context(), botID)
	if err != nil {
		if errors.Is(err, domain.ErrBotNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(errResponse(domain.ErrCodeNotFound, err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse(domain.ErrCodeInternal, "internal error"))
	}

	return c.JSON(fiber.Map{"data": resp})
}

// GetIdentity returns the minimal bot identity
// GET /v1/me/identity
func (h *BotHandler) GetIdentity(c *fiber.Ctx) error {
	botID, err := parseBotID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errResponse(domain.ErrCodeUnauthorized, "unauthorized"))
	}

	resp, err := h.botUC.GetIdentity(c.Context(), botID)
	if err != nil {
		if errors.Is(err, domain.ErrBotNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(errResponse(domain.ErrCodeNotFound, err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse(domain.ErrCodeInternal, "internal error"))
	}

	return c.JSON(fiber.Map{"data": resp})
}

// GetBootstrap returns the full context dump
// GET /v1/me/bootstrap
func (h *BotHandler) GetBootstrap(c *fiber.Ctx) error {
	botID, err := parseBotID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errResponse(domain.ErrCodeUnauthorized, "unauthorized"))
	}

	resp, err := h.botUC.GetBootstrap(c.Context(), botID)
	if err != nil {
		if errors.Is(err, domain.ErrBotNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(errResponse(domain.ErrCodeNotFound, err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse(domain.ErrCodeInternal, "internal error"))
	}

	return c.JSON(fiber.Map{"data": resp})
}

