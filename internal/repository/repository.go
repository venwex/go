package repository

import (
	m "example/test/internal/models"
	"example/test/internal/repository/postgres"
	"example/test/internal/repository/postgres/tasks"
	"example/test/internal/repository/postgres/users"
)

type UserRepository interface {
	GetUsers() ([]m.User, error)
	GetUserByID(id int) (m.User, error)
	CreateUser(user m.User) (m.User, error)
	UpdateUser(id int, name, email string) (m.User, error)
	DeleteUser(id int) (m.User, error)
}

type TaskRepository interface {
	GetTask(id int) (m.Task, error)
	GetTasks() ([]m.Task, error)
	CreateTask(title string) (m.Task, error)
	MarkDoneTask(id int, done bool) error
	DeleteTask(id int) (m.Task, error)
}

type Repositories struct {
	User UserRepository
	Task TaskRepository
}

func NewRepositories(db *postgres.Dialect) *Repositories {
	return &Repositories{
		User: users.NewRepository(db),
		Task: tasks.NewRepository(db),
	}
}
