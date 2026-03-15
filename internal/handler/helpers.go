package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// parseBotID extracts bot_id from fiber Locals (set by APIKeyAuth middleware)
func parseBotID(c *fiber.Ctx) (uuid.UUID, error) {
	val, ok := c.Locals("botId").(string)
	if !ok || val == "" {
		return uuid.Nil, fiber.ErrUnauthorized
	}
	return uuid.Parse(val)
}

// parsePathUUID parses a UUID path parameter by key
func parsePathUUID(c *fiber.Ctx, key string) (uuid.UUID, error) {
	return uuid.Parse(c.Params(key))
}

// errResponse builds a standard error response body
func errResponse(code, message string) fiber.Map {
	return fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": message,
		},
	}
}
