package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/minhducta/agent-core-service/internal/domain"
)

// BotResolver resolves a bot identity from a raw API key
type BotResolver interface {
	ResolveByAPIKey(ctx context.Context, rawKey string) (*domain.Bot, error)
}

// APIKeyAuth returns a middleware that validates Bearer API keys per-bot
func APIKeyAuth(resolver BotResolver) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    domain.ErrCodeUnauthorized,
					"message": "Authorization header is required",
				},
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    domain.ErrCodeUnauthorized,
					"message": "Invalid authorization header format. Expected: Bearer <api_key>",
				},
			})
		}

		rawKey := parts[1]
		bot, err := resolver.ResolveByAPIKey(c.Context(), rawKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    domain.ErrCodeUnauthorized,
					"message": "Invalid or inactive API key",
				},
			})
		}

		// Store bot_id in context locals for downstream handlers
		c.Locals("botId", bot.ID.String())
		return c.Next()
	}
}
