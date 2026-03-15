package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// BotRepository defines data access for the bots table
type BotRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Bot, error)
	GetByAPIKeyHash(ctx context.Context, hash string) (*Bot, error)
	UpdateLastSeen(ctx context.Context, id uuid.UUID) error
	ExecTx(ctx context.Context, fn func(*sqlx.Tx) error) error
}

// BotMemoryRepository defines data access for the bot_memories table
type BotMemoryRepository interface {
	ListByBotID(ctx context.Context, botID uuid.UUID) ([]BotMemory, error)
	Create(ctx context.Context, memory *BotMemory) error
	Delete(ctx context.Context, id uuid.UUID, botID uuid.UUID) error
}

// BotSkillRepository defines data access for the bot_skills table
type BotSkillRepository interface {
	ListByBotID(ctx context.Context, botID uuid.UUID) ([]BotSkill, error)
}

// BotPolicyRepository defines data access for the bot_policies table
type BotPolicyRepository interface {
	ListByBotID(ctx context.Context, botID uuid.UUID) ([]BotPolicy, error)
}

// TodoRepository defines data access for the todos table
type TodoRepository interface {
	ListByBotID(ctx context.Context, botID uuid.UUID, page, limit int) ([]Todo, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Todo, error)
	Update(ctx context.Context, todo *Todo) error
}

// TodoChecklistRepository defines data access for the todo_checklist_items table
type TodoChecklistRepository interface {
	ListByTodoID(ctx context.Context, todoID uuid.UUID) ([]TodoChecklistItem, error)
	GetByID(ctx context.Context, itemID uuid.UUID) (*TodoChecklistItem, error)
	Update(ctx context.Context, item *TodoChecklistItem) error
}

// HeartbeatRepository defines data access for the heartbeats table
type HeartbeatRepository interface {
	Create(ctx context.Context, hb *Heartbeat) error
	GetStatusByBotID(ctx context.Context, botID uuid.UUID) (*HeartbeatStatusResponse, error)
}
