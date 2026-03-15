package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/minhducta/agent-core-service/internal/domain"
	"github.com/minhducta/agent-core-service/internal/usecase"
)

// TodoHandler handles /v1/todos routes
type TodoHandler struct {
	todoUC *usecase.TodoUsecase
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(todoUC *usecase.TodoUsecase) *TodoHandler {
	return &TodoHandler{todoUC: todoUC}
}

// ListTodos returns paginated todos for the calling bot
// GET /v1/todos
func (h *TodoHandler) ListTodos(c *fiber.Ctx) error {
	botID, err := parseBotID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errResponse(domain.ErrCodeUnauthorized, "unauthorized"))
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	resp, err := h.todoUC.ListTodos(c.Context(), botID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse(domain.ErrCodeInternal, "failed to list todos"))
	}

	return c.JSON(resp)
}

// UpdateTodo applies a partial update to a todo
// PATCH /v1/todos/:id
func (h *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	_, err := parseBotID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errResponse(domain.ErrCodeUnauthorized, "unauthorized"))
	}

	todoID, err := parsePathUUID(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errResponse(domain.ErrCodeValidation, "invalid todo id"))
	}

	var req domain.UpdateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errResponse(domain.ErrCodeValidation, "invalid request body"))
	}

	todo, err := h.todoUC.UpdateTodo(c.Context(), todoID, req)
	if err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(errResponse(domain.ErrCodeNotFound, err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse(domain.ErrCodeInternal, "failed to update todo"))
	}

	return c.JSON(fiber.Map{"data": todo})
}

// GetChecklist returns checklist items for a todo
// GET /v1/todos/:id/checklist
func (h *TodoHandler) GetChecklist(c *fiber.Ctx) error {
	_, err := parseBotID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errResponse(domain.ErrCodeUnauthorized, "unauthorized"))
	}

	todoID, err := parsePathUUID(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errResponse(domain.ErrCodeValidation, "invalid todo id"))
	}

	items, err := h.todoUC.GetChecklist(c.Context(), todoID)
	if err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(errResponse(domain.ErrCodeNotFound, err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse(domain.ErrCodeInternal, "failed to get checklist"))
	}

	return c.JSON(fiber.Map{"data": items})
}

// UpdateChecklistItem applies a partial update to a checklist item
// PATCH /v1/todos/:id/checklist/:item_id
func (h *TodoHandler) UpdateChecklistItem(c *fiber.Ctx) error {
	_, err := parseBotID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errResponse(domain.ErrCodeUnauthorized, "unauthorized"))
	}

	todoID, err := parsePathUUID(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errResponse(domain.ErrCodeValidation, "invalid todo id"))
	}

	itemID, err := parsePathUUID(c, "item_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errResponse(domain.ErrCodeValidation, "invalid checklist item id"))
	}

	var req domain.UpdateChecklistItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errResponse(domain.ErrCodeValidation, "invalid request body"))
	}

	item, err := h.todoUC.UpdateChecklistItem(c.Context(), todoID, itemID, req)
	if err != nil {
		if errors.Is(err, domain.ErrChecklistItemNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(errResponse(domain.ErrCodeNotFound, err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errResponse(domain.ErrCodeInternal, "failed to update checklist item"))
	}

	return c.JSON(fiber.Map{"data": item})
}
