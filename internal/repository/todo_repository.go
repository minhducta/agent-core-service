package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/minhducta/agent-core-service/internal/domain"
	"github.com/minhducta/agent-core-service/pkg/database"
)

// TodoRepository implements domain.TodoRepository
type TodoRepository struct {
	db *database.DB
}

// NewTodoRepository creates a new TodoRepository
func NewTodoRepository(db *database.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// ListByBotID retrieves paginated todos assigned to or by a bot
func (r *TodoRepository) ListByBotID(ctx context.Context, botID uuid.UUID, page, limit int) ([]domain.Todo, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM todos WHERE assigned_to = $1 OR assigned_by = $1`
	if err := r.db.GetContext(ctx, &total, countQuery, botID); err != nil {
		return nil, 0, fmt.Errorf("repo.ListByBotID: count: %w", err)
	}

	offset := (page - 1) * limit
	query := `
		SELECT id, title, description, status, priority, result, due_date,
		       assigned_to, assigned_by, dependency_id, created_at, updated_at
		FROM todos
		WHERE assigned_to = $1 OR assigned_by = $1
		ORDER BY priority DESC, created_at DESC
		LIMIT $2 OFFSET $3`

	var todos []domain.Todo
	if err := r.db.SelectContext(ctx, &todos, query, botID, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("repo.ListByBotID: %w", err)
	}

	return todos, total, nil
}

// GetByID retrieves a single todo
func (r *TodoRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	var todo domain.Todo
	query := `
		SELECT id, title, description, status, priority, result, due_date,
		       assigned_to, assigned_by, dependency_id, created_at, updated_at
		FROM todos
		WHERE id = $1`

	if err := r.db.GetContext(ctx, &todo, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("repo.GetByID: %w", err)
	}

	return &todo, nil
}

// Update persists changes to a todo record
func (r *TodoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	query := `
		UPDATE todos
		SET status = $2, priority = $3, result = $4, due_date = $5, updated_at = $6
		WHERE id = $1`

	todo.UpdatedAt = time.Now()
	if _, err := r.db.ExecContext(ctx, query,
		todo.ID, todo.Status, todo.Priority, todo.Result, todo.DueDate, todo.UpdatedAt,
	); err != nil {
		return fmt.Errorf("repo.Update: %w", err)
	}

	return nil
}

// TodoChecklistRepository implements domain.TodoChecklistRepository
type TodoChecklistRepository struct {
	db *database.DB
}

// NewTodoChecklistRepository creates a new TodoChecklistRepository
func NewTodoChecklistRepository(db *database.DB) *TodoChecklistRepository {
	return &TodoChecklistRepository{db: db}
}

// ListByTodoID retrieves all checklist items for a todo ordered by order_index
func (r *TodoChecklistRepository) ListByTodoID(ctx context.Context, todoID uuid.UUID) ([]domain.TodoChecklistItem, error) {
	query := `
		SELECT id, todo_id, content, is_checked, order_index, created_at, updated_at
		FROM todo_checklist_items
		WHERE todo_id = $1
		ORDER BY order_index ASC`

	var items []domain.TodoChecklistItem
	if err := r.db.SelectContext(ctx, &items, query, todoID); err != nil {
		return nil, fmt.Errorf("repo.ListByTodoID: %w", err)
	}

	return items, nil
}

// GetByID retrieves a single checklist item
func (r *TodoChecklistRepository) GetByID(ctx context.Context, itemID uuid.UUID) (*domain.TodoChecklistItem, error) {
	var item domain.TodoChecklistItem
	query := `
		SELECT id, todo_id, content, is_checked, order_index, created_at, updated_at
		FROM todo_checklist_items
		WHERE id = $1`

	if err := r.db.GetContext(ctx, &item, query, itemID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("repo.GetByID: %w", err)
	}

	return &item, nil
}

// Update persists changes to a checklist item
func (r *TodoChecklistRepository) Update(ctx context.Context, item *domain.TodoChecklistItem) error {
	query := `
		UPDATE todo_checklist_items
		SET content = $2, is_checked = $3, updated_at = $4
		WHERE id = $1`

	item.UpdatedAt = time.Now()
	if _, err := r.db.ExecContext(ctx, query, item.ID, item.Content, item.IsChecked, item.UpdatedAt); err != nil {
		return fmt.Errorf("repo.Update: %w", err)
	}

	return nil
}
