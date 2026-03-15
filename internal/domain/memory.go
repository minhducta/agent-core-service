package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// MemoryType represents the category of a bot memory
type MemoryType string

const (
	MemoryTypeFact        MemoryType = "fact"
	MemoryTypePreference  MemoryType = "preference"
	MemoryTypeInstruction MemoryType = "instruction"
	MemoryTypeContext     MemoryType = "context"
)

// BotMemory represents a long-term memory entry for a bot
type BotMemory struct {
	ID         uuid.UUID      `db:"id"         json:"id"`
	BotID      uuid.UUID      `db:"bot_id"     json:"botId"`
	Type       MemoryType     `db:"type"       json:"type"`
	Content    string         `db:"content"    json:"content"`
	Tags       pq.StringArray `db:"tags"       json:"tags"`
	Importance int            `db:"importance" json:"importance"`
	ExpiresAt  *time.Time     `db:"expires_at" json:"expiresAt,omitempty"`
	CreatedAt  time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time      `db:"updated_at" json:"updatedAt"`
}

// CreateMemoryRequest is the request body for POST /v1/me/memories
type CreateMemoryRequest struct {
	Type       MemoryType `json:"type"`
	Content    string     `json:"content"`
	Tags       []string   `json:"tags"`
	Importance int        `json:"importance"`
	ExpiresAt  *time.Time `json:"expiresAt,omitempty"`
}
