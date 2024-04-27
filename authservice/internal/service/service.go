package service

import (
	"context"
	"main/internal/consumer"
	"main/internal/entity"
)

type Servicer interface {
	SaveUser(user entity.User) (entity.User, error)
	CheckUser(Email, password string) (entity.User, error)
}
type Service struct {
	consumer.UserConsumer
}

func NewService(serverurl string) *Service {
	return &Service{
		consumer.NewUserConsumer(serverurl),
	}
}

func (s *Service) SaveUser(us entity.User) (entity.User, error) {
	_, err := s.UserConsumer.CreateUser(context.Background(), us)
	if err != nil {
		return entity.User{}, err
	}
	return us, nil
}

func (s *Service) CheckUser(Email, password string) (entity.User, error) {
	return s.UserConsumer.CheckUser(context.Background(), Email, password)
}
