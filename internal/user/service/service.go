package service

import (
	"mygoboilerplate/internal/geolocation/repository"
	"mygoboilerplate/internal/user/entity"
)

type Servicer interface {
	GetUser(id int) (entity.User, error)
	ListUsers() []entity.User
}
type Service struct {
	repository.Repositorer
}
