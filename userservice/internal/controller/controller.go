package controller

import (
	"encoding/json"
	"main/internal/service"
	"net/http"
)

type Controllerer interface {
	GetProfile(w http.ResponseWriter, r *http.Request) error
	ListUsers(w http.ResponseWriter, r *http.Request) error
}
type Controller struct {
	service.Servicer
}

func NewController() Controllerer {
	return &Controller{
		service.NewService(),
	}
}

// GetProfile godoc
//
// @Summary Get user profile by email and password
// @Description Retrieves user profile based on provided email and password
// @Tags users
// @Accept json
// @Produce json
// @Param email body string true "User's email"
// @Param password body string true "User's password"
// @Success 200 {object} User "Successful operation"
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /profile [post]
func (c *Controller) GetProfile(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	user, err := c.Servicer.GetUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	return nil
}

// ListUsers godoc
//
// @Summary List all users
// @Description Retrieves a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} User "Successful operation"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /users [get]
func (c *Controller) ListUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := c.Servicer.ListUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	return nil
}
