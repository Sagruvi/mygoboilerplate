package app

import (
	"log"
	_ "main/docs"
	"main/internal/provider"
	"main/internal/service"
	"os"
)

func Run() {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret"
	}
	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8080"
	}
	serverurl := os.Getenv("SERVER_URL")
	if serverurl == "" {
		serverurl = "userservice:15002"
	}
	provider.Run(port)
	service.NewService(serverurl)
	log.Println("Сервис авторизации запущен на порту " + port)
}
