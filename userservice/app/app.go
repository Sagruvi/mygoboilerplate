package app

import (
	"log"
	_ "main/docs"
	"os"
)

func Run() {
	port := os.Getenv("USER_SERVICE_PORT")
	log.Println("Сервис user запущен на порту: " + port)
}
