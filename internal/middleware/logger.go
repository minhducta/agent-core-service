package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// RequestLogger returns a middleware that logs every request using Zap.
// It is nil-safe — if the logger is nil, it skips logging and calls c.Next().
func RequestLogger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if logger == nil {
			return c.Next()
		}

		start := time.Now()

		if err := c.Next(); err != nil {
			return err
		}

		logger.Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", c.IP()),
		)

		return nil
	}
}
