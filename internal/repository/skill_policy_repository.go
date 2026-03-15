package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/minhducta/agent-core-service/internal/domain"
	"github.com/minhducta/agent-core-service/pkg/database"
)

// BotSkillRepository implements domain.BotSkillRepository
type BotSkillRepository struct {
	db *database.DB
}

// NewBotSkillRepository creates a new BotSkillRepository
func NewBotSkillRepository(db *database.DB) *BotSkillRepository {
	return &BotSkillRepository{db: db}
}

// ListByBotID retrieves all skills for a bot
func (r *BotSkillRepository) ListByBotID(ctx context.Context, botID uuid.UUID) ([]domain.BotSkill, error) {
	query := `
		SELECT id, bot_id, name, description, usage_guide, created_at, updated_at
		FROM bot_skills
		WHERE bot_id = $1
		ORDER BY name ASC`

	var skills []domain.BotSkill
	if err := r.db.SelectContext(ctx, &skills, query, botID); err != nil {
		return nil, fmt.Errorf("repo.ListByBotID: %w", err)
	}

	return skills, nil
}

// BotPolicyRepository implements domain.BotPolicyRepository
type BotPolicyRepository struct {
	db *database.DB
}

// NewBotPolicyRepository creates a new BotPolicyRepository
func NewBotPolicyRepository(db *database.DB) *BotPolicyRepository {
	return &BotPolicyRepository{db: db}
}

// ListByBotID retrieves all policies for a bot
func (r *BotPolicyRepository) ListByBotID(ctx context.Context, botID uuid.UUID) ([]domain.BotPolicy, error) {
	query := `
		SELECT id, bot_id, action, effect, conditions, created_at, updated_at
		FROM bot_policies
		WHERE bot_id = $1
		ORDER BY action ASC`

	var policies []domain.BotPolicy
	if err := r.db.SelectContext(ctx, &policies, query, botID); err != nil {
		return nil, fmt.Errorf("repo.ListByBotID: %w", err)
	}

	return policies, nil
}
