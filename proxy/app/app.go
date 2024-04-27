package app

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
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
	newjwt := jwtauth.New("HS256", []byte("secret"), nil)
	con := gateway.NewGateway(userport, geoport)
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Post("/api/auth/register", con.Register)
		r.Post("/api/auth/login", con.Login)
	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(newjwt))
		r.Use(jwtauth.Authenticator)
		r.Get("/api/user/profile", con.Profile)
		r.Get("/api/user/list", con.List)
	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(newjwt))
		r.Use(jwtauth.Authenticator)
		r.Get("/api/geo/geocode", con.Geocode)
		r.Get("/api/geo/search", con.Search)
	})

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
