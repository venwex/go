package store

import (
	e "example/test/internal/errors"
	m "example/test/internal/models"
	"fmt"
	"sync"
)

type Store struct {
	mu     sync.Mutex
	nextID int
	m      map[int]m.Task
}

func NewStore() *Store {
	return &Store{
		nextID: 0,
		m:      make(map[int]m.Task),
	}
}

func (s *Store) GetTask(id int) (m.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id <= 0 {
		return m.Task{}, e.ErrInvalidID
	}

	task, ok := s.m[id]
	if !ok {
		return m.Task{}, e.ErrTaskNotFound
	}

	return task, nil
}

func (s *Store) GetTasks() []m.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	var tasks []m.Task
	for _, v := range s.m {
		tasks = append(tasks, v)
	}

	return tasks
}

func (s *Store) CreateTask(title string) (m.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	id := s.nextID

	_, ok := s.m[id]
	if ok {
		return m.Task{}, fmt.Errorf("%d id already exists", id)
	}

	task := m.NewTask(id, title)
	s.m[id] = task

	return task, nil
}

func (s *Store) MarkDoneTask(id int, done bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.m[id]
	if !ok {
		return e.ErrTaskNotFound
	}

	task.Done = done
	s.m[id] = task

	return nil
}

func (s *Store) DeleteTask(id int) (m.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.m[id]
	if !ok {
		return m.Task{}, e.ErrTaskNotFound
	}

	delete(s.m, id)

	return task, nil
}
