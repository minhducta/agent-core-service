package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/minhducta/agent-core-service/internal/domain"
	"github.com/minhducta/agent-core-service/pkg/database"
)

// BotMemoryRepository implements domain.BotMemoryRepository
type BotMemoryRepository struct {
	db *database.DB
}

// NewBotMemoryRepository creates a new BotMemoryRepository
func NewBotMemoryRepository(db *database.DB) *BotMemoryRepository {
	return &BotMemoryRepository{db: db}
}

// ListByBotID retrieves all memories for a bot ordered by importance desc
func (r *BotMemoryRepository) ListByBotID(ctx context.Context, botID uuid.UUID) ([]domain.BotMemory, error) {
	query := `
		SELECT id, bot_id, type, content, tags, importance, expires_at, created_at, updated_at
		FROM bot_memories
		WHERE bot_id = $1
		  AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY importance DESC, created_at DESC`

	var memories []domain.BotMemory
	if err := r.db.SelectContext(ctx, &memories, query, botID); err != nil {
		return nil, fmt.Errorf("repo.ListByBotID: %w", err)
	}

	return memories, nil
}

// Create inserts a new memory record
func (r *BotMemoryRepository) Create(ctx context.Context, memory *domain.BotMemory) error {
	query := `
		INSERT INTO bot_memories (id, bot_id, type, content, tags, importance, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	now := time.Now()
	memory.CreatedAt = now
	memory.UpdatedAt = now

	if _, err := r.db.ExecContext(ctx, query,
		memory.ID, memory.BotID, memory.Type, memory.Content,
		pq.Array(memory.Tags), memory.Importance, memory.ExpiresAt,
		memory.CreatedAt, memory.UpdatedAt,
	); err != nil {
		return fmt.Errorf("repo.Create: %w", err)
	}

	return nil
}

// Delete removes a memory by id, scoped to the owning bot
func (r *BotMemoryRepository) Delete(ctx context.Context, id uuid.UUID, botID uuid.UUID) error {
	query := `DELETE FROM bot_memories WHERE id = $1 AND bot_id = $2`

	result, err := r.db.ExecContext(ctx, query, id, botID)
	if err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repo.Delete: rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrMemoryNotFound
	}

	return nil
}
