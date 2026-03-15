package domain

import (
	"time"

	"github.com/google/uuid"
)

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
