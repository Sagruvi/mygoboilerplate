package app

import (
	"log"
	_ "main/docs"
	"main/internal/gateway"
	"main/internal/service"
	"os"
)

func Run() {
	port := os.Getenv("GEO_SERVICE_PORT")
	if port == "" {
		port = "15003"
	}
	err := gateway.NewGRPCGateway(service.NewService()).Run(port)
	if err != nil {
		panic(err)
	}
	log.Println("Геосервис запущен на порту " + port)
}
