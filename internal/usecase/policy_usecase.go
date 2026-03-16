package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/minhducta/agent-core-service/internal/domain"
	"go.uber.org/zap"
)

// PolicyUsecase handles bot_policies business logic
type PolicyUsecase struct {
	policyRepo domain.BotPolicyRepository
	logger     *zap.Logger
}

// NewPolicyUsecase creates a new PolicyUsecase
func NewPolicyUsecase(policyRepo domain.BotPolicyRepository, logger *zap.Logger) *PolicyUsecase {
	return &PolicyUsecase{policyRepo: policyRepo, logger: logger}
}

// ListPolicies returns all policies for the calling bot
func (uc *PolicyUsecase) ListPolicies(ctx context.Context, botID uuid.UUID) ([]domain.BotPolicy, error) {
	return uc.policyRepo.ListByBotID(ctx, botID)
}
