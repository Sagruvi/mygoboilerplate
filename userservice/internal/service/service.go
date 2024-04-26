package service

import (
	"main/internal/entity"
	"main/internal/repository"
)

type Servicer interface {
	CreateUser(user entity.User) error
	GetUser(email, password string) (entity.User, error)
	ListUsers() ([]entity.User, error)
}
type Service struct {
	repository.Repositorer
}

func NewService() Service {
	return Service{
		repository.NewRepository(),
	}
}
func (s Service) GetUser(email string, password string) (entity.User, error) {
	user, err := s.Repositorer.GetUser(email, password)
	return user, err
}
func (s Service) ListUsers() ([]entity.User, error) {
	users, err := s.Repositorer.ListUsers()
	if err != nil {
		return []entity.User{}, err
	}
	return users, nil
}
func (s Service) CreateUser(user entity.User) error {
	return s.Repositorer.CreateUser(user)
}
