package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/minhducta/agent-core-service/internal/domain"
)

// BotUsecase defines the business logic interface for bot identity operations
type BotUsecase interface {
	GetProfile(ctx context.Context, botID uuid.UUID) (*domain.BotProfileResponse, error)
	GetIdentity(ctx context.Context, botID uuid.UUID) (*domain.BotIdentityResponse, error)
	GetBootstrap(ctx context.Context, botID uuid.UUID) (*domain.BootstrapResponse, error)
}

// MemoryUsecase defines the business logic interface for memory operations
type MemoryUsecase interface {
	ListMemories(ctx context.Context, botID uuid.UUID) ([]domain.BotMemory, error)
	CreateMemory(ctx context.Context, botID uuid.UUID, req domain.CreateMemoryRequest) (*domain.BotMemory, error)
	DeleteMemory(ctx context.Context, botID uuid.UUID, memoryID uuid.UUID) error
}

// SkillUsecase defines the business logic interface for skill operations
type SkillUsecase interface {
	ListSkills(ctx context.Context, botID uuid.UUID) ([]domain.BotSkill, error)
}

// PolicyUsecase defines the business logic interface for policy operations
type PolicyUsecase interface {
	ListPolicies(ctx context.Context, botID uuid.UUID) ([]domain.BotPolicy, error)
}

// TodoUsecase defines the business logic interface for todo operations
type TodoUsecase interface {
	ListTodos(ctx context.Context, botID uuid.UUID, page, limit int) (*domain.TodoListResponse, error)
	UpdateTodo(ctx context.Context, todoID uuid.UUID, req domain.UpdateTodoRequest) (*domain.Todo, error)
	GetChecklist(ctx context.Context, todoID uuid.UUID) ([]domain.TodoChecklistItem, error)
	UpdateChecklistItem(ctx context.Context, todoID uuid.UUID, itemID uuid.UUID, req domain.UpdateChecklistItemRequest) (*domain.TodoChecklistItem, error)
}

// HeartbeatUsecase defines the business logic interface for heartbeat operations
type HeartbeatUsecase interface {
	RecordHeartbeat(ctx context.Context, botID uuid.UUID, req domain.HeartbeatRequest) (*domain.Heartbeat, error)
	GetStatus(ctx context.Context, botID uuid.UUID) (*domain.HeartbeatStatusResponse, error)
}
