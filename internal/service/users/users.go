package users

import (
	m "example/test/internal/models"
	"example/test/internal/repository"
	"example/test/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ServiceGetUsers(query m.UserQuery) (m.PaginatedResponse, error) {
	return s.repo.GetUsers(query)
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

func (s *UserService) ServiceGetCommonFriends(u1, u2 int) ([]m.User, error) {
	return s.repo.GetCommonFriends(u1, u2)
}

func (s *UserService) RegisterUser(user m.User) (m.User, string, error) {
	user, err := s.repo.CreateUser(user)
	if err != nil {
		return m.User{}, "", err
	}

	sessionID := uuid.New().String()
	return user, sessionID, nil
}

func (s *UserService) SignIn(login m.LoginUserDTO) (string, error) {
	user, err := s.repo.GetUserByEmail(login.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) PromoteUser(userID int) error {
	return s.repo.PromoteUser(userID)
}
