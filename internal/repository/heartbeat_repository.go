package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/minhducta/agent-core-service/internal/domain"
	"github.com/minhducta/agent-core-service/pkg/database"
)

// HeartbeatRepository implements domain.HeartbeatRepository
type HeartbeatRepository struct {
	db *database.DB
}

// NewHeartbeatRepository creates a new HeartbeatRepository
func NewHeartbeatRepository(db *database.DB) *HeartbeatRepository {
	return &HeartbeatRepository{db: db}
}

// Create inserts a heartbeat audit record
func (r *HeartbeatRepository) Create(ctx context.Context, hb *domain.Heartbeat) error {
	query := `
		INSERT INTO heartbeats (id, bot_id, status, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	hb.CreatedAt = time.Now()
	if _, err := r.db.ExecContext(ctx, query,
		hb.ID, hb.BotID, hb.Status, hb.Metadata, hb.CreatedAt,
	); err != nil {
		return fmt.Errorf("repo.Create: %w", err)
	}

	return nil
}

// GetStatusByBotID retrieves the latest heartbeat status for a bot
func (r *HeartbeatRepository) GetStatusByBotID(ctx context.Context, botID uuid.UUID) (*domain.HeartbeatStatusResponse, error) {
	query := `
		SELECT
			bot_id,
			status,
			created_at AS last_seen_at,
			(SELECT COUNT(*) FROM heartbeats WHERE bot_id = $1) AS total_count
		FROM heartbeats
		WHERE bot_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	type row struct {
		BotID      uuid.UUID              `db:"bot_id"`
		Status     domain.HeartbeatStatus `db:"status"`
		LastSeenAt time.Time              `db:"last_seen_at"`
		TotalCount int                    `db:"total_count"`
	}

	var r2 row
	if err := r.db.GetContext(ctx, &r2, query, botID); err != nil {
		return nil, fmt.Errorf("repo.GetStatusByBotID: %w", err)
	}

	lastSeen := r2.LastSeenAt
	return &domain.HeartbeatStatusResponse{
		BotID:      r2.BotID,
		LastStatus: r2.Status,
		LastSeenAt: &lastSeen,
		TotalCount: r2.TotalCount,
	}, nil
}
