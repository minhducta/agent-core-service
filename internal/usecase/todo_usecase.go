package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/minhducta/agent-core-service/internal/domain"
	"github.com/minhducta/agent-core-service/pkg/kafka"
	"go.uber.org/zap"
)

// TodoUsecase handles todo business logic
type TodoUsecase struct {
	todoRepo      domain.TodoRepository
	checklistRepo domain.TodoChecklistRepository
	producer      *kafka.Producer
	logger        *zap.Logger
}

// NewTodoUsecase creates a new TodoUsecase
func NewTodoUsecase(
	todoRepo domain.TodoRepository,
	checklistRepo domain.TodoChecklistRepository,
	producer *kafka.Producer,
	logger *zap.Logger,
) *TodoUsecase {
	return &TodoUsecase{
		todoRepo:      todoRepo,
		checklistRepo: checklistRepo,
		producer:      producer,
		logger:        logger,
	}
}

// ListTodos returns paginated todos for the calling bot
func (uc *TodoUsecase) ListTodos(ctx context.Context, botID uuid.UUID, page, limit int) (*domain.TodoListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	todos, total, err := uc.todoRepo.ListByBotID(ctx, botID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	return &domain.TodoListResponse{
		Data: todos,
		Meta: domain.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateTodo applies a partial update to a todo
func (uc *TodoUsecase) UpdateTodo(ctx context.Context, todoID uuid.UUID, req domain.UpdateTodoRequest) (*domain.Todo, error) {
	todo, err := uc.todoRepo.GetByID(ctx, todoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}
	if todo == nil {
		return nil, domain.ErrTodoNotFound
	}

	if req.Status != nil {
		todo.Status = *req.Status
	}
	if req.Priority != nil {
		todo.Priority = *req.Priority
	}
	if req.Result != nil {
		todo.Result = *req.Result
	}
	if req.DueDate != nil {
		todo.DueDate = req.DueDate
	}

	if err := uc.todoRepo.Update(ctx, todo); err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	if uc.producer != nil {
		_ = uc.producer.Publish(domain.EventTodoUpdated, map[string]interface{}{
			"todoId": todoID,
			"status": todo.Status,
		})
	}

	return todo, nil
}

// GetChecklist returns checklist items for a todo
func (uc *TodoUsecase) GetChecklist(ctx context.Context, todoID uuid.UUID) ([]domain.TodoChecklistItem, error) {
	// Verify todo exists
	todo, err := uc.todoRepo.GetByID(ctx, todoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}
	if todo == nil {
		return nil, domain.ErrTodoNotFound
	}

	return uc.checklistRepo.ListByTodoID(ctx, todoID)
}

// UpdateChecklistItem applies a partial update to a checklist item
func (uc *TodoUsecase) UpdateChecklistItem(ctx context.Context, todoID uuid.UUID, itemID uuid.UUID, req domain.UpdateChecklistItemRequest) (*domain.TodoChecklistItem, error) {
	item, err := uc.checklistRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get checklist item: %w", err)
	}
	if item == nil || item.TodoID != todoID {
		return nil, domain.ErrChecklistItemNotFound
	}

	if req.IsChecked != nil {
		item.IsChecked = *req.IsChecked
	}
	if req.Content != nil {
		item.Content = *req.Content
	}

	if err := uc.checklistRepo.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to update checklist item: %w", err)
	}

	return item, nil
}
