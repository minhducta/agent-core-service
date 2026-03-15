package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/minhducta/agent-core-service/internal/domain"
	"github.com/minhducta/agent-core-service/pkg/kafka"
	"go.uber.org/zap"
)

// HeartbeatUsecase handles heartbeat business logic
type HeartbeatUsecase struct {
	heartbeatRepo domain.HeartbeatRepository
	botRepo       domain.BotRepository
	producer      *kafka.Producer
	logger        *zap.Logger
}

// NewHeartbeatUsecase creates a new HeartbeatUsecase
func NewHeartbeatUsecase(
	heartbeatRepo domain.HeartbeatRepository,
	botRepo domain.BotRepository,
	producer *kafka.Producer,
	logger *zap.Logger,
) *HeartbeatUsecase {
	return &HeartbeatUsecase{
		heartbeatRepo: heartbeatRepo,
		botRepo:       botRepo,
		producer:      producer,
		logger:        logger,
	}
}

// RecordHeartbeat persists a heartbeat and updates bot last_seen_at
func (uc *HeartbeatUsecase) RecordHeartbeat(ctx context.Context, botID uuid.UUID, req domain.HeartbeatRequest) (*domain.Heartbeat, error) {
	var metaBytes []byte
	if req.Metadata != nil {
		b, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metaBytes = b
	}

	hb := &domain.Heartbeat{
		ID:       uuid.New(),
		BotID:    botID,
		Status:   req.Status,
		Metadata: metaBytes,
	}

	if err := uc.heartbeatRepo.Create(ctx, hb); err != nil {
		return nil, fmt.Errorf("failed to record heartbeat: %w", err)
	}

	// Update bot's last_seen_at asynchronously (best effort)
	go func() {
		bgCtx := context.Background()
		if err := uc.botRepo.UpdateLastSeen(bgCtx, botID); err != nil {
			uc.logger.Warn("failed to update last_seen_at", zap.Error(err))
		}
	}()

	_ = uc.producer.Publish(domain.EventHeartbeatReceived, map[string]interface{}{
		"botId":  botID,
		"status": req.Status,
	})

	uc.logger.Info("heartbeat recorded", zap.String("botId", botID.String()), zap.String("status", string(req.Status)))

	return hb, nil
}

// GetStatus returns the latest heartbeat status for the calling bot
func (uc *HeartbeatUsecase) GetStatus(ctx context.Context, botID uuid.UUID) (*domain.HeartbeatStatusResponse, error) {
	resp, err := uc.heartbeatRepo.GetStatusByBotID(ctx, botID)
	if err != nil {
		return nil, fmt.Errorf("failed to get heartbeat status: %w", err)
	}

	return resp, nil
}
