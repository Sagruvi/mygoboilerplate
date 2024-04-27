package app

import (
	"fmt"
	_ "main/docs"
	"main/internal/userprovider"
	"time"
)

func Run() {
	//port := os.Getenv("USER_SERVICE_PORT")
	provider := userprovider.NewUserProvider()
	go func() {
		err := provider.Run("15002")
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(20 * time.Second)
	fmt.Println("Сервис user запущен на порту: " + "15002")
	select {}
}
