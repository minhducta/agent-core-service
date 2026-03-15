package middleware

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/minhducta/agent-core-service/internal/domain"
	"github.com/minhducta/agent-core-service/internal/usecase"
)

// APIKeyAuth returns a middleware that validates Bearer API keys per-bot
func APIKeyAuth(botUC *usecase.BotUsecase) fiber.Handler {
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
		bot, err := botUC.ResolveByAPIKey(c.Context(), rawKey)
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

// hashAPIKey computes SHA-256 hex of the raw key (shared with usecase)
func hashAPIKey(rawKey string) string {
	h := sha256.Sum256([]byte(rawKey))
	return fmt.Sprintf("%x", h)
}
