package app

import (
	"log"
	_ "main/docs"
	"main/internal/provider"
	"os"
)

func Run() {
	clientPort := os.Getenv("USER_SERVICE_PORT")
	if clientPort == "" {
		clientPort = "15002"
	}
	serverPort := os.Getenv("AUTH_SERVICE_PORT")
	if serverPort == "" {
		serverPort = "15001"
	}
	prov := provider.NewProvider(clientPort, serverPort)
	prov.Run()
	log.Println("Сервис авторизации запущен на порту " + serverPort)
}
