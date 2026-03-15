package domain

import (
	"time"

	"github.com/google/uuid"
)

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
