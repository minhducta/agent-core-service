package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/minhducta/agent-core-service/internal/domain"
	"go.uber.org/zap"
)

// SkillUsecase handles bot_skills business logic
type SkillUsecase struct {
	skillRepo domain.BotSkillRepository
	logger    *zap.Logger
}

// NewSkillUsecase creates a new SkillUsecase
func NewSkillUsecase(skillRepo domain.BotSkillRepository, logger *zap.Logger) *SkillUsecase {
	return &SkillUsecase{skillRepo: skillRepo, logger: logger}
}

// ListSkills returns all skills for the calling bot
func (uc *SkillUsecase) ListSkills(ctx context.Context, botID uuid.UUID) ([]domain.BotSkill, error) {
	return uc.skillRepo.ListByBotID(ctx, botID)
}
