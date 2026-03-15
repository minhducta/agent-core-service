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

// BotSkill represents a skill/tool available to a bot
type BotSkill struct {
	ID          uuid.UUID `db:"id"          json:"id"`
	BotID       uuid.UUID `db:"bot_id"      json:"botId"`
	Name        string    `db:"name"        json:"name"`
	Description string    `db:"description" json:"description"`
	UsageGuide  string    `db:"usage_guide" json:"usageGuide"`
	CreatedAt   time.Time `db:"created_at"  json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at"  json:"updatedAt"`
}

// PolicyEffect represents the effect of a policy rule
type PolicyEffect string

const (
	PolicyEffectAllow PolicyEffect = "ALLOW"
	PolicyEffectDeny  PolicyEffect = "DENY"
)

// BotPolicy represents a permission rule for a bot
type BotPolicy struct {
	ID         uuid.UUID    `db:"id"         json:"id"`
	BotID      uuid.UUID    `db:"bot_id"     json:"botId"`
	Action     string       `db:"action"     json:"action"`
	Effect     PolicyEffect `db:"effect"     json:"effect"`
	Conditions []byte       `db:"conditions" json:"conditions,omitempty"`
	CreatedAt  time.Time    `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time    `db:"updated_at" json:"updatedAt"`
}

// BootstrapResponse is the full context dump for GET /v1/me/bootstrap
type BootstrapResponse struct {
	Bot      BotResponse  `json:"bot"`
	Memories []BotMemory  `json:"memories"`
	Skills   []BotSkill   `json:"skills"`
	Policies []BotPolicy  `json:"policies"`
}
