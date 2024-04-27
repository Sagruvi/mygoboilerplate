package app

import (
	"github.com/go-chi/chi"
	"main/internal/gateway"
	"net/http"
	"os"
)

func Run() {
	port := os.Getenv("PROXY_PORT")
	if port == "" {
		port = "8080"
	}
	userport := os.Getenv("USER_SERVICE_PORT")
	if userport == "" {
		userport = "userservice:15002"
	}
	geoport := os.Getenv("GEO_SERVICE_PORT")
	if geoport == "" {
		geoport = "geoservice:15003"
	}
	authport := os.Getenv("AUTH_SERVICE_PORT")
	if authport == "" {
		authport = "authservice:15001"
	}
	//newjwt := jwtauth.New("HS256", []byte("secret"), nil)
	con := gateway.NewGateway(userport, geoport, authport)
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Post("/api/auth/register", con.Register)
		r.Post("/api/auth/login", con.Login)
	})
	r.Group(func(r chi.Router) {
		r.Get("/api/user/profile", con.Profile)
		r.Get("/api/user/list", con.List)
	})
	r.Group(func(r chi.Router) {
		r.Get("/api/geo/geocode", con.Geocode)
		r.Get("/api/geo/search", con.Search)
	})

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
