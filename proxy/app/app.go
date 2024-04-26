package app

import (
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"main/internal/consumer"
	"net/http"

	"os"
)

func Run() {
	port := os.Getenv("PROXY_PORT")
	if port == "" {
		port = "8080"
	}
	geoport := os.Getenv("GEOSERVICE_PORT")
	if geoport == "" {
		geoport = "1234"
	}
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://proxy:"+port+"/swagger/doc.json"),
		))
	})
	r.Group(func(r chi.Router) {
		r.Get("/geoservice/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://geoservice:"+geoport+"/swagger/doc.json"),
		))
	})
	http.Handle("/", r)

	log.Printf("Starting proxy server on port %s\n", port)
	consumer.NewGeoConsumer(geoport)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start proxy server: %v", err)
	}
}
