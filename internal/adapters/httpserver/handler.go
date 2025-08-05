package httpserver

import (
	"github.com/mozhdekzm/heli-task/internal/application"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// TodoHandler struct holds the service dependency
type TodoHandler struct {
	svc application.TodoService
}

func NewTodoHandler(svc application.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

func (h *TodoHandler) Create(c echo.Context) error {
	var req CreateTodoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid date format, use YYYY-MM-DD"})
	}

	todo, err := h.svc.Create(req.Description, dueDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	res := GetTodoResponse(todo)
	return c.JSON(http.StatusCreated, res)
}

func (h *TodoHandler) List(c echo.Context) error {
	todos, err := h.svc.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch todos"})
	}
	return c.JSON(http.StatusOK, todos)
}
