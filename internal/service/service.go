package service

import (
	m "example/test/internal/models"
	"example/test/internal/store"
)

type Service struct {
	repo *store.Store
}

func NewService(repo *store.Store) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ServiceGetTask(id int) (m.Task, error) { // мб для валидации данных
	return s.repo.GetTask(id)
}

func (s *Service) ServiceGetTasks() []m.Task {
	return s.repo.GetTasks()
}

func (s *Service) ServiceCreateTask(title string) (m.Task, error) {
	return s.repo.CreateTask(title)
}

func (s *Service) ServiceMarkDoneTask(id int, done bool) error {
	return s.repo.MarkDoneTask(id, done)
}

func (s *Service) ServiceDeleteTask(id int) (m.Task, error) {
	return s.repo.DeleteTask(id)
}

/*
короче говоря, есть 3 слоя
store -> service -> handlers -> main (optional)

store: там описывается структура бд и все методы для взаимодействия с ним

service: создается структура сервис которая как раз хранит наш бд, этот слой служит мостом 
между store (database) and handlers (api). всю логику добавляения таска можно описать тут, а в хэндлере
использовать уже готовую функцию. в хэндлере мы лишь описываем тонкости на уровне http и ничего лишнего

handlers: тут тоже создается структура которая хранит наш сервис (и уже с помощью него взаимодействуем с 
сервисом, а через это к бд). все про http, rest api взаимодействия, самый верхний уровень. 
по идее не должен хранить в себе тяжелую или логику не относящуюся к этому уровню. 
в общем, я наконец понял про dependancy injection (DI) на практике лол.
*/