package domain

import (
	"time"

	"github.com/google/uuid"
)

// HeartbeatStatus represents the status reported in a heartbeat
type HeartbeatStatus string

const (
	HeartbeatStatusOK      HeartbeatStatus = "ok"
	HeartbeatStatusDegraded HeartbeatStatus = "degraded"
	HeartbeatStatusError   HeartbeatStatus = "error"
)

// Heartbeat is an audit log entry for bot liveness/status reports
type Heartbeat struct {
	ID        uuid.UUID       `db:"id"        json:"id"`
	BotID     uuid.UUID       `db:"bot_id"    json:"botId"`
	Status    HeartbeatStatus `db:"status"    json:"status"`
	Metadata  []byte          `db:"metadata"  json:"metadata,omitempty"`
	CreatedAt time.Time       `db:"created_at" json:"createdAt"`
}

// HeartbeatRequest is the body for POST /v1/heartbeat
type HeartbeatRequest struct {
	Status   HeartbeatStatus        `json:"status"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// HeartbeatStatusResponse is the response for GET /v1/heartbeat/status
type HeartbeatStatusResponse struct {
	BotID      uuid.UUID       `json:"botId"`
	LastStatus HeartbeatStatus `json:"lastStatus"`
	LastSeenAt *time.Time      `json:"lastSeenAt,omitempty"`
	TotalCount int             `json:"totalCount"`
}
