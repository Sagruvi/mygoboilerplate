package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/jwtauth"
	"mygoboilerplate/internal/auth/entity"
	"mygoboilerplate/internal/auth/repository"
	"mygoboilerplate/internal/auth/service"
	"net/http"
)

type GeoServicer interface {
	Login(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
}
type Controller struct {
	service   *service.Service
	TokenAuth *jwtauth.JWTAuth
}

func NewController(secret string) *Controller {
	return &Controller{
		service:   service.NewService(repository.NewRepository()),
		TokenAuth: jwtauth.New("HS256", []byte(secret), nil),
	}
}

// Login godoc
// @Summary login users
// @Description login users using JWT tokens
// @Tags auth
// @Accept  json
// @Produce  json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Success 200 {string} string "valid JWT token"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "error creating token"
// @Router /login [get]
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !c.service.CheckUser(user.Username, user.Password) {
		http.Error(w, "entity is not registered", http.StatusUnauthorized)
		return
	}
	_, tokenString, err := c.TokenAuth.
		Encode(map[string]interface{}{"username": user.Username, "password": user.Password})
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	} else {
		w.Write([]byte("token " + tokenString))
		w.WriteHeader(http.StatusOK)
		return
	}
}

// Register godoc
// @Summary register users
// @Description Register users using JWT tokens
// @Tags auth
// @Accept  json
// @Produce  json
//
//	@Param			Username and Password				body		string	true	"Username and Password"
//
// @Success 200 {string} string "entity is registered"
// @Failure 401 {string} string "error taking a claims"
// @Failure 500 {string} string "entity is already exists"
// @Failure 500 {string} string "error hashing password"
// @Router /register [get]
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = c.service.SaveUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("entity %s is registered: ", user.Username)))
}
