package httpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mozhdekzm/heli-task/config"
	"github.com/mozhdekzm/heli-task/internal/application"
	"net/http"
)

type Server struct {
	config      config.Config
	todoService application.TodoService
}

func New(cfg config.Config, todoService application.TodoService) Server {
	return Server{
		config:      cfg,
		todoService: todoService,
	}
}

func (s Server) Serv() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Health Check Route
	e.GET("/health/check", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Business Routes
	todoHandler := NewTodoHandler(s.todoService)
	e.POST("/todo", todoHandler.Create)
	e.GET("/todo", todoHandler.List)

	e.Logger.Fatal(e.Start(":" + s.config.ServerPort))
}
