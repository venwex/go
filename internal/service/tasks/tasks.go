package tasks

import (
	m "example/test/internal/models"
	"example/test/internal/repository"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) ServiceGetTask(id int) (m.Task, error) {
	return s.repo.GetTask(id)
}

func (s *TaskService) ServiceGetTasks() ([]m.Task, error) {
	return s.repo.GetTasks()
}

func (s *TaskService) ServiceCreateTask(title string) (m.Task, error) {
	return s.repo.CreateTask(title)
}

func (s *TaskService) ServiceMarkDoneTask(id int, done bool) error {
	return s.repo.MarkDoneTask(id, done)
}

func (s *TaskService) ServiceDeleteTask(id int) (m.Task, error) {
	return s.repo.DeleteTask(id)
}
