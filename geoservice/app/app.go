package app

import (
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	_ "main/docs"
	"main/internal/gateway"
	"main/internal/service"
	"os"
)

func Run() {
	port := os.Getenv("GEO_SERVICE_PORT")
	if port == "" {
		port = "1234"
	}

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://geoservice:"+port+"/swagger/doc.json"),
		))
	})
	err := gateway.NewGRPCGateway(service.NewService()).Run(port)
	if err != nil {
		panic(err)
	}
	log.Println("Геосервис запущен на порту " + port)
}
