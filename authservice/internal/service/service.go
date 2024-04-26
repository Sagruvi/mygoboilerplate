package service

import (
	"main/internal/entity"
	"main/internal/userconsumer"
	pb "main/proto"
)

type Servicer interface {
	SaveUser(Email, password string) (entity.User, error)
	CheckUser(Email, password string) error
}
type Service struct {
	userconsumer.UserConsumer
}

func NewService(serverurl string) *Service {
	return &Service{
		userconsumer.NewUserConsumer(serverurl),
	}
}

func (s *Service) SaveUser(Email, password string) (entity.User, error) {
	u := pb.User{
		Email:    Email,
		Password: password,
	}
	user, err := s.UserConsumer.CreateUser(&u)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (s *Service) CheckUser(Email, password string) error {
	u := pb.AuthOrLogin{
		Email:    Email,
		Password: password,
	}
	return s.UserConsumer.CheckUser(&u)
}
