package service

import (
	"main/internal/consumer"
	"main/internal/entity"
)

type Servicer interface {
	Login(email, password string) (entity.User, error)
	Geocode(address entity.GeocodeRequest) (string, error)
	Search(address string) ([]entity.Address, error)
	Get(email, password string) (entity.User, error)
	Register(user entity.User) (entity.User, error)
	List() ([]entity.User, error)
}
type Service struct {
	user consumer.UserConsumer
	geo  consumer.GeoConsumer
	auth consumer.AuthConsumer
}

func NewService(userConsumerPort, geoConsumerPort, authConsumerPort string) Servicer {
	return Service{
		user: consumer.NewUserConsumer(userConsumerPort),
		geo:  consumer.NewGeoConsumer(geoConsumerPort),
		auth: consumer.NewAuthConsumer(authConsumerPort),
	}
}

func (s Service) Login(email, password string) (entity.User, error) {
	return s.auth.Login(email, password)
}
func (s Service) Geocode(address entity.GeocodeRequest) (string, error) {
	return s.geo.Geocode(address)
}
func (s Service) Search(address string) ([]entity.Address, error) {
	return s.geo.AddressSearch(address)
}
func (s Service) Get(email, password string) (entity.User, error) {
	return s.user.Get(email, password)
}
func (s Service) Register(user entity.User) (entity.User, error) {
	return s.auth.Register(user)

}
func (s Service) List() ([]entity.User, error) {
	return s.user.List()
}
