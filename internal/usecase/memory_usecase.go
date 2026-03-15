package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/minhducta/agent-core-service/internal/domain"
	"github.com/minhducta/agent-core-service/pkg/cache"
	"github.com/minhducta/agent-core-service/pkg/kafka"
	"go.uber.org/zap"
)

// MemoryUsecase handles bot_memories business logic
type MemoryUsecase struct {
	memoryRepo domain.BotMemoryRepository
	cache      *cache.Cache
	producer   *kafka.Producer
	logger     *zap.Logger
}

// NewMemoryUsecase creates a new MemoryUsecase
func NewMemoryUsecase(
	memoryRepo domain.BotMemoryRepository,
	appCache *cache.Cache,
	producer *kafka.Producer,
	logger *zap.Logger,
) *MemoryUsecase {
	return &MemoryUsecase{
		memoryRepo: memoryRepo,
		cache:      appCache,
		producer:   producer,
		logger:     logger,
	}
}

// ListMemories returns all non-expired memories for the calling bot
func (uc *MemoryUsecase) ListMemories(ctx context.Context, botID uuid.UUID) ([]domain.BotMemory, error) {
	memories, err := uc.memoryRepo.ListByBotID(ctx, botID)
	if err != nil {
		return nil, fmt.Errorf("failed to list memories: %w", err)
	}

	return memories, nil
}

// CreateMemory creates a new memory for the calling bot
func (uc *MemoryUsecase) CreateMemory(ctx context.Context, botID uuid.UUID, req domain.CreateMemoryRequest) (*domain.BotMemory, error) {
	memory := &domain.BotMemory{
		ID:         uuid.New(),
		BotID:      botID,
		Type:       req.Type,
		Content:    req.Content,
		Tags:       req.Tags,
		Importance: req.Importance,
		ExpiresAt:  req.ExpiresAt,
	}

	if err := uc.memoryRepo.Create(ctx, memory); err != nil {
		return nil, fmt.Errorf("failed to create memory: %w", err)
	}

	_ = uc.cache.InvalidateMemories(ctx, botID.String())

	_ = uc.producer.Publish(domain.EventMemoryCreated, map[string]interface{}{
		"botId":    botID,
		"memoryId": memory.ID,
		"type":     memory.Type,
	})

	uc.logger.Info("memory created", zap.String("botId", botID.String()), zap.String("memoryId", memory.ID.String()))

	return memory, nil
}

// DeleteMemory removes a memory, scoped to the calling bot
func (uc *MemoryUsecase) DeleteMemory(ctx context.Context, botID uuid.UUID, memoryID uuid.UUID) error {
	if err := uc.memoryRepo.Delete(ctx, memoryID, botID); err != nil {
		return fmt.Errorf("failed to delete memory: %w", err)
	}

	_ = uc.cache.InvalidateMemories(ctx, botID.String())

	_ = uc.producer.Publish(domain.EventMemoryDeleted, map[string]interface{}{
		"botId":    botID,
		"memoryId": memoryID,
	})

	uc.logger.Info("memory deleted", zap.String("botId", botID.String()), zap.String("memoryId", memoryID.String()))

	return nil
}

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
