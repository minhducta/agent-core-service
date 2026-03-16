package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/minhducta/agent-core-service/internal/domain"
	"github.com/minhducta/agent-core-service/pkg/cache"
	"github.com/minhducta/agent-core-service/pkg/kafka"
	"go.uber.org/zap"
)

// BotUsecase handles bot identity business logic
type BotUsecase struct {
	botRepo    domain.BotRepository
	memoryRepo domain.BotMemoryRepository
	skillRepo  domain.BotSkillRepository
	policyRepo domain.BotPolicyRepository
	cache      *cache.Cache
	producer   *kafka.Producer
	logger     *zap.Logger
}

// NewBotUsecase creates a new BotUsecase
func NewBotUsecase(
	botRepo domain.BotRepository,
	memoryRepo domain.BotMemoryRepository,
	skillRepo domain.BotSkillRepository,
	policyRepo domain.BotPolicyRepository,
	appCache *cache.Cache,
	producer *kafka.Producer,
	logger *zap.Logger,
) *BotUsecase {
	return &BotUsecase{
		botRepo:    botRepo,
		memoryRepo: memoryRepo,
		skillRepo:  skillRepo,
		policyRepo: policyRepo,
		cache:      appCache,
		producer:   producer,
		logger:     logger,
	}
}

// ResolveByAPIKey validates the Bearer token and returns the bot
func (uc *BotUsecase) ResolveByAPIKey(ctx context.Context, rawKey string) (*domain.Bot, error) {
	hash := hashAPIKey(rawKey)

	// Check cache first
	if uc.cache != nil {
		var cachedBotID string
		if err := uc.cache.Get(ctx, fmt.Sprintf("agent:apikey:%s", hash), &cachedBotID); err == nil && cachedBotID != "" {
			id, err := uuid.Parse(cachedBotID)
			if err == nil {
				bot, err := uc.botRepo.GetByID(ctx, id)
				if err == nil && bot != nil && bot.Status == domain.BotStatusActive {
					return bot, nil
				}
			}
		}
	}

	bot, err := uc.botRepo.GetByAPIKeyHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve api key: %w", err)
	}
	if bot == nil {
		return nil, domain.ErrInvalidAPIKey
	}
	if bot.Status != domain.BotStatusActive {
		return nil, domain.ErrBotInactive
	}

	// Cache the mapping
	if uc.cache != nil {
		_ = uc.cache.SetBotByAPIKeyHash(ctx, hash, bot.ID.String())
	}

	return bot, nil
}

// GetProfile returns a bot profile with ref_links
func (uc *BotUsecase) GetProfile(ctx context.Context, botID uuid.UUID) (*domain.BotProfileResponse, error) {
	bot, err := uc.botRepo.GetByID(ctx, botID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot: %w", err)
	}
	if bot == nil {
		return nil, domain.ErrBotNotFound
	}

	// Touch last_seen
	_ = uc.botRepo.UpdateLastSeen(ctx, botID)

	refLinks := map[string]string{
		"identity":  "/v1/me/identity",
		"bootstrap": "/v1/me/bootstrap",
		"memories":  "/v1/me/memories",
		"skills":    "/v1/me/skills",
		"policies":  "/v1/me/policies",
		"todos":     "/v1/todos",
		"heartbeat": "/v1/heartbeat",
	}

	return &domain.BotProfileResponse{
		Bot:      bot.ToResponse(),
		RefLinks: refLinks,
	}, nil
}

// GetIdentity returns the minimal identity for a bot
func (uc *BotUsecase) GetIdentity(ctx context.Context, botID uuid.UUID) (*domain.BotIdentityResponse, error) {
	bot, err := uc.botRepo.GetByID(ctx, botID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot: %w", err)
	}
	if bot == nil {
		return nil, domain.ErrBotNotFound
	}

	return &domain.BotIdentityResponse{
		ID:    bot.ID,
		Name:  bot.Name,
		Role:  bot.Role,
		Vibe:  bot.Vibe,
		Emoji: bot.Emoji,
	}, nil
}

// GetBootstrap returns the full context dump for a bot
func (uc *BotUsecase) GetBootstrap(ctx context.Context, botID uuid.UUID) (*domain.BootstrapResponse, error) {
	bot, err := uc.botRepo.GetByID(ctx, botID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot: %w", err)
	}
	if bot == nil {
		return nil, domain.ErrBotNotFound
	}

	memories, err := uc.memoryRepo.ListByBotID(ctx, botID)
	if err != nil {
		if uc.logger != nil {
			uc.logger.Warn("failed to load memories for bootstrap", zap.Error(err))
		}
		memories = []domain.BotMemory{}
	}

	skills, err := uc.skillRepo.ListByBotID(ctx, botID)
	if err != nil {
		if uc.logger != nil {
			uc.logger.Warn("failed to load skills for bootstrap", zap.Error(err))
		}
		skills = []domain.BotSkill{}
	}

	policies, err := uc.policyRepo.ListByBotID(ctx, botID)
	if err != nil {
		if uc.logger != nil {
			uc.logger.Warn("failed to load policies for bootstrap", zap.Error(err))
		}
		policies = []domain.BotPolicy{}
	}

	return &domain.BootstrapResponse{
		Bot:      bot.ToResponse(),
		Memories: memories,
		Skills:   skills,
		Policies: policies,
	}, nil
}

// hashAPIKey returns the hex-encoded SHA-256 of the raw API key
func hashAPIKey(rawKey string) string {
	rawKey = strings.TrimPrefix(rawKey, "Bearer ")
	h := sha256.Sum256([]byte(rawKey))
	return fmt.Sprintf("%x", h)
}
