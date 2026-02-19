package handlers

import (
	"example/test/internal/handlers/tasks"
	"example/test/internal/handlers/users"
	"example/test/internal/service"
)

type Handlers struct {
	Task *tasks.TaskHandler
	User *users.UserHandler
}

func NewHandlers(service *service.Services) *Handlers {
	return &Handlers{
		Task: tasks.NewTaskHandler(service.Task),
		User: users.NewUserHandler(service.User),
	}
}
