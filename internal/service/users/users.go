package users

import (
	m "example/test/internal/models"
	"example/test/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ServiceGetUsers() ([]m.User, error) {
	return s.repo.GetUsers()
}

func (s *UserService) ServiceGetUser(id int) (m.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) ServiceCreateUser(user m.User) (m.User, error) {
	return s.repo.CreateUser(user)
}

func (s *UserService) ServiceUpdateUser(id int, name, email string) (m.User, error) {
	return s.repo.UpdateUser(id, name, email)
}

func (s *UserService) ServiceDeleteUser(id int) (m.User, error) {
	return s.repo.DeleteUser(id)
}
