package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhducta/agent-core-service/internal/domain"
)

// PolicyHandler handles /v1/me/policies routes
type PolicyHandler struct {
	policyUC PolicyUsecase
}

// NewPolicyHandler creates a new PolicyHandler
func NewPolicyHandler(policyUC PolicyUsecase) *PolicyHandler {
	return &PolicyHandler{policyUC: policyUC}
}

// ListPolicies returns policies for the calling bot
// GET /v1/me/policies
func (h *PolicyHandler) ListPolicies(c *fiber.Ctx) error {
	botID, err := parseBotID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errResponse(domain.ErrCodeUnauthorized, "unauthorized"))
	}

	policies, err := h.policyUC.ListPolicies(c.Context(), botID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse(domain.ErrCodeInternal, "failed to list policies"))
	}

	return c.JSON(fiber.Map{"data": policies})
}
