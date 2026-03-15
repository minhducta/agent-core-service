package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/minhducta/agent-core-service/internal/domain"
	"github.com/minhducta/agent-core-service/pkg/database"
)

// BotRepository implements domain.BotRepository
type BotRepository struct {
	db *database.DB
}

// NewBotRepository creates a new BotRepository
func NewBotRepository(db *database.DB) *BotRepository {
	return &BotRepository{db: db}
}

// GetByID retrieves a bot by its UUID
func (r *BotRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Bot, error) {
	var bot domain.Bot
	query := `
		SELECT id, name, role, vibe, emoji, avatar_url, api_key_hash, last_seen_at, status, created_at, updated_at
		FROM bots
		WHERE id = $1`

	if err := r.db.GetContext(ctx, &bot, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("repo.GetByID: %w", err)
	}

	return &bot, nil
}

// GetByAPIKeyHash retrieves a bot by the SHA-256 hash of its API key
func (r *BotRepository) GetByAPIKeyHash(ctx context.Context, hash string) (*domain.Bot, error) {
	var bot domain.Bot
	query := `
		SELECT id, name, role, vibe, emoji, avatar_url, api_key_hash, last_seen_at, status, created_at, updated_at
		FROM bots
		WHERE api_key_hash = $1`

	if err := r.db.GetContext(ctx, &bot, query, hash); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("repo.GetByAPIKeyHash: %w", err)
	}

	return &bot, nil
}

// UpdateLastSeen sets last_seen_at to now for the given bot
func (r *BotRepository) UpdateLastSeen(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE bots SET last_seen_at = $2, updated_at = $2 WHERE id = $1`
	now := time.Now()
	if _, err := r.db.ExecContext(ctx, query, id, now); err != nil {
		return fmt.Errorf("repo.UpdateLastSeen: %w", err)
	}
	return nil
}

// ExecTx executes a function within a database transaction
func (r *BotRepository) ExecTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("repo.ExecTx: begin: %w", err)
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("repo.ExecTx: tx err: %v, rollback err: %w", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
