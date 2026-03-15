package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhducta/agent-core-service/internal/domain"
)

// SkillHandler handles /v1/me/skills routes
type SkillHandler struct {
	skillUC SkillUsecase
}

// NewSkillHandler creates a new SkillHandler
func NewSkillHandler(skillUC SkillUsecase) *SkillHandler {
	return &SkillHandler{skillUC: skillUC}
}

// ListSkills returns skills for the calling bot
// GET /v1/me/skills
func (h *SkillHandler) ListSkills(c *fiber.Ctx) error {
	botID, err := parseBotID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errResponse(domain.ErrCodeUnauthorized, "unauthorized"))
	}

	skills, err := h.skillUC.ListSkills(c.Context(), botID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse(domain.ErrCodeInternal, "failed to list skills"))
	}

	return c.JSON(fiber.Map{"data": skills})
}
