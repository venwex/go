package service

import (
	"example/test/internal/repository"
	"example/test/internal/service/tasks"
	"example/test/internal/service/users"
)

type Services struct {
	Task *tasks.TaskService
	User *users.UserService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Task: tasks.NewService(repos.Task),
		User: users.NewService(repos.User),
	}
}
