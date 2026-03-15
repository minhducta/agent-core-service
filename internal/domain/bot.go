package domain

import (
	"time"

	"github.com/google/uuid"
)

// BotStatus represents the operational status of a bot
type BotStatus string

const (
	BotStatusActive   BotStatus = "active"
	BotStatusInactive BotStatus = "inactive"
	BotStatusBanned   BotStatus = "banned"
)

// Bot represents a bot identity in the system
type Bot struct {
	ID         uuid.UUID  `db:"id"           json:"id"`
	Name       string     `db:"name"         json:"name"`
	Role       string     `db:"role"         json:"role"`
	Vibe       string     `db:"vibe"         json:"vibe"`
	Emoji      string     `db:"emoji"        json:"emoji"`
	AvatarURL  string     `db:"avatar_url"   json:"avatarUrl"`
	APIKeyHash string     `db:"api_key_hash" json:"-"`
	LastSeenAt *time.Time `db:"last_seen_at" json:"lastSeenAt,omitempty"`
	Status     BotStatus  `db:"status"       json:"status"`
	CreatedAt  time.Time  `db:"created_at"   json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at"   json:"updatedAt"`
}

// BotResponse is a safe response DTO (no api_key_hash)
type BotResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Role       string     `json:"role"`
	Vibe       string     `json:"vibe"`
	Emoji      string     `json:"emoji"`
	AvatarURL  string     `json:"avatarUrl"`
	LastSeenAt *time.Time `json:"lastSeenAt,omitempty"`
	Status     BotStatus  `json:"status"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// BotProfileResponse is the full response for GET /v1/me — includes ref_links
type BotProfileResponse struct {
	Bot      BotResponse       `json:"bot"`
	RefLinks map[string]string `json:"refLinks"`
}

// BotIdentityResponse is the minimal response for GET /v1/me/identity
type BotIdentityResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Role  string    `json:"role"`
	Vibe  string    `json:"vibe"`
	Emoji string    `json:"emoji"`
}

// ToResponse converts Bot to BotResponse (omits sensitive fields)
func (b *Bot) ToResponse() BotResponse {
	return BotResponse{
		ID:         b.ID,
		Name:       b.Name,
		Role:       b.Role,
		Vibe:       b.Vibe,
		Emoji:      b.Emoji,
		AvatarURL:  b.AvatarURL,
		LastSeenAt: b.LastSeenAt,
		Status:     b.Status,
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  b.UpdatedAt,
	}
}
