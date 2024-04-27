package gateway

import (
	"encoding/json"
	"main/internal/entity"
	"main/internal/service"
	"net/http"
)

type Gateway interface {
	Geocode(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
	Profile(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}
type Controller struct {
	service.Servicer
}

func NewGateway(userPort, geoPort string) Gateway {
	return &Controller{service.NewService(userPort, geoPort)}
}
func (c *Controller) Geocode(w http.ResponseWriter, r *http.Request) {
	var address entity.GeocodeRequest
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	res, err := c.Servicer.Geocode(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(res)
}
func (c *Controller) Search(w http.ResponseWriter, r *http.Request) {
	var address entity.SearchRequest
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	res, err := c.Servicer.Search(address.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(res)
}
func (c *Controller) Profile(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	res, err := c.Servicer.Get(user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(res)
}
func (c *Controller) List(w http.ResponseWriter, r *http.Request) {
	res, err := c.Servicer.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(res)
}
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	res, err := c.Servicer.Register(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(res)
}
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	res, err := c.Servicer.Login(user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(res)
}
