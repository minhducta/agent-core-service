package domain

import (
	"time"

	"github.com/google/uuid"
)

// TodoStatus represents the lifecycle state of a todo task
type TodoStatus string

const (
	TodoStatusPending    TodoStatus = "pending"
	TodoStatusInProgress TodoStatus = "in_progress"
	TodoStatusDone       TodoStatus = "done"
	TodoStatusCancelled  TodoStatus = "cancelled"
)

// TodoPriority represents the priority level of a task
type TodoPriority string

const (
	TodoPriorityLow    TodoPriority = "low"
	TodoPriorityMedium TodoPriority = "medium"
	TodoPriorityHigh   TodoPriority = "high"
	TodoPriorityUrgent TodoPriority = "urgent"
)

// Todo represents a task assigned to or by a bot
type Todo struct {
	ID           uuid.UUID    `db:"id"            json:"id"`
	Title        string       `db:"title"         json:"title"`
	Description  string       `db:"description"   json:"description"`
	Status       TodoStatus   `db:"status"        json:"status"`
	Priority     TodoPriority `db:"priority"      json:"priority"`
	Result       string       `db:"result"        json:"result"`
	DueDate      *time.Time   `db:"due_date"      json:"dueDate,omitempty"`
	AssignedTo   *uuid.UUID   `db:"assigned_to"   json:"assignedTo,omitempty"`
	AssignedBy   *uuid.UUID   `db:"assigned_by"   json:"assignedBy,omitempty"`
	DependencyID *uuid.UUID   `db:"dependency_id" json:"dependencyId,omitempty"`
	CreatedAt    time.Time    `db:"created_at"    json:"createdAt"`
	UpdatedAt    time.Time    `db:"updated_at"    json:"updatedAt"`
}

// UpdateTodoRequest is the request body for PATCH /v1/todos/:id
type UpdateTodoRequest struct {
	Status   *TodoStatus   `json:"status,omitempty"`
	Priority *TodoPriority `json:"priority,omitempty"`
	Result   *string       `json:"result,omitempty"`
	DueDate  *time.Time    `json:"dueDate,omitempty"`
}

// TodoChecklistItem represents a checklist item within a todo
type TodoChecklistItem struct {
	ID         uuid.UUID `db:"id"          json:"id"`
	TodoID     uuid.UUID `db:"todo_id"     json:"todoId"`
	Content    string    `db:"content"     json:"content"`
	IsChecked  bool      `db:"is_checked"  json:"isChecked"`
	OrderIndex int       `db:"order_index" json:"orderIndex"`
	CreatedAt  time.Time `db:"created_at"  json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at"  json:"updatedAt"`
}

// UpdateChecklistItemRequest is the request body for PATCH /v1/todos/:id/checklist/:item_id
type UpdateChecklistItemRequest struct {
	IsChecked *bool   `json:"isChecked,omitempty"`
	Content   *string `json:"content,omitempty"`
}

// TodoListResponse is a paginated list of todos
type TodoListResponse struct {
	Data []Todo         `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
